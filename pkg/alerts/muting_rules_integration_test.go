//go:build integration
// +build integration

package alerts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v3/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

func TestIntegrationAlertsMutingRules(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)
	accountID := testhelpers.IntegrationTestAccountID

	// Schedule fields
	startTime, err1 := time.Parse(time.RFC3339, "2021-07-08T12:30:00Z")
	if err1 != nil {
		t.Fatal(err1)
	}
	endTime, err2 := time.Parse(time.RFC3339, "2021-07-08T14:30:00Z")
	if err2 != nil {
		t.Fatal(err2)
	}
	repeatCount := 10

	// Create a muting rule to work with in this test
	rule := MutingRuleCreateInput{
		Name:        testhelpers.RandSeq(5),
		Description: "Mute host-1 violations",
		Enabled:     true,
		Schedule: &MutingRuleScheduleCreateInput{
			EndRepeat:        nil,
			EndTime:          &NaiveDateTime{endTime},
			Repeat:           &MutingRuleScheduleRepeatTypes.DAILY,
			RepeatCount:      &repeatCount,
			StartTime:        &NaiveDateTime{startTime},
			TimeZone:         "America/Los_Angeles",
			WeeklyRepeatDays: nil,
		},
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
	testIntegrationMutingRuleNewName := testhelpers.RandSeq(5)
	updateRule := MutingRuleUpdateInput{}
	updateRule.Name = testIntegrationMutingRuleNewName

	updateResult, err := a.UpdateMutingRule(accountID, createResult.ID, updateRule)
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	err = a.DeleteMutingRule(accountID, createResult.ID)
	require.NoError(t, err)

	// Test: Not found
	getResult, err = a.GetMutingRule(accountID, createResult.ID)
	require.Error(t, err)
	require.Nil(t, getResult)
	_, ok := err.(*errors.NotFound)
	require.True(t, ok)
}
