package alerts

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// TODO: This is used by incidents_test.go still, need to refactor
// nolint
func newTestClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	return New(config.Config{
		AdminAPIKey:           "abc123",
		BaseURL:               ts.URL,
		InfrastructureBaseURL: ts.URL,
		UserAgent:             "newrelic/newrelic-client-go",
	})
}

// nolint
func newMockResponse(
	t *testing.T,
	mockJSONResponse string,
	statusCode int,
) Alerts {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)

	return New(config.Config{
		AdminAPIKey:           "abc123",
		BaseURL:               ts.URL,
		InfrastructureBaseURL: ts.URL,
		UserAgent:             "newrelic/newrelic-client-go",
	})
}

// nolint
func newIntegrationTestClient(t *testing.T) Alerts {
	personalAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	adminAPIKey := os.Getenv("NEW_RELIC_ADMIN_API_KEY")

	if personalAPIKey == "" && adminAPIKey == "" {
		t.Skipf("acceptance testing requires NEW_RELIC_API_KEY and NEW_RELIC_ADMIN_API_KEY")
	}

	return New(config.Config{
		AdminAPIKey:    adminAPIKey,
		PersonalAPIKey: personalAPIKey,
		LogLevel:       "debug",
	})
}
