package organization

import (
	"fmt"
	"regexp"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
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

	organizationCreateResponse, err := client.OrganizationCreate(
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

	// expected behaviour
	// require.Error(t, err)
	// require.Regexp(t, regexp.MustCompile(fmt.Sprintf("%s not found.", unitTestMockCustomerId)), err.Error())

	// current behaviour
	// the API isn't throwing an error, the error is instead, embedded in the response

	require.NoError(t, err)
	require.True(t, matchOrganizationUnauthorizedErrorRegex(organizationCreateResponse.JobId))
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

	organizationUpdateResponse, _ := client.OrganizationUpdate(
		OrganizationUpdateInput{
			Name: organizationNameUpdated,
		},
		unitTestMockOrganizationOneId,
	)

	// expected behaviour
	// require.Error(t, err)

	// actual behaviour
	// the error thrown is embedded as a field inside the response of the mutation

	require.NotNil(t, organizationUpdateResponse.Errors)
	require.True(t, matchOrganizationUnauthorizedErrorRegex(organizationUpdateResponse.Errors[0].Message))
}
