// +build integration

package alerts

import (
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
)

func TestAlertsMutingRules(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	// DTK terraform account
	accountID := 2520528

	// Create a policy to work with in this test
	rule := MutingRuleCreateInput{
		Name:        nr.RandSeq(5),
		Description: "Mute host-1 violations",
		Enabled:     true,
	}
	condition := MutingRuleCondition{
		Attribute: "tag.host",
		Operator:  "EQUALS",
		Values:    []string{"host-1"},
	}
	rule.Condition.Operator = "AND"
	rule.Condition.Conditions = append(rule.Condition.Conditions, condition)

	// Test: Create
	createResult, err := a.CreateMutingRule(accountID, rule)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	getResult, err := a.GetMutingRule(accountID, createResult.ID)
	require.NoError(t, err)
	require.NotNil(t, getResult)

	// Test: List
	listResult, err := a.ListMutingRules(accountID)
	require.NoError(t, err)
	require.True(t, len(listResult) > 0)

	// Test: Update
	testIntegrationMutingRuleNewName := nr.RandSeq(5)
	updateRule := MutingRuleUpdateInput{}
	updateRule.Name = testIntegrationMutingRuleNewName

	updateResult, err := a.UpdateMutingRule(accountID, createResult.ID, updateRule)
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: delete
	err = a.DeleteMutingRule(accountID, createResult.ID)
	require.NoError(t, err)

}
