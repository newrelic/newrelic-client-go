package customeradministration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCustomerAdministrationGetAccountsResponseJSON = `{
		"data": {
			"customerAdministration": {
				"accounts": {
					"items": [
						{
							"id": ` + unitTestMockAccountOneId + `
							"name": "` + unitTestMockAccountOneName + `",
							"regionCode": "us01",
							"status": "active"
						},
						{
							"id": ` + unitTestMockAccountTwoId + `,
							"name": "` + unitTestMockAccountTwoName + `",
							"regionCode": "us01",
							"status": "active"
						}
					],
					"nextCursor": "` + unitTestMockNextCursor + `",
					"totalCount": 2
				}
			}
		}
	}`
)

func TestUnit_CustomerAdministration_GetAccounts(t *testing.T) {
	t.Parallel()

	customeradministration := newMockResponse(t, testCustomerAdministrationGetAccountsResponseJSON, http.StatusOK)

	t.Skipf("Skipping this test as this test needs a fix, based on current API behaviour")

	expected := OrganizationAccountCollection{
		Items: []OrganizationAccount{
			{
				ID:         unitTestMockAccountOneIdAsInt,
				Name:       unitTestMockAccountOneName,
				RegionCode: "us01",
				Status:     "active",
			},
			{
				ID:         unitTestMockAccountTwoIdAsInt,
				Name:       unitTestMockAccountTwoName,
				RegionCode: "us01",
				Status:     "active",
			},
		},
		NextCursor: unitTestMockNextCursor,
		TotalCount: 2,
	}

	actual, err := customeradministration.GetAccounts(
		unitTestMockNextCursor,
		OrganizationAccountFilterInput{
			Name:           OrganizationAccountNameFilterInput{},
			OrganizationId: OrganizationAccountOrganizationIdFilterInput{unitTestMockOrganizationId},
			SharingMode:    OrganizationAccountSharingModeFilterInput{},
			Status:         OrganizationAccountStatusFilterInput{},
		},
		[]OrganizationAccountSortInput{},
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
