//go:build unit || integration
// +build unit integration

package federatedlogs

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newMockClient(t *testing.T, mockJSONResponse string, statusCode int) FederatedLogs {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

func newIntegrationTestClient(t *testing.T) FederatedLogs {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
