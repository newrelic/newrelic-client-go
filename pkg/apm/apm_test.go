package apm

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
)

// nolint
func newTestClient(t *testing.T, handler http.Handler) APM {
	ts := httptest.NewServer(handler)
	tc := mock.NewTestConfig(t, ts)

	c := New(tc)

	return c
}

// nolint
func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) APM {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

// nolint
func newIntegrationTestClient(t *testing.T) APM {
	cfg := mock.NewIntegrationTestConfig(t)
	client := New(cfg)

	return client
}
