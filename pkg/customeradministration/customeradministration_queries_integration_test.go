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

func TestIntegrationCustomerAdministration_GetUser(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	getUserResponse, err := client.GetUser()
	require.NoError(t, err)
	require.NotNil(t, getUserResponse.ID)
}

func TestIntegrationCustomerAdministration_GetUsers(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	getUsersResponse, err := client.GetUsers(
		"",
		MultiTenantIdentityUserFilterInput{
			// Authentication Domain ID is required in this filter
			AuthenticationDomainId: &MultiTenantIdentityAuthenticationDomainIdInput{Eq: "WORK-IN-PROGRESS"},
			Email:                  MultiTenantIdentityUserEmailInput{},
			EmailVerificationState: MultiTenantIdentityEmailVerificationStateInput{},
			ID:                     MultiTenantIdentityUserIdInput{},
			Name:                   MultiTenantIdentityUserNameInput{},
			PendingUpgradeRequest:  MultiTenantIdentityPendingUpgradeRequestInput{},
		},
		[]MultiTenantIdentityUserSortInput{},
	)

	require.NoError(t, err)
	require.Greater(t, getUsersResponse.TotalCount, 0)
}

// This is currently throwing an unauthorized error, need to check

//func TestIntegrationCustomerAdministration_GetOrganizations(t *testing.T) {
//	t.Parallel()
//	_, err := mock.GetTestAccountID()
//	if err != nil {
//		t.Skipf("%s", err)
//	}
//
//	client := newIntegrationTestClient(t)
//	getOrganizationsResponse, err := client.GetOrganizations(
//		"",
//		OrganizationCustomerOrganizationFilterInput{},
//	)
//
//	require.NoError(t, err)
//	require.Greater(t, len(getOrganizationsResponse.Items), 0)
//}

// This is currently failing - needs to be fixed as this is caused at a schema level

//func TestIntegrationCustomerAdministration_GetGroups(t *testing.T) {
//	t.Parallel()
//	_, err := mock.GetTestAccountID()
//	if err != nil {
//		t.Skipf("%s", err)
//	}
//
//	client := newIntegrationTestClient(t)
//	getAccountsResponse, err := client.GetGroups(
//		"",
//		MultiTenantIdentityGroupFilterInput{
//			AllowsCapability:       MultiTenantIdentityAllowsCapabilityInput{},
//			AuthenticationDomainId: MultiTenantIdentityAuthenticationDomainIdInput{},
//			ID:                     MultiTenantIdentityGroupIdInput{},
//			Members:                MultiTenantIdentityGroupMemberIdInput{Contains: nil},
//			Name:                   MultiTenantIdentityGroupNameInput{},
//			OrganizationId:         MultiTenantIdentityOrganizationIdInput{"WORK-IN-PROGRESS"},
//		},
//		[]MultiTenantIdentityGroupSortInput{},
//	)
//
//	require.NoError(t, err)
//	fmt.Println(getAccountsResponse)
//	require.Greater(t, getAccountsResponse.TotalCount, 0)
//}
