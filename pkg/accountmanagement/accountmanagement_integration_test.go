package accountmanagement

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCreateAccount(t *testing.T) {
	t.Skipf("skipping create account test case as there is no delete account API")
	t.Parallel()
	name := "Test-" + mock.RandSeq(5)
	createAccountInput := AccountManagementCreateInput{
		Name:       name,
		RegionCode: "us01",
	}
	acctMgmt := newAccountManagementTestClient(t)
	actual, err := acctMgmt.AccountManagementCreateAccount(createAccountInput)
	require.Nil(t, err)
	require.NotNil(t, actual.ManagedAccount.ID)
	require.Equal(t, createAccountInput.RegionCode, actual.ManagedAccount.RegionCode)
	require.Equal(t, createAccountInput.Name, actual.ManagedAccount.Name)
}

func TestIntegrationUpdateAccount(t *testing.T) {
	t.Parallel()
	name := "Dont Delete-" + mock.RandSeq(5)
	updateAccountInput := AccountManagementUpdateInput{
		ID:   3833494,
		Name: name,
	}
	acctMgmt := newAccountManagementTestClient(t)
	actual, err := acctMgmt.AccountManagementUpdateAccount(updateAccountInput)
	require.Nil(t, err)
	require.NotNil(t, actual.ManagedAccount.RegionCode)
	require.Equal(t, updateAccountInput.ID, actual.ManagedAccount.ID)
	require.Equal(t, updateAccountInput.Name, actual.ManagedAccount.Name)
}

func TestIntegrationCreateAccountError(t *testing.T) {
	t.Parallel()
	name := "Test-" + mock.RandSeq(5)
	createAccountInput := AccountManagementCreateInput{
		Name:       name,
		RegionCode: "test",
	}
	acctMgmt := newAccountManagementTestClient(t)
	actual, err := acctMgmt.AccountManagementCreateAccount(createAccountInput)
	require.Nil(t, actual)
	require.NotNil(t, err)
}

func TestIntegrationUpdateAccountError(t *testing.T) {
	t.Parallel()
	name := "Dont Delete-" + mock.RandSeq(5)
	updateAccountInput := AccountManagementUpdateInput{
		ID:   123345,
		Name: name,
	}
	acctMgmt := newAccountManagementTestClient(t)
	actual, err := acctMgmt.AccountManagementUpdateAccount(updateAccountInput)
	require.Nil(t, actual)
	require.NotNil(t, err)
}
