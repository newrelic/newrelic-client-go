package accountmanagement

import (
	"fmt"
	"log"
	"testing"
	"time"

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

func TestIntegrationAccountManagement_CreateUpdateCancelAccount(t *testing.T) {
	t.Parallel()
	accountManagementClient := newAccountManagementTestClient(t)

	// Create Account
	name := "client-go-test-account-" + mock.RandSeq(5)
	createAccountInput := AccountManagementCreateInput{
		Name:       name,
		RegionCode: "us01",
	}
	createAccountResponse, err := accountManagementClient.AccountManagementCreateAccount(createAccountInput)

	require.Nil(t, err)
	require.NotNil(t, createAccountResponse.ManagedAccount.ID)
	require.Equal(t, createAccountInput.RegionCode, createAccountResponse.ManagedAccount.RegionCode)
	require.Equal(t, createAccountInput.Name, createAccountResponse.ManagedAccount.Name)
	time.Sleep(time.Second * 2)

	// Update Account
	updateAccountInput := AccountManagementUpdateInput{
		ID:   createAccountResponse.ManagedAccount.ID,
		Name: name + "-updated",
	}
	updateAccountResponse, err := accountManagementClient.AccountManagementUpdateAccount(updateAccountInput)
	fmt.Println("updateAccountResponse", updateAccountResponse)

	require.Nil(t, err)
	require.NotNil(t, updateAccountResponse.ManagedAccount.ID)
	require.Equal(t, updateAccountResponse.ManagedAccount.ID, createAccountResponse.ManagedAccount.ID)
	require.Equal(t, updateAccountInput.Name, updateAccountResponse.ManagedAccount.Name)
	time.Sleep(time.Second * 3)

	// Get Account
	getAccountResponse, err := accountManagementClient.GetManagedAccounts()
	fmt.Println("getAccountResponse", getAccountResponse)
	require.Nil(t, err)
	require.NotNil(t, getAccountResponse)
	foundAccountInGetResponse := false

	for _, account := range *getAccountResponse {
		if account.ID == updateAccountResponse.ManagedAccount.ID {
			foundAccountInGetResponse = true
			break
		}
	}

	require.True(t, foundAccountInGetResponse)
	time.Sleep(time.Second * 3)

	// Cancel Account
	cancelAccountResponse, err := accountManagementClient.AccountManagementCancelAccount(createAccountResponse.ManagedAccount.ID)

	require.Nil(t, err)
	require.NotNil(t, cancelAccountResponse)
	time.Sleep(time.Second * 2)

	// Get Account to Confirm Account Cancellation based on the value of `isCanceled`
	isCancelled := true
	getAccountResponse, err = accountManagementClient.GetManagedAccountsWithAdditionalArguments(&isCancelled)

	require.Nil(t, err)
	require.NotNil(t, getAccountResponse)
	foundAccountInGetResponse = false

	for _, account := range *getAccountResponse {
		if account.ID == updateAccountResponse.ManagedAccount.ID {
			foundAccountInGetResponse = true
			require.True(t, account.IsCanceled)
			break
		}
	}

	require.True(t, foundAccountInGetResponse)

}

func TestIntegrationGetManagedAccounts(t *testing.T) {
	t.Parallel()
	accountManagementClient := newAccountManagementTestClient(t)

	actual, _ := accountManagementClient.GetManagedAccounts()

	log.Println(actual)
	require.NotNil(t, actual)
	require.NotZero(t, len(*actual))
}

func TestIntegrationGetManagedAccountsModified_CanceledAccounts(t *testing.T) {
	t.Parallel()
	accountManagementClient := newAccountManagementTestClient(t)

	cancelled := true
	actual, _ := accountManagementClient.GetManagedAccountsWithAdditionalArguments(&cancelled)
	log.Println(actual)
	require.NotNil(t, actual)
	require.NotZero(t, len(*actual))
}

func TestIntegrationGetManagedAccountsModified_NonCanceledAccounts(t *testing.T) {
	t.Parallel()
	accountManagementClient := newAccountManagementTestClient(t)

	cancelled := false
	actual, _ := accountManagementClient.GetManagedAccountsWithAdditionalArguments(&cancelled)
	log.Println(actual)
	require.NotNil(t, actual)
	require.NotZero(t, len(*actual))
}

func TestIntegrationGetManagedAccountsModified_AllCancellationStatuses(t *testing.T) {
	t.Parallel()
	accountManagementClient := newAccountManagementTestClient(t)

	actual, _ := accountManagementClient.GetManagedAccountsWithAdditionalArguments(nil)

	require.NotNil(t, actual)
	require.NotZero(t, len(*actual))

	foundCancelledAccount := false
	foundUncancelledAccount := false

	for _, acct := range *actual {
		if foundUncancelledAccount == true && foundCancelledAccount == true {
			break
		}
		if acct.IsCanceled == true {
			foundCancelledAccount = true
		}
		if acct.IsCanceled == false {
			foundUncancelledAccount = true
		}
	}

	require.True(t, foundCancelledAccount)
	require.True(t, foundUncancelledAccount)
}
