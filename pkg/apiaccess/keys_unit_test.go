//go:build unit
// +build unit

package apiaccess

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// newMockAPIAccess creates an APIAccess client backed by a single-response mock server.
func newMockAPIAccess(t *testing.T, mockJSONResponse string, statusCode int) APIAccess {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	t.Cleanup(ts.Close)
	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

// multiResponseServer creates a mock server that returns responses from the provided
// slice in order. Once the slice is exhausted it panics to signal unexpected calls.
func multiResponseServer(t *testing.T, responses []string) *httptest.Server {
	t.Helper()
	callIndex := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if callIndex >= len(responses) {
			t.Fatalf("unexpected request #%d — only %d responses were configured", callIndex+1, len(responses))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(responses[callIndex]))
		require.NoError(t, err)
		callIndex++
	}))
	t.Cleanup(ts.Close)
	return ts
}

// requestCapturingServer creates a mock server that records all request bodies and
// returns the same response on every call.
func requestCapturingServer(t *testing.T, mockJSONResponse string) (*httptest.Server, *[][]byte) {
	t.Helper()
	var bodies [][]byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		bodies = append(bodies, body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(mockJSONResponse))
		require.NoError(t, err)
	}))
	t.Cleanup(ts.Close)
	return ts, &bodies
}

// --- TestGetPreferredIngestLicenseKey_FindsByName ---

func TestGetPreferredIngestLicenseKey_FindsByName(t *testing.T) {
	t.Parallel()

	mockResponse := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{
								"id": "key-id-1",
								"key": "NRAL-preferred-key-value",
								"name": "Installer Ingest License Key",
								"createdAt": 1000,
								"type": "INGEST",
								"ingestType": "LICENSE"
							}
						],
						"nextCursor": ""
					}
				}
			}
		}
	}`

	client := newMockAPIAccess(t, mockResponse, http.StatusOK)
	keyValue, err := client.GetPreferredIngestLicenseKey(context.Background(), 12345)
	require.NoError(t, err)
	assert.Equal(t, "NRAL-preferred-key-value", keyValue)
}

// --- TestGetPreferredIngestLicenseKey_FallsBackToOldest ---

func TestGetPreferredIngestLicenseKey_FallsBackToOldest(t *testing.T) {
	t.Parallel()

	// First call: name search returns no keys (preferred name not found)
	nameSearchResponse := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [],
						"nextCursor": ""
					}
				}
			}
		}
	}`

	// Second call: full search returns 3 keys with different creation times
	fullSearchResponse := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{
								"id": "key-id-newest",
								"key": "NRAL-newest",
								"name": "Key C",
								"createdAt": 3000,
								"type": "INGEST",
								"ingestType": "LICENSE"
							},
							{
								"id": "key-id-oldest",
								"key": "NRAL-oldest",
								"name": "Key A",
								"createdAt": 1000,
								"type": "INGEST",
								"ingestType": "LICENSE"
							},
							{
								"id": "key-id-middle",
								"key": "NRAL-middle",
								"name": "Key B",
								"createdAt": 2000,
								"type": "INGEST",
								"ingestType": "LICENSE"
							}
						],
						"nextCursor": ""
					}
				}
			}
		}
	}`

	ts := multiResponseServer(t, []string{nameSearchResponse, fullSearchResponse})
	tc := mock.NewTestConfig(t, ts)
	client := New(tc)

	keyValue, err := client.GetPreferredIngestLicenseKey(context.Background(), 12345)
	require.NoError(t, err)
	assert.Equal(t, "NRAL-oldest", keyValue)
}

// --- TestGetPreferredIngestLicenseKey_ReturnsEmptyWhenNoKeys ---

func TestGetPreferredIngestLicenseKey_ReturnsEmptyWhenNoKeys(t *testing.T) {
	t.Parallel()

	emptyResponse := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [],
						"nextCursor": ""
					}
				}
			}
		}
	}`

	// Both name search and full search return empty results
	ts := multiResponseServer(t, []string{emptyResponse, emptyResponse})
	tc := mock.NewTestConfig(t, ts)
	client := New(tc)

	keyValue, err := client.GetPreferredIngestLicenseKey(context.Background(), 12345)
	require.NoError(t, err)
	assert.Equal(t, "", keyValue)
}

// --- TestSearchAPIAccessKeysWithContext_FollowsNextCursor ---

type graphQLReqBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func TestSearchAPIAccessKeysWithContext_FollowsNextCursor(t *testing.T) {
	t.Parallel()

	page1Response := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{
								"id": "key-id-1",
								"key": "NRAL-key-1",
								"name": "Key 1",
								"createdAt": 1000,
								"type": "INGEST",
								"ingestType": "LICENSE"
							}
						],
						"nextCursor": "cursor-for-page-2"
					}
				}
			}
		}
	}`

	page2Response := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{
								"id": "key-id-2",
								"key": "NRAL-key-2",
								"name": "Key 2",
								"createdAt": 2000,
								"type": "INGEST",
								"ingestType": "LICENSE"
							}
						],
						"nextCursor": ""
					}
				}
			}
		}
	}`

	var capturedBodies [][]byte
	callIndex := 0
	responses := []string{page1Response, page2Response}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		capturedBodies = append(capturedBodies, body)

		if callIndex >= len(responses) {
			t.Fatalf("unexpected request #%d", callIndex+1)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(responses[callIndex]))
		require.NoError(t, err)
		callIndex++
	}))
	t.Cleanup(ts.Close)

	tc := mock.NewTestConfig(t, ts)
	client := New(tc)

	keys, err := client.SearchAPIAccessKeysWithContext(context.Background(), APIAccessKeySearchQuery{
		Scope: APIAccessKeySearchScope{
			AccountIDs:  []int{12345},
			IngestTypes: []APIAccessIngestKeyType{APIAccessIngestKeyTypeTypes.LICENSE},
		},
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.INGEST},
	})

	require.NoError(t, err)
	require.Len(t, keys, 2, "expected keys from both pages to be collected")
	assert.Equal(t, "NRAL-key-1", keys[0].Key)
	assert.Equal(t, "NRAL-key-2", keys[1].Key)

	// Verify the second request sent the cursor from page 1
	require.Len(t, capturedBodies, 2, "expected exactly 2 requests to be made")

	var req1, req2 graphQLReqBody
	require.NoError(t, json.Unmarshal(capturedBodies[0], &req1))
	require.NoError(t, json.Unmarshal(capturedBodies[1], &req2))

	// First request cursor should be nil (not set)
	assert.Nil(t, req1.Variables["cursor"])

	// Second request cursor should be the value from page 1's nextCursor
	assert.Equal(t, "cursor-for-page-2", req2.Variables["cursor"])

	// Verify the query includes the cursor parameter
	assert.Contains(t, req1.Query, "$cursor: String")
	assert.Contains(t, req1.Query, "cursor: $cursor")
}
