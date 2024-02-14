//go:build unit
// +build unit

package usermanagement

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCreateUserResponseJSON = `{
    "data":{
        "userManagementCreateUser":{
            "createdUser":{
                "authenticationDomainId":"` + mockAuthenticationDomainId + `",
                "email":"` + mockUserEmail + `",
				"id": "` + mockUserId + `",
				"name":"` + userName + `",
                "type":{
                    "displayName":"Basic",
					"id":"1"
                }
            }
        }
    }
}`

	testUpdateUserResponseJSON = `{
	"data": {
		"userManagementUpdateUser": {
			"user": {
				"email": "` + mockUserEmailUpdated + `",
				"emailVerificationState": "Pending",
				"groups": {
					"nextCursor": null,
					"totalCount": 0
				},
				"id": "` + mockUserId + `",
				"lastActive": null,
				"name": "` + userNameUpdated + `",
				"pendingUpgradeRequest": null,
				"timeZone": "Etc/UTC",
				"type": {
					"displayName": "Core",
					"id": "2"
				}
			}
		}
	}
}`
	testDeleteUserResponseJson = `{
    "data":{
        "userManagementDeleteUser":{
            "deletedUser":{
                "id":"` + mockUserId + `"
            }
        }
    }
}`

	testGetUserResponseJSON = `{
    "data":{
        "actor":{
            "organization":{
                "userManagement":{
                    "authenticationDomains":{
                        "authenticationDomains":[
                            {
                                "id":"` + mockAuthenticationDomainId + `",
                                "name":"` + mockAuthenticationDomainName + `",
                                "users":{
                                    "users":[
                                        {
                                            "email":"` + mockUserEmail + `",
                                            "emailVerificationState":"Pending",
                                            "groups":{
                                                "groups":[]
                                            },
                                            "id":"` + mockUserId + `",
                                            "lastActive":null,
                                            "name":"` + userName + `",
                                            "pendingUpgradeRequest":null,
                                            "timeZone":"Etc/UTC",
                                            "type":{
                                                "displayName":"Basic",
                                                "id":"1"
                                            }
                                        }
                                    ]
                                }
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

func TestUnitCreateUser(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testCreateUserResponseJSON, http.StatusCreated)
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: mockAuthenticationDomainId,
		Name:                   userName,
		Email:                  mockUserEmail,
		UserType:               UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	}

	expected := &UserManagementCreateUserPayload{
		CreatedUser: UserManagementCreatedUser{
			AuthenticationDomainId: mockAuthenticationDomainId,
			Email:                  mockUserEmail,
			ID:                     mockUserId,
			Name:                   userName,
			Type: UserManagementUserType{
				DisplayName: "Basic",
				ID:          "1",
			},
		},
	}

	actual, err := user.UserManagementCreateUser(createUserInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitUpdateUser(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testUpdateUserResponseJSON, http.StatusCreated)
	updateUserInput := UserManagementUpdateUser{
		Name:     userNameUpdated,
		Email:    mockUserEmailUpdated,
		UserType: UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	}

	expected := &UserManagementUpdateUserPayload{
		User: UserManagementUser{
			Email:                  mockUserEmailUpdated,
			EmailVerificationState: "Pending",
			TimeZone:               "Etc/UTC",
			Groups:                 UserManagementUserGroups{},
			ID:                     mockUserId,
			Name:                   userNameUpdated,
			PendingUpgradeRequest:  UserManagementPendingUpgradeRequest{},
			Type: UserManagementUserType{
				DisplayName: "Core",
				ID:          "2",
			},
		},
	}

	actual, err := user.UserManagementUpdateUser(updateUserInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
func TestUnitDeleteUser(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testDeleteUserResponseJson, http.StatusCreated)
	deleteUserInput := UserManagementDeleteUser{
		ID: mockUserId,
	}

	expected := &UserManagementDeleteUserPayload{
		DeletedUser: UserManagementDeletedUser{
			ID: mockUserId,
		},
	}

	actual, err := user.UserManagementDeleteUser(deleteUserInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitGetUser(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testGetUserResponseJSON, http.StatusCreated)

	expected := &UserManagementAuthenticationDomains{
		AuthenticationDomains: []UserManagementAuthenticationDomain{
			{
				ID:   mockAuthenticationDomainId,
				Name: mockAuthenticationDomainName,
				Users: UserManagementUsers{
					Users: []UserManagementUser{
						{
							ID:                     mockUserId,
							Name:                   userName,
							Email:                  mockUserEmail,
							EmailVerificationState: "Pending",
							TimeZone:               "Etc/UTC",
							PendingUpgradeRequest:  UserManagementPendingUpgradeRequest{},
							Type: UserManagementUserType{
								DisplayName: "Basic",
								ID:          "1",
							},
							Groups: UserManagementUserGroups{
								Groups: []UserManagementUserGroup{},
							},
						},
					},
				},
			},
		},
		NextCursor: "",
		TotalCount: 1,
	}

	actual, err := user.UserManagementGetUsers(
		[]string{mockAuthenticationDomainId},
		[]string{mockUserId},
		"",
		"",
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
