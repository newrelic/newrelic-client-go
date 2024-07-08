package organization

import (
	"fmt"
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

	testUpdateSharedAccountResponseJSON = `{
		"data":  {
			"organizationUpdateSharedAccount":  {
				"sharedAccount":  {
					"accountId":  ` + fmt.Sprint(unitTestMockAccountOneId) + `,
					"id":  "` + unitTestMockOrganizationOneId + `",
					"limitingRoleId":  ` + fmt.Sprint(unitTestMockLimitingRoleId) + `,
					"name":  "` + unitTestMockOrganizationOneName + `",
					"sourceOrganizationId":  "` + unitTestMockOrganizationOneId + `",
					"sourceOrganizationName":  "` + unitTestMockOrganizationOneName + `"
				}
			}
		}
	}`

	testRevokeSharedAccountResponseJSON = `{
		"data":  {
			"organizationRevokeSharedAccount":  {
				"sharedAccount":  {
					"accountId":  ` + fmt.Sprint(unitTestMockAccountOneId) + `,
					"id":  "` + unitTestMockOrganizationOneId + `",
					"limitingRoleId":  ` + fmt.Sprint(unitTestMockLimitingRoleId) + `,
					"name":  "` + unitTestMockOrganizationOneName + `",
					"sourceOrganizationId":  "` + unitTestMockOrganizationOneId + `",
					"sourceOrganizationName":  "` + unitTestMockOrganizationOneName + `"
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

func TestUnitOrganizationUpdateSharedAccount(t *testing.T) {
	t.Parallel()

	organization := newMockResponse(t, testUpdateSharedAccountResponseJSON, http.StatusOK)

	expected := &OrganizationUpdateSharedAccountResponse{
		SharedAccount: OrganizationSharedAccount{
			AccountID:              unitTestMockAccountOneId,
			ID:                     unitTestMockOrganizationOneId,
			LimitingRoleId:         unitTestMockLimitingRoleId,
			Name:                   unitTestMockOrganizationOneName,
			SourceOrganizationId:   unitTestMockOrganizationOneId,
			SourceOrganizationName: unitTestMockOrganizationOneName,
		},
	}

	actual, err := organization.OrganizationUpdateSharedAccount(
		OrganizationUpdateSharedAccountInput{
			ID:             fmt.Sprint(unitTestMockAccountOneId),
			LimitingRoleId: unitTestMockLimitingRoleId,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUnitOrganizationRevokeSharedAccount(t *testing.T) {
	t.Parallel()

	organization := newMockResponse(t, testRevokeSharedAccountResponseJSON, http.StatusOK)

	expected := &OrganizationRevokeSharedAccountResponse{
		SharedAccount: OrganizationSharedAccount{
			AccountID:              unitTestMockAccountOneId,
			ID:                     unitTestMockOrganizationOneId,
			LimitingRoleId:         unitTestMockLimitingRoleId,
			Name:                   unitTestMockOrganizationOneName,
			SourceOrganizationId:   unitTestMockOrganizationOneId,
			SourceOrganizationName: unitTestMockOrganizationOneName,
		},
	}

	actual, err := organization.OrganizationRevokeSharedAccount(
		OrganizationRevokeSharedAccountInput{
			ID: fmt.Sprint(unitTestMockAccountOneId),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
