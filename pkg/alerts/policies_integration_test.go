// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
)

func TestAlertsPolicy_Legacy(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	testIntegrationPolicyNameRandStr := nr.RandSeq(5)
	policy := Policy{
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               fmt.Sprintf("test-alert-policy-%s", testIntegrationPolicyNameRandStr),
	}

	// Test: Create
	createResult, err := a.CreatePolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Read
	readResult, err := a.GetPolicy(createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test: Update
	createResult.Name = fmt.Sprintf("test-alert-policy-updated-%s", testIntegrationPolicyNameRandStr)
	updateResult, err := a.UpdatePolicy(*createResult)

	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	deleteResult, err := a.DeletePolicy(updateResult.ID)

	require.NoError(t, err)
	require.NotNil(t, *deleteResult)
}

func TestAlertsQueryPolicy_GraphQL_Enabled(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	// DTK terraform account
	accountID := 2520528

	// Create a policy to work with in this test
	testIntegrationPolicyNameRandStr := nr.RandSeq(5)
	policy := AlertsPolicyInput{}
	policy.IncidentPreference = PER_POLICY
	policy.Name = fmt.Sprintf("test-alert-policy-%s", testIntegrationPolicyNameRandStr)

	// Test: Create
	createResult, err := a.CreatePolicyMutation(accountID, policy)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Query for the policy we policy we just created
	queryResult, err := a.QueryPolicy(accountID, createResult.ID)
	require.NoError(t, err)
	require.NotNil(t, queryResult)

	// Search
	searchResults, err := a.QueryPolicySearch(accountID, AlertsPoliciesSearchCriteriaInput{})
	require.NoError(t, err)
	require.NotNil(t, searchResults)

	// Test: Update
	updatePolicy := AlertsPolicyUpdateInput{}
	updatePolicy.Name = fmt.Sprintf("test-alert-policy-updated-%s", testIntegrationPolicyNameRandStr)
	updatePolicy.IncidentPreference = createResult.IncidentPreference

	updateResult, err := a.UpdatePolicyMutation(accountID, createResult.ID, updatePolicy)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	assert.Equal(t, updateResult.Name, updatePolicy.Name)
	assert.Equal(t, updateResult.IncidentPreference, updatePolicy.IncidentPreference)

	// Test: Delete
	deleteResult, err := a.DeletePolicyMutation(accountID, createResult.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)

	// Expect that the thing we just deleted does not still exist
	queryResult, err = a.QueryPolicy(accountID, createResult.ID)
	require.Error(t, err)
	require.Nil(t, queryResult)
}
