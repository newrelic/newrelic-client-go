// +build integration

package workloads

import (
	"os"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

var (
	testWorkloadName        = "testWorkload"
	testUpdatedWorkloadName = testWorkloadName + "Updated"
	testWorkloadQuery       = "(name like 'tf_test' or id = 'tf_test' or domainId = 'tf_test')"
	testAccountID           = 2508259
	testCreateInput         = CreateInput{
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
	testUpdateInput = UpdateInput{
		Name: testUpdatedWorkloadName,
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

func TestIntegrationWorkload(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.CreateWorkload(testAccountID, &testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)

	// Test: Get
	workload, err := client.GetWorkload(testAccountID, *created.ID)

	require.NoError(t, err)
	require.NotNil(t, workload)

	// Test: List
	workloads, err := client.ListWorkloads(testAccountID)

	require.NoError(t, err)
	require.Greater(t, len(workloads), 0)

	// Test: Update
	time.Sleep(time.Second) // Updates are failing intermittently without this
	updated, err := client.UpdateWorkload(*created.GUID, &testUpdateInput)

	require.NoError(t, err)
	require.NotNil(t, workload)
	require.Equal(t, testUpdateInput.Name, *updated.Name)

	// Test: Duplicate
	duplicateInput := DuplicateInput{
		Name: "duplicateWorkload",
	}
	duplicate, err := client.DuplicateWorkload(testAccountID, *created.GUID, &duplicateInput)

	require.NoError(t, err)
	require.NotNil(t, duplicate)
	require.Equal(t, "duplicateWorkload", *duplicate.Name)

	// Test: Delete
	deleted, err := client.DeleteWorkload(*created.GUID)

	require.NoError(t, err)
	require.NotNil(t, deleted)

	deleted, err = client.DeleteWorkload(*duplicate.GUID)

	require.NoError(t, err)
	require.NotNil(t, deleted)
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
