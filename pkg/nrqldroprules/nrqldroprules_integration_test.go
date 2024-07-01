//go:build integration
// +build integration

package nrqldroprules

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
	"strconv"
)

func TestIntegrationDropRules(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand                     = mock.RandSeq(5)
		testRuleDescription      = "testRuleDescription_" + rand
		testOtherRuleDescription = "testRuleOtherDescription_" + rand
		testRuleNrql             = "SELECT * FROM Log WHERE container_name = 'noise'"
		testCreateInput          = []NRQLDropRulesCreateDropRuleInput{
			{
				Description: testRuleDescription,
				NRQL:        testRuleNrql,
				Action:      NRQLDropRulesActionTypes.DROP_DATA,
			},
			{
				Description: testOtherRuleDescription,
				NRQL:        testRuleNrql,
				Action:      NRQLDropRulesActionTypes.DROP_DATA,
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.NRQLDropRulesCreate(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 2, len(created.Successes))

	// Test: List
	rules, err := client.GetList(testAccountID)
	require.NoError(t, err)
	require.Greater(t, len(rules.Rules), 0)

	// Test: Rule Exist
	dropRuleID, _ := strconv.Atoi(created.Successes[0].ID)
	rule, err := client.GetDropRuleByID(testAccountID, dropRuleID)
	require.NoError(t, err)
	require.NotNil(t, rule)

	// Test: Delete
	testDeleteInput := []string{created.Successes[0].ID, created.Successes[1].ID}
	deleted, err := client.NRQLDropRulesDelete(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.NotEmpty(t, deleted)
	require.Equal(t, 2, len(deleted.Successes))
}

func TestIntegrationDropRules_Fail(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand                     = mock.RandSeq(5)
		testRuleDescription      = "testRuleDescription_" + rand
		testOtherRuleDescription = "testRuleOtherDescription_" + rand
		testRuleNrql             = "not a query"
		testCreateInput          = []NRQLDropRulesCreateDropRuleInput{
			{
				Description: testRuleDescription,
				NRQL:        testRuleNrql,
				Action:      NRQLDropRulesActionTypes.DROP_DATA,
			},
			{
				Description: testOtherRuleDescription,
				NRQL:        testRuleNrql,
				Action:      NRQLDropRulesActionTypes.DROP_DATA,
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.NRQLDropRulesCreate(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 2, len(created.Failures))

}

func newIntegrationTestClient(t *testing.T) Nrqldroprules {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
