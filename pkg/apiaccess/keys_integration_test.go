//go:build integration
// +build integration

package apiaccess

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
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
