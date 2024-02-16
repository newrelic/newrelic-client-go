package customeradministration

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCustomerAdministration_GetAccounts(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	getAccountsResponse, err := client.GetAccounts(
		"",
		OrganizationAccountFilterInput{
			ID:             OrganizationAccountIdFilterInput{Eq: mock.IntegrationTestAccountID},
			Name:           OrganizationAccountNameFilterInput{},
			OrganizationId: OrganizationAccountOrganizationIdFilterInput{"WORK-IN-PROGRESS"},
			SharingMode:    OrganizationAccountSharingModeFilterInput{},
			Status:         OrganizationAccountStatusFilterInput{},
		},
		[]OrganizationAccountSortInput{},
	)
	require.NoError(t, err)
	require.Equal(t, getAccountsResponse.Items[0].ID, mock.IntegrationTestAccountID)
}
