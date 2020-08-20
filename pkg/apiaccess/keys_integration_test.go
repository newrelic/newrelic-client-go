// +build integration

package apiaccess

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationAPIAccess_IngestKeys(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Setup
	createInput := APIAccessCreateInput{
		Ingest: []APIAccessCreateIngestKeyInput{
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
	updateResult, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{
			{
				KeyID: createResult[0].ID,
				Name:  createResult[0].Name,
				Notes: "testing notes update",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	deleteResult, err := client.DeleteAPIAccessKey(APIAccessDeleteInput{
		IngestKeyIDs: []string{createResult[0].ID},
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
	createInput := APIAccessCreateInput{
		User: []APIAccessCreateUserKeyInput{
			{
				AccountID: testhelpers.TestAccountID,
				Name:      "test-integration-api-access",
				Notes:     "This user key was created by an integration test.",
				UserID:    userID,
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
	searchResult, err := client.SearchAPIAccessKeys(APIAccessKeySearchQuery{
		Scope: APIAccessKeySearchScope{
			AccountIDs: []int{testhelpers.TestAccountID},
		},
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.USER},
	})
	require.NoError(t, err)
	require.Greater(t, len(searchResult), 0)

	// Test: Update
	updateResult, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		User: []APIAccessUpdateUserKeyInput{
			{
				KeyID: createResult[0].ID,
				Name:  createResult[0].Name,
				Notes: "testing notes update",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, "testing notes update", updateResult[0].Notes)

	// Test: Delete
	deleteResult, err := client.DeleteAPIAccessKey(APIAccessDeleteInput{
		UserKeyIDs: []string{createResult[0].ID},
	})
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

func newIntegrationTestClient(t *testing.T) APIAccess {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
