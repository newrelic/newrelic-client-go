//go:build integration
// +build integration

package agentapplications

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// GUID for Dummy App in integration test account
const INTEGRATION_TEST_APM_APPLICATION_GUID = "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTczNDgyNjM4"

func TestIntegrationAgentApplicationBrowser_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	settings := AgentApplicationBrowserSettingsInput{}

	// Create
	createResult, err := client.AgentApplicationCreateBrowser(testAccountID, appName, settings)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.Equal(t, appName, createResult.Name)

	// Update
	updateSettings := AgentApplicationSettingsUpdateInput{}
	updateResult, err := client.AgentApplicationSettingsUpdate(createResult.GUID, updateSettings)
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Delete
	deleteResult, err := client.AgentApplicationDelete(createResult.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.True(t, deleteResult.Success)
}

func TestIntegrationAgentApplicationBrowser_EnableThenDisableSettings(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	cookiesEnabled := true
	distributedTracingEnabled := true
	settings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            &cookiesEnabled,
		DistributedTracingEnabled: &distributedTracingEnabled,
		LoaderType:                AgentApplicationBrowserLoaderTypes.LITE,
	}

	// Create
	createResult, err := client.AgentApplicationCreateBrowser(testAccountID, appName, settings)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.Equal(t, appName, createResult.Name)

	cookiesEnabled = false
	// Update
	updateSettings := AgentApplicationSettingsUpdateInput{
		BrowserMonitoring: &AgentApplicationSettingsBrowserMonitoringInput{
			Loader: &AgentApplicationSettingsBrowserLoaderInputTypes.PRO,
			DistributedTracing: &AgentApplicationSettingsBrowserDistributedTracingInput{
				Enabled: false,
			},
			Privacy: &AgentApplicationSettingsBrowserPrivacyInput{
				CookiesEnabled: &cookiesEnabled,
			},
		},
	}
	updateResult, err := client.AgentApplicationSettingsUpdate(createResult.GUID, updateSettings)
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Delete
	deleteResult, err := client.AgentApplicationDelete(createResult.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.True(t, deleteResult.Success)
}

func TestIntegrationAgentApplicationBrowser_DisableThenEnableSettings(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	cookiesEnabled := false
	distributedTracingEnabled := false
	settings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            &cookiesEnabled,
		DistributedTracingEnabled: &distributedTracingEnabled,
		LoaderType:                AgentApplicationBrowserLoaderTypes.LITE,
	}

	// Create
	createResult, err := client.AgentApplicationCreateBrowser(testAccountID, appName, settings)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.Equal(t, appName, createResult.Name)

	cookiesEnabled = true
	// Update
	updateSettings := AgentApplicationSettingsUpdateInput{
		BrowserMonitoring: &AgentApplicationSettingsBrowserMonitoringInput{
			Loader: &AgentApplicationSettingsBrowserLoaderInputTypes.PRO,
			DistributedTracing: &AgentApplicationSettingsBrowserDistributedTracingInput{
				Enabled: true,
			},
			Privacy: &AgentApplicationSettingsBrowserPrivacyInput{
				CookiesEnabled: &cookiesEnabled,
			},
		},
	}
	updateResult, err := client.AgentApplicationSettingsUpdate(createResult.GUID, updateSettings)
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Delete
	deleteResult, err := client.AgentApplicationDelete(createResult.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.True(t, deleteResult.Success)
}

func TestIntegrationAgentApplicationBrowser_InvalidLoaderTypeInput(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	cookiesEnabled := true
	distributedTracingEnabled := true
	settings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            &cookiesEnabled,
		DistributedTracingEnabled: &distributedTracingEnabled,
		LoaderType:                AgentApplicationBrowserLoader("INVALID"),
	}

	// Create
	result, err := client.AgentApplicationCreateBrowser(testAccountID, appName, settings)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestIntegrationAgentApplicationEnableAPMBrowser_Basic(t *testing.T) {
	t.Parallel()

	client := newAgentApplicationIntegrationTestClient(t)
	settings := AgentApplicationBrowserSettingsInput{}

	// Enable
	result, err := client.AgentApplicationEnableApmBrowser(common.EntityGUID(INTEGRATION_TEST_APM_APPLICATION_GUID), settings)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntegrationAgentApplicationEnableAPMBrowser_WithSettings(t *testing.T) {
	t.Parallel()

	cookiesEnabled := true
	distributedTracingEnabled := true
	client := newAgentApplicationIntegrationTestClient(t)
	settings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            &cookiesEnabled,
		DistributedTracingEnabled: &distributedTracingEnabled,
		LoaderType:                AgentApplicationBrowserLoaderTypes.PRO,
	}

	// Enable
	result, err := client.AgentApplicationEnableApmBrowser(common.EntityGUID(INTEGRATION_TEST_APM_APPLICATION_GUID), settings)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntegrationAgentApplicationAPM_Basic(t *testing.T) {
	t.Parallel()
	client := newAgentApplicationIntegrationTestClient(t)

	aliasName := testhelpers.IntegrationTestApplicationEntityNameNew
	// updating an existing application setting
	// this is expected to throw no error, and successfully updating application setting
	applicationSettingTestResult, err := client.AgentApplicationSettingsUpdate(
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		AgentApplicationSettingsUpdateInput{
			Alias: aliasName,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, applicationSettingTestResult)
	require.Equal(t, aliasName, applicationSettingTestResult.Alias)
}

func TestIntegrationAgentApplicationAPM_WithSettings(t *testing.T) {
	t.Parallel()
	client := newAgentApplicationIntegrationTestClient(t)
	aliasName := testhelpers.IntegrationTestApplicationEntityNameNew
	// updating an existing application setting
	// this is expected to throw no error, and successfully updating application setting
	AgentApplicationsResult, err := client.AgentApplicationSettingsUpdate(
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		AgentApplicationSettingsUpdateInput{
			ApmConfig: &AgentApplicationSettingsApmConfigInput{
				ApdexTarget: 0.5,
			},
		},
	)

	require.NoError(t, err)
	require.NotNil(t, AgentApplicationsResult)
	require.Equal(t, aliasName, AgentApplicationsResult.Alias)
	require.Equal(t, AgentApplicationsResult.ApmSettings.ApmConfig.ApdexTarget, 0.5)
}

func TestIntegrationAgentApplicationAPM_WithSettingsError(t *testing.T) {
	t.Parallel()
	client := newAgentApplicationIntegrationTestClient(t)

	txnValue := 0.5
	// updating an existing application setting
	// this is expected to throw error
	_, err := client.AgentApplicationSettingsUpdate(
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		AgentApplicationSettingsUpdateInput{
			TransactionTracer: &AgentApplicationSettingsTransactionTracerInput{
				TransactionThresholdValue: &txnValue,
			},
		},
	)

	require.Error(t, err)
}

func newAgentApplicationIntegrationTestClient(t *testing.T) AgentApplications {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
