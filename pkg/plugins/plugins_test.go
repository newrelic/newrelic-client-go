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

var (
	listPluginsResponseJSON = `{ "plugins": [] }` // todo: add more realistic response data
)

func TestListPlugins(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, listPluginsResponseJSON, http.StatusOK)

	expected := []*Plugin{}

	actual, err := client.ListPlugins(nil)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestListPluginsWithParams(t *testing.T) {
	t.Parallel()

	guidFilter := "test.newrelic_redis_plugin"
	idsFilter := "123"

	client := newTestPluginsClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[guid]")
		require.Equal(t, guidFilter, name)

		ids := values.Get("filter[ids]")
		require.Equal(t, idsFilter, ids)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(listPluginsResponseJSON))

		require.NoError(t, err)
	}))

	params := ListPluginsParams{
		GUID: guidFilter,
		IDs:  []int{123},
	}

	expected := []*Plugin{}

	actual, err := client.ListPlugins(&params)

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
