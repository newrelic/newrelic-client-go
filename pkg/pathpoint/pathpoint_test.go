//go:build unit || integration
// +build unit integration

package pathpoint

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Pathpoint {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

func newIntegrationTestClient(t *testing.T) Pathpoint {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

const (
	testFlowGUID EntityGUID = "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA"
	testFlowName            = "Supply Chain Flow"
)

var testPathpointAccountID = mock.IntegrationTestAccountID
