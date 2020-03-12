// +build integration

package workloads

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

var (
	testWorkloadName  = "testWorkload"
	testWorkloadQuery = "(name like 'tf_test' or id = 'tf_test' or domainId = 'tf_test')"
	testAccountID     = 2508259
	testCreateInput   = CreateInput{
		Name: &testWorkloadName,
		ScopeAccountsInput: ScopeAccountsInput{
			AccountIDs: []*int{&testAccountID},
		},
		EntitySearchQueries: []*EntitySearchQueryInput{
			{
				Name:  "testQuery",
				Query: &testWorkloadQuery,
			},
		},
	}
)

func TestIntegrationGetWorkload(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	actual, err := client.GetWorkload(testAccountID, 791)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

func TestIntegrationListWorkloads(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	actual, err := client.ListWorkloads(2508259)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

func TestIntegrationCreateWorkload(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	actual, err := client.CreateWorkload(testAccountID, &testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

// nolint
func newIntegrationTestClient(t *testing.T) Workloads {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing for graphql requires your personal API key")
	}

	return New(config.Config{
		PersonalAPIKey: apiKey,
		UserAgent:      "newrelic/newrelic-client-go",
		LogLevel:       "debug",
	})
}
