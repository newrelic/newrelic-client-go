//go:build unit
// +build unit

package apiaccess

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

const (
	testUpdateKeysMockResponse = `{
		"data": {
			"apiAccessUpdateKeys": {
				"updatedKeys": [{
					"id": "A1B2C3",
					"name": "test-key",
					"type": "INGEST",
					"enabled": false,
					"ingestType": "LICENSE",
					"accountId": 1,
					"owningTeamIds": []
				}],
				"errors": []
			}
		}
	}`

	testUpdateUserKeysMockResponse = `{
		"data": {
			"apiAccessUpdateKeys": {
				"updatedKeys": [{
					"id": "D4E5F6",
					"name": "test-key",
					"type": "USER",
					"enabled": false,
					"accountId": 1,
					"userId": 1
				}],
				"errors": []
			}
		}
	}`

	testCreateKeysMockResponse = `{
		"data": {
			"apiAccessCreateKeys": {
				"createdKeys": [{
					"id": "A1B2C3",
					"type": "INGEST",
					"enabled": true,
					"ingestType": "LICENSE",
					"accountId": 1,
					"owningTeamIds": []
				}],
				"errors": []
			}
		}
	}`

	testGetKeyMockResponse = `{
		"data": {
			"actor": {
				"apiAccess": {
					"key": {
						"id": "A1B2C3",
						"name": "test-key",
						"type": "INGEST",
						"enabled": false,
						"ingestType": "LICENSE",
						"accountId": 1,
						"owningTeamIds": ["team-guid-1", "team-guid-2"]
					}
				}
			}
		}
	}`
)

type graphQLRequestBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// newMockResponseCapturingRequest returns a client pointed at a mock server, plus the
// decoded body of the last request it received, so tests can assert on the GraphQL
// variables actually sent over the wire.
func newMockResponseCapturingRequest(t *testing.T, mockJSONResponse string) (APIAccess, *graphQLRequestBody) {
	lastRequest := &graphQLRequestBody{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, lastRequest))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte(mockJSONResponse))
		require.NoError(t, err)
	}))
	t.Cleanup(ts.Close)

	return New(mock.NewTestConfig(t, ts)), lastRequest
}

func boolPtr(b bool) *bool { return &b }

// updateKeyVars returns the variables sent for the single key in the "ingest" or "user" update list.
func updateKeyVars(t *testing.T, req *graphQLRequestBody, keyList string) map[string]interface{} {
	keys, ok := req.Variables["keys"].(map[string]interface{})
	require.True(t, ok)
	list, ok := keys[keyList].([]interface{})
	require.True(t, ok)
	require.Len(t, list, 1)
	return list[0].(map[string]interface{})
}

func TestUnitAPIAccess_UpdateKeys_SendsEnabledFalse(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testUpdateKeysMockResponse)

	keys, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: "A1B2C3", Enabled: boolPtr(false)}},
	})
	require.NoError(t, err)

	vars := updateKeyVars(t, req, "ingest")
	assert.Equal(t, false, vars["enabled"])
	require.Len(t, keys, 1)
	assert.False(t, keys[0].Enabled)
}

func TestUnitAPIAccess_UpdateKeys_UnsetEnabledOmittedAndOwningTeamIDsNull(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testUpdateKeysMockResponse)

	_, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: "A1B2C3", Name: "renamed"}},
	})
	require.NoError(t, err)

	vars := updateKeyVars(t, req, "ingest")
	assert.Equal(t, "renamed", vars["name"])
	assert.NotContains(t, vars, "enabled")
	// A nil OwningTeamIDs is sent as null, which leaves the key's teams unchanged.
	require.Contains(t, vars, "owningTeamIds")
	assert.Nil(t, vars["owningTeamIds"])
}

func TestUnitAPIAccess_UpdateKeys_SendsEmptyOwningTeamIDsToClear(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testUpdateKeysMockResponse)

	_, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		Ingest: []APIAccessUpdateIngestKeyInput{{KeyID: "A1B2C3", OwningTeamIDs: []string{}}},
	})
	require.NoError(t, err)

	vars := updateKeyVars(t, req, "ingest")
	require.Contains(t, vars, "owningTeamIds")
	assert.Equal(t, []interface{}{}, vars["owningTeamIds"])
}

func TestUnitAPIAccess_UpdateUserKeys_SendsEnabledFalse(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testUpdateUserKeysMockResponse)

	keys, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		User: []APIAccessUpdateUserKeyInput{{KeyID: "D4E5F6", Enabled: boolPtr(false)}},
	})
	require.NoError(t, err)

	vars := updateKeyVars(t, req, "user")
	assert.Equal(t, false, vars["enabled"])
	require.Len(t, keys, 1)
	assert.False(t, keys[0].Enabled)
}

func TestUnitAPIAccess_UpdateUserKeys_UnsetEnabledOmitted(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testUpdateUserKeysMockResponse)

	_, err := client.UpdateAPIAccessKeys(APIAccessUpdateInput{
		User: []APIAccessUpdateUserKeyInput{{KeyID: "D4E5F6", Name: "renamed"}},
	})
	require.NoError(t, err)

	vars := updateKeyVars(t, req, "user")
	assert.Equal(t, "renamed", vars["name"])
	assert.NotContains(t, vars, "enabled")
}

func TestUnitAPIAccess_CreateKeys_UnsetOwningTeamIDsNull(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testCreateKeysMockResponse)

	_, err := client.CreateAPIAccessKeys(APIAccessCreateInput{
		Ingest: []APIAccessCreateIngestKeyInput{{AccountID: 1, IngestType: APIAccessIngestKeyTypeTypes.LICENSE}},
	})
	require.NoError(t, err)

	keys, ok := req.Variables["keys"].(map[string]interface{})
	require.True(t, ok)
	ingest := keys["ingest"].([]interface{})[0].(map[string]interface{})
	// A nil OwningTeamIDs is sent as null, which attaches no teams.
	require.Contains(t, ingest, "owningTeamIds")
	assert.Nil(t, ingest["owningTeamIds"])
}

func TestUnitAPIAccess_GetKey_DefaultOmitsIncludeDisabled(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testGetKeyMockResponse)

	_, err := client.GetAPIAccessKey("A1B2C3", APIAccessKeyTypeTypes.INGEST)
	require.NoError(t, err)

	assert.NotContains(t, req.Query, "includeDisabled")
}

func TestUnitAPIAccess_GetKeyIncludingDisabled(t *testing.T) {
	t.Parallel()
	client, req := newMockResponseCapturingRequest(t, testGetKeyMockResponse)

	key, err := client.GetAPIAccessKeyIncludingDisabled("A1B2C3", APIAccessKeyTypeTypes.INGEST)
	require.NoError(t, err)

	assert.Contains(t, req.Query, "includeDisabled: true")
	assert.False(t, key.Enabled)
	assert.Equal(t, []string{"team-guid-1", "team-guid-2"}, key.OwningTeamIDs)
}
