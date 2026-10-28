//go:build unit
// +build unit

package fleetcontrol

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

const testEntitySearchMockResponse = `{
	"data": {
		"actor": {
			"entityManagement": {
				"entitySearch": {
					"entities": [],
					"nextCursor": "next-page-cursor"
				}
			}
		}
	}
}`

// newMockResponseCapturingRequest is like newMockResponse but also captures the
// raw request body of the last request sent to the mock server, so tests can
// assert on the GraphQL variables that were actually sent over the wire.
func newMockResponseCapturingRequest(t *testing.T, mockJSONResponse string, statusCode int) (Fleetcontrol, *[]byte) {
	var lastRequestBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		lastRequestBody = body

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err = w.Write([]byte(mockJSONResponse))
		require.NoError(t, err)
	}))
	t.Cleanup(ts.Close)

	tc := mock.NewTestConfig(t, ts)

	return New(tc), &lastRequestBody
}

type graphQLRequestBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func TestUnitGetEntitySearch_SendsCursorVariable(t *testing.T) {
	t.Parallel()

	client, lastRequestBody := newMockResponseCapturingRequest(t, testEntitySearchMockResponse, http.StatusOK)

	cursor := "current-page-cursor"
	result, err := client.GetEntitySearch(&cursor, "type = 'FLEET'")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "next-page-cursor", result.NextCursor)

	var reqBody graphQLRequestBody
	require.NoError(t, json.Unmarshal(*lastRequestBody, &reqBody))

	assert.Contains(t, reqBody.Query, "$cursor: String")
	assert.Contains(t, reqBody.Query, "cursor: $cursor")
	assert.Equal(t, cursor, reqBody.Variables["cursor"])
	assert.Equal(t, "type = 'FLEET'", reqBody.Variables["query"])
}

func TestUnitGetEntitySearch_NilCursorSendsNull(t *testing.T) {
	t.Parallel()

	client, lastRequestBody := newMockResponseCapturingRequest(t, testEntitySearchMockResponse, http.StatusOK)

	_, err := client.GetEntitySearch(nil, "type = 'FLEET'")
	require.NoError(t, err)

	var reqBody graphQLRequestBody
	require.NoError(t, json.Unmarshal(*lastRequestBody, &reqBody))

	assert.Nil(t, reqBody.Variables["cursor"])
}
