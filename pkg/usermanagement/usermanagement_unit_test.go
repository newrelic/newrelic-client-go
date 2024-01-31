//go:build unit
// +build unit

package usermanagement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	timestampString = "2022-07-25T12:08:07.179638Z"
	timestamp       = nrtime.DateTime(timestampString)
	user            = "test-user"
	accountId       = 10867072
	channelId       = "0d11fd42-5919-4767-8cf5-e07cb71c1b04"
	id              = "03bd4929-3d86-4447-a077-a901b5d511ff"

	testCreateUserRequestJSON = `{
    "userManagementCreateUser": {
      "createdUser": {
        "authenticationDomainId": "0cc21d98-8dc2-484a-bb26-258e17ede584",
        "email": "unittest@tf.com",
        "name": "unit test tf",
        "type": {
          "displayName": "Core"
        }
      }
    }
  }`

	testUpdateUserRequestJSON = `{
    "userManagementUpdateUser": {
      "user": {
        "email": "unittest@tf.com",
        "name": "test unit",
        "type": {
          "displayName": "Basic"
        }
      }
    }
  }`
	testDeletedResponseJson = `{
    "userManagementDeleteUser": {
      "deletedUser": {
        "id": "1000081687"
      }
    }
  }`
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Usermanagement {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateUserRequestJSON)
	user := newMockResponse(t, respJSON, http.StatusCreated)
	createUserInput := UserManagementCreateUser{
		AuthenticationDomainId: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		Name:                   "unit test tf",
		Email:                  "unittest@tf.com",
		UserType:               UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	}
	userManagementUserType := UserManagementUserType{
		DisplayName: "Core",
	}
	createdUser := UserManagementCreatedUser{
		Name:                   createUserInput.Name,
		AuthenticationDomainId: createUserInput.AuthenticationDomainId,
		Email:                  createUserInput.Email,
		Type:                   userManagementUserType,
	}
	expected := &UserManagementCreateUserPayload{
		CreatedUser: createdUser,
	}
	actual, err := user.UserManagementCreateUser(createUserInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateUserRequestJSON)
	user := newMockResponse(t, respJSON, http.StatusCreated)
	updateUserInput := UserManagementUpdateUser{
		Name:     "test unit",
		Email:    "unittest@tf.com",
		UserType: UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	}
	userManagementUserType := UserManagementUserType{
		DisplayName: "Basic",
	}
	updatedUser := UserManagementUser{
		Name:  updateUserInput.Name,
		Email: updateUserInput.Email,
		Type:  userManagementUserType,
	}
	expected := &UserManagementUpdateUserPayload{
		User: updatedUser,
	}
	actual, err := user.UserManagementUpdateUser(updateUserInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
func TestDeleteUser(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeletedResponseJson)
	user := newMockResponse(t, respJSON, http.StatusCreated)
	deleteUserInput := UserManagementDeleteUser{
		ID: "1000081687",
	}
	deletedUser := UserManagementDeletedUser{
		ID: deleteUserInput.ID,
	}
	expected := &UserManagementDeleteUserPayload{
		DeletedUser: deletedUser,
	}

	actual, err := user.UserManagementDeleteUser(deleteUserInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
