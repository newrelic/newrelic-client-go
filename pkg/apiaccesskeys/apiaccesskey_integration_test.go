// +build integration

package apiaccesskeys

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAPIAccessKeys(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Setup
	createInput := APIAccessCreateKeysInput{
		Ingest: []APIAccessCreateIngestKeyInput{
			{
				AccountID:  testhelpers.TestAccountID,
				IngestType: "BROWSER",
				Name:       "test-integration-api-access",
				Notes:      "This ingest key was created by an integration test.",
			},
		},
	}

	deleteInput := APIAccessDeleteInput{
		IngestKeyIds: []string{},
	}

	// Test: Create
	createResult, err := client.CreateAPIAccessKeysMutation(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get
	getResult, err := client.GetAPIAccessKeyMutation(APIAccessGetInput{
		ID:      createResult[0].ID,
		KeyType: createResult[0].Type,
	})
	require.NoError(t, err)
	require.NotNil(t, getResult)

	// Test: Update
	updateResult, err := client.UpdateAPIAccessKeyMutation(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateKeyInput{
			{
				ID:    createResult[0].ID,
				Name:  createResult[0].Name,
				Notes: "testing notes update",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	deleteInput.IngestKeyIds = []string{createResult[0].ID}
	err = client.DeleteAPIAccessKeyMutation(deleteInput)
	require.NoError(t, err)
}

func newIntegrationTestClient(t *testing.T) APIAccessKeys {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
