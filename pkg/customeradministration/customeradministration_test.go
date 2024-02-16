package customeradministration

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Customeradministration {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

// WORK IN PROGRESS
//func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Customeradministration {
//	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
//	tc := mock.NewTestConfig(t, ts)
//
//	return New(tc)
//}
