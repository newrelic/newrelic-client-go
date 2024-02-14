//go:build unit
// +build unit

package usermanagement

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testGetAuthenticationDomainsRequestJSON = `{
	"data": {
		"actor": {
			"organization": {
				"userManagement": {
					"authenticationDomains": {
						"authenticationDomains": [{
							"id": "` + mockAuthenticationDomainId + `",
							"name": "` + mockAuthenticationDomainName + `",
							"provisioningType": "manual"
						}, {
							"id": "` + mockAuthenticationDomainIdTwo + `",
							"name": "` + mockAuthenticationDomainNameTwo + `",
							"provisioningType": "manual"
						}],
						"nextCursor": null,
						"totalCount": 2
					}
				}
			}
		}
	}
}`
)

func TestUnitGetAuthenticationDomains(t *testing.T) {
	t.Parallel()
	user := newMockResponse(t, testGetAuthenticationDomainsRequestJSON, http.StatusCreated)

	expected := &UserManagementAuthenticationDomains{
		[]UserManagementAuthenticationDomain{
			{ID: mockAuthenticationDomainId, Name: mockAuthenticationDomainName, ProvisioningType: "manual"},
			{ID: mockAuthenticationDomainIdTwo, Name: mockAuthenticationDomainNameTwo, ProvisioningType: "manual"},
		},
		"",
		2,
	}

	actual, err := user.GetAuthenticationDomains("", []string{""})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
