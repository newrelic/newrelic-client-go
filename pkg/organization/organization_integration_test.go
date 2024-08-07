package organization

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationOrganizationCreate_CustomerIdNotFoundError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	_, err = client.OrganizationCreate(
		unitTestMockCustomerId,
		&OrganizationNewManagedAccountInput{
			Name:       "Some Random Managed Account",
			RegionCode: OrganizationRegionCodeEnumTypes.US01,
		},
		OrganizationCreateOrganizationInput{
			Name: "Some Random Organization",
		},
		&OrganizationSharedAccountInput{
			AccountID:      0000,
			LimitingRoleId: 0000,
		},
	)

	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile(fmt.Sprintf("%s not found.", unitTestMockCustomerId)), err.Error())
}

func TestIntegrationOrganizationCreate_AccessDeniedError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	_, err = client.OrganizationCreate(
		"",
		&OrganizationNewManagedAccountInput{
			Name:       "Some Random Managed Account",
			RegionCode: OrganizationRegionCodeEnumTypes.US01,
		},
		OrganizationCreateOrganizationInput{
			Name: "Some Random Organization",
		},
		&OrganizationSharedAccountInput{
			AccountID:      0000,
			LimitingRoleId: 0000,
		},
	)

	// commenting this out since the following commented code was written on the basis of previous API behaviour
	// an error would not directly be thrown, instead, the error was embedded in the response previously

	//require.NoError(t, err)
	//require.True(t, matchOrganizationUnauthorizedErrorRegex(organizationCreateResponse.JobId))

	require.Error(t, err)
	require.True(t, matchOrganizationUnauthorizedErrorRegex(err.Error()))
}

func TestIntegrationOrganizationUpdate(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	if organizationId == "" {
		t.Errorf("cannot run integration test as the environment variable INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID is missing, or has an empty value")
	}

	client := newIntegrationTestClient(t)

	organizationUpdateResponse, _ := client.OrganizationUpdate(
		OrganizationUpdateInput{
			Name: organizationNameUpdated,
		},
		organizationId,
	)

	require.NotNil(t, organizationUpdateResponse.OrganizationInformation)
	require.Equal(t, organizationUpdateResponse.OrganizationInformation.Name, organizationNameUpdated)
}

func TestIntegrationOrganizationUpdate_AccessDeniedError(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	_, err = client.OrganizationUpdate(
		OrganizationUpdateInput{
			Name: organizationNameUpdated,
		},
		unitTestMockOrganizationOneId,
	)

	// commenting this out since the following commented code was written on the basis of previous API behaviour
	// an error would not directly be thrown, instead, the error was embedded in the response previously
	//require.NotNil(t, organizationUpdateResponse.Errors)
	//require.True(t, matchOrganizationUnauthorizedErrorRegex(organizationUpdateResponse.Errors[0].Message))

	require.Error(t, err)
	require.True(t, matchOrganizationUnauthorizedErrorRegex(err.Error()))
}

func TestIntegrationOrganizationRevokeSharedAccount_Error(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	response, _ := client.OrganizationRevokeSharedAccount(
		OrganizationRevokeSharedAccountInput{ID: fmt.Sprint(unitTestMockAccountOneId)},
	)

	// current API behaviour appears to be returning zero against the Shared Account ID if the operation is unsuccessful
	require.Zero(t, response.SharedAccount.AccountID)

	//require.Error(t, err)
	//require.True(t, matchOrganizationUnauthorizedErrorRegex(err.Error()))
}
