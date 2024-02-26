//go:build integration
// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

func TestIntegrationConditions(t *testing.T) {
	t.Skipf("This a legacy API which will be deprecated soon")

	var (
		testConditionRandStr = nr.RandSeq(5)
		testConditionPolicy  = Policy{
			Name: fmt.Sprintf("test-integration-alert-conditions-%s",
				testConditionRandStr),
			IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		}
		testCondition = Condition{
			Type:       ConditionTypes.APMApplicationMetric,
			Name:       "Adpex (High)",
			Enabled:    true,
			Entities:   []string{},
			Metric:     MetricTypes.Apdex,
			RunbookURL: "",
			Terms: []ConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    0.9,
					TimeFunction: TimeFunctionTypes.All,
				},
				{
					Duration:     5,
					Operator:     "equal",
					Priority:     "warning",
					Threshold:    0.95,
					TimeFunction: TimeFunctionTypes.All,
				},
			},
			UserDefined: ConditionUserDefined{
				Metric:        "",
				ValueFunction: "",
			},
			Scope:               "application",
			GCMetric:            "",
			ViolationCloseTimer: 0,
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	policy, err := client.CreatePolicy(testConditionPolicy)

	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	createResult, err := client.CreateCondition(policy.ID, testCondition)

	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get
	listResult, err := client.ListConditions(policy.ID)

	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	readResult, err := client.GetCondition(policy.ID, createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test: Update
	createResult.Name = "Apdex Update Test"
	updateResult, err := client.UpdateCondition(*createResult)

	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, "Apdex Update Test", updateResult.Name)

	// Test: Delete
	result, err := client.DeleteCondition(updateResult.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
}
