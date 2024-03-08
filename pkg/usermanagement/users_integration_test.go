//go:build integration
// +build integration

package usermanagement

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
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
	err = UserManagementUserCleanupForIntegrationTests(client, authenticationDomainId)
	require.NoError(t, err)

	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: authenticationDomainId,
		Name:                   userName,
		// userEmailUpdated is not a typo
		Email:    userEmailUpdated,
		UserType: UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	}

	createUserResponse, err := client.UserManagementCreateUser(createUserInput)

	require.NoError(t, err)
	require.NotNil(t, createUserResponse.CreatedUser.ID)

	getUserResponse, err := client.UserManagementGetUsers([]string{authenticationDomainId}, []string{}, "", "")
	require.NoError(t, err)
	require.NotNil(t, getUserResponse)

	updateUserInput := UserManagementUpdateUser{
		Name:     userNameUpdated,
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
	getUserResponse, err := client.UserManagementGetUsers([]string{mockAuthenticationDomainId}, []string{}, "", "")
	require.NoError(t, err)
	require.Zero(t, len(getUserResponse.AuthenticationDomains))
}

func UserManagementUserCleanupForIntegrationTests(client Usermanagement, authenticationDomainId string) error {
	getUsersResponse, err := client.UserManagementGetUsers([]string{authenticationDomainId}, []string{}, userNamePrefix, "")
	if err != nil {
		return err
	}

	for _, authDomain := range getUsersResponse.AuthenticationDomains {
		if authDomain.ID == authenticationDomainId {
			for _, u := range authDomain.Users.Users {
				if strings.Contains(u.Name, userNamePrefix) {
					_, err := client.UserManagementDeleteUser(UserManagementDeleteUser{ID: u.ID})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
