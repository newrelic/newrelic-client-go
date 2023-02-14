//go:build integration
// +build integration

package agentapplication

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationAgentApplicationBrowser_DefaultSettings(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newAgentApplicationIntegrationTestClient(t)

	appName := testhelpers.GenerateRandomName(10)
	applicationSettings := AgentApplicationBrowserSettingsInput{}
	actual, err := client.AgentApplicationCreateBrowser(testAccountID, appName, applicationSettings)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

func TestIntegrationAgentApplicationBrowser_BasicSettings(t *testing.T) {
	t.Parallel()

	client := newAgentApplicationIntegrationTestClient(t)
	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skip("Skipping TestIntegrationAgentApplicationBrowser_Basic. Account ID required to run this integration test.")
	}

	appName := testhelpers.GenerateRandomName(10)
	applicationSettings := AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            true,
		DistributedTracingEnabled: true,
		LoaderType:                AgentApplicationBrowserLoaderTypes.LITE,
	}
	actual, err := client.AgentApplicationCreateBrowser(testAccountID, appName, applicationSettings)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

func newAgentApplicationIntegrationTestClient(t *testing.T) AgentApplication {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
