//go:build integration
// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

func TestIntegrationListInfrastructureConditions(t *testing.T) {
	t.Parallel()

	var (
		testIntegrationInfrastructureConditionRandStr = nr.RandSeq(5)
		testIntegrationInfrastructureConditionPolicy  = Policy{
			Name: fmt.Sprintf("test-integration-infrastructure-conditions-%s",
				testIntegrationInfrastructureConditionRandStr),
			IncidentPreference: "PER_POLICY",
		}
		thresholdZeroValue                              = float64(0)
		testIntegrationInfrastructureConditionThreshold = InfrastructureConditionThreshold{
			Duration: 6,
			Value:    &thresholdZeroValue,
		}

		testIntegrationInfrastructureCondition = InfrastructureCondition{
			Comparison:   "equal",
			Critical:     &testIntegrationInfrastructureConditionThreshold,
			Enabled:      true,
			Name:         "Java is running",
			ProcessWhere: "(commandName = 'java')",
			Type:         "infra_process_running",
			Where:        "(hostname LIKE '%cassandra%')",
			Description:  "Mozzarella halloumi the big cheese cottage cheese cheese and biscuits cheeseburger fromage frais roquefort.",
		}
	)

	alerts := newIntegrationTestClient(t)

	// Setup
	policy, err := alerts.CreatePolicy(testIntegrationInfrastructureConditionPolicy)

	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := alerts.DeletePolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	testIntegrationInfrastructureCondition.PolicyID = policy.ID
	created, err := alerts.CreateInfrastructureCondition(testIntegrationInfrastructureCondition)

	require.NoError(t, err)
	require.NotZero(t, created)

	// Test: List
	conditions, err := alerts.ListInfrastructureConditions(policy.ID)

	require.NoError(t, err)
	require.Greater(t, len(conditions), 0)

	// Test: Get
	condition, err := alerts.GetInfrastructureCondition(created.ID)

	require.NoError(t, err)
	require.NotZero(t, condition)

	// Test: Update
	created.Name = "Updated"
	created.Description = ""
	updated, err := alerts.UpdateInfrastructureCondition(*created)

	require.NoError(t, err)
	require.NotZero(t, updated)
	require.Equal(t, "", updated.Description)

	// Test: Delete
	err = alerts.DeleteInfrastructureCondition(created.ID)

	require.NoError(t, err)
}
