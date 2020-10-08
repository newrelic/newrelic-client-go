package cloud

import (
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// Create or modify a cloud integration.
//
// For details and mutation examples visit
// [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/manage-your-aws-azure-google-cloud-integrations-graphql-api).
func (a *Cloud) CloudConfigureIntegration(
	accountID int,
	integrations CloudIntegrationsInput,
) (*CloudConfigureIntegrationPayload, error) {

	resp := CloudConfigureIntegrationResponse{}
	vars := map[string]interface{}{
		"accountId":    accountID,
		"integrations": integrations,
	}

	if err := a.client.NerdGraphQuery(CloudConfigureIntegrationMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CloudConfigureIntegrationPayload, nil
}

type CloudConfigureIntegrationResponse struct {
	CloudConfigureIntegrationPayload CloudConfigureIntegrationPayload `json:"CloudConfigureIntegration"`
}

const CloudConfigureIntegrationMutation = `mutation(
	$accountId: Int!,
	$integrations: CloudIntegrationsInput!,
) { CloudConfigureIntegration(
	accountId: $accountId,
	integrations: $integrations,
) {
	
} }`

// Disable a cloud integration. Stops collecting data for the specified integration.
//
// For details and mutation examples visit
// [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/manage-your-aws-azure-google-cloud-integrations-graphql-api).
func (a *Cloud) CloudDisableIntegration(
	accountID int,
	integrations CloudDisableIntegrationsInput,
) (*CloudDisableIntegrationPayload, error) {

	resp := CloudDisableIntegrationResponse{}
	vars := map[string]interface{}{
		"accountId":    accountID,
		"integrations": integrations,
	}

	if err := a.client.NerdGraphQuery(CloudDisableIntegrationMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CloudDisableIntegrationPayload, nil
}

type CloudDisableIntegrationResponse struct {
	CloudDisableIntegrationPayload CloudDisableIntegrationPayload `json:"CloudDisableIntegration"`
}

const CloudDisableIntegrationMutation = `mutation(
	$accountId: Int,
	$integrations: CloudDisableIntegrationsInput,
) { CloudDisableIntegration(
	accountId: $accountId,
	integrations: $integrations,
) {
	
} }`

// Link a cloud provider account to a New Relic Account.
//
// For details and mutation examples visit
// [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/manage-your-aws-azure-google-cloud-integrations-graphql-api).
func (a *Cloud) CloudLinkAccount(
	accountID int,
	accounts CloudLinkCloudAccountsInput,
) (*CloudLinkAccountPayload, error) {

	resp := CloudLinkAccountResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"accounts":  accounts,
	}

	if err := a.client.NerdGraphQuery(CloudLinkAccountMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CloudLinkAccountPayload, nil
}

type CloudLinkAccountResponse struct {
	CloudLinkAccountPayload CloudLinkAccountPayload `json:"CloudLinkAccount"`
}

const CloudLinkAccountMutation = `mutation(
	$accountId: Int,
	$accounts: CloudLinkCloudAccountsInput,
) { CloudLinkAccount(
	accountId: $accountId,
	accounts: $accounts,
) {
	
} }`

// Rename one or more linked cloud provider accounts.
//
// For details and mutation examples visit
// [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/manage-your-aws-azure-google-cloud-integrations-graphql-api).
func (a *Cloud) CloudRenameAccount(
	accountID int,
	accounts CloudRenameAccountsInput,
) (*CloudRenameAccountPayload, error) {

	resp := CloudRenameAccountResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"accounts":  accounts,
	}

	if err := a.client.NerdGraphQuery(CloudRenameAccountMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CloudRenameAccountPayload, nil
}

type CloudRenameAccountResponse struct {
	CloudRenameAccountPayload CloudRenameAccountPayload `json:"CloudRenameAccount"`
}

const CloudRenameAccountMutation = `mutation(
	$accountId: Int,
	$accounts: [CloudRenameAccountsInput],
) { CloudRenameAccount(
	accountId: $accountId,
	accounts: $accounts,
) {
	
} }`

// Unlink one or more cloud provider accounts.
// Stops collecting data for all the associated integrations.
//
// For details and mutation examples visit
// [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/manage-your-aws-azure-google-cloud-integrations-graphql-api).
func (a *Cloud) CloudUnlinkAccount(
	accountID int,
	accounts CloudUnlinkAccountsInput,
) (*CloudUnlinkAccountPayload, error) {

	resp := CloudUnlinkAccountResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"accounts":  accounts,
	}

	if err := a.client.NerdGraphQuery(CloudUnlinkAccountMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.CloudUnlinkAccountPayload, nil
}

type CloudUnlinkAccountResponse struct {
	CloudUnlinkAccountPayload CloudUnlinkAccountPayload `json:"CloudUnlinkAccount"`
}

const CloudUnlinkAccountMutation = `mutation(
	$accountId: Int,
	$accounts: [CloudUnlinkAccountsInput],
) { CloudUnlinkAccount(
	accountId: $accountId,
	accounts: $accounts,
) {
	
} }`

// Get all linked cloud provider accounts scoped to the Actor.
func (a *Cloud) GetLinkedAccounts(
	provider string,
) (*[]CloudLinkedAccount, error) {

	resp := linkedAccountsResponse{}
	vars := map[string]interface{}{
		"provider": provider,
	}

	if err := a.client.NerdGraphQuery(getLinkedAccountsQuery, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.Actor.Cloud.LinkedAccounts) == 0 {
		return nil, errors.NewNotFound("")
	}

	return &resp.Actor.Cloud.LinkedAccounts, nil
}

const getLinkedAccountsQuery = `query(
	$provider: String,
) { actor { cloud { linkedAccounts(
	provider: $provider,
) {
	authLabel
	createdAt
	disabled
	externalId
	id
	integrations {
		__typename
		createdAt
		id
		name
		nrAccountId
		updatedAt
		... on CloudAlbIntegration {
			__typename
			awsRegions
			fetchExtendedInventory
			fetchTags
			inventoryPollingInterval
			loadBalancerPrefixes
			metricsPollingInterval
			tagKey
			tagValue
		}
		... on CloudApigatewayIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
			stagePrefixes
			tagKey
			tagValue
		}
		... on CloudAutoscalingIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsAppsyncIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsAthenaIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsDirectconnectIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsDocdbIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsGlueIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsMqIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsMskIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsQldbIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsStatesIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAwsWafIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudAzureApimanagementIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureAppserviceIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureCosmosdbIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureCostmanagementIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			tagKeys
		}
		... on CloudAzureFunctionsIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureLoadbalancerIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureMariadbIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureMysqlIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzurePostgresqlIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureRediscacheIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureServicebusIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureSqlIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureSqlmanagedIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureStorageIntegration {
			__typename
			id
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureVirtualmachineIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureVirtualnetworksIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudAzureVmsIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
			resourceGroups
		}
		... on CloudBaseIntegration {
			__typename
		}
		... on CloudBillingIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudCloudfrontIntegration {
			__typename
			fetchLambdasAtEdge
			fetchTags
			inventoryPollingInterval
			metricsPollingInterval
			tagKey
			tagValue
		}
		... on CloudCloudtrailIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudDynamodbIntegration {
			__typename
			awsRegions
			fetchExtendedInventory
			fetchTags
			inventoryPollingInterval
			metricsPollingInterval
			tagKey
			tagValue
		}
		... on CloudKinesisIntegration {
			__typename
			awsRegions
			fetchShards
			fetchTags
			inventoryPollingInterval
			metricsPollingInterval
			tagKey
			tagValue
		}
		... on CloudLambdaIntegration {
			__typename
			awsRegions
			fetchTags
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudRdsIntegration {
			__typename
			awsRegions
			fetchTags
			inventoryPollingInterval
			metricsPollingInterval
			tagKey
			tagValue
		}
		... on CloudRedshiftIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
			tagKey
			tagValue
		}
		... on CloudRoute53Integration {
			__typename
			fetchExtendedInventory
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudS3Integration {
			__typename
			fetchExtendedInventory
			fetchTags
			inventoryPollingInterval
			metricsPollingInterval
			tagKey
			tagValue
		}
		... on CloudSesIntegration {
			__typename
			awsRegions
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudSnsIntegration {
			__typename
			awsRegions
			fetchExtendedInventory
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudSqsIntegration {
			__typename
			awsRegions
			fetchExtendedInventory
			fetchTags
			inventoryPollingInterval
			metricsPollingInterval
			queuePrefixes
			tagKey
			tagValue
		}
		... on CloudTrustedadvisorIntegration {
			__typename
			inventoryPollingInterval
			metricsPollingInterval
		}
		... on CloudVpcIntegration {
			__typename
			awsRegions
			fetchNatGateway
			fetchVpn
			inventoryPollingInterval
			metricsPollingInterval
			tagKey
			tagValue
		}
	}
	name
	nrAccountId
	provider {
		__typename
		createdAt
		icon
		id
		name
		slug
		updatedAt
		... on CloudAwsGovCloudProvider {
			__typename
			awsAccountId
		}
		... on CloudAwsProvider {
			__typename
			roleAccountId
			roleExternalId
		}
		... on CloudBaseProvider {
			__typename
		}
		... on CloudGcpProvider {
			__typename
			serviceAccountId
		}
	}
	updatedAt
} } } }`
