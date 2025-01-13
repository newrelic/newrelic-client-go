//go:build integration
// +build integration

package usermanagement

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCreateGroup(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createGroupInput := UserManagementCreateGroup{
		AuthenticationDomainId: authenticationDomainId,
		DisplayName:            groupName,
	}

	createGroupResponse, err := client.UserManagementCreateGroup(createGroupInput)

	require.NoError(t, err)
	require.NotNil(t, createGroupResponse.Group.ID)

	deleteGroupInput := UserManagementDeleteGroup{ID: createGroupResponse.Group.ID}
	deleteGroupResponse, err := client.UserManagementDeleteGroup(deleteGroupInput)

	require.NoError(t, err)
	require.NotNil(t, deleteGroupResponse)
}

func TestIntegrationCreateGroupError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createGroupInput := UserManagementCreateGroup{
		AuthenticationDomainId: mockAuthenticationDomainId,
		DisplayName:            groupName,
	}

	_, err = client.UserManagementCreateGroup(createGroupInput)

	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationUpdateGroupError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	updateGroupInput := UserManagementUpdateGroup{
		DisplayName: fmt.Sprintf("%s-updated", groupName),
		ID:          mockUserId,
	}

	_, err = client.UserManagementUpdateGroup(updateGroupInput)

	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationDeleteGroupError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	deleteGroupInput := UserManagementDeleteGroup{
		ID: mockUserId,
	}

	_, err = client.UserManagementDeleteGroup(deleteGroupInput)

	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationGroupManagementWithoutUsers(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	fmt.Println("client", client)
	err = UserManagementGroupCleanupForIntegrationTests(client, authenticationDomainId)
	fmt.Println("err", err)
	require.NoError(t, err)

	displayName := fmt.Sprintf("%s-new", groupName)
	createGroupInput := UserManagementCreateGroup{
		AuthenticationDomainId: authenticationDomainId,
		DisplayName:            displayName,
	}

	createGroupResponse, err := client.UserManagementCreateGroup(createGroupInput)
	fmt.Println("createGroupResponse", createGroupResponse, err)

	require.NoError(t, err)
	require.NotNil(t, createGroupResponse.Group.ID)

	displayNameUpdated := fmt.Sprintf("%s-updated", displayName)

	updateGroupInput := UserManagementUpdateGroup{
		DisplayName: displayNameUpdated,
		ID:          createGroupResponse.Group.ID,
	}

	updateGroupResponse, err := client.UserManagementUpdateGroup(updateGroupInput)
	fmt.Println("updateGroupResponse", updateGroupResponse, err)

	require.NoError(t, err)
	require.NotNil(t, updateGroupResponse.Group.ID)
	require.Equal(t, updateGroupResponse.Group.DisplayName, displayNameUpdated)

	deleteGroupInput := UserManagementDeleteGroup{ID: createGroupResponse.Group.ID}
	deleteGroupResponse, err := client.UserManagementDeleteGroup(deleteGroupInput)

	require.NoError(t, err)
	require.NotNil(t, deleteGroupResponse)

}

func UserManagementGroupCleanupForIntegrationTests(client Usermanagement, authenticationDomainId string) error {
	getGroupsResponse, err := client.UserManagementGetGroupsWithUsers(
		[]string{authenticationDomainId},
		[]string{},
		"",
	)

	if err != nil {
		return err
	}

	for _, a := range getGroupsResponse.AuthenticationDomains {
		if a.ID == authenticationDomainId {
			for _, g := range a.Groups.Groups {
				if strings.Contains(g.DisplayName, groupNamePrefix) {
					_, err := client.UserManagementDeleteGroup(UserManagementDeleteGroup{ID: g.ID})
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func TestIntegrationAddUsersToGroupsAndRemove(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// fetch a user and a group to link the user to
	groupSearchTerm := "Integration Test Group 1 DO NOT DELETE"
	userSearchTerm := "Integration Test User 1 DO NOT DELETE"

	getGroupsResponse, err := client.UserManagementGetGroupsWithUsers(
		[]string{authenticationDomainId},
		[]string{},
		groupSearchTerm,
	)
	require.NoError(t, err)

	getUsersResponse, err := client.UserManagementGetUsers(
		[]string{authenticationDomainId},
		[]string{},
		userSearchTerm,
		"",
	)
	require.NoError(t, err)

	groupID := ""
	userID := ""

	for _, a := range getGroupsResponse.AuthenticationDomains {
		if a.ID == authenticationDomainId {
			for _, g := range a.Groups.Groups {
				if strings.Contains(g.DisplayName, groupSearchTerm) {
					groupID = g.ID
				}
			}
		}
	}

	for _, a := range getUsersResponse.AuthenticationDomains {
		if a.ID == authenticationDomainId {
			for _, u := range a.Users.Users {
				if strings.Contains(u.Name, userSearchTerm) {
					userID = u.ID
				}
			}
		}
	}

	require.NotEmpty(t, groupID)
	require.NotEmpty(t, userID)

	addUsersToGroupInput := UserManagementUsersGroupsInput{
		GroupIds: []string{groupID},
		UserIDs:  []string{userID},
	}

	addUsersToGroupResponse, err := client.UserManagementAddUsersToGroups(addUsersToGroupInput)
	require.NoError(t, err)
	require.NotNil(t, addUsersToGroupResponse)

	removeUsersFromGroupResponse, err := client.UserManagementRemoveUsersFromGroups(addUsersToGroupInput)
	require.NoError(t, err)
	require.NotNil(t, removeUsersFromGroupResponse)
}

func TestIntegrationAddUsersToGroupsError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	addUsersToGroupInput := UserManagementUsersGroupsInput{
		GroupIds: []string{mockGroupId},
		UserIDs:  []string{mockUserId},
	}

	_, err = client.UserManagementAddUsersToGroups(addUsersToGroupInput)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationRemoveUsersFromGroupsError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	removeUsersFromGroupInput := UserManagementUsersGroupsInput{
		GroupIds: []string{mockGroupId},
		UserIDs:  []string{mockUserId},
	}

	_, err = client.UserManagementRemoveUsersFromGroups(removeUsersFromGroupInput)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Could not find the target or you are unauthorized.")
}

func TestIntegrationGroupManagementWithUsers(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	err = UserManagementGroupCleanupForIntegrationTests(client, authenticationDomainId)
	require.NoError(t, err)

	displayName := fmt.Sprintf("%s-new", groupName)
	createGroupInput := UserManagementCreateGroup{
		AuthenticationDomainId: authenticationDomainId,
		DisplayName:            displayName,
	}

	createGroupResponse, err := client.UserManagementCreateGroup(createGroupInput)

	require.NoError(t, err)
	require.NotNil(t, createGroupResponse.Group.ID)

	getUserResponse, err := client.UserManagementGetUsers([]string{authenticationDomainId}, []string{}, "Integration Test User 1 DO NOT DELETE", "")
	require.NoError(t, err)
	require.NotNil(t, getUserResponse)

	userID := ""

	for _, authDomain := range getUserResponse.AuthenticationDomains {
		if authDomain.ID == authenticationDomainId {
			for _, u := range authDomain.Users.Users {
				if strings.Contains(u.Name, "Integration Test User 1 DO NOT DELETE") {
					userID = u.ID
				}
			}
		}
	}

	require.NotEmpty(t, userID)

	addUsersToGroupInput := UserManagementUsersGroupsInput{
		GroupIds: []string{createGroupResponse.Group.ID},
		UserIDs:  []string{userID},
	}

	addUsersToGroupResponse, err := client.UserManagementAddUsersToGroups(addUsersToGroupInput)
	require.NoError(t, err)
	require.NotNil(t, addUsersToGroupResponse)

	displayNameUpdated := fmt.Sprintf("%s-updated", displayName)

	updateGroupInput := UserManagementUpdateGroup{
		DisplayName: displayNameUpdated,
		ID:          createGroupResponse.Group.ID,
	}

	updateGroupResponse, err := client.UserManagementUpdateGroup(updateGroupInput)

	require.NoError(t, err)
	require.NotNil(t, updateGroupResponse.Group.ID)
	require.Equal(t, updateGroupResponse.Group.DisplayName, displayNameUpdated)

	removeUsersFromGroupResponse, err := client.UserManagementRemoveUsersFromGroups(addUsersToGroupInput)
	require.NoError(t, err)
	require.NotNil(t, removeUsersFromGroupResponse)

	deleteGroupInput := UserManagementDeleteGroup{ID: createGroupResponse.Group.ID}
	deleteGroupResponse, err := client.UserManagementDeleteGroup(deleteGroupInput)

	require.NoError(t, err)
	require.NotNil(t, deleteGroupResponse)

}
