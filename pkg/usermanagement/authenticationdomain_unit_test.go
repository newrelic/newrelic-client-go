package usermanagement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testGetAllAuthDomains = `{
    "actor": {
      "organization": {
        "userManagement": {
          "authenticationDomains": {
            "authenticationDomains": [
              {
                "id": "0cc21d98-8dc2-484a-bb26-258e17ede584",
                "name": "Default",
                "provisioningType": "manual"
              },
              {
                "id": "a8e96cbe-b430-436a-bc1f-9b27875cabab",
                "name": "test_new_auth_domain",
                "provisioningType": "manual"
              },
              {
                "id": "d589c6fc-7f6a-4a0e-8539-18c492f7bb2d",
                "name": "test_saml_auth",
                "provisioningType": "scim"
              }
            ],
			"nextCursor": null,
 			"totalCount": 3
          }
        }
      }
    }
  }`

	testGetAuthDomain = `{
    "actor": {
      "organization": {
        "userManagement": {
          "authenticationDomains": {
            "authenticationDomains": [
              {
                "id": "0cc21d98-8dc2-484a-bb26-258e17ede584",
                "name": "Default",
                "provisioningType": "manual"
              }
            ],
			"nextCursor": null,
 			"totalCount": 1
          }
        }
      }
    }
  }`
)


func TestGetAllAuthDomains(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetAllAuthDomains)
	authDomains := newMockResponse(t, respJSON, http.StatusCreated)

	authdomain1 :=UserManagementAuthenticationDomain{
		ID: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		ProvisioningType: "manual",
		Name: "Default",
	}
	authdomain2 :=UserManagementAuthenticationDomain{
		ID: "a8e96cbe-b430-436a-bc1f-9b27875cabab",
		ProvisioningType: "manual",
		Name: "test_new_auth_domain",
	}
	authdomain3 :=UserManagementAuthenticationDomain{
		ID: "d589c6fc-7f6a-4a0e-8539-18c492f7bb2d",
		ProvisioningType: "scim",
		Name: "test_saml_auth",
	}
	authDomainsList:=[]UserManagementAuthenticationDomain{authdomain1,authdomain2,authdomain3}
	expected := &UserManagementAuthenticationDomains{
		AuthenticationDomains: authDomainsList,
		TotalCount: 3,
		NextCursor: "",
	}
	actual, err := authDomains.GetAuthenticationDomains("", nil)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetAuthDomain(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetAuthDomain)
	authDomain := newMockResponse(t, respJSON, http.StatusCreated)

	authdomain1 :=UserManagementAuthenticationDomain{
		ID: "0cc21d98-8dc2-484a-bb26-258e17ede584",
		ProvisioningType: "manual",
		Name: "Default",
	}

	authDomainsList:=[]UserManagementAuthenticationDomain{authdomain1}
	expected := &UserManagementAuthenticationDomains{
		AuthenticationDomains: authDomainsList,
		TotalCount: 1,
		NextCursor: "",
	}
	actual, err := authDomain.GetAuthenticationDomains("", []string{"0cc21d98-8dc2-484a-bb26-258e17ede584"})
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}


