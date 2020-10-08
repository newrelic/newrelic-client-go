// +build integration

package cloud

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestCloud_RenameAccounts(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	// DTK terraform account
	// accountID := 2508259
	//
	// input := CloudRenameAccountsInput{}
	// input.LinkedAccountId = 48552
	// input.Name = "NEW-DTK-NAME"
	//
	// response, err := a.CloudRenameAccount(accountID, []CloudRenameAccountsInput{input})
	// require.NoError(t, err)
	// t.Logf("response: %+v", response)

	response, err := a.GetLinkedAccounts("aws")
	require.NoError(t, err)
	require.NotNil(t, response)

}

func newIntegrationTestClient(t *testing.T) Cloud {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
