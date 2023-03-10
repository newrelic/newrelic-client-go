package accountmanagement

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newAccountManagementTestClient(t *testing.T) Accountmanagement {
	cfg := mock.NewIntegrationTestConfig(t)
	client := New(cfg)

	return client
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Accountmanagement {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}
