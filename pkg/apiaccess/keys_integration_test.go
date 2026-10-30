//go:build integration
// +build integration

package apiaccess

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationAPIAccess_IngestKeys(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Setup
	createInput := APIAccessCreateInput{
		Ingest: []APIAccessCreateIngestKeyInput{
			{
				AccountID:  testAccountID,
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

// getTestTeamGUIDs returns the comma-separated team entity GUIDs in
// INTEGRATION_TESTING_NEW_RELIC_TEAM_GUIDS, or nil when unset.
func getTestTeamGUIDs() []string {
	v := os.Getenv("INTEGRATION_TESTING_NEW_RELIC_TEAM_GUIDS")
	if v == "" {
		return nil
	}
	return strings.Split(v, ",")
}

func TestIntegrationAPIAccess_IngestKeyEnabledAndOwningTeams(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	teamGUIDs := getTestTeamGUIDs()

	// Test: Create (with owning teams when configured; otherwise OwningTeamIDs is nil and sent as null)
	createResult, err := client.CreateAPIAccessKeys(APIAccessCreateInput{
		Ingest: []APIAccessCreateIngestKeyInput{{
			AccountID:     testAccountID,
			IngestType:    "BROWSER",
			Name:          "test-integration-api-access-enabled-" + mock.RandSeq(5),
			OwningTeamIDs: teamGUIDs,
		}},
	})
	require.NoError(t, err)
	require.Len(t, createResult, 1)
	keyID := createResult[0].ID
	require.True(t, createResult[0].Enabled)
	require.ElementsMatch(t, teamGUIDs, createResult[0].OwningTeamIDs)

	if len(teamGUIDs) > 0 {
		// Test: Replace the owning team set
		replaced := teamGUIDs[:1]
		updateResult, updateErr := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
			Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: keyID, OwningTeamIDs: replaced}},
		})
		require.NoError(t, updateErr)
		require.Equal(t, replaced, updateResult[0].OwningTeamIDs)

		// Test: An update without OwningTeamIDs (sent as null) leaves the owning teams unchanged
		updateResult, updateErr = client.UpdateAPIAccessKeys(APIAccessUpdateInput{
			Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: keyID, Notes: "teams should be unchanged"}},
		})
		require.NoError(t, updateErr)
		require.Equal(t, replaced, updateResult[0].OwningTeamIDs)

		// Test: Clear all owning teams
		updateResult, updateErr = client.UpdateAPIAccessKeys(APIAccessUpdateInput{
			Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: keyID, OwningTeamIDs: []string{}}},
		})
		require.NoError(t, updateErr)
		require.Empty(t, updateResult[0].OwningTeamIDs)
	}

	// Test: Disable
	disabled := false
	updateResult, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: keyID, Enabled: &disabled}},
	})
	require.NoError(t, err)
	require.False(t, updateResult[0].Enabled)

	// Test: A disabled key is not returned by default...
	_, err = client.GetAPIAccessKey(keyID, APIAccessKeyTypeTypes.INGEST)
	require.ErrorContains(t, err, "Key not found")

	// ...but is returned when including disabled keys
	getResult, err := client.GetAPIAccessKeyIncludingDisabled(keyID, APIAccessKeyTypeTypes.INGEST)
	require.NoError(t, err)
	require.False(t, getResult.Enabled)

	// Test: Re-enable
	enabled := true
	updateResult, err = client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: keyID, Enabled: &enabled}},
	})
	require.NoError(t, err)
	require.True(t, updateResult[0].Enabled)

	// Test: Delete a disabled key
	_, err = client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: keyID, Enabled: &disabled}},
	})
	require.NoError(t, err)

	deleteResult, err := client.DeleteAPIAccessKey(APIAccessDeleteInput{IngestKeyIDs: []string{keyID}})
	require.NoError(t, err)
	require.Len(t, deleteResult, 1)
}

func TestIntegrationAPIAccess_UserKeyEnabled(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	testUserID, err := mock.GetTestUserID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Setup
	createResult, err := client.CreateAPIAccessKeys(APIAccessCreateInput{
		User: []APIAccessCreateUserKeyInput{{
			AccountID: testAccountID,
			UserID:    testUserID,
			Name:      "test-integration-api-access-enabled-" + mock.RandSeq(5),
		}},
	})
	require.NoError(t, err)
	require.Len(t, createResult, 1)
	keyID := createResult[0].ID

	// Test: Disable
	disabled := false
	updateResult, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		User: []APIAccessUpdateUserKeyInput{{KeyID: keyID, Enabled: &disabled}},
	})
	require.NoError(t, err)
	require.False(t, updateResult[0].Enabled)

	getResult, err := client.GetAPIAccessKeyIncludingDisabled(keyID, APIAccessKeyTypeTypes.USER)
	require.NoError(t, err)
	require.False(t, getResult.Enabled)

	// Test: Delete a disabled key
	deleteResult, err := client.DeleteAPIAccessKey(APIAccessDeleteInput{UserKeyIDs: []string{keyID}})
	require.NoError(t, err)
	require.Len(t, deleteResult, 1)
}

func TestIntegrationAPIAccess_UserKeys(t *testing.T) {
	t.Parallel()

	userID, err := mock.GetTestUserID()
	if err != nil {
		t.Skipf("Skipping `TestIntegrationAPIAccess_UserKeys` integration test due error: %v", err)
		return
	}

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Setup
	createInput := APIAccessCreateInput{
		User: []APIAccessCreateUserKeyInput{
			{
				AccountID: testAccountID,
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
			AccountIDs: []int{testAccountID},
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

func TestIntegrationAPIAccess_UpdateIngestKeyError(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	mockIngestKeyID, _ := mock.GetNonExistentIDs()
	_, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{
			{
				KeyID: mockIngestKeyID,
				Name:  "Lorem Ipsum",
				Notes: "Lorem Ipsum",
			},
		},
	})
	require.Error(t, err)
	require.Equal(t, validateAPIAccessKeyError(err, "INGEST"), true)
}

func TestIntegrationAPIAccess_UpdateUserKeyError(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	_, mockUserKeyID := mock.GetNonExistentIDs()
	_, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		User: []APIAccessUpdateUserKeyInput{
			{
				KeyID: mockUserKeyID,
				Name:  "Lorem Ipsum",
				Notes: "Lorem Ipsum",
			},
		},
	})
	require.Error(t, err)
	require.Equal(t, validateAPIAccessKeyError(err, "USER"), true)
}

func TestIntegrationAPIAccess_DeleteIngestKeyError(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	mockIngestKeyIDOne, mockIngestKeyIDTwo := mock.GetNonExistentIDs()
	_, err := client.DeleteAPIAccessKey(APIAccessDeleteInput{
		IngestKeyIDs: []string{
			mockIngestKeyIDOne,
			mockIngestKeyIDTwo,
		},
	})
	require.Error(t, err)
	require.Equal(t, validateAPIAccessKeyError(err, "INGEST"), true)
}

func TestIntegrationAPIAccess_DeleteUserKeyError(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	_, mockUserKeyID := mock.GetNonExistentIDs()
	_, err := client.DeleteAPIAccessKey(APIAccessDeleteInput{
		UserKeyIDs: []string{
			mockUserKeyID,
		},
	})
	require.Error(t, err)
	require.Equal(t, validateAPIAccessKeyError(err, "USER"), true)
}

var possibleIngestKeyErrors = []string{
	string(APIAccessIngestKeyErrorTypeTypes.NOT_FOUND),
	string(APIAccessIngestKeyErrorTypeTypes.INVALID),
	string(APIAccessIngestKeyErrorTypeTypes.FORBIDDEN),
}

var possibleUserKeyErrors = []string{
	string(APIAccessUserKeyErrorTypeTypes.NOT_FOUND),
	string(APIAccessUserKeyErrorTypeTypes.INVALID),
	string(APIAccessUserKeyErrorTypeTypes.FORBIDDEN),
}

func validateAPIAccessKeyError(err error, errorType string) (status bool) {
	errorMessage := err.Error()
	listOfErrorMessages := strings.Split(errorMessage, "\n")
	listOfErrorMessages = listOfErrorMessages[1 : len(listOfErrorMessages)-1]
	validatedErrorMessages := 0

	if errorType == "USER" {
		for errorIndex := 0; errorIndex < len(listOfErrorMessages); errorIndex++ {
			for listIndex := 0; listIndex < len(possibleUserKeyErrors); listIndex++ {
				match, _ := regexp.MatchString(possibleUserKeyErrors[listIndex], listOfErrorMessages[errorIndex])
				if match {
					validatedErrorMessages += 1
					break
				}
			}
		}
	} else if errorType == "INGEST" {
		for errorIndex := 0; errorIndex < len(listOfErrorMessages); errorIndex++ {
			for listIndex := 0; listIndex < len(possibleIngestKeyErrors); listIndex++ {
				match, _ := regexp.MatchString(possibleIngestKeyErrors[listIndex], listOfErrorMessages[errorIndex])
				if match {
					validatedErrorMessages += 1
					break
				}
			}
		}
	}

	if validatedErrorMessages == len(listOfErrorMessages) {
		return true
	}

	return false
}
