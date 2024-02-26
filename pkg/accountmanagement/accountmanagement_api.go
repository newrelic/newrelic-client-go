// Code generated by tutone: DO NOT EDIT
package accountmanagement

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v3/pkg/errors"
)

// Creates an organization-scoped account.
func (a *Accountmanagement) AccountManagementCreateAccount(
	managedAccount AccountManagementCreateInput,
) (*AccountManagementCreateResponse, error) {
	return a.AccountManagementCreateAccountWithContext(context.Background(),
		managedAccount,
	)
}

// Creates an organization-scoped account.
func (a *Accountmanagement) AccountManagementCreateAccountWithContext(
	ctx context.Context,
	managedAccount AccountManagementCreateInput,
) (*AccountManagementCreateResponse, error) {

	resp := AccountManagementCreateAccountQueryResponse{}
	vars := map[string]interface{}{
		"managedAccount": managedAccount,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AccountManagementCreateAccountMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AccountManagementCreateResponse, nil
}

type AccountManagementCreateAccountQueryResponse struct {
	AccountManagementCreateResponse AccountManagementCreateResponse `json:"AccountManagementCreateAccount"`
}

const AccountManagementCreateAccountMutation = `mutation(
	$managedAccount: AccountManagementCreateInput!,
) { accountManagementCreateAccount(
	managedAccount: $managedAccount,
) {
	managedAccount {
		id
		name
		regionCode
	}
} }`

// Updates an account.
func (a *Accountmanagement) AccountManagementUpdateAccount(
	managedAccount AccountManagementUpdateInput,
) (*AccountManagementUpdateResponse, error) {
	return a.AccountManagementUpdateAccountWithContext(context.Background(),
		managedAccount,
	)
}

// Updates an account.
func (a *Accountmanagement) AccountManagementUpdateAccountWithContext(
	ctx context.Context,
	managedAccount AccountManagementUpdateInput,
) (*AccountManagementUpdateResponse, error) {

	resp := AccountManagementUpdateAccountQueryResponse{}
	vars := map[string]interface{}{
		"managedAccount": managedAccount,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AccountManagementUpdateAccountMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AccountManagementUpdateResponse, nil
}

type AccountManagementUpdateAccountQueryResponse struct {
	AccountManagementUpdateResponse AccountManagementUpdateResponse `json:"AccountManagementUpdateAccount"`
}

const AccountManagementUpdateAccountMutation = `mutation(
	$managedAccount: AccountManagementUpdateInput!,
) { accountManagementUpdateAccount(
	managedAccount: $managedAccount,
) {
	managedAccount {
		id
		name
		regionCode
	}
} }`

// Admin-level info about the accounts in an organization.
func (a *Accountmanagement) GetManagedAccounts() (*[]AccountManagementManagedAccount, error) {
	return a.GetManagedAccountsWithContext(context.Background())
}

// Admin-level info about the accounts in an organization.
func (a *Accountmanagement) GetManagedAccountsWithContext(
	ctx context.Context,
) (*[]AccountManagementManagedAccount, error) {

	resp := managedAccountsResponse{}
	vars := map[string]interface{}{}

	if err := a.client.NerdGraphQueryWithContext(ctx, getManagedAccountsQuery, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.Actor.Organization.AccountManagement.ManagedAccounts) == 0 {
		return nil, errors.NewNotFound("")
	}

	return &resp.Actor.Organization.AccountManagement.ManagedAccounts, nil
}

const getManagedAccountsQuery = `query { actor { organization { accountManagement { managedAccounts {
	id
	name
	regionCode
} } } } }`
