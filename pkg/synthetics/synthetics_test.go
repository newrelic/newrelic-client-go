// +build unit

package synthetics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
)

// TestError checks that messages are reported in the correct
// order by going through priorities backwards
func TestError(t *testing.T) {
	t.Parallel()

	// Default empty
	e := ErrorResponse{}
	assert.Equal(t, "", e.Error())

	// 3rd Messages concat
	e.Messages = []ErrorDetail{
		{Message: "detail message"},
		{Message: "another detail"},
	}
	assert.Equal(t, "detail message, another detail", e.Error())

	// 2nd Message
	e.Message = "message"
	assert.Equal(t, "message", e.Error())

	// 1st Server Error Message
	e.ServerErrorMessage = "server message"
	assert.Equal(t, "server message", e.Error())

}

func newTestClient(t *testing.T, handler http.Handler) Synthetics {
	ts := httptest.NewServer(handler)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Synthetics {
	return newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJSONResponse))

		require.NoError(t, err)
	}))
}
