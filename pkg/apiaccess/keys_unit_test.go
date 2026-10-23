//go:build unit
// +build unit

package apiaccess

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// keySearchRequestBody is used to parse the NerdGraph request body in tests.
type keySearchRequestBody struct {
	Variables map[string]interface{} `json:"variables"`
}

func newAPIAccessClient(t *testing.T, ts *httptest.Server) APIAccess {
	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

func TestSearchAPIAccessKeysWithContext_FollowsNextCursor(t *testing.T) {
	t.Parallel()

	page1 := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{"id": "key-1", "name": "Key One", "type": "INGEST"},
							{"id": "key-2", "name": "Key Two", "type": "INGEST"}
						],
						"nextCursor": "page-2"
					}
				}
			}
		}
	}`

	page2 := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{"id": "key-3", "name": "Key Three", "type": "INGEST"}
						],
						"nextCursor": null
					}
				}
			}
		}
	}`

	requestCount := 0
	var capturedCursors []interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req keySearchRequestBody
		require.NoError(t, json.Unmarshal(body, &req))

		capturedCursors = append(capturedCursors, req.Variables["cursor"])

		w.Header().Set("Content-Type", "application/json")
		if requestCount == 1 {
			_, _ = w.Write([]byte(page1))
		} else {
			_, _ = w.Write([]byte(page2))
		}
	}))
	defer ts.Close()

	client := newAPIAccessClient(t, ts)

	params := APIAccessKeySearchQuery{
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.INGEST},
	}

	keys, err := client.SearchAPIAccessKeys(params)

	require.NoError(t, err)
	assert.Len(t, keys, 3)
	assert.Equal(t, "key-1", keys[0].ID)
	assert.Equal(t, "key-2", keys[1].ID)
	assert.Equal(t, "key-3", keys[2].ID)

	assert.Equal(t, 2, requestCount, "expected exactly 2 requests")

	// First request: cursor should be nil
	assert.Nil(t, capturedCursors[0])

	// Second request: cursor should be "page-2"
	assert.Equal(t, "page-2", capturedCursors[1])
}

func TestSearchAPIAccessKeysWithContext_SinglePage(t *testing.T) {
	t.Parallel()

	singlePage := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{"id": "key-only", "name": "Only Key", "type": "INGEST"}
						]
					}
				}
			}
		}
	}`

	requestCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(singlePage))
	}))
	defer ts.Close()

	client := newAPIAccessClient(t, ts)

	params := APIAccessKeySearchQuery{
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.INGEST},
	}

	keys, err := client.SearchAPIAccessKeys(params)

	require.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, "key-only", keys[0].ID)
	assert.Equal(t, 1, requestCount, "expected exactly 1 request")
}

func TestSearchAPIAccessKeysWithContext_NameFilter(t *testing.T) {
	t.Parallel()

	response := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{"id": "installer-key", "name": "Installer Ingest License Key", "type": "INGEST"}
						]
					}
				}
			}
		}
	}`

	var capturedBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		capturedBody = string(body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	defer ts.Close()

	client := newAPIAccessClient(t, ts)

	params := APIAccessKeySearchQuery{
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.INGEST},
		Scope: APIAccessKeySearchScope{
			Name: &APIAccessKeyNameQuery{
				Contains: "Installer",
			},
		},
	}

	keys, err := client.SearchAPIAccessKeys(params)

	require.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, "installer-key", keys[0].ID)

	assert.True(t, strings.Contains(capturedBody, `"contains":"Installer"`),
		"expected request body to contain name filter, got: %s", capturedBody)
}

func TestSearchAPIAccessKeysWithContext_ErrorOnPage(t *testing.T) {
	t.Parallel()

	page1 := `{
		"data": {
			"actor": {
				"apiAccess": {
					"keySearch": {
						"keys": [
							{"id": "key-1", "name": "Key One", "type": "INGEST"}
						],
						"nextCursor": "page-2"
					}
				}
			}
		}
	}`

	page2Error := `{
		"errors": [
			{"message": "internal server error on page 2"}
		]
	}`

	requestCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		if requestCount == 1 {
			_, _ = w.Write([]byte(page1))
		} else {
			_, _ = w.Write([]byte(page2Error))
		}
	}))
	defer ts.Close()

	client := newAPIAccessClient(t, ts)

	params := APIAccessKeySearchQuery{
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.INGEST},
	}

	keys, err := client.SearchAPIAccessKeys(params)

	require.Error(t, err, "expected an error when a page returns errors")
	assert.Nil(t, keys, "expected no partial results on error")
	assert.Equal(t, 2, requestCount, "expected 2 requests before encountering the error")
}
