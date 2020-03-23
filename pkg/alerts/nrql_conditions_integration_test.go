// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
)

var (
	testNrqlConditionRandomString = nr.RandSeq(5)
	nrqlConditionBase             = NrqlConditionBase{
		Description: "test description",
		Enabled:     true,
		Name:        fmt.Sprintf("test-nrql-condition-%s", testNrqlConditionRandomString),
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
	}
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

func TestIntegrationNrqlConditions_Baseline(t *testing.T) {
	t.Parallel()

	var (
		randStr             = nr.RandSeq(5)
		createBaselineInput = NrqlConditionInput{
			NrqlConditionBase: nrqlConditionBase,
			BaselineDirection: &NrqlBaselineDirections.LowerOnly,
		}
		updateBaselineInput = NrqlConditionInput{
			NrqlConditionBase: nrqlConditionBase,
			BaselineDirection: &NrqlBaselineDirections.LowerOnly,
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

	// Test: Create (baseline condition)
	created, err := client.CreateNrqlConditionBaselineMutation(nr.TestAccountID, policy.ID, createBaselineInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.ID)
	require.NotNil(t, created.PolicyID)
	require.Equal(t, NrqlConditionType("BASELINE"), created.Type)

	// Test: Get (baseline condition)
	readResult, err := client.GetNrqlConditionQuery(nr.TestAccountID, created.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, NrqlConditionType("BASELINE"), readResult.Type)
	require.Equal(t, "test description", readResult.Description)

	// Test: Update (baseline condition)
	// There is currently a timing issue with this test.
	// TODO: Once the upstream is fixed, test the updated fields to ensure the this worked
	updateBaselineInput.Description = "test description updated"
	_, err = client.UpdateNrqlConditionBaselineMutation(nr.TestAccountID, created.ID, updateBaselineInput)
	require.NoError(t, err)

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
		randStr           = nr.RandSeq(5)
		createStaticInput = NrqlConditionInput{
			NrqlConditionBase: nrqlConditionBase,
			ValueFunction:     &NrqlConditionValueFunctions.SingleValue,
		}
		updateStaticInput = NrqlConditionInput{
			NrqlConditionBase: nrqlConditionBase,
			ValueFunction:     &NrqlConditionValueFunctions.Sum,
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

	// Test: Create (static condition)
	createdStatic, err := client.CreateNrqlConditionStaticMutation(nr.TestAccountID, policy.ID, createStaticInput)
	require.NoError(t, err)
	require.NotNil(t, createdStatic)
	require.NotNil(t, createdStatic.ID)
	require.NotNil(t, createdStatic.PolicyID)
	require.Equal(t, NrqlConditionType("STATIC"), createdStatic.Type)

	// Test: Get (static condition)
	readResult, err := client.GetNrqlConditionQuery(nr.TestAccountID, createdStatic.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, NrqlConditionType("STATIC"), readResult.Type)
	require.Equal(t, "test description", readResult.Description)

	// Test: Update (static condition)
	// There is currently a timing issue with this test.
	// TODO: Once the upstream is fixed, test the updated fields to ensure the this worked
	updateStaticInput.Description = "test description updated"
	_, err = client.UpdateNrqlConditionStaticMutation(nr.TestAccountID, readResult.ID, updateStaticInput)
	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_ErrorScenarios(t *testing.T) {
	t.Parallel()

	var (
		randStr = nr.RandSeq(5)

		// Invalid NrqlConditionInput (Baseline and ValueFunction cannot exist together)
		testInvalidMutationInput = NrqlConditionInput{
			NrqlConditionBase: nrqlConditionBase,

			// Having both of the following fields should result in an error returned from the API
			BaselineDirection: &NrqlBaselineDirections.LowerOnly,
			ValueFunction:     &NrqlConditionValueFunctions.SingleValue,
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

	// Test: Create Invalid (should result in an error)
	createdBaseline, err := client.CreateNrqlConditionBaselineMutation(nr.TestAccountID, policy.ID, testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, createdBaseline)

	// Test: Update Invalid (should result in an error)
	updatedBaseline, err := client.UpdateNrqlConditionBaselineMutation(nr.TestAccountID, "8675309", testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, updatedBaseline)

	// Test: Create Invalid (should result in an error)
	createdStatic, err := client.CreateNrqlConditionStaticMutation(nr.TestAccountID, policy.ID, testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, createdStatic)

	// Test: Update Invalid (should result in an error)
	updatedStatic, err := client.UpdateNrqlConditionStaticMutation(nr.TestAccountID, "8675309", testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, updatedStatic)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_Search(t *testing.T) {
	t.Parallel()

	var (
		randStr            = nr.RandSeq(5)
		conditionName      = fmt.Sprintf("test-nrql-condition-%s", randStr)
		testConditionInput = NrqlConditionInput{
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
			BaselineDirection: &NrqlBaselineDirections.LowerOnly,
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

	condition, err := client.CreateNrqlConditionBaselineMutation(nr.TestAccountID, policy.ID, testConditionInput)
	require.NoError(t, err)
	require.NotNil(t, condition)

	// Test: Search
	searchResults, err := client.SearchNrqlConditionsQuery(nr.TestAccountID, searchCriteria)
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
