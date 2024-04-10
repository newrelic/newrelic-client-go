package customeradministration

import (
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Customeradministration {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

var (
	organizationId = os.Getenv("INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID")

	authenticationDomainName = "Test-Auth-Domain DO NOT DELETE"
	authenticationDomainId   = os.Getenv("INTEGRATION_TESTING_NEW_RELIC_AUTHENTICATION_DOMAIN_ID")

	roleName = "Integration Test Role 1 DO NOT DELETE"
	roleId   = "38236"
)

// WORK IN PROGRESS
//func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Customeradministration {
//	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
//	tc := mock.NewTestConfig(t, ts)
//
//	return New(tc)
//}
