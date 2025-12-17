//go:build integration
// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCompoundConditions_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr       = mock.RandSeq(5)
		conditionName = fmt.Sprintf("test-compound-condition-%s", randStr)
	)

	client := newIntegrationTestClient(t)

	// Setup: Create a policy
	testPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, testPolicy)
	require.NoError(t, err)

	// Setup: Create two NRQL conditions to use as component conditions
	nrqlConditionInput1 := NrqlConditionCreateInput{
		NrqlConditionCreateBase: NrqlConditionCreateBase{
			Enabled: true,
			Name:    fmt.Sprintf("test-nrql-condition-1-%s", randStr),
			Nrql: NrqlConditionCreateQuery{
				Query: "SELECT count(*) FROM Transaction",
			},
			Terms: []NrqlConditionTerm{
				{
					Threshold:            floatPtr(1.0),
					ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
					Priority:             NrqlConditionPriorities.Critical,
				},
			},
			ViolationTimeLimitSeconds: 3600,
		},
	}

	nrqlConditionInput2 := NrqlConditionCreateInput{
		NrqlConditionCreateBase: NrqlConditionCreateBase{
			Enabled: true,
			Name:    fmt.Sprintf("test-nrql-condition-2-%s", randStr),
			Nrql: NrqlConditionCreateQuery{
				Query: "SELECT average(duration) FROM Transaction",
			},
			Terms: []NrqlConditionTerm{
				{
					Threshold:            floatPtr(0.5),
					ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
					Priority:             NrqlConditionPriorities.Warning,
				},
			},
			ViolationTimeLimitSeconds: 3600,
		},
	}

	condition1, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, nrqlConditionInput1)
	require.NoError(t, err)
	require.NotNil(t, condition1)

	condition2, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, nrqlConditionInput2)
	if err != nil {
		t.Logf("Error creating NRQL condition 2: %+v", err)
	}
	require.NoError(t, err)
	require.NotNil(t, condition2)

	// Test: Create compound condition
	createInput := CompoundConditionCreateInput{
		Name:                  conditionName,
		Enabled:               true,
		FacetMatchingBehavior: stringPtr(string(AlertsFacetMatchingBehaviorTypes.FACETS_IGNORED)),
		RunbookURL:            stringPtr("https://example.com/runbook"),
		ThresholdDuration:     intPtr(60),
		TriggerExpression:     "A AND B",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    condition1.ID,
				Alias: "A",
			},
			{
				ID:    condition2.ID,
				Alias: "B",
			},
		},
	}

	created, err := client.CreateCompoundCondition(testAccountID, policy.ID, createInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.ID)
	require.Equal(t, conditionName, created.Name)
	require.Equal(t, true, created.Enabled)
	require.Equal(t, string(AlertsFacetMatchingBehaviorTypes.FACETS_IGNORED), created.FacetMatchingBehavior)
	require.Equal(t, "https://example.com/runbook", created.RunbookURL)
	require.Equal(t, 60, created.ThresholdDuration)
	require.Equal(t, "A AND B", created.TriggerExpression)
	require.Len(t, created.ComponentConditions, 2)

	// Test: Search compound conditions (search all, then filter in code)
	searchResults, err := client.SearchCompoundConditions(testAccountID, nil, nil, nil)
	require.NoError(t, err)
	require.Greater(t, len(searchResults), 0)

	var foundCondition *CompoundCondition
	for _, c := range searchResults {
		if c.ID == created.ID {
			foundCondition = c
			break
		}
	}
	require.NotNil(t, foundCondition)
	require.Equal(t, conditionName, foundCondition.Name)

	// Test: Update compound condition
	updateInput := CompoundConditionUpdateInput{
		Name:                  fmt.Sprintf("%s-updated", conditionName),
		Enabled:               boolPtr(false),
		FacetMatchingBehavior: stringPtr(string(AlertsFacetMatchingBehaviorTypes.FACETS_IGNORED)),
		RunbookURL:            stringPtr("https://example.com/updated-runbook"),
		ThresholdDuration:     intPtr(60),
		TriggerExpression:     "A OR B",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    condition1.ID,
				Alias: "A",
			},
			{
				ID:    condition2.ID,
				Alias: "B",
			},
		},
	}

	updated, err := client.UpdateCompoundCondition(testAccountID, created.ID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, fmt.Sprintf("%s-updated", conditionName), updated.Name)
	require.Equal(t, false, updated.Enabled)
	require.Equal(t, string(AlertsFacetMatchingBehaviorTypes.FACETS_IGNORED), updated.FacetMatchingBehavior)
	require.Equal(t, "https://example.com/updated-runbook", updated.RunbookURL)
	require.Equal(t, 60, updated.ThresholdDuration)
	require.Equal(t, "A OR B", updated.TriggerExpression)

	// Test: Delete compound condition
	deletedID, err := client.DeleteCompoundCondition(testAccountID, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, deletedID)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationCompoundConditions_Search(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr        = mock.RandSeq(5)
		conditionName1 = fmt.Sprintf("test-compound-condition-1-%s", randStr)
		conditionName2 = fmt.Sprintf("test-compound-condition-2-%s", randStr)
	)

	client := newIntegrationTestClient(t)

	// Setup: Create a policy
	testPolicy := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-%s", randStr),
	}
	policy, err := client.CreatePolicyMutation(testAccountID, testPolicy)
	require.NoError(t, err)

	// Setup: Create two NRQL conditions
	nrqlConditionInput1 := NrqlConditionCreateInput{
		NrqlConditionCreateBase: NrqlConditionCreateBase{
			Enabled: true,
			Name:    fmt.Sprintf("test-nrql-condition-1-%s", randStr),
			Nrql: NrqlConditionCreateQuery{
				Query: "SELECT count(*) FROM Transaction",
			},
			Terms: []NrqlConditionTerm{
				{
					Threshold:            floatPtr(1.0),
					ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
					Priority:             NrqlConditionPriorities.Critical,
				},
			},
			ViolationTimeLimitSeconds: 3600,
		},
	}

	nrqlConditionInput2 := NrqlConditionCreateInput{
		NrqlConditionCreateBase: NrqlConditionCreateBase{
			Enabled: true,
			Name:    fmt.Sprintf("test-nrql-condition-2-%s", randStr),
			Nrql: NrqlConditionCreateQuery{
				Query: "SELECT average(duration) FROM Transaction",
			},
			Terms: []NrqlConditionTerm{
				{
					Threshold:            floatPtr(0.5),
					ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
					Priority:             NrqlConditionPriorities.Warning,
				},
			},
			ViolationTimeLimitSeconds: 3600,
		},
	}

	condition1, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, nrqlConditionInput1)
	require.NoError(t, err)

	condition2, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy.ID, nrqlConditionInput2)
	require.NoError(t, err)

	// Create compound condition 1
	createInput1 := CompoundConditionCreateInput{
		Name:              conditionName1,
		Enabled:           true,
		ThresholdDuration: intPtr(60),
		TriggerExpression: "A AND B",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    condition1.ID,
				Alias: "A",
			},
			{
				ID:    condition2.ID,
				Alias: "B",
			},
		},
	}

	created1, err := client.CreateCompoundCondition(testAccountID, policy.ID, createInput1)
	require.NoError(t, err)
	require.NotNil(t, created1)

	// Create compound condition 2
	createInput2 := CompoundConditionCreateInput{
		Name:              conditionName2,
		Enabled:           true,
		ThresholdDuration: intPtr(60),
		TriggerExpression: "A OR B",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    condition1.ID,
				Alias: "a",
			},
			{
				ID:    condition2.ID,
				Alias: "b",
			},
		},
	}

	created2, err := client.CreateCompoundCondition(testAccountID, policy.ID, createInput2)
	require.NoError(t, err)
	require.NotNil(t, created2)

	// Test: Search with filter for specific ID
	filter := &AlertsCompoundConditionFilterInput{
		Id: &AlertsCompoundConditionIDFilter{
			Eq: &created1.ID,
		},
	}
	searchResults, err := client.SearchCompoundConditions(testAccountID, filter, nil, nil)
	require.NoError(t, err)
	require.Greater(t, len(searchResults), 0)

	// Verify we got the correct condition
	var foundCondition *CompoundCondition
	for _, c := range searchResults {
		if c.ID == created1.ID {
			foundCondition = c
			break
		}
	}
	require.NotNil(t, foundCondition)
	require.Equal(t, conditionName1, foundCondition.Name)

	// Test: Search with multiple IDs using IN operator
	filterIn := &AlertsCompoundConditionFilterInput{
		Id: &AlertsCompoundConditionIDFilter{
			In: []string{created1.ID, created2.ID},
		},
	}
	searchResultsIn, err := client.SearchCompoundConditions(testAccountID, filterIn, nil, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(searchResultsIn), 2)

	// Test: Search with sort
	sort := []AlertsCompoundConditionSortInput{
		{
			Key:       string(AlertsCompoundConditionSortKeyTypes.NAME),
			Direction: string(AlertsCompoundConditionSortDirectionTypes.ASCENDING),
		},
	}
	searchResultsSort, err := client.SearchCompoundConditions(testAccountID, nil, sort, nil)
	require.NoError(t, err)
	require.Greater(t, len(searchResultsSort), 0)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy.ID, policy.Name, err)
		}
	}()
}

func TestIntegrationCompoundConditions_UpdatePolicyID(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		randStr       = mock.RandSeq(5)
		conditionName = fmt.Sprintf("test-compound-condition-%s", randStr)
	)

	client := newIntegrationTestClient(t)

	// Setup: Create two policies
	testPolicy1 := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-1-%s", randStr),
	}
	policy1, err := client.CreatePolicyMutation(testAccountID, testPolicy1)
	require.NoError(t, err)

	testPolicy2 := AlertsPolicyInput{
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		Name:               fmt.Sprintf("test-alert-policy-2-%s", randStr),
	}
	policy2, err := client.CreatePolicyMutation(testAccountID, testPolicy2)
	require.NoError(t, err)

	// Setup: Create two NRQL conditions
	nrqlConditionInput1 := NrqlConditionCreateInput{
		NrqlConditionCreateBase: NrqlConditionCreateBase{
			Enabled: true,
			Name:    fmt.Sprintf("test-nrql-condition-1-%s", randStr),
			Nrql: NrqlConditionCreateQuery{
				Query: "SELECT count(*) FROM Transaction",
			},
			Terms: []NrqlConditionTerm{
				{
					Threshold:            floatPtr(1.0),
					ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
					Priority:             NrqlConditionPriorities.Critical,
				},
			},
			ViolationTimeLimitSeconds: 3600,
		},
	}

	nrqlConditionInput2 := NrqlConditionCreateInput{
		NrqlConditionCreateBase: NrqlConditionCreateBase{
			Enabled: true,
			Name:    fmt.Sprintf("test-nrql-condition-2-%s", randStr),
			Nrql: NrqlConditionCreateQuery{
				Query: "SELECT average(duration) FROM Transaction",
			},
			Terms: []NrqlConditionTerm{
				{
					Threshold:            floatPtr(0.5),
					ThresholdOccurrences: ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             AlertsNRQLConditionTermsOperatorTypes.ABOVE,
					Priority:             NrqlConditionPriorities.Warning,
				},
			},
			ViolationTimeLimitSeconds: 3600,
		},
	}

	condition1, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy1.ID, nrqlConditionInput1)
	require.NoError(t, err)

	condition2, err := client.CreateNrqlConditionStaticMutation(testAccountID, policy1.ID, nrqlConditionInput2)
	require.NoError(t, err)

	// Test: Create compound condition in policy1
	createInput := CompoundConditionCreateInput{
		Name:              conditionName,
		Enabled:           true,
		ThresholdDuration: intPtr(60),
		TriggerExpression: "A AND B",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    condition1.ID,
				Alias: "A",
			},
			{
				ID:    condition2.ID,
				Alias: "B",
			},
		},
	}

	created, err := client.CreateCompoundCondition(testAccountID, policy1.ID, createInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, policy1.ID, created.PolicyID)

	// Test: Update compound condition to move to policy2
	updateInput := CompoundConditionUpdateInput{
		Name:              conditionName,
		Enabled:           boolPtr(true),
		PolicyID:          stringPtr(policy2.ID),
		ThresholdDuration: intPtr(60),
		TriggerExpression: "A AND B",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    condition1.ID,
				Alias: "A",
			},
			{
				ID:    condition2.ID,
				Alias: "B",
			},
		},
	}

	updated, err := client.UpdateCompoundCondition(testAccountID, created.ID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, policy2.ID, updated.PolicyID)

	// Deferred teardown
	defer func() {
		_, err := client.DeletePolicyMutation(testAccountID, policy1.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy1.ID, policy1.Name, err)
		}
		_, err = client.DeletePolicyMutation(testAccountID, policy2.ID)
		if err != nil {
			t.Logf("error cleaning up alert policy %s (%s): %s", policy2.ID, policy2.Name, err)
		}
	}()
}

// Helper functions to create pointers
func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
