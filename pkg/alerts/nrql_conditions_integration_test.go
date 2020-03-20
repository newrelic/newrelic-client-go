// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
)

func TestIntegrationNrqlConditions(t *testing.T) {
	t.Parallel()

	var (
		randomString = nr.RandSeq(5)
		alertPolicy  = Policy{
			Name:               fmt.Sprintf("test-integration-nrql-policy-%s", randomString),
			IncidentPreference: "PER_POLICY",
		}
		nrqlConditionName        = fmt.Sprintf("test-integration-nrql-condition-%s", randomString)
		nrqlConditionNameUpdated = fmt.Sprintf("test-integration-nrql-condition-updated-%s", randomString)
		nrqlCondition            = NrqlCondition{
			Nrql: NrqlQuery{
				Query:      "SELECT count(*) FROM Transactions",
				SinceValue: "3",
			},
			Terms: []ConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    1,
					TimeFunction: "all",
				},
			},
			Type:                "static",
			Name:                nrqlConditionName,
			RunbookURL:          "https://www.example.com",
			ValueFunction:       "single_value",
			ViolationCloseTimer: 3600,
			Enabled:             true,
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	policy, err := client.CreatePolicy(alertPolicy)

	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)

		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()

	// Test: Create
	createResult, err := client.CreateNrqlCondition(policy.ID, nrqlCondition)

	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: List
	listResult, err := client.ListNrqlConditions(policy.ID)

	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	readResult, err := client.GetNrqlCondition(policy.ID, createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test: Update
	createResult.Name = nrqlConditionNameUpdated
	updateResult, err := client.UpdateNrqlCondition(*createResult)

	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, nrqlConditionNameUpdated, updateResult.Name)

	// Test: Delete
	result, err := client.DeleteNrqlCondition(updateResult.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntegrationNrqlConditions_Search(t *testing.T) {
	t.Parallel()

	var (
		testAccountID      = 2520528
		randStr            = nr.RandSeq(5)
		conditionName      = fmt.Sprintf("test-nrql-condition-%s", randStr)
		testConditionInput = NrqlConditionBaselineInput{
			NrqlConditionBase: NrqlConditionBase{
				Description: "test description",
				Enabled:     true,
				Name:        conditionName,
				Nrql: NrqlConditionQuery{
					Query:            "SELECT uniqueCount(host) from Transaction where appName='Dummy App'",
					EvaluationOffset: 3,
				},
				RunbookURL: "test.com",
				Terms: []NrqlConditionTerms{
					{
						Threshold:            1,
						ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             NrqlConditionOperators.Above,
						Priority:             NrqlConditionPriorities.Critical,
					},
				},
				ViolationTimeLimit: NrqlConditionViolationTimeLimits.OneHour,
			},
			BaselineDirection: NrqlBaselineDirections.LowerOnly,
		}
		searchCriteria = NrqlConditionsSearchCriteria{
			NameLike: conditionName,
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	setupPolicy := Policy{
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicy(setupPolicy)
	require.NoError(t, err)

	condition, err := client.CreateNrqlConditionBaselineMutation(testAccountID, policy.ID, testConditionInput)
	require.NoError(t, err)
	require.NotNil(t, condition)

	// Test: Search
	searchResults, err := client.SearchNrqlConditionsQuery(testAccountID, searchCriteria)
	require.NoError(t, err)
	require.Greater(t, len(searchResults), 0)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_Baseline(t *testing.T) {
	t.Parallel()

	var (
		randStr       = nr.RandSeq(5)
		testAccountID = 2520528

		testCreateInput = NrqlConditionBaselineInput{
			NrqlConditionBase: NrqlConditionBase{
				Description: "test description",
				Enabled:     true,
				Name:        fmt.Sprintf("test-nrql-condition-%s", randStr),
				Nrql: NrqlConditionQuery{
					Query:            "SELECT uniqueCount(host) from Transaction where appName='Dummy App'",
					EvaluationOffset: 3,
				},
				RunbookURL: "test.com",
				Terms: []NrqlConditionTerms{
					{
						Threshold:            1,
						ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             NrqlConditionOperators.Above,
						Priority:             NrqlConditionPriorities.Critical,
					},
				},
				ViolationTimeLimit: NrqlConditionViolationTimeLimits.OneHour,
			},
			BaselineDirection: NrqlBaselineDirections.LowerOnly,
		}
	)

	// Setup
	client := newIntegrationTestClient(t)

	testPolicy := Policy{
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}

	policy, err := client.CreatePolicy(testPolicy)

	require.NoError(t, err)

	// Test: Create
	created, err := client.CreateNrqlConditionBaselineMutation(testAccountID, policy.ID, testCreateInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.ID)
	require.NotNil(t, created.PolicyID)

	// Test: Get
	readResult, err := client.GetNrqlConditionQuery(testAccountID, created.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, NrqlConditionType("BASELINE"), readResult.Type)
	require.Equal(t, "test description", readResult.Description)

	// Test: Update
	// There is currently a timing issue with this test.
	// TODO: re-enable once fixed in the upstream API
	// updated, err := client.UpdateNrqlConditionBaselineMutation(testAccountID, created.ID, testUpdateInput)
	// require.NoError(t, err)
	// require.NotNil(t, updated)
	// require.Equal(t, "test description updated", updated.Description)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_Static(t *testing.T) {
	t.Parallel()

	var (
		testAccountID         = 2520528
		randStr               = nr.RandSeq(5)
		testCreateStaticInput = NrqlConditionStaticInput{
			NrqlConditionBase: NrqlConditionBase{
				Description: "test description",
				Enabled:     true,
				Name:        fmt.Sprintf("test-nrql-condition-%s", randStr),
				Nrql: NrqlConditionQuery{
					Query:            "SELECT uniqueCount(host) from Transaction where appName='Dummy App'",
					EvaluationOffset: 3,
				},
				RunbookURL: "test.com",
				Terms: []NrqlConditionTerms{
					{
						Threshold:            1,
						ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             NrqlConditionOperators.Above,
						Priority:             NrqlConditionPriorities.Critical,
					},
				},
				ViolationTimeLimit: NrqlConditionViolationTimeLimits.OneHour,
			},
			ValueFunction: NrqlConditionValueFunctions.SingleValue,
		}
	)

	// Setup
	client := newIntegrationTestClient(t)

	testPolicy := Policy{
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}

	policy, err := client.CreatePolicy(testPolicy)

	require.NoError(t, err)

	// Test: Create
	created, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, testCreateStaticInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.ID)
	require.NotNil(t, created.PolicyID)

	// Test: Get
	readResult, err := client.GetNrqlConditionQuery(testAccountID, created.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, NrqlConditionType("BASELINE"), readResult.Type)
	require.Equal(t, "test description", readResult.Description)

	// Test: Update
	// There is currently a timing issue with this test.
	// TODO: re-enable once fixed in the upstream API
	// updated, err := client.UpdateNrqlConditionStaticMutation(testAccountID, created.ID, testUpdateStaticInput)
	// require.NoError(t, err)
	// require.NotNil(t, updated)
	// require.Equal(t, "test description updated", updated.Description)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}
