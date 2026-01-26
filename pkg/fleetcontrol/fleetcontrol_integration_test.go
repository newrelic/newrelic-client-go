//go:build integration
// +build integration

package fleetcontrol

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCreateConfigurationAndVersion(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	configurationVersionOneBody := "Body for Version 1"
	configurationVersionTwoBody := "Body for Version 2"

	createConfigurationResponse, err := client.FleetControlCreateConfiguration(
		configurationVersionOneBody,
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": "{\"name\": \"Random Build v100 Test\", \"agentType\": \"NRInfra\", \"managedEntityType\": \"KUBERNETESCLUSTER\"}",
			},
		},
		testOrganizationID,
	)

	require.NoError(t, err)
	fmt.Println(createConfigurationResponse.ConfigurationEntityGUID)
	fmt.Println(createConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)

	require.NotNil(t, createConfigurationResponse.ConfigurationEntityGUID)
	require.NotNil(t, createConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.Equal(t, createConfigurationResponse.ConfigurationVersion.ConfigurationVersionNumber, 1)

	addVersionToConfigurationResponse, err := client.FleetControlCreateConfiguration(
		configurationVersionTwoBody,
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": fmt.Sprintf("{\"agentConfiguration\": \"%s\"}", createConfigurationResponse.ConfigurationEntityGUID),
			},
		},
		testOrganizationID,
	)

	require.NoError(t, err)
	fmt.Println(addVersionToConfigurationResponse.ConfigurationEntityGUID)
	fmt.Println(addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.NotNil(t, addVersionToConfigurationResponse.ConfigurationEntityGUID)
	require.NotNil(t, addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.NotEqual(t, createConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID, addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.Equal(t, addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionNumber, 2)

	time.Sleep(time.Second * 10)

	getConfigurationResponse, err := client.FleetControlGetConfiguration(
		createConfigurationResponse.ConfigurationEntityGUID,
		testOrganizationID,
		GetConfigurationModeTypes.ConfigEntity,
		-1,
	)

	require.NoError(t, err)
	require.Equal(t, *getConfigurationResponse, GetConfigurationResponse(configurationVersionTwoBody))

	getConfigurationVersionOneResponse, err := client.FleetControlGetConfiguration(
		createConfigurationResponse.ConfigurationEntityGUID,
		testOrganizationID,
		GetConfigurationModeTypes.ConfigEntity,
		1,
	)

	require.NoError(t, err)
	require.Equal(t, *getConfigurationVersionOneResponse, GetConfigurationResponse(configurationVersionOneBody))

	getConfigurationVersionTwoResponse, err := client.FleetControlGetConfiguration(
		addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID,
		testOrganizationID,
		GetConfigurationModeTypes.ConfigVersionEntity,
		-1,
	)

	require.NoError(t, err)
	require.Equal(t, *getConfigurationVersionTwoResponse, GetConfigurationResponse(configurationVersionTwoBody))

	// Test getting all configuration versions
	getConfigurationVersionsResponse, err := client.FleetControlGetConfigurationVersions(
		createConfigurationResponse.ConfigurationEntityGUID,
		testOrganizationID,
	)

	require.NoError(t, err)
	require.NotNil(t, getConfigurationVersionsResponse)
	require.NotNil(t, getConfigurationVersionsResponse.Versions)
	require.GreaterOrEqual(t, len(getConfigurationVersionsResponse.Versions), 2)

	// Verify that we have at least 2 versions
	var foundVersionOne, foundVersionTwo bool
	for _, version := range getConfigurationVersionsResponse.Versions {
		if version.Version == "1" {
			foundVersionOne = true
		}
		if version.Version == "2" {
			foundVersionTwo = true
		}
	}
	require.True(t, foundVersionOne, "Expected to find version 1")
	require.True(t, foundVersionTwo, "Expected to find version 2")

}

func TestIntegrationDeleteBlob(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createBlobResponse, err := client.FleetControlDeleteConfiguration(
		"NDgyOTY3M3xOR0VQfEFHRU5UX0NPTkZJR1VSQVRJT058MDE5YjBkMWUtMzBiNS03NGYwLTk2M2EtMjk1NzZjNWUwNjEx",
		testOrganizationID,
	)

	fmt.Println(createBlobResponse)
	require.NoError(t, err)

	// require.NotNil(t, createUserResponse.CreatedUser.ID)
}

func TestIntegrationDeleteConfigurationVersion(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	err = client.FleetControlDeleteConfigurationVersion(
		"NDgyOTY3M3xOR0VQfEFHRU5UX0NPTkZJR1VSQVRJT05fVkVSU0lPTnwwMTliZTljZC1jNDljLTdjZTgtOWJjOS03Y2UyNTVjYWIzMjI",
		testOrganizationID,
	)

	require.NoError(t, err)

	// require.NotNil(t, createUserResponse.CreatedUser.ID)
}

func TestIntegrationGetEntity(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	entityID := "NDgyOTY3M3xOR0VQfEZMRUVUfDAxOWJmOTRmLTAwY2MtNzBjNy1iNzA1LWYzNTdlNjJlZGNjNA"

	entity, err := client.GetEntity(entityID)
	require.NoError(t, err)
	require.NotNil(t, entity)

	// Type assert to EntityManagementFleetEntity since the ID indicates it's a FLEET entity
	fleetEntity, ok := (*entity).(*EntityManagementFleetEntity)
	require.True(t, ok, "Expected entity to be of type EntityManagementFleetEntity")
	require.NotNil(t, fleetEntity)
	require.Equal(t, entityID, fleetEntity.ID)
	require.NotEmpty(t, fleetEntity.Name)
	require.NotEmpty(t, fleetEntity.Type)

	fmt.Printf("Successfully retrieved entity: %s (Type: %s, Name: %s)\n", fleetEntity.ID, fleetEntity.Type, fleetEntity.Name)
}

func TestIntegrationGetEntitySearch(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Search for FLEET type entities
	query := "type = 'FLEET'"

	searchResult, err := client.GetEntitySearch("", query)
	require.NoError(t, err)
	require.NotNil(t, searchResult)
	require.NotNil(t, searchResult.Entities)

	fmt.Printf("Found %d entities matching query '%s'\n", len(searchResult.Entities), query)

	// Verify we have at least one entity
	require.GreaterOrEqual(t, len(searchResult.Entities), 1, "Expected at least one FLEET entity")

	// Check the first entity to verify it's properly unmarshaled
	if len(searchResult.Entities) > 0 {
		firstEntity := searchResult.Entities[0]

		// Try to type assert to EntityManagementFleetEntity
		fleetEntity, ok := firstEntity.(*EntityManagementFleetEntity)
		require.True(t, ok, "Expected entity to be of type EntityManagementFleetEntity")
		require.NotNil(t, fleetEntity)
		require.NotEmpty(t, fleetEntity.ID)
		require.NotEmpty(t, fleetEntity.Name)
		require.NotEmpty(t, fleetEntity.Type)
		require.Equal(t, "FLEET", fleetEntity.Type)

		fmt.Printf("First entity: ID=%s, Name=%s, Type=%s\n", fleetEntity.ID, fleetEntity.Name, fleetEntity.Type)
	}
}

// doesn't work yet, because the fleet deploy part is not yet figured out
func TestIntegrationFleetDeploymentCreateAndUpdate(t *testing.T) {
	t.Skipf("TBD")

	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Step 1: Create a test fleet first (required for deployment)
	fleetName := fmt.Sprintf("Test Fleet for Deployment %d", time.Now().Unix())
	createFleetInput := FleetControlFleetEntityCreateInput{
		Name:              fleetName,
		Description:       "Test fleet for deployment integration test",
		ManagedEntityType: FleetControlManagedEntityTypeTypes.HOST,
		Product:           "Infrastructure",
		Scope: FleetControlScopedReferenceInput{
			ID:   testOrganizationID,
			Type: FleetControlEntityScopeTypes.ORGANIZATION,
		},
		Tags: []FleetControlTagInput{
			{
				Key:    "environment",
				Values: []string{"test"},
			},
		},
	}

	createFleetResponse, err := client.FleetControlCreateFleet(createFleetInput)
	require.NoError(t, err)
	require.NotNil(t, createFleetResponse)
	require.NotNil(t, createFleetResponse.Entity.ID)
	require.Equal(t, fleetName, createFleetResponse.Entity.Name)

	fleetID := createFleetResponse.Entity.ID
	fmt.Printf("Created test fleet with ID: %s\n", fleetID)

	// Step 2: Create a fleet deployment
	deploymentName := fmt.Sprintf("Test Deployment %d", time.Now().Unix())
	createDeploymentInput := FleetControlFleetDeploymentCreateInput{
		FleetId:     fleetID,
		Name:        deploymentName,
		Description: "Test deployment for integration test",
		Scope: FleetControlScopedReferenceInput{
			ID:   testOrganizationID,
			Type: FleetControlEntityScopeTypes.ORGANIZATION,
		},
		Tags: []FleetControlTagInput{
			{
				Key:    "test-type",
				Values: []string{"integration"},
			},
		},
	}

	createDeploymentResponse, err := client.FleetControlCreateFleetDeployment(createDeploymentInput)
	require.NoError(t, err)
	require.NotNil(t, createDeploymentResponse)
	require.NotNil(t, createDeploymentResponse.Entity.ID)
	require.Equal(t, fleetID, createDeploymentResponse.Entity.FleetId)
	require.Equal(t, deploymentName, createDeploymentResponse.Entity.Name)
	require.NotEmpty(t, createDeploymentResponse.Entity.Phase)

	deploymentID := createDeploymentResponse.Entity.ID
	fmt.Printf("Created deployment with ID: %s\n", deploymentID)

	// Verify deployment metadata
	require.NotNil(t, createDeploymentResponse.Entity.Metadata)
	require.NotZero(t, createDeploymentResponse.Entity.Metadata.CreatedAt)
	require.NotEmpty(t, createDeploymentResponse.Entity.Metadata.CreatedBy.ID)

	// Verify tags were set
	require.NotEmpty(t, createDeploymentResponse.Entity.Tags)
	foundTestTag := false
	for _, tag := range createDeploymentResponse.Entity.Tags {
		if tag.Key == "test-type" {
			foundTestTag = true
			require.Contains(t, tag.Values, "integration")
		}
	}
	require.True(t, foundTestTag, "Expected to find test-type tag")

	// Step 3: Update the fleet deployment
	updatedDeploymentName := fmt.Sprintf("Updated Test Deployment %d", time.Now().Unix())
	updateDeploymentInput := FleetControlFleetDeploymentUpdateInput{
		Name:        updatedDeploymentName,
		Description: "Updated description for integration test",
		Tags: []FleetControlTagInput{
			{
				Key:    "test-type",
				Values: []string{"integration", "updated"},
			},
			{
				Key:    "status",
				Values: []string{"modified"},
			},
		},
	}

	updateDeploymentResponse, err := client.FleetControlUpdateFleetDeployment(updateDeploymentInput, deploymentID)
	require.NoError(t, err)
	require.NotNil(t, updateDeploymentResponse)
	require.Equal(t, deploymentID, updateDeploymentResponse.Entity.ID)
	require.Equal(t, updatedDeploymentName, updateDeploymentResponse.Entity.Name)
	require.Equal(t, "Updated description for integration test", updateDeploymentResponse.Entity.Description)

	// Verify update metadata
	require.NotNil(t, updateDeploymentResponse.Entity.Metadata)
	require.NotZero(t, updateDeploymentResponse.Entity.Metadata.UpdatedAt)
	require.NotEmpty(t, updateDeploymentResponse.Entity.Metadata.UpdatedBy.ID)

	// Verify updated tags
	require.NotEmpty(t, updateDeploymentResponse.Entity.Tags)
	foundUpdatedTag := false
	foundStatusTag := false
	for _, tag := range updateDeploymentResponse.Entity.Tags {
		if tag.Key == "test-type" {
			foundUpdatedTag = true
			require.Contains(t, tag.Values, "updated")
		}
		if tag.Key == "status" {
			foundStatusTag = true
			require.Contains(t, tag.Values, "modified")
		}
	}
	require.True(t, foundUpdatedTag, "Expected to find updated test-type tag")
	require.True(t, foundStatusTag, "Expected to find status tag")

	fmt.Printf("Successfully created and updated deployment: %s\n", deploymentID)
}
