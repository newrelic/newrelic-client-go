package organization

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCreateOrganizationResponseJSON = `{
		"data": {
			"organizationCreate": {
				"jobId": "` + unitTestMockOrganizationCreateJobId + `"
			}
		}
	}`

	testUpdateOrganizationResponseJSON = `{
		"data":  {
			"organizationUpdate":  {
				"organizationInformation":  {
					"id":  "` + unitTestMockOrganizationOneId + `",
					"name":  "` + unitTestMockOrganizationOneName + `"
				}
			}
		}
	}`
)

func TestUnitCreateOrganization(t *testing.T) {
	t.Parallel()

	organization := newMockResponse(t, testCreateOrganizationResponseJSON, http.StatusCreated)

	expected := &OrganizationCreateOrganizationResponse{
		JobId: unitTestMockOrganizationCreateJobId,
	}

	actual, err := organization.OrganizationCreate(
		unitTestMockCustomerId,
		&OrganizationNewManagedAccountInput{
			Name: "Test Account",
		},
		OrganizationCreateOrganizationInput{
			Name: unitTestMockOrganizationOneName,
		},
		&OrganizationSharedAccountInput{},
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitUpdateOrganization(t *testing.T) {
	t.Parallel()

	organization := newMockResponse(t, testUpdateOrganizationResponseJSON, http.StatusCreated)

	expected := &OrganizationUpdateResponse{
		OrganizationInformation: OrganizationInformation{
			ID:   unitTestMockOrganizationOneId,
			Name: unitTestMockOrganizationOneName,
		},
	}

	actual, err := organization.OrganizationUpdate(
		OrganizationUpdateInput{
			Name: unitTestMockOrganizationOneName,
		},
		unitTestMockOrganizationOneId,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
