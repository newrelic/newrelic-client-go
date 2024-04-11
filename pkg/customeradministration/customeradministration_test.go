package customeradministration

import (
	"os"
	"strconv"
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

	unitTestMockAccountOneId         = "9999999"
	unitTestMockAccountOneIdAsInt, _ = strconv.Atoi(unitTestMockAccountOneId)
	unitTestMockAccountOneName       = "customerAdministration getAccounts Unit Test Mock Account 1"
	unitTestMockAccountTwoId         = "8888888"
	unitTestMockAccountTwoIdAsInt, _ = strconv.Atoi(unitTestMockAccountTwoId)
	unitTestMockAccountTwoName       = "customerAdministration getAccounts Unit Test Mock Account 2"
	unitTestMockNextCursor           = "=ExXzxTX1XDN0XXX2EXXyXXX"

	unitTestMockOrganizationId = "58a5a9b8-158c-4189-85ea-e08281c58c98"
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Customeradministration {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}
