//go:build integration
// +build integration

package applicationsettings

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrationAgentApplicationSettings_All(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)

	aliasName := testhelpers.IntegrationTestApplicationEntityNameNew
	// updating an existing application setting
	// this is expected to throw no error, and successfully updating application setting
	applicationSettingTestResult, err := client.AgentApplicationSettingsUpdate(
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		AgentApplicationSettingsUpdateInput{
			Alias: &aliasName,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, applicationSettingTestResult)
	require.Equal(t, aliasName, applicationSettingTestResult.Alias)

	// updating an existing application setting
	// this is expected to throw no error, and successfully updating application setting
	applicationSettingTestResult, err = client.AgentApplicationSettingsUpdate(
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		AgentApplicationSettingsUpdateInput{
			ApmConfig: &AgentApplicationSettingsApmConfigInput{
				ApdexTarget: 0.5,
			},
		},
	)

	require.NoError(t, err)
	require.NotNil(t, applicationSettingTestResult)
	require.Equal(t, aliasName, applicationSettingTestResult.Alias)
	require.Equal(t, applicationSettingTestResult.ApmSettings.ApmConfig.ApdexTarget, 0.5)

}

func TestIntegrationAgentApplicationSettingsError(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)

	TransactionValue := 0.5
	// updating an existing application setting
	// this is expected to throw no error, and successfully updating application setting
	_, err := client.AgentApplicationSettingsUpdate(
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		AgentApplicationSettingsUpdateInput{
			TransactionTracer: &AgentApplicationSettingsTransactionTracerInput{
				TransactionThresholdValue: &TransactionValue,
			},
		},
	)

	require.Error(t, err)
}
