package usermanagement

import (
	"fmt"
	"strings"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

var (
	authenticationDomainId     = "84cb286a-8eb0-4478-b469-cdf2ccfef553"
	mockAuthenticationDomainId = "fae55e6b-b1ce-4a0f-83b2-ee774798f2cc"
	userName                   = fmt.Sprintf("newrelic-client-go-integration-test-mock-user-%s", mock.RandSeq(5))
	userEmail                  = fmt.Sprintf("developer-toolkit+%s@newrelic.com", mock.RandSeq(5))
	userNameUpdated            = fmt.Sprintf("%s_updated", userName)
	mockUserId                 = "9999999999"
)

func TestIntegrationCreateUser(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: authenticationDomainId,
		Name:                   userName,
		Email:                  userEmail,
		UserType:               UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	}

	createUserResponse, err := client.UserManagementCreateUser(createUserInput)

	require.NoError(t, err)
	require.NotNil(t, createUserResponse.CreatedUser.ID)
}

func TestIntegrationCreateUserError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: mockAuthenticationDomainId,
		Name:                   userName,
		Email:                  userEmail,
		UserType:               UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	}

	_, err = client.UserManagementCreateUser(createUserInput)

	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationUpdateUserError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	updateUserInput := UserManagementUpdateUser{
		Name:     userNameUpdated,
		Email:    userEmail,
		UserType: UserManagementRequestedTierNameTypes.CORE_USER_TIER,
		ID:       mockUserId,
	}

	_, err = client.UserManagementUpdateUser(updateUserInput)

	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationDeleteUserError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	deleteUserInput := UserManagementDeleteUser{
		ID: mockUserId,
	}

	_, err = client.UserManagementDeleteUser(deleteUserInput)

	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationUserManagement(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: authenticationDomainId,
		Name:                   userName,
		Email:                  userEmail,
		UserType:               UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	}

	createUserResponse, err := client.UserManagementCreateUser(createUserInput)

	require.NoError(t, err)
	require.NotNil(t, createUserResponse.CreatedUser.ID)

	getUserResponse, err := client.GetUsers([]string{authenticationDomainId}, []string{}, "", "")
	require.NoError(t, err)
	require.NotNil(t, getUserResponse)

	updateUserInput := UserManagementUpdateUser{
		Name:     userNameUpdated,
		Email:    userEmail,
		UserType: UserManagementRequestedTierNameTypes.CORE_USER_TIER,
		ID:       createUserResponse.CreatedUser.ID,
	}

	updateUserResponse, err := client.UserManagementUpdateUser(updateUserInput)

	require.NoError(t, err)
	require.NotNil(t, updateUserResponse.User.ID)
	require.Equal(t, updateUserResponse.User.ID, createUserResponse.CreatedUser.ID)
	require.Equal(t, updateUserResponse.User.Name, userNameUpdated)
	require.Equal(t, updateUserResponse.User.Type.DisplayName, "Core")

	deleteUserInput := UserManagementDeleteUser{ID: updateUserResponse.User.ID}

	deleteUserResponse, err := client.UserManagementDeleteUser(deleteUserInput)

	require.NoError(t, err)
	require.NotNil(t, deleteUserResponse.DeletedUser.ID)

}

func TestIntegrationGetUsersError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	getUserResponse, err := client.GetUsers([]string{mockAuthenticationDomainId}, []string{}, "", "")
	require.NoError(t, err)
	require.Zero(t, len(getUserResponse.AuthenticationDomains))
}

func TestIntegrationGetUsersTestAndCleanup(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	getUsersResponse, err := client.GetUsers([]string{authenticationDomainId}, []string{}, "newrelic-client-go", "")
	require.NoError(t, err)

	for _, authDomain := range getUsersResponse.AuthenticationDomains {
		if authDomain.ID == authenticationDomainId {
			for _, u := range authDomain.Users.Users {
				if strings.Contains(u.Name, "newrelic-client-go") {
					_, err := client.UserManagementDeleteUser(UserManagementDeleteUser{ID: u.ID})
					require.NoError(t, err)
				}
			}
		}
	}
}

func newIntegrationTestClient(t *testing.T) Usermanagement {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
