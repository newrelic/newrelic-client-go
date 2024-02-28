//go:build unit
// +build unit

package usermanagement

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCreateGroupResponseJSON = `{
	"data": {
		"userManagementCreateGroup": {
			"group": {
				"displayName": "` + groupName + `",
				"id": "` + mockGroupId + `",
				"users": {
					"nextCursor": null,
					"totalCount": 0
				}
			}
		}
	}
}`

	testAddUsersToGroupsJSON = `{
	"data": {
		"userManagementAddUsersToGroups": {
			"groups": [{
				"displayName": "` + groupName + `",
				"id": "` + mockGroupId + `",
				"users": {
					"nextCursor": null,
					"totalCount": 2
				}
			}]
		}
	}
}`

	testUpdateGroupResponseJSON = `{
	"data": {
		"userManagementUpdateGroup": {
			"group": {
				"displayName": "` + groupName + `",
				"id": "` + mockGroupId + `",
				"users": {
					"nextCursor": null,
					"totalCount": 2
				}
			}
		}
	}
}`

	testRemoveUsersFromGroupsJSON = `{
	"data": {
		"userManagementRemoveUsersFromGroups": {
			"groups": [{
				"displayName": "` + groupName + `",
				"id": "` + mockGroupId + `",
				"users": {
					"nextCursor": null,
					"totalCount": 0
				}
			}]
		}
	}
}`

	testDeleteGroupResponseJson = `{
	"data": {
		"userManagementDeleteGroup": {
			"group": {
				"displayName": "` + groupName + `",
				"id": "` + mockGroupId + `",
				"users": {
					"nextCursor": null,
					"totalCount": 0
				}
			}
		}
	}
}`

	testGetUsersInGroupsResponseJSON = `{
    "data":{
        "actor":{
            "organization":{
                "userManagement":{
                    "authenticationDomains":{
                        "authenticationDomains":[
                            {
                                "groups":{
                                    "groups":[
                                        {
                                            "displayName":"` + groupName + `",
                                            "id":"` + mockGroupId + `",
                                            "users":{
                                                "users":[
                                                    {
                                                        "email":"` + mockUserEmail + `",
                                                        "id":"` + mockUserId + `",
                                                        "name":"` + userName + `",
                                                        "timeZone":"Etc/UTC"
                                                    }
                                                ]
                                            }
                                        }
                                    ]
                                },
                                "id":"` + mockAuthenticationDomainId + `",
                                "name":"` + mockAuthenticationDomainName + `"
                            }
                        ],
                        "nextCursor":null,
                        "totalCount":1
                    }
                }
            }
        }
    }
}`
)

func TestUnitCreateGroup(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testCreateGroupResponseJSON, http.StatusCreated)
	createGroupInput := UserManagementCreateGroup{
		AuthenticationDomainId: mockAuthenticationDomainId,
		DisplayName:            groupName,
	}

	expected := &UserManagementCreateGroupPayload{Group: UserManagementGroup{
		DisplayName: groupName,
		ID:          mockGroupId,
		Users: UserManagementGroupUsers{
			NextCursor: "",
			TotalCount: 0,
			Users:      nil,
		},
	}}

	actual, err := user.UserManagementCreateGroup(createGroupInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitAddUsersToGroups(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testAddUsersToGroupsJSON, http.StatusCreated)

	addUsersToGroupsInput := UserManagementUsersGroupsInput{
		GroupIds: []string{mockGroupId},
		UserIDs:  []string{mockUserId, mockUserIdTwo},
	}

	expected := &UserManagementAddUsersToGroupsPayload{Groups: []UserManagementGroup{
		{
			DisplayName: groupName,
			ID:          mockGroupId,
			Users: UserManagementGroupUsers{
				NextCursor: "",
				TotalCount: 2,
				Users:      nil,
			},
		},
	}}

	actual, err := user.UserManagementAddUsersToGroups(addUsersToGroupsInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitUpdateGroup(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testUpdateGroupResponseJSON, http.StatusCreated)
	updateGroupInput := UserManagementUpdateGroup{
		DisplayName: groupName,
	}

	expected := &UserManagementUpdateGroupPayload{Group: UserManagementGroup{
		DisplayName: groupName,
		ID:          mockGroupId,
		Users: UserManagementGroupUsers{
			NextCursor: "",
			TotalCount: 2,
			Users:      nil,
		},
	}}

	actual, err := user.UserManagementUpdateGroup(updateGroupInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitRemoveUsersFromGroups(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testRemoveUsersFromGroupsJSON, http.StatusCreated)

	addUsersToGroupsInput := UserManagementUsersGroupsInput{
		GroupIds: []string{mockGroupId},
		UserIDs:  []string{mockUserId, mockUserIdTwo},
	}

	expected := &UserManagementRemoveUsersFromGroupsPayload{Groups: []UserManagementGroup{
		{
			DisplayName: groupName,
			ID:          mockGroupId,
			Users: UserManagementGroupUsers{
				NextCursor: "",
				TotalCount: 0,
				Users:      nil,
			},
		},
	}}

	actual, err := user.UserManagementRemoveUsersFromGroups(addUsersToGroupsInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitDeleteGroup(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testDeleteGroupResponseJson, http.StatusCreated)
	deleteGroupInput := UserManagementDeleteGroup{
		ID: mockGroupId,
	}

	expected := &UserManagementDeleteGroupPayload{Group: UserManagementGroup{
		DisplayName: groupName,
		ID:          mockGroupId,
		Users: UserManagementGroupUsers{
			NextCursor: "",
			TotalCount: 0,
			Users:      nil,
		},
	}}

	actual, err := user.UserManagementDeleteGroup(deleteGroupInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitGetUsersInGroups(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testGetUsersInGroupsResponseJSON, http.StatusCreated)

	expected := &UserManagementAuthenticationDomains{
		AuthenticationDomains: []UserManagementAuthenticationDomain{
			{
				ID:   mockAuthenticationDomainId,
				Name: mockAuthenticationDomainName,
				Groups: UserManagementGroups{
					Groups: []UserManagementGroup{
						{
							DisplayName: groupName,
							ID:          mockGroupId,
							Users: UserManagementGroupUsers{
								Users: []UserManagementGroupUser{
									{
										Email:    mockUserEmail,
										Name:     userName,
										ID:       mockUserId,
										TimeZone: "Etc/UTC",
									},
								},
							},
						},
					},
				},
			},
		},
		NextCursor: "",
		TotalCount: 1,
	}

	actual, err := user.UserManagementGetGroupsWithUsers(
		[]string{mockAuthenticationDomainId},
		[]string{mockGroupId},
		"",
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
