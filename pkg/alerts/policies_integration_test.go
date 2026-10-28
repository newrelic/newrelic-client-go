//go:build integration
// +build integration

package alerts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestAlertsPolicy_Legacy(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	testIntegrationPolicyNameRandStr := testhelpers.RandSeq(5)
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

	// Test: Not found
	readResult, err = a.GetPolicy(createResult.ID)
	require.Error(t, err)
	require.Nil(t, readResult)
	_, ok := err.(*errors.NotFound)
	require.True(t, ok)
}

func TestAlertsQueryPolicy_GraphQL_Enabled(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	// DTK terraform account
	accountID := testhelpers.IntegrationTestAccountID

	// Create a policy to work with in this test
	testIntegrationPolicyNameRandStr := testhelpers.RandSeq(5)
	policy := AlertsPolicyInput{}
	policy.IncidentPreference = AlertsIncidentPreferenceTypes.PER_POLICY
	policy.Name = fmt.Sprintf("test-alert-policy-%s", testIntegrationPolicyNameRandStr)

	// Test: Create
	createResult, err := a.CreatePolicyMutation(accountID, policy)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.NotEmpty(t, createResult.EntityGuid, "EntityGuid should be returned on create")

	// Query for the policy we policy we just created
	queryResult, err := a.QueryPolicy(accountID, createResult.ID)
	require.NoError(t, err)
	require.NotNil(t, queryResult)
	require.NotEmpty(t, queryResult.EntityGuid, "EntityGuid should be returned on query")
	assert.Equal(t, createResult.EntityGuid, queryResult.EntityGuid, "EntityGuid should match between create and query")

	// Search: all policies (no criteria)
	searchResults, err := a.QueryPolicySearch(accountID, AlertsPoliciesSearchCriteriaInput{})
	require.NoError(t, err)
	require.NotNil(t, searchResults)
	// Verify entityGuid is present in search results
	if len(searchResults) > 0 {
		require.NotEmpty(t, searchResults[0].EntityGuid, "EntityGuid should be returned in search results")
	}

	// Search: exact name match
	nameResults, err := a.QueryPolicySearch(accountID, AlertsPoliciesSearchCriteriaInput{Name: createResult.Name})
	require.NoError(t, err)
	require.NotNil(t, nameResults)
	require.Len(t, nameResults, 1, "exact name search should return exactly one policy")
	require.Equal(t, createResult.ID, nameResults[0].ID)
	require.Equal(t, createResult.Name, nameResults[0].Name)

	// Search: partial name match (nameLike)
	nameLikeResults, err := a.QueryPolicySearch(accountID, AlertsPoliciesSearchCriteriaInput{NameLike: fmt.Sprintf("test-alert-policy-%s", testIntegrationPolicyNameRandStr)})
	require.NoError(t, err)
	require.NotNil(t, nameLikeResults)
	require.GreaterOrEqual(t, len(nameLikeResults), 1, "nameLike search should return at least one policy")
	found := false
	for _, p := range nameLikeResults {
		if p.ID == createResult.ID {
			found = true
			break
		}
	}
	require.True(t, found, "created policy should appear in nameLike search results")

	// Test: Update
	updatePolicy := AlertsPolicyUpdateInput{}
	updatePolicy.Name = fmt.Sprintf("test-alert-policy-updated-%s", testIntegrationPolicyNameRandStr)
	updatePolicy.IncidentPreference = createResult.IncidentPreference

	updateResult, err := a.UpdatePolicyMutation(accountID, createResult.ID, updatePolicy)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	assert.Equal(t, updateResult.Name, updatePolicy.Name)
	assert.Equal(t, updateResult.IncidentPreference, updatePolicy.IncidentPreference)
	require.NotEmpty(t, updateResult.EntityGuid, "EntityGuid should be returned on update")
	assert.Equal(t, createResult.EntityGuid, updateResult.EntityGuid, "EntityGuid should remain the same after update")

	// Test: Delete
	deleteResult, err := a.DeletePolicyMutation(accountID, createResult.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)

	// Expect that the thing we just deleted does not still exist
	queryResult, err = a.QueryPolicy(accountID, createResult.ID)
	require.Error(t, err)
	require.Nil(t, queryResult)
	_, ok := err.(*errors.NotFound)
	require.True(t, ok)
}
