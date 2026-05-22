//go:build integration
// +build integration

// Package-level coverage note
//
// These integration tests use a dedicated fleet account (NEW_RELIC_FLEET_TEST_*)
// and are organised into three focused workflows:
//
//  1. Fleet lifecycle  – create, update, read, search, delete a fleet.
//  2. Configuration lifecycle – create a config, add versions, read by version,
//     list all versions, delete a version, delete the config.
//  3. Deployment & managed-entity search – read-only searches only.
//
// What is intentionally NOT covered and why:
//
//   - Adding/removing managed entities (fleet members): a managed entity is a
//     host that has the NR Infrastructure agent installed and actively reporting
//     to the fleet account. There is no programmatic way to enrol a host from a
//     test — it requires a real agent to be running and assigned to the fleet.
//     Hardcoding entity GUIDs from real hosts is fragile because those hosts
//     eventually stop reporting and the GUIDs become stale.
//
//   - Creating or updating fleet deployments: a deployment targets a fleet that
//     has at least one managed entity enrolled in a ring. Without controllable
//     managed entities (see above) any create/update call will fail with
//     "Error occurred while adding managed entities to the fleet ring". The
//     deployment APIs are therefore only exercised via read-only entity search.
//
//   - Triggering a deployment (FleetControlDeploy): depends on a valid
//     deployment in a non-terminal phase, which in turn requires managed
//     entities. Out of scope for the same reason.

package fleetcontrol

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// TestIntegrationFleetLifecycle covers the full fleet CRUD flow plus entity
// lookup operations (GetEntity, GetEntitySearch) that are driven by the fleet
// created within the same test — no hardcoded GUIDs.
func TestIntegrationFleetLifecycle(t *testing.T) {
	t.Parallel()

	_, err := mock.GetFleetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// --- Create ---
	fleetName := fmt.Sprintf("integration-test-fleet-%d", time.Now().Unix())
	createInput := FleetControlFleetEntityCreateInput{
		Name:              fleetName,
		Description:       "Fleet created by integration test",
		ManagedEntityType: FleetControlManagedEntityTypeTypes.HOST,
		Product:           "Infrastructure",
		Scope: FleetControlScopedReferenceInput{
			ID:   testOrganizationID,
			Type: FleetControlEntityScopeTypes.ORGANIZATION,
		},
		OperatingSystem: &FleetControlOperatingSystemCreateInput{
			Type: FleetControlOperatingSystemTypeTypes.LINUX,
		},
		Tags: []FleetControlTagInput{
			{Key: "environment", Values: []string{"test"}},
		},
	}

	createResp, err := client.FleetControlCreateFleet(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotEmpty(t, createResp.Entity.ID)
	require.Equal(t, fleetName, createResp.Entity.Name)

	fleetID := createResp.Entity.ID

	deleted := false
	defer func() {
		if deleted {
			return
		}
		_, _ = client.FleetControlDeleteFleet(fleetID) // best-effort cleanup
	}()

	// --- Update ---
	updatedName := fmt.Sprintf("integration-test-fleet-updated-%d", time.Now().Unix())
	updateResp, err := client.FleetControlUpdateFleet(FleetControlUpdateFleetEntityInput{
		Name:        updatedName,
		Description: "Updated by integration test",
		Tags: []FleetControlTagInput{
			{Key: "environment", Values: []string{"test"}},
			{Key: "updated", Values: []string{"true"}},
		},
	}, fleetID)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, fleetID, updateResp.Entity.ID)
	require.Equal(t, updatedName, updateResp.Entity.Name)

	// --- GetEntity (by ID from create response — no hardcoded GUID) ---
	entity, err := client.GetEntity(fleetID)
	require.NoError(t, err)
	require.NotNil(t, entity)

	fleetEntity, ok := (*entity).(*EntityManagementFleetEntity)
	require.True(t, ok, "expected *EntityManagementFleetEntity")
	require.Equal(t, fleetID, fleetEntity.ID)
	require.Equal(t, updatedName, fleetEntity.Name)
	require.Equal(t, "FLEET", fleetEntity.Type)

	// --- GetEntitySearch (search for FLEET entities, verify ours is present) ---
	searchResp, err := client.GetEntitySearch("", "type = 'FLEET'")
	require.NoError(t, err)
	require.NotNil(t, searchResp)
	require.GreaterOrEqual(t, len(searchResp.Entities), 1, "expected at least the fleet we created")

	found := false
	for _, e := range searchResp.Entities {
		if fe, ok := e.(*EntityManagementFleetEntity); ok && fe.ID == fleetID {
			found = true
			require.Equal(t, "FLEET", fe.Type)
		}
	}
	require.True(t, found, "created fleet not found in entity search results")

	// --- Delete ---
	deleteResp, err := client.FleetControlDeleteFleet(fleetID)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
	require.Equal(t, fleetID, deleteResp.ID)
	deleted = true
}

// TestIntegrationConfigurationLifecycle covers the full configuration CRUD
// flow: create config (v1), add a second version (v2), get each version by
// config entity GUID, get all versions, delete v1, then delete the config.
// Everything is driven by GUIDs returned from the API — no hardcoded values.
func TestIntegrationConfigurationLifecycle(t *testing.T) {
	t.Parallel()

	_, err := mock.GetFleetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	const (
		bodyV1 = "Body for Version 1"
		bodyV2 = "Body for Version 2"
	)

	// --- Create config (v1) ---
	createV1Resp, err := client.FleetControlCreateConfiguration(
		bodyV1,
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": `{"name": "integration-test-config", "agentType": "NRInfra", "managedEntityType": "KUBERNETESCLUSTER"}`,
			},
		},
		testOrganizationID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, createV1Resp.ConfigurationEntityGUID)
	require.NotEmpty(t, createV1Resp.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.Equal(t, 1, createV1Resp.ConfigurationVersion.ConfigurationVersionNumber)

	configGUID := createV1Resp.ConfigurationEntityGUID
	v1GUID := createV1Resp.ConfigurationVersion.ConfigurationVersionEntityGUID

	configDeleted := false
	defer func() {
		if configDeleted {
			return
		}
		_, _ = client.FleetControlDeleteConfiguration(configGUID, testOrganizationID) // best-effort cleanup
	}()

	// --- Add version (v2) ---
	createV2Resp, err := client.FleetControlCreateConfiguration(
		bodyV2,
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": fmt.Sprintf(`{"agentConfiguration": "%s"}`, configGUID),
			},
		},
		testOrganizationID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, createV2Resp.ConfigurationEntityGUID)
	require.NotEmpty(t, createV2Resp.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.Equal(t, 2, createV2Resp.ConfigurationVersion.ConfigurationVersionNumber)
	require.NotEqual(t, v1GUID, createV2Resp.ConfigurationVersion.ConfigurationVersionEntityGUID)

	v2GUID := createV2Resp.ConfigurationVersion.ConfigurationVersionEntityGUID

	// Give the blob service a moment to index the new versions
	time.Sleep(5 * time.Second)

	// The blob service stores the POST body as-is. Since PostWithContext JSON-marshals
	// its argument, a plain Go string is stored with outer quotes (e.g. `"Body for Version 2"`).
	// We strip those quotes when asserting to compare against the original content.
	configBody := func(r *GetConfigurationResponse) string {
		return strings.Trim(string(*r), "\"")
	}

	// --- Get latest (should be v2) ---
	latestResp, err := client.FleetControlGetConfiguration(
		configGUID, testOrganizationID, GetConfigurationModeTypes.ConfigEntity, -1,
	)
	require.NoError(t, err)
	require.Equal(t, bodyV2, configBody(latestResp))

	// --- Get v1 by version number ---
	v1Resp, err := client.FleetControlGetConfiguration(
		configGUID, testOrganizationID, GetConfigurationModeTypes.ConfigEntity, 1,
	)
	require.NoError(t, err)
	require.Equal(t, bodyV1, configBody(v1Resp))

	// --- Get v2 directly via version entity GUID ---
	v2Resp, err := client.FleetControlGetConfiguration(
		v2GUID, testOrganizationID, GetConfigurationModeTypes.ConfigVersionEntity, -1,
	)
	require.NoError(t, err)
	require.Equal(t, bodyV2, configBody(v2Resp))

	// --- List all versions ---
	versionsResp, err := client.FleetControlGetConfigurationVersions(configGUID, testOrganizationID)
	require.NoError(t, err)
	require.NotNil(t, versionsResp)
	require.GreaterOrEqual(t, len(versionsResp.Versions), 2)

	foundV1, foundV2 := false, false
	for _, v := range versionsResp.Versions {
		if v.Version == "1" {
			foundV1 = true
		}
		if v.Version == "2" {
			foundV2 = true
		}
	}
	require.True(t, foundV1, "expected version 1 in list")
	require.True(t, foundV2, "expected version 2 in list")

	// --- Delete v1 (version-level delete) ---
	err = client.FleetControlDeleteConfigurationVersion(v1GUID, testOrganizationID)
	require.NoError(t, err)

	// --- Delete the config entirely ---
	_, err = client.FleetControlDeleteConfiguration(configGUID, testOrganizationID)
	require.NoError(t, err)
	configDeleted = true
}

// TestIntegrationDeploymentAndManagedEntitySearch exercises the search/read APIs
// for fleet deployments and managed (host) entities.
//
// Create/update operations for both deployments and managed-entity membership
// are deliberately excluded here — see the package-level note at the top of
// this file for the full rationale. In short: both require a real NR
// Infrastructure agent to be running and enrolled in the fleet, which cannot
// be orchestrated from an integration test.
func TestIntegrationDeploymentAndManagedEntitySearch(t *testing.T) {
	t.Parallel()

	_, err := mock.GetFleetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Search for any fleet deployment entities in the org
	deploymentSearch, err := client.GetEntitySearch("", "type = 'FLEET_DEPLOYMENT'")
	require.NoError(t, err)
	require.NotNil(t, deploymentSearch)
	// Zero results is acceptable in a fresh account — we just verify the call works
	for _, e := range deploymentSearch.Entities {
		dep, ok := e.(*EntityManagementFleetDeploymentEntity)
		require.True(t, ok, "expected *EntityManagementFleetDeploymentEntity for FLEET_DEPLOYMENT result")
		require.NotEmpty(t, dep.ID)
	}

	// Search for any existing fleet entities (read-only sanity check)
	fleetSearch, err := client.GetEntitySearch("", "type = 'FLEET'")
	require.NoError(t, err)
	require.NotNil(t, fleetSearch)
	for _, e := range fleetSearch.Entities {
		fe, ok := e.(*EntityManagementFleetEntity)
		require.True(t, ok, "expected *EntityManagementFleetEntity for FLEET result")
		require.NotEmpty(t, fe.ID)
		require.Equal(t, "FLEET", fe.Type)
	}

	// If any fleets exist, verify GetFleetMembers works (cursor=nil means first page)
	if len(fleetSearch.Entities) > 0 {
		firstFleet := fleetSearch.Entities[0].(*EntityManagementFleetEntity)
		membersResp, err := client.GetFleetMembers(nil, &FleetControlFleetMembersFilterInput{
			FleetId: firstFleet.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, membersResp)
		// Zero members is fine; we just confirm the API round-trips cleanly
	}
}
