// +build integration

package cloud

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestCloudAccount_Basic(t *testing.T) {
	t.Skipf("Skipping this test tdue to an upstream API failure")

	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	testARN := os.Getenv("INTEGRATION_TESTING_AWS_ARN")
	if testARN == "" {
		t.Skip("an AWS ARN is required to run cloud account tests")
		return
	}

	a := newIntegrationTestClient(t)
	// Reset everything
	getResponse, err := a.GetLinkedAccounts("aws")
	require.NoError(t, err)

	for _, linkedAccount := range *getResponse {
		if linkedAccount.NrAccountId == testAccountID {
			a.CloudUnlinkAccount(testAccountID, []CloudUnlinkAccountsInput{
				{
					LinkedAccountId: linkedAccount.ID,
				},
			})
		}
	}

	// Link the account
	linkResponse, err := a.CloudLinkAccount(testAccountID, CloudLinkCloudAccountsInput{
		Aws: []CloudAwsLinkAccountInput{
			{
				Arn:  testARN,
				Name: "DTK Integration Testing",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, linkResponse)

	// Get the linked account
	getResponse, err = a.GetLinkedAccounts("aws")
	require.NoError(t, err)

	var linkedAccountID int
	for _, linkedAccount := range *getResponse {
		if linkedAccount.NrAccountId == testAccountID {
			linkedAccountID = linkedAccount.ID
			break
		}
	}
	require.NotZero(t, linkedAccountID)

	// Rename the account
	newName := "NEW-DTK-NAME"
	renameResponse, err := a.CloudRenameAccount(testAccountID, []CloudRenameAccountsInput{
		{
			LinkedAccountId: linkedAccountID,
			Name:            newName,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, renameResponse)

	// Unlink the account
	// unlinkResponse, err := a.CloudUnlinkAccount(testAccountID, CloudUnlinkAccountsInput{linkedAccountID})
	// require.NoError(t, err)
	// require.NotNil(t, unlinkResponse)
}

func newIntegrationTestClient(t *testing.T) Cloud {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
