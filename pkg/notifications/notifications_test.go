//go:build unit || integration
// +build unit integration

package notifications

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Notifications {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func newIntegrationTestClient(t *testing.T) Notifications {
	cfg := mock.NewIntegrationTestConfig(t)
	client := New(cfg)

	return client
}
