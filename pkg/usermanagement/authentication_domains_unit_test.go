//go:build unit
// +build unit

package usermanagement

import (
	"net/http"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) UserManagement {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

var (
	testGetAuthenticationDomainsRequestJSON = `{
	"data": {
		"actor": {
			"organization": {
				"userManagement": {
					"authenticationDomains": {
						"authenticationDomains": [{
							"id": "fae55e6b-b1ce-4a0f-83b2-ee774798f2cc",
							"name": "Mock Authentication Domain 1",
							"provisioningType": "manual"
						}, {
							"id": "841acb6f-b14b-47d0-b811-3209ce97aa07",
							"name": "Mock Authentication Domain 2",
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
			{ID: "fae55e6b-b1ce-4a0f-83b2-ee774798f2cc", Name: "Mock Authentication Domain 1", ProvisioningType: "manual"},
			{ID: "841acb6f-b14b-47d0-b811-3209ce97aa07", Name: "Mock Authentication Domain 2", ProvisioningType: "manual"},
		},
		"",
		2,
	}

	actual, err := user.GetAuthenticationDomains("", []string{""})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
