//go:build unit || integration
// +build unit integration

package usermanagement

import (
	"fmt"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Usermanagement {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Usermanagement {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

var (
	authenticationDomainId = "84cb286a-8eb0-4478-b469-cdf2ccfef553"

	userNamePrefix  = "newrelic-client-go-integration-test-mock-user"
	userName        = fmt.Sprintf("%s-%s", userNamePrefix, mock.RandSeq(5))
	groupNamePrefix = "newrelic-client-go-integration-test-mock-group"
	groupName       = fmt.Sprintf("%s-%s", groupNamePrefix, mock.RandSeq(5))

	userEmail        = fmt.Sprintf("developer-toolkit+%s@newrelic.com", mock.RandSeq(5))
	userEmailUpdated = fmt.Sprintf("developer-toolkit+%s@newrelic.com", mock.RandSeq(10))
	userNameUpdated  = fmt.Sprintf("%s-updated", userName)

	mockUserId                      = "9999999999"
	mockAuthenticationDomainId      = "fae55e6b-b1ce-4a0f-83b2-ee774798f2cc"
	mockGroupId                     = "e5f30c8d-acf6-481c-a492-ef516957b479"
	mockAuthenticationDomainName    = "Mock Authentication Domain 1"
	mockAuthenticationDomainIdTwo   = "b5aff74c-f56c-420a-980e-4d97e7278861"
	mockAuthenticationDomainNameTwo = "Mock Authentication Domain Two"
	mockUserEmail                   = "mock@mock.mock"
	mockUserEmailUpdated            = fmt.Sprintf("updated_%s", mockUserEmail)
)
