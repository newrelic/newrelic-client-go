// +build integration

package apiaccess

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAPIAccess_IngestKeys(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Setup
	createInput := ApiAccessCreateInput{
		Ingest: []ApiAccessCreateIngestKeyInput{
			{
				AccountID:  testhelpers.TestAccountID,
				IngestType: "BROWSER",
				Name:       "test-integration-api-access",
				Notes:      "This ingest key was created by an integration test.",
			},
		},
	}

	// Test: Create
	createResult, err := client.CreateAPIAccessKeys(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get
	getResult, err := client.GetAPIAccessKey(createResult[0].ID, createResult[0].Type)
	require.NoError(t, err)
	require.NotNil(t, getResult)

	// Test: Update
	updateResult, err := client.UpdateAPIAccessKeys(ApiAccessUpdateInput{
		Ingest: []ApiAccessUpdateIngestKeyInput{
			{
				KeyId: createResult[0].ID,
				Name:  createResult[0].Name,
				Notes: "testing notes update",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	deleteResult, err := client.DeleteAPIAccessKey(ApiAccessDeleteInput{
		IngestKeyIds: []string{createResult[0].ID},
	})
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

func TestIntegrationAPIAccess_UserKeys(t *testing.T) {
	t.Parallel()

	userID, err := testhelpers.GetTestUserID()
	if err != nil {
		t.Skipf("Skipping `TestIntegrationAPIAccess_UserKeys` integration test due error: %v", err)
		return
	}

	client := newIntegrationTestClient(t)

	// Setup
	createInput := ApiAccessCreateInput{
		User: []ApiAccessCreateUserKeyInput{
			{
				AccountID: testhelpers.TestAccountID,
				Name:      "test-integration-api-access",
				Notes:     "This user key was created by an integration test.",
				UserId:    userID,
			},
		},
	}

	// Test: Create
	createResult, err := client.CreateAPIAccessKeys(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get
	getResult, err := client.GetAPIAccessKey(createResult[0].ID, createResult[0].Type)
	require.NoError(t, err)
	require.NotNil(t, getResult)

	// Test: Search
	searchResult, err := client.SearchAPIAccessKeys(ApiAccessKeySearchQuery{
		Scope: ApiAccessKeySearchScope{
			AccountIds: []int{testhelpers.TestAccountID},
		},
		Types: []ApiAccessKeyType{ApiAccessKeyTypeTypes.INGEST},
	})
	require.NoError(t, err)
	require.Greater(t, len(searchResult), 0)

	// Test: Update
	updateResult, err := client.UpdateAPIAccessKeys(ApiAccessUpdateInput{
		User: []ApiAccessUpdateUserKeyInput{
			{
				KeyId: createResult[0].ID,
				Name:  createResult[0].Name,
				Notes: "testing notes update",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, "testing notes update", updateResult[0].Notes)

	// Test: Delete
	deleteResult, err := client.DeleteAPIAccessKey(ApiAccessDeleteInput{
		UserKeyIds: []string{createResult[0].ID},
	})
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

func newIntegrationTestClient(t *testing.T) APIAccess {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
