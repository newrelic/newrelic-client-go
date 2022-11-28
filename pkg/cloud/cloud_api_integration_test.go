//go:build integration
// +build integration

package cloud

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
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
