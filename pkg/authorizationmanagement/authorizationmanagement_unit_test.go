package authorizationmanagement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testAuthorizationManagementCreateResponse = `{
		"authorizationManagementGrantAccess": {
			"roles": [
			  {
				"accountId": 1111111,
				"displayName": "All Product Admin",
				"name": "all_product_admin",
				"organizationId": null,
				"roleId": 1254,
				"type": "Role::V2::Standard"
			  },
			  {
				"accountId": 1111112,
				"displayName": "All Product Admin",
				"name": "all_product_admin",
				"organizationId": null,
				"roleId": 1254,
				"type": "Role::V2::Standard"
			  }
			]
		}
  }`
)

func TestCreateAuthorizationManagement(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testAuthorizationManagementCreateResponse)
	authorizationManagement := newMockResponse(t, respJSON, http.StatusCreated)

	roleId := "1254"
	accountId := 1111111
	groupId := "1111112"

	authorizationManagementGrantAccessInput := AuthorizationManagementGrantAccess{
		AccountAccessGrants: []AuthorizationManagementAccountAccessGrant{
			{
				AccountID: accountId,
				RoleId:    roleId,
			},
		},
		GroupId: groupId,
	}

	actual, err := authorizationManagement.AuthorizationManagementGrantAccess(authorizationManagementGrantAccessInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.GreaterOrEqual(t, len(actual.Roles), 1)

	roleIdFound := false
	accountIdFound := false

	for _, role := range actual.Roles {
		if fmt.Sprint(role.RoleId) == roleId {
			roleIdFound = true
		}
		if role.AccountID == accountId {
			accountIdFound = true
		}
	}

	assert.True(t, roleIdFound, "Role not found")
	assert.True(t, accountIdFound, "Account not found")
}
