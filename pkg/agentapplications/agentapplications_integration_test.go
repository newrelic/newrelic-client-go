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

func TestIntegrationAgentApplicationBrowser_WithSettings(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	cookiesEnabled := true
	settings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            &cookiesEnabled,
		DistributedTracingEnabled: true,
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

func TestIntegrationAgentApplicationBrowser_InvalidLoaderTypeInput(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)
	appName := testhelpers.GenerateRandomName(10)
	cookiesEnabled := true
	settings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            &cookiesEnabled,
		DistributedTracingEnabled: true,
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
	client := newAgentApplicationIntegrationTestClient(t)
	settings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            &cookiesEnabled,
		DistributedTracingEnabled: true,
		LoaderType:                AgentApplicationBrowserLoaderTypes.PRO,
	}

	// Enable
	result, err := client.AgentApplicationEnableApmBrowser(common.EntityGUID(INTEGRATION_TEST_APM_APPLICATION_GUID), settings)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func newAgentApplicationIntegrationTestClient(t *testing.T) AgentApplications {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
