package alerts

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TODO: This is used by incidents_test.go still, need to refactor
// nolint
func newTestClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		AdminAPIKey:           "abc123",
		BaseURL:               ts.URL,
		InfrastructureBaseURL: ts.URL,
		UserAgent:             "newrelic/newrelic-client-go",
	})

	return c
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

	client := New(config.Config{
		AdminAPIKey:    adminAPIKey,
		PersonalAPIKey: personalAPIKey,
		LogLevel:       "debug",
	})

	return client
}

func TestSetBaseURL(t *testing.T) {
	a := New(config.Config{
		BaseURL: "http://localhost",
	})

	assert.Equal(t, "http://localhost", a.client.Config.BaseURL)
}

func TestSetInfrastructureBaseURL(t *testing.T) {
	a := New(config.Config{
		InfrastructureBaseURL: "http://localhost",
	})

	assert.Equal(t, "http://localhost", a.infraClient.Config.BaseURL)
}
