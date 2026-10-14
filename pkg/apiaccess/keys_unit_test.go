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
	testKeySearchFirstPage = `{
	"data": {
		"actor": {
			"apiAccess": {
				"keySearch": {
					"keys": [
						{"id": "key-1", "key": "aaa", "name": "First", "type": "INGEST", "ingestType": "LICENSE"},
						{"id": "key-2", "key": "bbb", "name": "Second", "type": "INGEST", "ingestType": "LICENSE"}
					],
					"nextCursor": "page-2"
				}
			}
		}
	}
}`

	testKeySearchLastPage = `{
	"data": {
		"actor": {
			"apiAccess": {
				"keySearch": {
					"keys": [
						{"id": "key-3", "key": "ccc", "name": "Third", "type": "INGEST", "ingestType": "LICENSE"}
					],
					"nextCursor": null
				}
			}
		}
	}
}`
)

type keySearchRequestBody struct {
	Variables map[string]interface{} `json:"variables"`
}

// newPagedKeySearchClient serves pages keyed by the cursor variable on each
// request and records the cursor sent with every request.
func newPagedKeySearchClient(t *testing.T, pages map[string]string) (APIAccess, *[]interface{}) {
	var cursors []interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req keySearchRequestBody
		require.NoError(t, json.Unmarshal(body, &req))

		cursor := req.Variables["cursor"]
		cursors = append(cursors, cursor)

		key := ""
		if c, ok := cursor.(string); ok {
			key = c
		}

		page, ok := pages[key]
		require.True(t, ok, "unexpected cursor %v", cursor)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(page))
		require.NoError(t, err)
	}))
	t.Cleanup(ts.Close)

	return New(mock.NewTestConfig(t, ts)), &cursors
}

func TestUnitSearchAPIAccessKeys_FollowsNextCursor(t *testing.T) {
	t.Parallel()

	client, cursors := newPagedKeySearchClient(t, map[string]string{
		"":       testKeySearchFirstPage,
		"page-2": testKeySearchLastPage,
	})

	keys, err := client.SearchAPIAccessKeys(APIAccessKeySearchQuery{
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.INGEST},
	})

	require.NoError(t, err)
	require.Len(t, keys, 3)
	assert.Equal(t, "key-1", keys[0].ID)
	assert.Equal(t, "key-3", keys[2].ID)
	assert.Equal(t, []interface{}{nil, "page-2"}, *cursors)
}

func TestUnitSearchAPIAccessKeys_SinglePage(t *testing.T) {
	t.Parallel()

	client, cursors := newPagedKeySearchClient(t, map[string]string{
		"": testKeySearchLastPage,
	})

	keys, err := client.SearchAPIAccessKeys(APIAccessKeySearchQuery{
		Types: []APIAccessKeyType{APIAccessKeyTypeTypes.INGEST},
	})

	require.NoError(t, err)
	require.Len(t, keys, 1)
	assert.Len(t, *cursors, 1)
}
