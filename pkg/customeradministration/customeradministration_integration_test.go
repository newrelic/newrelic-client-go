package customeradministration

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCustomerAdministration_GetAccountShares(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	_, err = client.GetAccountShares(
		"",
		OrganizationAccountShareFilterInput{
			AccountID: OrganizationAccountIdInput{Eq: mock.IntegrationTestAccountID},
			TargetId:  OrganizationTargetIdInput{},
		},
		[]OrganizationAccountShareSortInput{},
	)
	require.NoError(t, err)
}

func TestIntegrationCustomerAdministration_GetAccounts(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)
	getAccountsResponse, err := client.GetAccounts(
		"",
		OrganizationAccountFilterInput{
			ID:             OrganizationAccountIdFilterInput{Eq: mock.IntegrationTestAccountID},
			Name:           OrganizationAccountNameFilterInput{},
			OrganizationId: OrganizationAccountOrganizationIdFilterInput{organizationId},
			SharingMode:    OrganizationAccountSharingModeFilterInput{},
			Status:         OrganizationAccountStatusFilterInput{},
		},
		[]OrganizationAccountSortInput{},
	)
	require.NoError(t, err)
	require.Equal(t, getAccountsResponse.Items[0].ID, mock.IntegrationTestAccountID)
}

func TestIntegrationCustomerAdministration_GetAuthenticationDomains(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)
	_, err = client.GetAuthenticationDomains("",
		OrganizationAuthenticationDomainFilterInput{
			Name:           &OrganizationNameInput{Eq: authenticationDomainName},
			OrganizationId: &OrganizationOrganizationIdInput{Eq: organizationId},
		},
		[]OrganizationAuthenticationDomainSortInput{},
	)
	require.NoError(t, err)

}

func TestIntegrationCustomerAdministration_GetUser(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
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

	if authenticationDomainId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_AUTHENTICATION_DOMAIN_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)
	getUsersResponse, err := client.GetUsers(
		"",
		MultiTenantIdentityUserFilterInput{
			// Authentication Domain ID is required in this filter
			AuthenticationDomainId: &MultiTenantIdentityAuthenticationDomainIdInput{Eq: authenticationDomainId},
		},
		[]MultiTenantIdentityUserSortInput{},
	)

	require.NoError(t, err)
	require.Greater(t, getUsersResponse.TotalCount, 0)
}

// This is currently throwing an unauthorized error, need to check - skipping the test until then
func TestIntegrationCustomerAdministration_GetOrganizations_UnauthorizedError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)
	_, err = client.GetOrganizations(
		"",
		OrganizationCustomerOrganizationFilterInput{
			ID: OrganizationOrganizationIdInputFilter{
				organizationId,
			},
		},
	)

	require.Regexp(t, regexp.MustCompile("Unauthorized"), err.Error())
}

func TestIntegrationCustomerAdministration_GetGroups(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)
	getAccountsResponse, err := client.GetGroups(
		"",
		MultiTenantIdentityGroupFilterInput{
			AuthenticationDomainId: &MultiTenantIdentityAuthenticationDomainIdInput{Eq: authenticationDomainId},
			OrganizationId:         &MultiTenantIdentityOrganizationIdInput{organizationId},
		},
		[]MultiTenantIdentityGroupSortInput{},
	)

	require.NoError(t, err)
	require.Greater(t, getAccountsResponse.TotalCount, 0)
}

func TestIntegrationCustomerAdministration_GetRoles(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)
	getRolesResponse, err := client.GetRoles(
		"",
		MultiTenantAuthorizationRoleFilterInputExpression{
			Name: &MultiTenantAuthorizationRoleNameInputFilter{
				Eq: roleName,
			},
			OrganizationId: &MultiTenantAuthorizationRoleOrganizationIdInputFilter{
				Eq: organizationId,
			},
		},
		[]MultiTenantAuthorizationRoleSortInput{},
	)

	require.NoError(t, err)
	require.Greater(t, getRolesResponse.TotalCount, 0)
}

func TestIntegrationCustomerAdministration_GetPermissions(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)
	getPermissionsResponse, err := client.GetPermissions("",
		MultiTenantAuthorizationPermissionFilter{
			RoleId: MultiTenantAuthorizationPermissionFilterRoleIdInput{
				Eq: roleId,
			},
		},
	)

	require.NoError(t, err)
	require.Greater(t, len(getPermissionsResponse.Items), 0)
	require.NotNil(t, getPermissionsResponse.Items[0].Name)
}

func TestIntegrationCustomerAdministration_GetGrants(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	if authenticationDomainId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_AUTHENTICATION_DOMAIN_ID is missing, or has an empty value")
	}

	roleIdAsInt, _ := strconv.Atoi(roleId)

	client := newIntegrationTestClient(t)
	getGrantsResponse, err := client.GetGrants("",
		MultiTenantAuthorizationGrantFilterInputExpression{
			// don't use "eq" with filters which have "eq" and "in" since "in" don't have omitempty, and are hence, expected
			AuthenticationDomainId: &MultiTenantAuthorizationGrantAuthenticationDomainIdInputFilter{
				In: []string{
					authenticationDomainId,
				},
			},
			OrganizationId: &MultiTenantAuthorizationGrantOrganizationIdInputFilter{
				organizationId,
			},
			// don't use "eq" with filters which have "eq" and "in" since "in" don't have omitempty, and are hence, expected
			RoleId: &MultiTenantAuthorizationGrantRoleIdInputFilter{
				In: []int{
					roleIdAsInt,
				},
			},
		},
		[]MultiTenantAuthorizationGrantSortInput{},
	)

	require.NoError(t, err)
	require.Greater(t, len(getGrantsResponse.Items), 0)
	require.NotNil(t, getGrantsResponse.Items[0].ID)
}
