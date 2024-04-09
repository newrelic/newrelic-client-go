package authorizationmanagement

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegration_GrantAccess_InvalidGroupError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	_, err = client.AuthorizationManagementGrantAccess(
		AuthorizationManagementGrantAccess{
			GroupId: mockGroupId,
			AccountAccessGrants: []AuthorizationManagementAccountAccessGrant{
				{
					AccountID: mockAccountId,
					RoleId:    mockRoleId,
				},
			},
		},
	)

	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile(invalidGroupErrorRegularExpression), err.Error())
}

func TestIntegration_GrantAccess_AccountRole_InvalidAccountIDAndRoleID(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	userManagementClient := newIntegrationTestClientUserManagement(t)

	createdGroupID, err := UserManagementCreateDeleteTestGroupUtility(
		userManagementClient,
		false,
		"",
	)

	require.NotNil(t, createdGroupID)
	require.NoError(t, err)

	_, err = client.AuthorizationManagementGrantAccess(
		AuthorizationManagementGrantAccess{
			GroupId: createdGroupID,
			AccountAccessGrants: []AuthorizationManagementAccountAccessGrant{
				{
					AccountID: mockAccountId,
					RoleId:    mockRoleId,
				},
			},
		},
	)

	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile(invalidAccountIdAndRoleIdErrorRegularExpression), err.Error())

	deletedGroupID, err := UserManagementCreateDeleteTestGroupUtility(
		userManagementClient,
		true,
		createdGroupID,
	)

	require.NoError(t, err)
	require.NotNil(t, deletedGroupID)
}

func TestIntegration_GrantAccess_AccountRole_InvalidAccountID(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	userManagementClient := newIntegrationTestClientUserManagement(t)

	createdGroupID, err := UserManagementCreateDeleteTestGroupUtility(
		userManagementClient,
		false,
		"",
	)

	require.NotNil(t, createdGroupID)
	require.NoError(t, err)

	_, err = client.AuthorizationManagementGrantAccess(
		AuthorizationManagementGrantAccess{
			GroupId: createdGroupID,
			AccountAccessGrants: []AuthorizationManagementAccountAccessGrant{
				{
					AccountID: mockAccountId,
					RoleId:    roleId,
				},
			},
		},
	)

	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile(invalidAccountIdErrorRegularExpression), err.Error())

	deletedGroupID, err := UserManagementCreateDeleteTestGroupUtility(
		userManagementClient,
		true,
		createdGroupID,
	)

	require.NoError(t, err)
	require.NotNil(t, deletedGroupID)
}

func TestIntegration_GrantAccess_AccountRole_Success(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	userManagementClient := newIntegrationTestClientUserManagement(t)

	createdGroupID, err := UserManagementCreateDeleteTestGroupUtility(
		userManagementClient,
		false,
		"",
	)

	require.NotNil(t, createdGroupID)
	require.NoError(t, err)

	grantAccessResponse, err := client.AuthorizationManagementGrantAccess(
		AuthorizationManagementGrantAccess{
			GroupId: createdGroupID,
			AccountAccessGrants: []AuthorizationManagementAccountAccessGrant{
				{
					AccountID: mock.IntegrationTestAccountID,
					RoleId:    roleId,
				},
			},
		},
	)

	require.NoError(t, err)
	foundRole := false
	for _, v := range grantAccessResponse.Roles {
		if strconv.Itoa(v.RoleId) == roleId && v.AccountID == mock.IntegrationTestAccountID {
			foundRole = true
		}
	}

	require.Equal(t, foundRole, true)

	revokeAccessResponse, err := client.AuthorizationManagementRevokeAccess(
		AuthorizationManagementRevokeAccess{
			GroupId: createdGroupID,
			AccountAccessGrants: []AuthorizationManagementAccountAccessGrant{
				{
					AccountID: mock.IntegrationTestAccountID,
					RoleId:    roleId,
				},
			},
		},
	)

	require.NoError(t, err)

	// the following condition is actually supposed to be the block of code commented below
	// however, the API is currently returning a list of existent roles in the response of this
	// mutation, which is why the response is empty after the only role linked to the group
	// is removed. Hence, the condition to have an empty list of roles has been added below until this is fixed in the API.
	require.Equal(t, revokeAccessResponse.Roles, []AuthorizationManagementGrantedRole{})

	//foundRevokedRole := false
	//for _, v := range revokeAccessResponse.Roles {
	//	if strconv.Itoa(v.RoleId) == roleId && v.AccountID == mock.IntegrationTestAccountID {
	//		foundRevokedRole = true
	//	}
	//}
	//
	//require.Equal(t, foundRevokedRole, true)

	deletedGroupID, err := UserManagementCreateDeleteTestGroupUtility(
		userManagementClient,
		true,
		createdGroupID,
	)

	require.NoError(t, err)
	require.NotNil(t, deletedGroupID)
}

func TestIntegration_RevokeAccess_InvalidGroupError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	_, err = client.AuthorizationManagementRevokeAccess(
		AuthorizationManagementRevokeAccess{
			GroupId: mockGroupId,
			AccountAccessGrants: []AuthorizationManagementAccountAccessGrant{
				{
					AccountID: mockAccountId,
					RoleId:    mockRoleId,
				},
			},
		},
	)

	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile(invalidGroupErrorRegularExpression), err.Error())
}
