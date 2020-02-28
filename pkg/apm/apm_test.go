package apm

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// nolint
func newTestClient(handler http.Handler) APM {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		PersonalAPIKey: "abc123",
		BaseURL:        ts.URL,
		UserAgent:      "newrelic/newrelic-client-go",
		LogLevel:       "debug",
	})

	return c
}

// nolint
func newMockResponse(
	t *testing.T,
	mockJSONResponse string,
	statusCode int,
) APM {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)

	return New(config.Config{
		PersonalAPIKey: "abc123",
		BaseURL:        ts.URL,
		UserAgent:      "newrelic/newrelic-client-go",
	})
}

// nolint
func newIntegrationTestClient(t *testing.T) APM {
	personalAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	adminAPIKey := os.Getenv("NEW_RELIC_ADMIN_API_KEY")

	if personalAPIKey == "" && adminAPIKey == "" {
		t.Skipf("acceptance testing requires NEW_RELIC_API_KEY and NEW_RELIC_ADMIN_API_KEY")
	}

	client := New(config.Config{
		AdminAPIKey:    adminAPIKey,
		PersonalAPIKey: personalAPIKey,
		LogLevel:       "debug",
	})

	return client
}
