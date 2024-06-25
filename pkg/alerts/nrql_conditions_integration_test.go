//go:build integration
// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	testNrqlConditionRandomString       = mock.RandSeq(5)
	nrqlConditionBaseThreshold          = 1.0                                         // needed for setting pointer
	nrqlConditionBaseThresholdZeroValue = float64(0)                                  // needed for setting pointer
	nrqlConditionBaseSignalFillValue    = float64(0.1)                                // needed for setting pointer
	nrqlConditionBaseExpirationDuration = 1200                                        // needed for setting pointer
	nrqlConditionBaseEvalOffset         = 3                                           // needed for setting pointer
	nrqlConditionBaseAggWindow          = 60                                          // needed for setting pointer
	nrqlConditionBaseAggMethod          = NrqlConditionAggregationMethodTypes.Cadence // needed for setting pointer
	nrqlConditionBaseAggDelay           = 2                                           // needed for setting pointer
	nrqlConditionBaseAggTimer           = 5                                           // needed for setting pointer
	nrqlConditionBaseSlideBy            = 30                                          // needed for setting pointer
	nrqlConditionEvaluationDelay        = 60                                          // needed for setting pointer
	nrqlConditionCreateBase             = nrqlCreateFactory(ConditionArgs{})
	nrqlConditionUpdateBase = nrqlUpdateFactory(ConditionArgs{})
	nrqlConditionCreateWithStreamingMethods = nrqlCreateFactory(ConditionArgs{})
	nrqlConditionUpdateWithStreamingMethods = nrqlUpdateFactory(ConditionArgs{})
	nrqlConditionCreateWithSlideBy = nrqlCreateFactory(ConditionArgs{
		slideBy: nrqlConditionBaseSlideBy,
	})
	nrqlConditionUpdateWithSlideBy = nrqlUpdateFactory(ConditionArgs{
		slideBy: nrqlConditionBaseSlideBy,
	})
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
		createBaselineInput = NrqlConditionCreateInput{
			NrqlConditionCreateBase: nrqlConditionCreateWithSlideBy,
			BaselineDirection:       &NrqlBaselineDirections.LowerOnly,
		}
		updateBaselineInput = NrqlConditionUpdateInput{
			NrqlConditionUpdateBase: nrqlConditionUpdateWithSlideBy,
			BaselineDirection:       &NrqlBaselineDirections.LowerOnly,
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
		createStaticInput = NrqlConditionCreateInput{
			NrqlConditionCreateBase: nrqlConditionCreateBase,
		}
		updateStaticInput = NrqlConditionUpdateInput{
			NrqlConditionUpdateBase: nrqlConditionUpdateBase,
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

func TestIntegrationNrqlConditions_Search(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr            = mock.RandSeq(5)
		conditionName      = fmt.Sprintf("test-nrql-condition-%s", randStr)
		testConditionInput = NrqlConditionCreateInput{
			NrqlConditionCreateBase: nrqlCreateFactory(ConditionArgs{name: conditionName}),
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
		testConditionInput = NrqlConditionCreateInput{
			NrqlConditionCreateBase: NrqlConditionCreateBase{
				Description: "test description",
				Enabled:     true,
				Name:        conditionName,
				Nrql: NrqlConditionCreateQuery{
					Query:            "SELECT rate(sum(apm.service.cpu.usertime.utilization), 1 second) * 100 as cpuUsage FROM Metric WHERE appName like 'Dummy App'",
					EvaluationOffset: &nrqlConditionBaseEvalOffset,
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
		testConditionInput = NrqlConditionCreateInput{
			NrqlConditionCreateBase: nrqlCreateFactory(
				ConditionArgs{threshold: nrqlConditionBaseThresholdZeroValue},
			),
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

func TestIntegrationNrqlConditions_StreamingMethods(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr                     = mock.RandSeq(5)
		createStreamingMethodsInput = NrqlConditionCreateInput{
			NrqlConditionCreateBase: nrqlConditionCreateWithStreamingMethods,
		}
		updateStreamingMethodsInput = NrqlConditionUpdateInput{
			NrqlConditionUpdateBase: nrqlConditionUpdateWithStreamingMethods,
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

	// Test: Create (static condition with streaming methods fields)
	createdStaticWithStreamingMethods, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, createStreamingMethodsInput)
	require.NoError(t, err)
	require.NotNil(t, createdStaticWithStreamingMethods)
	require.NotNil(t, createdStaticWithStreamingMethods.ID)
	require.NotNil(t, createdStaticWithStreamingMethods.PolicyID)
	require.NotNil(t, createdStaticWithStreamingMethods.Signal)
	require.NotNil(t, createdStaticWithStreamingMethods.Expiration)
	require.Equal(t, &nrqlConditionBaseAggMethod, createdStaticWithStreamingMethods.Signal.AggregationMethod)

	// Test: Get (static condition with streaming methods fields)
	readResult, err := client.GetNrqlConditionQuery(testAccountID, createdStaticWithStreamingMethods.ID)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.Equal(t, &nrqlConditionBaseAggMethod, readResult.Signal.AggregationMethod)
	require.Equal(t, &nrqlConditionBaseAggDelay, readResult.Signal.AggregationDelay)
	require.Equal(t, &nrqlConditionEvaluationDelay, readResult.Signal.EvaluationDelay)
	require.Nil(t, createdStaticWithStreamingMethods.Signal.AggregationTimer)

	// Test: Not found
	notFoundResult, err := client.GetNrqlConditionQuery(testAccountID, "1")
	require.Error(t, err)
	require.Nil(t, notFoundResult)
	_, ok := err.(*errors.NotFound)
	require.True(t, ok)

	// Test: Update (static condition with streaming methods fields modified)
	nrqlConditionBaseAggMethodUpdated := NrqlConditionAggregationMethodTypes.EventTimer // needed for setting pointer

	updateStreamingMethodsInput.Signal.AggregationMethod = &nrqlConditionBaseAggMethodUpdated
	updateStreamingMethodsInput.Signal.AggregationDelay = nil
	updateStreamingMethodsInput.Signal.AggregationTimer = &nrqlConditionBaseAggTimer

	updatedStaticWithStreamingMethods, err := client.UpdateNrqlConditionStaticMutation(testAccountID, readResult.ID, updateStreamingMethodsInput)
	require.NoError(t, err)
	require.Equal(t, &nrqlConditionBaseAggMethodUpdated, updatedStaticWithStreamingMethods.Signal.AggregationMethod)
	require.Equal(t, &nrqlConditionBaseAggTimer, updatedStaticWithStreamingMethods.Signal.AggregationTimer)
	require.Nil(t, updatedStaticWithStreamingMethods.Signal.AggregationDelay)

	// Test: Update (static condition without streaming methods fields)
	updateStreamingMethodsInput.Signal.AggregationMethod = nil
	updateStreamingMethodsInput.Signal.AggregationDelay = nil
	updateStreamingMethodsInput.Signal.AggregationTimer = nil
	updateStreamingMethodsInput.Signal.EvaluationOffset = &nrqlConditionBaseEvalOffset

	updatedStaticWithoutStreamingMethods, err := client.UpdateNrqlConditionStaticMutation(testAccountID, readResult.ID, updateStreamingMethodsInput)
	require.NoError(t, err)
	require.Nil(t, updatedStaticWithoutStreamingMethods.Signal.AggregationMethod)
	require.Nil(t, updatedStaticWithoutStreamingMethods.Signal.AggregationDelay)
	require.Nil(t, updatedStaticWithoutStreamingMethods.Signal.AggregationTimer)
	require.Equal(t, &nrqlConditionBaseEvalOffset, updatedStaticWithoutStreamingMethods.Signal.EvaluationOffset)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationNrqlConditions_IgnoreOnExpectedTermination(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr                     = mock.RandSeq(5)
		defaultExpectedTerm				  = true
		createExpectedTermination = NrqlConditionCreateInput{
			NrqlConditionCreateBase: nrqlCreateFactory(ConditionArgs{
				ignoreOnExpectedTermination: &defaultExpectedTerm,
			}),
		}
		updateExpectedTermination = NrqlConditionUpdateInput{
			NrqlConditionUpdateBase: nrqlUpdateFactory(ConditionArgs{
				ignoreOnExpectedTermination: &defaultExpectedTerm,
			}),
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

	// Test: Create (static condition with streaming methods fields)
	createdCondition, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, createExpectedTermination)
	require.NoError(t, err)
	require.NotNil(t, createdCondition)
	require.NotNil(t, createdCondition.ID)
	require.NotNil(t, createdCondition.PolicyID)
	require.NotNil(t, createdCondition.Expiration)
	require.Equal(t, true, createdCondition.Expiration.IgnoreOnExpectedTermination)

	// Test: Update (static condition with updated expected termination)
	updateExpectedTermination.Expiration.IgnoreOnExpectedTermination = false
	updatedCondition, err := client.UpdateNrqlConditionStaticMutation(testAccountID, createdCondition.ID, updateExpectedTermination)
	require.NoError(t, err)
	require.Equal(t, false, updatedCondition.Expiration.IgnoreOnExpectedTermination)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func buildFactoryWithDefaults(baseCondition ConditionArgs) ConditionArgs {
	if baseCondition.description == "" {
		baseCondition.description = "test description"
	}
	if baseCondition.enabled == nil {
		defaultEnabled := true
		baseCondition.enabled = &defaultEnabled
	}
	if baseCondition.name == "" {
		baseCondition.name = "default name"
	}
	if baseCondition.query == "" {
		baseCondition.query = "SELECT rate(sum(apm.service.cpu.usertime.utilization), 1 second) * 100 as cpuUsage FROM Metric WHERE appName like 'Dummy App'"
	}
	if baseCondition.runbookURL == "" {
		baseCondition.runbookURL = "test.com"
	}
	if baseCondition.threshold == 0.0 {
		baseCondition.threshold = nrqlConditionBaseThreshold
	}
	if baseCondition.thresholdOccurrences == "" {
		baseCondition.thresholdOccurrences = ThresholdOccurrences.AtLeastOnce
	}
	if baseCondition.thresholdDuration == 0 {
		baseCondition.thresholdDuration = 600
	}
	if baseCondition.operator == "" {
		baseCondition.operator = AlertsNRQLConditionTermsOperatorTypes.ABOVE
	}
	if baseCondition.priority == "" {
		baseCondition.priority = NrqlConditionPriorities.Critical
	}
	if baseCondition.violationTimeLimitSeconds == 0 {
		baseCondition.violationTimeLimitSeconds = 3600
	}
	if baseCondition.closeViolationsOnExpiration == nil {
		closeViolationsDefault := false
		baseCondition.closeViolationsOnExpiration = &closeViolationsDefault
	}
	if baseCondition.ignoreOnExpectedTermination == nil {
		expectedTermDefault := false
		baseCondition.ignoreOnExpectedTermination = &expectedTermDefault
	}
	if baseCondition.expirationDuration == 0 {
		baseCondition.expirationDuration = nrqlConditionBaseExpirationDuration
	}
	if baseCondition.openViolationOnExpiration == nil {
		openViolationDefault := false
		baseCondition.openViolationOnExpiration = &openViolationDefault
	}
	if baseCondition.aggregationWindow == 0 {
		baseCondition.aggregationWindow = nrqlConditionBaseAggWindow
	}
	if baseCondition.fillOption == "" {
		baseCondition.fillOption = AlertsFillOptionTypes.STATIC
	}
	if baseCondition.fillValue == 0.0 {
		baseCondition.fillValue = nrqlConditionBaseSignalFillValue
	}
	if baseCondition.aggregationMethod == "" {
		baseCondition.aggregationMethod = "average"
	}
	if baseCondition.aggregationDelay == 0 {
		baseCondition.aggregationDelay = 300
	}
	if baseCondition.evaluationDelay == 0 {
		baseCondition.evaluationDelay = nrqlConditionEvaluationDelay
	}
	if baseCondition.evaluationOffset == 0 {
		baseCondition.evaluationOffset = nrqlConditionBaseEvalOffset
	}
	if baseCondition.slideBy == 0 {
		baseCondition.slideBy = 0
	}
	return baseCondition
}

func nrqlUpdateFactory(args ConditionArgs) NrqlConditionUpdateBase {
	var factory = buildFactoryWithDefaults(args)
	return NrqlConditionUpdateBase{
		Description: factory.description,
		Enabled:     *factory.enabled,
		Name:        factory.name,
		Nrql: NrqlConditionUpdateQuery{
			Query: factory.query,
			EvaluationOffset: &factory.evaluationOffset,
		},
		RunbookURL: factory.runbookURL,
		Terms: []NrqlConditionTerm{
			{
				Threshold:            &factory.threshold,
				ThresholdOccurrences: factory.thresholdOccurrences,
				ThresholdDuration:    factory.thresholdDuration,
				Operator:             factory.operator,
				Priority:             factory.priority,
			},
		},
		ViolationTimeLimitSeconds: factory.violationTimeLimitSeconds,
		Expiration: &AlertsNrqlConditionExpiration{
			CloseViolationsOnExpiration: *factory.closeViolationsOnExpiration,
			ExpirationDuration:          &factory.expirationDuration,
			OpenViolationOnExpiration:   *factory.openViolationOnExpiration,
			IgnoreOnExpectedTermination: *factory.ignoreOnExpectedTermination,
		},
		Signal: &AlertsNrqlConditionUpdateSignal{
			AggregationWindow: &factory.aggregationWindow,
			FillOption:        &factory.fillOption,
			FillValue:         &factory.fillValue,
			AggregationMethod: &factory.aggregationMethod,
			AggregationDelay:  &factory.aggregationDelay,
			EvaluationDelay:   &factory.evaluationDelay,
			SlideBy:           &factory.slideBy,
		},
	}
}

func nrqlCreateFactory(args ConditionArgs) NrqlConditionCreateBase {
	var factory = buildFactoryWithDefaults(args)
	return NrqlConditionCreateBase{
		Description: factory.description,
		Enabled:     *factory.enabled,
		Name:        factory.name,
		Nrql: NrqlConditionCreateQuery{
			Query: factory.query,
			EvaluationOffset: &factory.evaluationOffset,
		},
		RunbookURL: factory.runbookURL,
		Terms: []NrqlConditionTerm{
			{
				Threshold:            &factory.threshold,
				ThresholdOccurrences: factory.thresholdOccurrences,
				ThresholdDuration:    factory.thresholdDuration,
				Operator:             factory.operator,
				Priority:             factory.priority,
			},
		},
		ViolationTimeLimitSeconds: factory.violationTimeLimitSeconds,
		Expiration: &AlertsNrqlConditionExpiration{
			CloseViolationsOnExpiration: *factory.closeViolationsOnExpiration,
			ExpirationDuration:          &factory.expirationDuration,
			OpenViolationOnExpiration:   *factory.openViolationOnExpiration,
			IgnoreOnExpectedTermination: *factory.ignoreOnExpectedTermination,
		},
		Signal: &AlertsNrqlConditionCreateSignal{
			AggregationWindow: &factory.aggregationWindow,
			FillOption:        &factory.fillOption,
			FillValue:         &factory.fillValue,
			AggregationMethod: &factory.aggregationMethod,
			AggregationDelay:  &factory.aggregationDelay,
			EvaluationDelay:   &factory.evaluationDelay,
			SlideBy:           &factory.slideBy,
		},
	}
}

type ConditionArgs struct {
	description string;
	enabled *bool;
	name string;
	query string;
	runbookURL string;
	threshold float64;
	thresholdOccurrences ThresholdOccurrence;
	thresholdDuration int;
	operator AlertsNRQLConditionTermsOperator;
	priority NrqlConditionPriority;
	violationTimeLimitSeconds int;
	closeViolationsOnExpiration *bool;
	expirationDuration int;
	openViolationOnExpiration *bool;
	ignoreOnExpectedTermination *bool;
	aggregationWindow int;
	fillOption AlertsFillOption;
	fillValue float64;
	aggregationMethod NrqlConditionAggregationMethod;
	aggregationDelay int;
	evaluationDelay int;
	evaluationOffset int;
	slideBy int;
}
