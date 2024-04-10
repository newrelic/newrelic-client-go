package authorizationmanagement

import (
	"fmt"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
)

func newIntegrationTestClient(t *testing.T) Authorizationmanagement {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

func newIntegrationTestClientUserManagement(t *testing.T) usermanagement.Usermanagement {
	tc := mock.NewIntegrationTestConfig(t)
	return usermanagement.New(tc)
}

var (
	authenticationDomainId = "84cb286a-8eb0-4478-b469-cdf2ccfef553"

	// do not add more than nine 9s here
	mockAccountId = 999999999
	mockGroupId   = "fake-group-id"
	mockRoleId    = "fake-role-id"

	roleId = "38236"

	invalidGroupErrorRegularExpression              = `An error occurred resolving this field`
	invalidAccountIdAndRoleIdErrorRegularExpression = `error granting access`
	invalidAccountIdErrorRegularExpression          = `Validation failed: granted_on entity .* do not belong to the same organization`
)

func UserManagementCreateDeleteTestGroupUtility(
	userManagementClient usermanagement.Usermanagement,
	delete bool,
	groupId string,
) (string, error) {
	if !delete {
		groupCreateResponse, err := userManagementClient.UserManagementCreateGroup(
			usermanagement.UserManagementCreateGroup{
				AuthenticationDomainId: authenticationDomainId,
				DisplayName:            fmt.Sprintf("mock-group-created-for-account-role-grant-%s", mock.RandSeq(5)),
			},
		)

		if err != nil {
			return "", err
		}

		return groupCreateResponse.Group.ID, nil

	} else {
		if groupId == "" {
			return "", fmt.Errorf("iud of the group to be deleted has not been specified")
		}
		groupDeleteResponse, err := userManagementClient.UserManagementDeleteGroup(
			usermanagement.UserManagementDeleteGroup{
				ID: groupId,
			},
		)

		if err != nil {
			return "", err
		}

		return groupDeleteResponse.Group.ID, nil
	}
}

// WORK IN PROGRESS
//func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Authorizationmanagement {
//	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
//	tc := mock.NewTestConfig(t, ts)
//
//	return New(tc)
//}
