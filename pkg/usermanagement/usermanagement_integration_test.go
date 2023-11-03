//go:build integration
// +build integration

package usermanagement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateUser(t *testing.T) {
	t.Parallel()

	// Test: Create
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		Name:                   "unit test tf",
		Email:                  "unittest@tf.com",
		UserType:               UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	}
	client := newIntegrationTestClient(t)
	created, err := client.UserManagementCreateUser(createUserInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	//defer deleteTestUser(t, created.CreatedUser.ID)
	var CreatedUser = created.CreatedUser

	require.Equal(t, CreatedUser.Name, createUserInput.Name)
	require.Equal(t, CreatedUser.Email, createUserInput.Email)
	require.Equal(t, CreatedUser.Type.DisplayName, "Core")
	require.NotNil(t, CreatedUser.ID)
}

func createTestUser(t *testing.T) *UserManagementCreatedUser {
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		Name:                   "unit test tf",
		Email:                  "unittest@tf.com",
		UserType:               UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	}
	client := newIntegrationTestClient(t)
	user, err := client.UserManagementCreateUser(createUserInput)
	require.NoError(t, err)

	return &user.CreatedUser
}

func deleteTestUser(t *testing.T, userId string) {
	deleteUserInput := UserManagementDeleteUser{
		ID: userId,
	}
	client := newIntegrationTestClient(t)
	deleted, err := client.UserManagementDeleteUser(deleteUserInput)
	require.NoError(t, err)
	require.NotNil(t, deleted)
}

func TestIntegrationUpdateUser(t *testing.T) {
	t.Parallel()

	createdUser := createTestUser(t)
	updateUserInput := UserManagementUpdateUser{
		Name:     "test unit update",
		Email:    "unittestupdated@tf.com",
		UserType: UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
		ID:       createdUser.ID,
	}
	client := newIntegrationTestClient(t)
	updated, err := client.UserManagementUpdateUser(updateUserInput)

	require.NoError(t, err)
	require.NotNil(t, updated)
	//defer deleteTestUser(t, updated.User.ID)
	var updatedUser = updated.User
	require.Equal(t, updatedUser.Name, updateUserInput.Name)
	require.Equal(t, updatedUser.Email, updateUserInput.Email)
	require.Equal(t, updatedUser.Type.DisplayName, "Basic")
	require.Equal(t, createdUser.ID, updatedUser.ID)
}

func TestIntegrationDeleteUser(t *testing.T) {
	t.Parallel()

	createdUser := createTestUser(t)
	deleteUserInput := UserManagementDeleteUser{
		ID: createdUser.ID,
	}
	client := newIntegrationTestClient(t)
	deleted, err := client.UserManagementDeleteUser(deleteUserInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.Equal(t, deleted.DeletedUser.ID, createdUser.ID)
}

func TestIntegrationDeleteInvalidUser(t *testing.T) {
	t.Parallel()

	deleteUserInput := UserManagementDeleteUser{
		ID: "testinvaliduserdelete",
	}
	client := newIntegrationTestClient(t)
	deleted, err := client.UserManagementDeleteUser(deleteUserInput)

	require.Error(t, err)
	require.Nil(t, deleted)
}

func TestIntegrationCreateUserInvalidAuthDomain(t *testing.T) {
	t.Parallel()

	// Test: Create
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: "0cc21d98-8dc2-484a-bb26",
		Name:                   "unit test tf",
		Email:                  "unittest@tf.com",
		UserType:               UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	}
	client := newIntegrationTestClient(t)
	created, err := client.UserManagementCreateUser(createUserInput)
	require.Error(t, err)
	require.Nil(t, created)

}

func TestIntegrationCreateUserInvalidEmail(t *testing.T) {
	t.Parallel()

	// Test: Create
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		Name:                   "unit test tf",
		Email:                  "unittest@tf.com@invalid",
		UserType:               UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	}
	client := newIntegrationTestClient(t)
	created, err := client.UserManagementCreateUser(createUserInput)
	require.Error(t, err)
	require.Nil(t, created)

}

func TestIntegrationCreateUserExisting(t *testing.T) {
	t.Parallel()

	createdUser := createTestUser(t)
	// Test: Create
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		Name:                   "unit test tf",
		Email:                  createdUser.Email,
		UserType:               UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	}
	client := newIntegrationTestClient(t)
	created, err := client.UserManagementCreateUser(createUserInput)
	defer deleteTestUser(t, createdUser.ID)
	require.Error(t, err)
	require.Nil(t, created)

}

func TestIntegrationCreateUserWithDefaultAuthDomain(t *testing.T) {
	t.Parallel()

	// Test: Create
	client := newIntegrationTestClient(t)
	authDomains, err := client.GetAuthenticationDomains("", nil)
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: authDomains.AuthenticationDomains[0].ID,
		Name:                   "unit test tf",
		Email:                  "unittest@tf.com",
		UserType:               UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	}

	created, err := client.UserManagementCreateUser(createUserInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	//defer deleteTestUser(t, created.CreatedUser.ID)
	var CreatedUser = created.CreatedUser

	require.Equal(t, CreatedUser.Name, createUserInput.Name)
	require.Equal(t, CreatedUser.Email, createUserInput.Email)
	require.Equal(t, CreatedUser.Type.DisplayName, "Core")
	require.NotNil(t, CreatedUser.ID)
}
