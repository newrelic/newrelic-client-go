//go:build integration
// +build integration

package logconfigurations

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationDataPartitionRule(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping `TestIntegrationDataPartitionRule` as the API endpoint is not recommended.")
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "Log_testName_" + rand
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsCreateDataPartitionRuleInput{
			Description:         testDescription,
			TargetDataPartition: LogConfigurationsLogDataPartitionName(testName),
			Enabled:             true,
			MatchingCriteria: &LogConfigurationsDataPartitionRuleMatchingCriteriaInput{
				AttributeName:      "hostname",
				MatchingExpression: "localhost",
				MatchingMethod:     "EQUALS",
			},
			RetentionPolicy: "SECONDARY",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateDataPartitionRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 0, len(created.Errors))

	// Test: Delete
	testDeleteInput := created.Rule.ID
	deleted, err := client.LogConfigurationsDeleteDataPartitionRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	// This is returning nil after successful delete
	require.Empty(t, deleted)
	require.Equal(t, 0, len(deleted.Errors))

}

// Create with invalid name
// It should always begin with Log_
func TestIntegrationDataPartitionRule_ValidName(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping `TestIntegrationDataPartitionRule_ValidName` as the API endpoint is not recommended.")

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "testName_" + rand
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsCreateDataPartitionRuleInput{
			Description:         testDescription,
			TargetDataPartition: LogConfigurationsLogDataPartitionName(testName),
			Enabled:             true,
			MatchingCriteria: &LogConfigurationsDataPartitionRuleMatchingCriteriaInput{
				AttributeName:      "hostname",
				MatchingExpression: "localhost",
				MatchingMethod:     "EQUALS",
			},
			RetentionPolicy: "SECONDARY",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateDataPartitionRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 1, len(created.Errors))

	// Test: Delete
	testDeleteInput := created.Rule.ID
	deleted, err := client.LogConfigurationsDeleteDataPartitionRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	// This is returning nil after successful delete
	require.NotEmpty(t, deleted)
	require.Equal(t, 1, len(deleted.Errors))
}

// Create with invalid name
// It should always be Unique
func TestIntegrationDataPartitionRule_DuplicateName(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping `TestIntegrationDataPartitionRule_DuplicateName` as the API endpoint is not recommended.")

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "Log_testName_" + rand
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsCreateDataPartitionRuleInput{
			Description:         testDescription,
			TargetDataPartition: LogConfigurationsLogDataPartitionName(testName),
			Enabled:             true,
			MatchingCriteria: &LogConfigurationsDataPartitionRuleMatchingCriteriaInput{
				AttributeName:      "hostname",
				MatchingExpression: "localhost",
				MatchingMethod:     "EQUALS",
			},
			RetentionPolicy: "SECONDARY",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	_, err = client.LogConfigurationsCreateDataPartitionRule(testAccountID, testCreateInput)
	created, err := client.LogConfigurationsCreateDataPartitionRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 1, len(created.Errors))

	// Test: Delete
	testDeleteInput := created.Rule.ID
	deleted, err := client.LogConfigurationsDeleteDataPartitionRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	// This is returning nil after successful delete
	require.NotEmpty(t, deleted)
	require.Equal(t, 1, len(deleted.Errors))
}

// Create with invalid attributeName
// It must be a combination of a-z A-Z 0-9
func TestIntegrationDataPartitionRule_attributeName(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping `TestIntegrationDataPartitionRule_attributeName` as the API endpoint is not recommended.")

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "Log_testName_" + rand
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsCreateDataPartitionRuleInput{
			Description:         testDescription,
			TargetDataPartition: LogConfigurationsLogDataPartitionName(testName),
			Enabled:             true,
			MatchingCriteria: &LogConfigurationsDataPartitionRuleMatchingCriteriaInput{
				AttributeName:      "Select * form None",
				MatchingExpression: "localhost",
				MatchingMethod:     "EQUALS",
			},
			RetentionPolicy: "SECONDARY",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateDataPartitionRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 1, len(created.Errors))

	// Test: Delete
	testDeleteInput := created.Rule.ID
	deleted, err := client.LogConfigurationsDeleteDataPartitionRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	// This is returning nil after successful delete
	require.NotEmpty(t, deleted)
	require.Equal(t, 1, len(deleted.Errors))
}

// Update
func TestIntegrationDataPartitionRuleUpdate(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping `TestIntegrationDataPartitionRuleUpdate` as the API endpoint is not recommended.")

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "Log_testName_" + rand
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsCreateDataPartitionRuleInput{
			Description:         testDescription,
			TargetDataPartition: LogConfigurationsLogDataPartitionName(testName),
			Enabled:             true,
			MatchingCriteria: &LogConfigurationsDataPartitionRuleMatchingCriteriaInput{
				AttributeName:      "hostname",
				MatchingExpression: "localhost",
				MatchingMethod:     "EQUALS",
			},
			RetentionPolicy: "SECONDARY",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateDataPartitionRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, 0, len(created.Errors))

	//Test: Update
	update, err := client.LogConfigurationsUpdateDataPartitionRule(testAccountID, LogConfigurationsUpdateDataPartitionRuleInput{
		Enabled: false,
		ID:      created.Rule.ID,
	})

	require.NoError(t, err)
	require.NotNil(t, update)
	require.NotEmpty(t, update)
	require.Equal(t, 0, len(update.Errors))

	// Test: Delete
	testDeleteInput := update.Rule.ID
	deleted, err := client.LogConfigurationsDeleteDataPartitionRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	// This is returning nil after successful delete
	require.Empty(t, deleted)
	require.Equal(t, 0, len(deleted.Errors))
}
