//go:build integration
// +build integration

package agentapplication

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// GUID for Dummy App in our test account
const INTEGRATION_TEST_APM_APPLICATION_GUID = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"

func TestIntegrationAgentApplicationBrowser_DefaultSettings(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	applicationSettings := AgentApplicationBrowserSettingsInput{}

	// Create
	createResult, err := client.AgentApplicationCreateBrowser(testAccountID, appName, applicationSettings)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Delete
	deleteResult, err := client.AgentApplicationDelete(common.EntityGUID(createResult.GUID))
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.True(t, deleteResult.Success)
}

func TestIntegrationAgentApplicationBrowser_BasicSettings(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	applicationSettings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            true,
		DistributedTracingEnabled: true,
		LoaderType:                AgentApplicationBrowserLoaderTypes.LITE,
	}

	// Create
	createResult, err := client.AgentApplicationCreateBrowser(testAccountID, appName, applicationSettings)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Delete
	deleteResult, err := client.AgentApplicationDelete(common.EntityGUID(createResult.GUID))
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.True(t, deleteResult.Success)
}

func TestIntegrationAgentApplicationBrowser_InvalidLoaderTypeError(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	applicationSettings := AgentApplicationBrowserSettingsInput{
		LoaderType: AgentApplicationBrowserLoader("INVALID"),
	}

	// Should result in an error
	result, err := client.AgentApplicationCreateBrowser(testAccountID, appName, applicationSettings)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestIntegrationAgentApplicationEnableAPMBrowser_DefaultSettings(t *testing.T) {
	t.Parallel()

	client := newAgentApplicationIntegrationTestClient(t)
	guid := common.EntityGUID(INTEGRATION_TEST_APM_APPLICATION_GUID)
	applicationSettings := AgentApplicationBrowserSettingsInput{}

	result, err := client.AgentApplicationEnableApmBrowser(guid, applicationSettings)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntegrationAgentApplicationEnableAPMBrowser_BasicSettings(t *testing.T) {
	t.Parallel()

	client := newAgentApplicationIntegrationTestClient(t)
	guid := common.EntityGUID(INTEGRATION_TEST_APM_APPLICATION_GUID)
	applicationSettings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            true,
		DistributedTracingEnabled: true,
		LoaderType:                AgentApplicationBrowserLoaderTypes.LITE,
	}

	result, err := client.AgentApplicationEnableApmBrowser(guid, applicationSettings)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func newAgentApplicationIntegrationTestClient(t *testing.T) AgentApplication {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}

// no-lint
func cleanIntegrationTestBrowserApplications(t *testing.T) error {
	entitiesClient := entities.New(testhelpers.NewIntegrationTestConfig(t))
	applicationClient := newAgentApplicationIntegrationTestClient(t)
	query := "domain = 'BROWSER' AND entityType = 'BROWSER_APPLICATION' AND (name LIKE '%nr-test-%' OR name LIKE '%nr_test_%')"

	fmt.Printf("[INFO] cleaning up any extraneous integration test resources...")

	for {
		matches, err := entitiesClient.GetEntitySearchByQuery(
			entities.EntitySearchOptions{},
			query,
			[]entities.EntitySearchSortCriteria{},
		)

		if err != nil {
			return fmt.Errorf("error cleaning up dangling synthetics resources: %s", err)
		}

		if matches != nil {
			resources := matches.Results.Entities
			for _, r := range resources {
				_, err := applicationClient.AgentApplicationDelete(common.EntityGUID(r.GetGUID()))
				if err != nil {
					log.Printf("[ERROR] error deleting dangling resource: %s", err)
				}
			}

			fmt.Printf("\n[INFO] deleted %d extraneous resources", len(resources))
		}

		if matches.Results.NextCursor == "" {
			break
		}
	}

	return nil
}
