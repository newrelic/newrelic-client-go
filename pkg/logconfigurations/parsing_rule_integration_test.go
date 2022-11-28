//go:build integration
// +build integration

package logconfigurations

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationParsingRule(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsParsingRuleConfiguration{
			Attribute:   "attribute",
			Description: testDescription,
			Enabled:     true,
			Grok:        "sampleattribute=%{NUMBER:test:int}",
			Lucene:      "logtype:linux_messages",
			NRQL:        "SELECT * FROM Log WHERE logtype = 'linux_messages'",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateParsingRule(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)

	// Get API
	rules, err := client.GetParsingRules(testAccountID)
	require.NoError(t, err)
	require.NotEmpty(t, rules)
	require.NotNil(t, rules)

	// Test: Update
	update, err := client.LogConfigurationsUpdateParsingRule(testAccountID, created.Rule.ID, LogConfigurationsParsingRuleConfiguration{

		Attribute:   "attribute",
		Description: testDescription + "_update",
		Enabled:     true,
		Grok:        "sampleattribute=%{NUMBER:test:int}",
		Lucene:      "logtype:linux_messages",
		NRQL:        "SELECT * FROM Log WHERE logtype = 'linux_messages'",
	})

	require.NoError(t, err)
	require.NotNil(t, update)
	require.NotEmpty(t, update)

	// Test: Delete
	testDeleteInput := created.Rule.ID
	deleted, err := client.LogConfigurationsDeleteParsingRule(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.Empty(t, deleted)

	rules, err = client.GetParsingRules(testAccountID)
	require.NoError(t, err)
	require.NotEmpty(t, rules)
	require.NotNil(t, rules)
	require.NotEqual(t, len(created.Rule), len(rules))
}
