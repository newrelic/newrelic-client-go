//go:build unit
// +build unit

package usermanagement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCreateGroupResponseJSON = `{
	"userManagementCreateGroup": {
      "group": {
        "displayName": "api-test",
        "id": "66cdd38c-45a2-4106-95a0-8198f5cff775",
        "users": {
          "totalCount": 0,
          "users": []
        }
      }
    }
}`

	testUpdateGroupResponseJSON = `{
	"userManagementUpdateGroup": {
      "group": {
        "displayName": "api-test-name-update",
        "id": "66cdd38c-45a2-4106-95a0-8198f5cff775",
        "users": {
          "nextCursor": null,
          "totalCount": 0,
          "users": []
        }
      }
    }
  }`
	testDeletedGroupResponseJson = `{
	"userManagementDeleteGroup": {
      "group": {
        "displayName": "api-test-name-update",
        "id": "66cdd38c-45a2-4106-95a0-8198f5cff775",
        "users": {
          "nextCursor": null,
          "totalCount": 0,
          "users": []
        }
      }
    }
}`
)

func TestCreateGroup(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateGroupResponseJSON)
	user := newMockResponse(t, respJSON, http.StatusCreated)
	id := "66cdd38c-45a2-4106-95a0-8198f5cff775"
	createInput := UserManagementCreateGroup{
		AuthenticationDomainId: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		DisplayName:            "api-test",
	}

	createdResp := UserManagementGroup{
		DisplayName: createInput.DisplayName,
		ID:          id,
		Users: UserManagementGroupUsers{
			TotalCount: 0,
			Users:      []UserManagementGroupUser{},
		},
	}
	expected := &UserManagementCreateGroupPayload{
		Group: createdResp,
	}
	actual, err := user.UserManagementCreateGroup(createInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

// temporarily commenting this test out, until the right query is identified and used here
// func TestReadGroup(t *testing.T) {
// 	t.Parallel()
// 	respJSON := fmt.Sprintf(`{ "data":%s }`, testReadGroupResponseJson)
// 	user := newMockResponse(t, respJSON, http.StatusCreated)
// 	id := "66cdd38c-45a2-4106-95a0-8198f5cff775"
// }

func TestUpdateGroup(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateGroupResponseJSON)
	user := newMockResponse(t, respJSON, http.StatusCreated)
	id := "66cdd38c-45a2-4106-95a0-8198f5cff775"

	updateInput := UserManagementUpdateGroup{
		DisplayName: "api-test-name-update",
		ID:          id,
	}

	updateResp := UserManagementGroup{
		DisplayName: updateInput.DisplayName,
		ID:          id,
		Users: UserManagementGroupUsers{
			TotalCount: 0,
			Users:      []UserManagementGroupUser{},
		},
	}
	expected := &UserManagementUpdateGroupPayload{
		Group: updateResp,
	}
	actual, err := user.UserManagementUpdateGroup(updateInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteGroup(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeletedGroupResponseJson)
	user := newMockResponse(t, respJSON, http.StatusCreated)
	id := "66cdd38c-45a2-4106-95a0-8198f5cff775"
	deleteInput := UserManagementDeleteGroup{
		ID: id,
	}
	deletedResp := UserManagementGroup{
		DisplayName: "api-test-name-update",
		ID:          deleteInput.ID,
		Users: UserManagementGroupUsers{
			TotalCount: 0,
			Users:      []UserManagementGroupUser{},
		},
	}
	expected := &UserManagementDeleteGroupPayload{
		Group: deletedResp,
	}

	actual, err := user.UserManagementDeleteGroup(deleteInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
