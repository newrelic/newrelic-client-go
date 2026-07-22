//go:build unit || integration
// +build unit integration

package scorecards

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newMockClient(t *testing.T, mockJSONResponse string, statusCode int) Scorecards {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

func newIntegrationTestClient(t *testing.T) Scorecards {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
