package organization

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// TODO: Uncomment when integration tests are added to this package
//func newIntegrationTestClient(t *testing.T) Organization {
//	tc := mock.NewIntegrationTestConfig(t)
//	return New(tc)
//}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Organization {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

var (
	unitTestMockOrganizationCreateJobId = "bec1a268-53b8-4dc5-8522-37e648fc9d38"
	unitTestMockCustomerId              = "CC-0000000000"
	unitTestMockOrganizationOneName     = "Mock Organization One"
	unitTestMockOrganizationOneId       = "e1fe1ff8-0032-43d5-935f-caf47567a71d"
)
