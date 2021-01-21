// +build integration

package alerts

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

var (
	testNrqlConditionRandomString       = mock.RandSeq(5)
	nrqlConditionBaseThreshold          = 1.0          // needed for setting pointer
	nrqlConditionBaseThresholdZeroValue = float64(0)   // needed for setting pointer
	nrqlConditionBaseSignalFillValue    = float64(0.1) // needed for setting pointer
	nrqlConditionBaseExpirationDuration = 1200         // needed for setting pointer
	nrqlConditionBaseEvalOffset         = 3            // needed for setting pointer
	nrqlConditionBaseAggWindow          = 60           // needed for setting pointer
	nrqlConditionBase                   = NrqlConditionBase{
		Description: "test description",
		Enabled:     true,
		Name:        fmt.Sprintf("test-nrql-condition-%s", testNrqlConditionRandomString),
		Nrql: NrqlConditionQuery{
			Query:            "SELECT uniqueCount(host) from Transaction where appName='Dummy App'",
			EvaluationOffset: 3,
		},
		RunbookURL: "test.com",
		Terms: []NrqlConditionTerm{
			{
				Threshold:            &nrqlConditionBaseThreshold,
				ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
				ThresholdDuration:    600,
				Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
				Priority:             NrqlConditionPriorities.Critical,
			},
		},
		ViolationTimeLimitSeconds: 3600,
		Expiration: &AlertsNrqlConditionExpiration{
			CloseViolationsOnExpiration: true,
			ExpirationDuration:          &nrqlConditionBaseExpirationDuration,
			OpenViolationOnExpiration:   false,
		},
		Signal: &AlertsNrqlConditionSignal{
			AggregationWindow: &nrqlConditionBaseAggWindow,
			EvaluationOffset:  &nrqlConditionBaseEvalOffset,
			FillOption:        &AlertsFillOptionTypes.STATIC,
			FillValue:         &nrqlConditionBaseSignalFillValue,
		},
	}
)

// REST API integration test (deprecated)
func TestIntegrationNrqlConditions(t *testing.T) {
	t.Parallel()

	var (
		randomString = mock.RandSeq(5)
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
					Threshold:    nrqlConditionBaseThreshold,
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

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr             = mock.RandSeq(5)
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
	testPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, testPolicy)
	require.NoError(t, err)

	// Test: Create (baseline condition)
	created, err := client.CreateNrqlConditionBaselineMutation(testAccountID, policy.ID, createBaselineInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.ID)
	require.NotNil(t, created.PolicyID)
	require.NotNil(t, created.Signal)
	require.NotNil(t, created.Expiration)
	require.Equal(t, NrqlConditionType("BASELINE"), created.Type)

	// Test: Get (baseline condition)
	require.NoError(t, err)
	readResult, err := client.GetNrqlConditionQuery(testAccountID, created.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, NrqlConditionType("BASELINE"), readResult.Type)
	require.Equal(t, "test description", readResult.Description)

	// Test: Update (baseline condition)
	// There is currently a timing issue with this test.
	// TODO: Once the upstream is fixed, test the updated fields to ensure the this worked
	updateBaselineInput.Description = "test description updated"
	_, err = client.UpdateNrqlConditionBaselineMutation(testAccountID, created.ID, updateBaselineInput)
	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_Static(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr           = mock.RandSeq(5)
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
	testPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, testPolicy)
	require.NoError(t, err)

	// Test: Create (static condition)
	createdStatic, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, createStaticInput)
	require.NoError(t, err)
	require.NotNil(t, createdStatic)
	require.NotNil(t, createdStatic.ID)
	require.NotNil(t, createdStatic.PolicyID)
	require.NotNil(t, createdStatic.Signal)
	require.NotNil(t, createdStatic.Expiration)
	require.Equal(t, NrqlConditionType("STATIC"), createdStatic.Type)

	// Test: Get (static condition)
	readResult, err := client.GetNrqlConditionQuery(testAccountID, createdStatic.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, NrqlConditionType("STATIC"), readResult.Type)
	require.Equal(t, "test description", readResult.Description)

	// Test: Update (static condition)
	// There is currently a timing issue with this test.
	// TODO: Once the upstream is fixed, test the updated fields to ensure the this worked
	updateStaticInput.Description = "test description updated"
	_, err = client.UpdateNrqlConditionStaticMutation(testAccountID, readResult.ID, updateStaticInput)
	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_Outlier(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		expectedGroups     = 1
		violationOverlap   = false
		randStr            = mock.RandSeq(5)
		thresholdCritical  = 0.1
		createOutlierInput = NrqlConditionInput{
			NrqlConditionBase: NrqlConditionBase{
				Description: "test description",
				Enabled:     true,
				Name:        fmt.Sprintf("test-nrql-condition-%s", randStr),
				Nrql: NrqlConditionQuery{
					Query:            "SELECT average(duration) FROM Transaction WHERE appName='Dummy App' FACET host",
					EvaluationOffset: 3,
				},
				RunbookURL: "http://example.com",
				Terms: []NrqlConditionTerm{
					{
						Threshold:            &thresholdCritical,
						ThresholdOccurrences: ThresholdOccurrences.All,
						ThresholdDuration:    120,
						Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
						Priority:             NrqlConditionPriorities.Critical,
					},
				},
				ViolationTimeLimitSeconds: 3600,
			},
			ExpectedGroups:              &expectedGroups,
			OpenViolationOnGroupOverlap: &violationOverlap,
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

	// Test: Create (outlier condition)
	createdOutlier, err := client.CreateNrqlConditionOutlierMutation(testAccountID, strconv.Itoa(policy.ID), createOutlierInput)
	require.NoError(t, err)
	require.NotNil(t, createdOutlier)
	require.NotNil(t, createdOutlier.ID)
	require.NotNil(t, createdOutlier.PolicyID)
	require.NotNil(t, createdOutlier.Signal)
	require.NotNil(t, createdOutlier.Expiration)
	require.Equal(t, NrqlConditionTypes.Outlier, createdOutlier.Type)

	// Test: Get (outlier condition)
	readResult, err := client.GetNrqlConditionQuery(testAccountID, createdOutlier.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, NrqlConditionTypes.Outlier, readResult.Type)
	require.Equal(t, "test description", readResult.Description)

	// Test: Update (outlier condition)
	createOutlierInput.Description = "test description updated"
	updated, err := client.UpdateNrqlConditionOutlierMutation(testAccountID, createdOutlier.ID, createOutlierInput)
	require.NoError(t, err)
	require.Equal(t, "test description updated", updated.Description)

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

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr = mock.RandSeq(5)

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
	testPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, testPolicy)
	require.NoError(t, err)

	// Test: Create Invalid (should result in an error)
	createdBaseline, err := client.CreateNrqlConditionBaselineMutation(testAccountID, policy.ID, testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, createdBaseline)

	// Test: Update Invalid (should result in an error)
	updatedBaseline, err := client.UpdateNrqlConditionBaselineMutation(testAccountID, "8675309", testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, updatedBaseline)

	// Test: Create Invalid (should result in an error)
	createdStatic, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, createdStatic)

	// Test: Update Invalid (should result in an error)
	updatedStatic, err := client.UpdateNrqlConditionStaticMutation(testAccountID, "8675309", testInvalidMutationInput)
	require.Error(t, err)
	require.Nil(t, updatedStatic)

	// Test: 'Not Found' error for non-existent condition
	_, err = client.GetNrqlConditionQuery(testAccountID, "999999999999999")
	require.Error(t, err)
	_, ok := err.(*errors.NotFound)
	require.True(t, ok)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_Search(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr            = mock.RandSeq(5)
		conditionName      = fmt.Sprintf("test-nrql-condition-%s", randStr)
		thresholdCritical  = 1.0
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
				Terms: []NrqlConditionTerm{
					{
						Threshold:            &thresholdCritical,
						ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
						Priority:             NrqlConditionPriorities.Critical,
					},
				},
				ViolationTimeLimitSeconds: 3600,
			},
			BaselineDirection: &NrqlBaselineDirections.LowerOnly,
		}
		searchCriteria = NrqlConditionsSearchCriteria{
			NameLike: conditionName,
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	setupPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, setupPolicy)
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
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_CreateStatic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr            = mock.RandSeq(5)
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
				Terms: []NrqlConditionTerm{
					{
						Threshold:            &nrqlConditionBaseThreshold,
						ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
						Priority:             NrqlConditionPriorities.Critical,
					},
					{
						Threshold:            &nrqlConditionBaseThresholdZeroValue,
						ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             AlertsNRQLConditionTermsOperatorTypes.EQUALS,
						Priority:             NrqlConditionPriorities.Warning,
					},
				},
				ViolationTimeLimitSeconds: 3600,
			},
			ValueFunction: &NrqlConditionValueFunctions.SingleValue,
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	setupPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, setupPolicy)
	require.NoError(t, err)

	condition, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, testConditionInput)
	require.NoError(t, err)
	require.NotNil(t, condition)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_ZeroValueThreshold(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr            = mock.RandSeq(5)
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
				Terms: []NrqlConditionTerm{
					{
						Threshold:            &nrqlConditionBaseThresholdZeroValue,
						ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
						Priority:             NrqlConditionPriorities.Critical,
					},
				},
				ViolationTimeLimitSeconds: 3600,
			},
			ValueFunction: &NrqlConditionValueFunctions.SingleValue,
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	setupPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, setupPolicy)
	require.NoError(t, err)

	condition, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, testConditionInput)
	require.NoError(t, err)
	require.NotNil(t, condition)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}
