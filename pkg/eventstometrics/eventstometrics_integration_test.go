// +build integration

package eventstometrics

import (
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationEventsToMetrics(t *testing.T) {
	t.Parallel()

	var (
		rand                = nr.RandSeq(5)
		testRuleName        = "testRule_" + rand
		testOtherRuleName   = "testRuleOther_" + rand
		testRuleDescription = "testRuleDescription"
		testRuleNrql        = "SELECT uniqueCount(account_id) AS `Transaction.account_id` FROM Transaction FACET appName, name"
		testAccountID       = nr.TestAccountID
		testCreateInput     = []EventsToMetricsCreateRuleInput{
			{
				AccountID:   testAccountID,
				Name:        testRuleName,
				Description: testRuleDescription,
				Nrql:        testRuleNrql,
			},
			{
				AccountID:   testAccountID,
				Name:        testOtherRuleName,
				Description: testRuleDescription,
				Nrql:        testRuleNrql,
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.CreateRules(testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 2, len(created))

	// Test: Get
	rule, err := client.GetRule(testAccountID, created[0].ID)

	require.NoError(t, err)
	require.NotNil(t, rule)

	// Test: Get Multiple
	rules, err := client.GetRules(testAccountID, []string{created[0].ID, created[1].ID})

	require.NoError(t, err)
	require.NotNil(t, rules)
	require.NotEmpty(t, rules)
	require.Equal(t, 2, len(rules))

	// Test: List
	rules, err = client.ListRules(testAccountID)
	require.NoError(t, err)
	require.Greater(t, len(rules), 0)

	// Test: Update
	testUpdateInput := []EventsToMetricsUpdateRuleInput{
		{
			AccountID: testAccountID,
			RuleId:    created[0].ID,
			Enabled:   false,
		},
	}

	updated, err := client.UpdateRules(testUpdateInput)

	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotEmpty(t, updated)
	require.Equal(t, 1, len(updated))
	require.Equal(t, testUpdateInput[0].Enabled, updated[0].Enabled)

	// Test: Delete
	testDeleteInput := []EventsToMetricsDeleteRuleInput{
		{
			AccountID: testAccountID,
			RuleId:    created[0].ID,
		},
		{
			AccountID: testAccountID,
			RuleId:    created[1].ID,
		},
	}
	deleted, err := client.DeleteRules(testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.NotEmpty(t, deleted)
	require.Equal(t, 2, len(deleted))
}

func newIntegrationTestClient(t *testing.T) EventsToMetrics {
	tc := nr.NewIntegrationTestConfig(t)

	return New(tc)
}
