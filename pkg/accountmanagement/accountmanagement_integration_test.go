package accountmanagement

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateAccount(t *testing.T) {
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
	updateAcountInput := AccountManagementUpdateInput{
		ID:   3833494,
		Name: name,
	}
	acctMgmt := newAccountManagementTestClient(t)
	actual, err := acctMgmt.AccountManagementUpdateAccount(updateAcountInput)
	require.Nil(t, err)
	require.NotNil(t, actual.ManagedAccount.RegionCode)
	require.Equal(t, updateAcountInput.ID, actual.ManagedAccount.ID)
	require.Equal(t, updateAcountInput.Name, actual.ManagedAccount.Name)
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
	updateAcountInput := AccountManagementUpdateInput{
		ID:   123345,
		Name: name,
	}
	acctMgmt := newAccountManagementTestClient(t)
	actual, err := acctMgmt.AccountManagementUpdateAccount(updateAcountInput)
	require.Nil(t, actual)
	require.NotNil(t, err)
}
