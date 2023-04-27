//go:build integration
// +build integration

package cloud

import (
	"os"
	"strings"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestCloudAccount_Basic(t *testing.T) {
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
	unlinkResponse, err := a.CloudUnlinkAccount(testAccountID, []CloudUnlinkAccountsInput{
		{
			LinkedAccountId: linkedAccountID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, unlinkResponse)
}

func TestCloudAccount_SingleLinkedAccount(t *testing.T) {
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

	if len(linkResponse.LinkedAccounts) == 0 {
		t.Skip("skipping TestCloudAccount_SingleLinkedAccount due to no linked accounts")
		return
	}

	// Get the linked account
	linkedAccountId := linkResponse.LinkedAccounts[0].ID
	getResponse, err := a.GetLinkedAccount(testAccountID, linkedAccountId)
	foundLinkedAccountId := getResponse.ID

	require.NoError(t, err)
	require.Equal(t, linkedAccountId, foundLinkedAccountId)

	// Rename the account
	newName := "NEW-DTK-NAME"
	renameResponse, err := a.CloudRenameAccount(testAccountID, []CloudRenameAccountsInput{
		{
			LinkedAccountId: linkedAccountId,
			Name:            newName,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, renameResponse)

	// Unlink the account
	unlinkResponse, err := a.CloudUnlinkAccount(testAccountID, []CloudUnlinkAccountsInput{
		{
			LinkedAccountId: linkedAccountId,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, unlinkResponse)
}

func TestCloudAccount_GCPIntegrations(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	testGCPProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	if testGCPProjectID == "" {
		t.Skip("A GCP Project ID is required to run test")
		return
	}

	a := newIntegrationTestClient(t)

	// Link the account
	linkResponse, linkErr := a.CloudLinkAccount(testAccountID, CloudLinkCloudAccountsInput{
		Gcp: []CloudGcpLinkAccountInput{
			{
				Name:      "OAC-Test-Account",
				ProjectId: testGCPProjectID,
			},
		},
	})

	require.NoError(t, linkErr)
	require.NotNil(t, linkResponse)

	if len(linkResponse.LinkedAccounts) == 0 {
		t.Skip("skipping TestCloudAccount_GCPIntegrations due to no linked accounts")
		return
	}

	linkedAccountId := linkResponse.LinkedAccounts[0].ID

	integrationsRes, integrationsErr := a.CloudConfigureIntegration(testAccountID, CloudIntegrationsInput{
		Gcp: CloudGcpIntegrationsInput{
			GcpAlloydb:          []CloudGcpAlloydbIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpAppengine:        []CloudGcpAppengineIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpBigquery:         []CloudGcpBigqueryIntegrationInput{{LinkedAccountId: linkedAccountId, MetricsPollingInterval: 400, FetchTags: true, FetchTableMetrics: true}},
			GcpBigtable:         []CloudGcpBigtableIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpComposer:         []CloudGcpComposerIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpDataflow:         []CloudGcpDataflowIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpDataproc:         []CloudGcpDataprocIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpDatastore:        []CloudGcpDatastoreIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirebasedatabase: []CloudGcpFirebasedatabaseIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirebasehosting:  []CloudGcpFirebasehostingIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirebasestorage:  []CloudGcpFirebasestorageIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirestore:        []CloudGcpFirestoreIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFunctions:        []CloudGcpFunctionsIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpInterconnect:     []CloudGcpInterconnectIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpKubernetes:       []CloudGcpKubernetesIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpLoadbalancing:    []CloudGcpLoadbalancingIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpMemcache:         []CloudGcpMemcacheIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpPubsub:           []CloudGcpPubsubIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpRedis:            []CloudGcpRedisIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpRouter:           []CloudGcpRouterIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpRun:              []CloudGcpRunIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpSpanner:          []CloudGcpSpannerIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpSql:              []CloudGcpSqlIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpStorage:          []CloudGcpStorageIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpVms:              []CloudGcpVmsIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpVpcaccess:        []CloudGcpVpcaccessIntegrationInput{{LinkedAccountId: linkedAccountId}},
		},
	})

	require.NoError(t, integrationsErr)
	require.NotNil(t, integrationsRes)
	require.Len(t, integrationsRes.Errors, 0)
	require.Greater(t, len(integrationsRes.Integrations), 0)

	disableRes, disableErr := a.CloudDisableIntegration(testAccountID, CloudDisableIntegrationsInput{
		Gcp: CloudGcpDisableIntegrationsInput{
			GcpAlloydb:          []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpAppengine:        []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpBigquery:         []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpVpcaccess:        []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpVms:              []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpStorage:          []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpSpanner:          []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpSql:              []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpRun:              []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpRouter:           []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpPubsub:           []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpMemcache:         []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpRedis:            []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpLoadbalancing:    []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpKubernetes:       []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpInterconnect:     []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFunctions:        []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirestore:        []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirebasestorage:  []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirebasehosting:  []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpDatastore:        []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpDataproc:         []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpDataflow:         []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpComposer:         []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpBigtable:         []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
			GcpFirebasedatabase: []CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountId}},
		}})

	require.NoError(t, disableErr)
	require.NotNil(t, disableRes.DisabledIntegrations)
	require.Len(t, disableRes.Errors, 0)
	require.Greater(t, len(integrationsRes.Integrations), 0)

	// Unlink the account
	unlinkResponse, err := a.CloudUnlinkAccount(testAccountID, []CloudUnlinkAccountsInput{
		{
			LinkedAccountId: linkedAccountId,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, unlinkResponse)
}

func newIntegrationTestClient(t *testing.T) Cloud {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}

func TestCloudAccount_AzureMonitorIntegration(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	azureCredentials := map[string]string{
		"INTEGRATION_TESTING_AZURE_APPLICATION_ID":   os.Getenv("INTEGRATION_TESTING_AZURE_APPLICATION_ID"),
		"INTEGRATION_TESTING_AZURE_CLIENT_SECRET_ID": os.Getenv("INTEGRATION_TESTING_AZURE_CLIENT_SECRET_ID"),
		"INTEGRATION_TESTING_AZURE_SUBSCRIPTION_ID":  os.Getenv("INTEGRATION_TESTING_AZURE_SUBSCRIPTION_ID"),
		"INTEGRATION_TESTING_AZURE_TENANT_ID":        os.Getenv("INTEGRATION_TESTING_AZURE_TENANT_ID"),
	}

	var credentialsNotFound []string

	for key, value := range azureCredentials {
		if value == "" {
			credentialsNotFound = append(credentialsNotFound, key)
		}
	}

	if len(credentialsNotFound) != 0 {
		t.Skipf("Skipping this test, as the following required Azure credentials do not exist in the environment: \n%s", strings.Join(credentialsNotFound[:], ", "))
	}

	// Reset everything - unlink the account if linked already.
	getResponse, err := client.GetLinkedAccounts("azure")
	require.NoError(t, err)

	for _, linkedAccount := range *getResponse {
		if linkedAccount.NrAccountId == testAccountID {
			client.CloudUnlinkAccount(testAccountID, []CloudUnlinkAccountsInput{
				{
					LinkedAccountId: linkedAccount.ID,
				},
			})
		}
	}

	// Link the account
	linkResponse, err := client.CloudLinkAccount(testAccountID, CloudLinkCloudAccountsInput{
		Azure: []CloudAzureLinkAccountInput{
			{
				Name:           "TEST_AZURE_ACCOUNT",
				ApplicationID:  azureCredentials["INTEGRATION_TESTING_AZURE_APPLICATION_ID"],
				ClientSecret:   SecureValue(azureCredentials["INTEGRATION_TESTING_AZURE_CLIENT_SECRET_ID"]),
				SubscriptionId: azureCredentials["INTEGRATION_TESTING_AZURE_SUBSCRIPTION_ID"],
				TenantId:       azureCredentials["INTEGRATION_TESTING_AZURE_TENANT_ID"],
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, linkResponse)

	// Get the linked account
	getResponse, err = client.GetLinkedAccounts("azure")
	require.NoError(t, err)

	var linkedAccountID int
	for _, linkedAccount := range *getResponse {
		if linkedAccount.NrAccountId == testAccountID {
			linkedAccountID = linkedAccount.ID
			break
		}
	}
	require.NotZero(t, linkedAccountID)

	// Create a new AzureMonitor Cloud Integration.
	azureMonitorIntegrationRes, azureMonitorIntegrationErr := client.CloudConfigureIntegration(testAccountID, CloudIntegrationsInput{
		Azure: CloudAzureIntegrationsInput{
			AzureMonitor: []CloudAzureMonitorIntegrationInput{{
				LinkedAccountId:        linkedAccountID,
				Enabled:                true,
				ExcludeTags:            []string{"env:staging", "env:testing"},
				IncludeTags:            []string{"env:production"},
				MetricsPollingInterval: 1200,
				ResourceTypes:          []string{"microsoft.datashare/accounts"},
				ResourceGroups:         []string{"resource_groups"},
			}},
		},
	})

	require.NoError(t, azureMonitorIntegrationErr)
	require.NotNil(t, azureMonitorIntegrationRes)
	require.Len(t, azureMonitorIntegrationRes.Errors, 0)
	require.Greater(t, len(azureMonitorIntegrationRes.Integrations), 0)

	// Delete the created AzureMonitor Cloud Integration.
	azureMonitorDisableIntegrationRes, azureMonitorDisableIntegrationErr := client.CloudDisableIntegration(testAccountID, CloudDisableIntegrationsInput{
		Azure: CloudAzureDisableIntegrationsInput{
			AzureMonitor: []CloudDisableAccountIntegrationInput{{
				LinkedAccountId: linkedAccountID,
			}},
		},
	})

	require.NoError(t, azureMonitorDisableIntegrationErr)
	require.NotNil(t, azureMonitorDisableIntegrationRes)
	require.Len(t, azureMonitorDisableIntegrationRes.Errors, 0)

	// Unlink the linked account.
	unlinkResponse, err := client.CloudUnlinkAccount(testAccountID, []CloudUnlinkAccountsInput{
		{
			LinkedAccountId: linkedAccountID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, unlinkResponse)
}
