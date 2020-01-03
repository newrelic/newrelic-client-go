// +build unit

package plugins

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPlugins(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{ "plugins": [] }`, http.StatusOK)

	expected := []*Plugin{}

	actual, err := client.ListPlugins()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

// nolint
func newTestPluginsClient(handler http.Handler) Plugins {
	ts := httptest.NewServer(handler)

	return New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})
}

// nolint
func newMockResponse(
	t *testing.T,
	mockJSONResponse string,
	statusCode int,
) Plugins {
	return newTestPluginsClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJSONResponse))

		require.NoError(t, err)
	}))
}
