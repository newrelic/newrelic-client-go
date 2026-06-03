//go:build unit || integration
// +build unit integration

package federatedlogs

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Federatedlogs {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func newIntegrationTestClient(t *testing.T) Federatedlogs {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}

// newFleetIntegrationTestClient builds a Federatedlogs client authenticated
// with the fleet-test credentials.
func newFleetIntegrationTestClient(t *testing.T) Federatedlogs {
	tc := mock.NewFleetIntegrationTestConfig(t)

	return New(tc)
}
