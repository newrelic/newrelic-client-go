//go:build unit
// +build unit

package organization

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testGetOrganizationResponseJSON = `{
		"data": {
			"customerAdministration": {
				"organizations": {
					"items": [
						{
							"contractId": "` + contractId + `",
							"customerId": "` + customerId + `",
							"id": "` + org1Id + `",
							"name": "` + org1Name + `"
						},
						{
							"contractId": "` + contractId + `",
							"customerId": "` + customerId + `",
							"id": "` + org2Id + `",
							"name": "` + org2Name + `"
						}
					],
					"nextCursor": null
				}
			}
		}
	}`

	testCreateOrganizationResponseJSON = `{
		"data": {
			"organizationCreate": {
				"jobId": "` + orgCreateJobId + `"
			}
		}
	}`
)

func TestUnitCreateOrganization(t *testing.T) {
	t.Parallel()

	organization := newMockResponse(t, testCreateOrganizationResponseJSON, http.StatusCreated)

	expected := &OrganizationCreateOrganizationResponse{
		JobId: orgCreateJobId,
	}

	actual, err := organization.OrganizationCreate(customerId, OrganizationNewManagedAccountInput{
		Name: "Test Account",
	}, OrganizationCreateOrganizationInput{
		Name: org1Name,
	}, OrganizationSharedAccountInput{})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)

}

func TestUnitGetOrganizations(t *testing.T) {
	t.Parallel()

	organization := newMockResponse(t, testGetOrganizationResponseJSON, http.StatusCreated)

	expected := &OrganizationCustomerOrganizationWrapper{
		[]OrganizationCustomerOrganization{
			{ContractId: contractId, CustomerId: customerId, ID: org1Id, Name: org1Name},
			{ContractId: contractId, CustomerId: customerId, ID: org2Id, Name: org2Name},
		},
		"",
	}

	actual, err := organization.OrganizationGetOrganizations(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
