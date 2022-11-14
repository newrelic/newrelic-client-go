//go:build integration
// +build integration

package logconfigurations

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationObfuscationRule(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "testName_" + rand
		testDescription = "testDescription_" + rand
		expression      = createTestObfuscationExpression(t, testAccountID)
		testCreateInput = LogConfigurationsCreateObfuscationRuleInput{
			Description: testDescription,
			Name:        testName,
			Filter:      "entity.guids='MjUyMDUyOHxJTkZSQXxOQXwzMjI2NzYxMDM5Njk4NjQ3MTM2'",
			Enabled:     true,
			Actions: []LogConfigurationsCreateObfuscationActionInput{
				{
					Attributes:   []string{"awsAccountId"},
					ExpressionId: expression.ID,
					Method:       "MASK",
				},
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateObfuscationRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)

	// Test: Delete
	testDeleteInput := created.ID
	deleted, err := client.LogConfigurationsDeleteObfuscationRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.NotEmpty(t, deleted)
}

//Update
func TestIntegrationObfuscationRuleUpdate(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "testName_" + rand
		testDescription = "testDescription_" + rand
		expression      = createTestObfuscationExpression(t, testAccountID)
		testCreateInput = LogConfigurationsCreateObfuscationRuleInput{
			Description: testDescription,
			Name:        testName,
			Filter:      "entity.guids='MjUyMDUyOHxJTkZSQXxOQXwzMjI2NzYxMDM5Njk4NjQ3MTM2'",
			Enabled:     true,
			Actions: []LogConfigurationsCreateObfuscationActionInput{
				{
					Attributes:   []string{"awsAccountId"},
					ExpressionId: expression.ID,
					Method:       "MASK",
				},
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateObfuscationRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)

	//Test: Update - fail
	//Actions attribute not given
	update, err := client.LogConfigurationsUpdateObfuscationRule(testAccountID, LogConfigurationsUpdateObfuscationRuleInput{
		Name: testName + "_update",
		ID:   created.ID,
	})

	require.Error(t, err)
	require.Nil(t, update)
	require.Empty(t, update)

	//Test: Update
	//Actions attribute given
	update, err = client.LogConfigurationsUpdateObfuscationRule(testAccountID, LogConfigurationsUpdateObfuscationRuleInput{
		Name: testName + "_update",
		ID:   created.ID,
		Actions: []LogConfigurationsUpdateObfuscationActionInput{
			{
				Attributes:   []string{"message"},
				ExpressionId: expression.ID,
				Method:       "MASK",
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, update)
	require.NotEmpty(t, update)

	// Test: Delete
	testDeleteInput := created.ID
	deleted, err := client.LogConfigurationsDeleteObfuscationRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.NotEmpty(t, deleted)
}

func createTestObfuscationExpression(t *testing.T, testAccountID int) *LogConfigurationsObfuscationExpression {

	var (
		rand            = mock.RandSeq(5)
		testName        = "testName_" + rand
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsCreateObfuscationExpressionInput{

			Description: testDescription,
			Name:        testName,
			Regex:       "(^http.*)",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateObfuscationExpression(testAccountID, testCreateInput)
	require.NoError(t, err)

	return created
}
