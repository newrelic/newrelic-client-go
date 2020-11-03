package infrastructure

type LinkCloudAccountInput struct {
	AccountID int                     `json:"accountId"`
	AWS       []LinkAWSAccountInput   `json:"aws"`
	Azure     []LinkAzureAccountInput `json:"azure"`
	GCP       []LinkGCPAccountInput   `json:"gcp"`
}

type LinkAWSAccountInput struct {
	Name string `json:"name"`
	ARN  string `json:"arn"`
}

type LinkAzureAccountInput struct {
	Name           string `json:"name"`
	ApplicationID  string `json:"applicationId"`
	ClientSecret   string `json:"clientSecret"`
	TenantID       string `json:"tenantId"`
	SubscriptionId string `json:"subscriptionId"`
}

type LinkGCPAccountInput struct {
	Name      string `json:"name"`
	ProjectId string `json:"projectId"`
}

type linkCloudAccountOutput struct {
	LinkedAccounts []LinkedCloudAccount `json:"linkedAccounts"`
}

type LinkedCloudAccount struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AuthLabel string `json:"authLabel"`
}

// LinkCloudAccount is the mutation to link a New Relic account to a Cloud account (AWS / Azure / GCP).
func (i *Infrastructure) LinkCloudAccount(account LinkCloudAccountInput) ([]LinkedCloudAccount, error) {
	vars := map[string]interface{}{
		"account": account,
	}

	var resp linkCloudAccountOutput
	if err := i.client.NerdGraphQuery(linkCloudAccountMutation, vars, &resp); err != nil {
		return nil, err
	}

	return resp.LinkedAccounts, nil
}

type UnlinkCloudAccountInput struct {
	Accounts []LinkedCloudAccountRef `json:"accounts"`
}

type LinkedCloudAccountRef struct {
	LinkedAccountId int `json:"linkedAccountId"`
}

type UnlinkedCloudAccount struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type unlinkCloudAccountOutput struct {
	UnlinkedAccounts []UnlinkedCloudAccount `json:"unlinkedAccounts"`
}

// UnlinkCloudAccount is the mutation to unlink a New Relic account from a Cloud account (AWS / Azure / GCP).
func (i *Infrastructure) UnlinkCloudAccount(accountID int, accounts UnlinkCloudAccountInput) ([]UnlinkedCloudAccount, error) {
	vars := map[string]interface{}{
		"accountID": accountID,
		"accounts":  accounts,
	}

	resp := unlinkCloudAccountOutput{}
	if err := i.client.NerdGraphQuery(unlinkCloudAccountMutation, vars, &resp); err != nil {
		return nil, err
	}

	return resp.UnlinkedAccounts, nil
}

const (
	graphqlErrors = `
    errors {
      type
      message
    }
    `

	linkCloudAccountMutation = `
		mutation($accounts: LinkCloudAccountInput!) {
			cloudLinkAccount(accounts: $accounts) {
				linkedAccounts {
					id
					name
					authLabel
				}` +
		graphqlErrors +
		`} }`

	unlinkCloudAccountMutation = `
		mutation($accountID: Int!, $accounts: UnlinkCloudAccountInput!) {
			cloudUnlinkAccount(accountId: $accountID, accounts: $accounts) {
				unlinkedAccounts {
					id
					name
				}` +
		graphqlErrors +
		`} }`
)
