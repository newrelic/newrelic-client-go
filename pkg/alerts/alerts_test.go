//go:build unit || integration
// +build unit integration

package alerts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// TODO: This is used by incidents_test.go still, need to refactor
func newTestClient(t *testing.T, handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Alerts {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func newIntegrationTestClient(t *testing.T) Alerts {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
