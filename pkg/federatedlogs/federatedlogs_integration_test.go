//go:build integration
// +build integration

package federatedlogs

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

const orgID = "fb33fea3-4d7e-4736-9701-acb59a634fdf"

// getFleetTestOrganizationID returns the org used for entity scope when
// running the Setup/Partition tests under fleet-test credentials.
func getFleetTestOrganizationID() string {
	if id := os.Getenv("NEW_RELIC_FLEET_TEST_ORGANIZATION_ID"); id != "" {
		return id
	}
	return "b961cf81-d62b-4359-8822-7b1d6dadd374"
}

func TestIntegrationFederatedLogs_AwsConnection(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	ctx := context.Background()

	connectionID, cleanup := createTestAwsConnection(t, client, testAccountID, orgID)
	t.Cleanup(cleanup)

	// Read the connection back via the polymorphic GetEntity call.
	getResp, err := client.GetEntityWithContext(ctx, connectionID)
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.NotNil(t, *getResp, "GetEntity returned nil interface for newly-created connection")

	entity, ok := (*getResp).(*EntityManagementAwsConnectionEntity)
	require.True(t, ok, "expected *EntityManagementAwsConnectionEntity, got %T", *getResp)
	require.Equal(t, connectionID, entity.ID)

	enabled := false
	updateInput := EntityManagementAwsConnectionEntityUpdateInput{
		Description: "updated by integration test",
		Enabled:     &enabled,
		Region:      "us-west-2",
		Credential: &EntityManagementAwsCredentialsUpdateInput{
			AssumeRole: EntityManagementAwsAssumeRoleConfigUpdateInput{
				RoleArn: "arn:aws:iam::123456789012:role/nr-test-role-rotated",
			},
		},
	}
	// EntityManagement updates use optimistic concurrency — pull the current
	// metadata.version from the entity we just fetched.
	updateResp, err := client.EntityManagementUpdateAwsConnectionWithContext(ctx, updateInput, connectionID, entity.Metadata.Version)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, connectionID, updateResp.Entity.ID)
}

func TestIntegrationFederatedLogs_Setup(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetFleetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newFleetIntegrationTestClient(t)
	ctx := context.Background()

	fleetID, cleanupFleet := createTestFleet(t, ctx)
	t.Cleanup(cleanupFleet)

	scopeOrgID := getFleetTestOrganizationID()
	connectionID, cleanupConn := createTestAwsConnection(t, client, testAccountID, scopeOrgID)
	t.Cleanup(cleanupConn)

	setupName := mock.GenerateRandomName(0) + "-setup"
	createInput := FederatedLogsCreateSetupInput{
		Name:        setupName,
		Description: "Created by integration test",
		Storage: FederatedLogsSetupStorageInput{
			DataLocationBucket:     "nr-test-fed-logs-bucket",
			Database:               "nr_test_fed_logs_db",
			DataIngestConnectionId: connectionID,
			QueryConnectionId:      connectionID,
			CloudProviderConfiguration: FederatedLogsCloudProviderConfigurationInput{
				Provider: FederatedLogsCloudProviderTypes.AWS,
				Region:   "us-east-1",
			},
		},
		DefaultPartition: FederatedLogsDefaultPartitionInput{
			Storage: FederatedLogsPartitionStorageInput{
				Table:           "nr_test_default_partition",
				DataLocationUri: "s3://nr-test-fed-logs-bucket/nr_test_default_partition",
			},
		},
		Forwarder: &FederatedLogsForwarderInput{
			Type: FederatedLogsForwarderTypeTypes.PIPELINE_CONTROL,
			PipelineControl: &FederatedLogsPipelineControlConfigurationInput{
				FleetId: fleetID,
			},
		},
	}

	createResp, err := client.FederatedLogsCreateSetupWithContext(ctx, createInput)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotEmpty(t, createResp.Setup.ID)
	require.Equal(t, setupName, createResp.Setup.Name)

	setupID := createResp.Setup.ID

	// Cleanup — soft-delete via lifecycle transition.
	t.Cleanup(func() {
		deleteInput := FederatedLogsUpdateSetupInput{
			LifecycleStatus: &FederatedLogsLifecycleStatusInput{
				Status: FederatedLogsLifecycleStateTypes.DELETING,
			},
		}
		if _, err := client.FederatedLogsUpdateSetupWithContext(ctx, setupID, deleteInput); err != nil {
			t.Logf("cleanup: failed to soft-delete setup %s: %v", setupID, err)
		}
	})

	// Read back via the federatedLogs.setup query.
	getResp, err := client.GetSetupWithContext(ctx, setupID)
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, setupID, getResp.ID)
	require.Equal(t, setupName, getResp.Name)

	updatedDescription := "updated by integration test"
	updateInput := FederatedLogsUpdateSetupInput{
		Description: updatedDescription,
		Forwarder: &FederatedLogsForwarderInput{
			Type: FederatedLogsForwarderTypeTypes.PIPELINE_CONTROL,
			PipelineControl: &FederatedLogsPipelineControlConfigurationInput{
				FleetId: fleetID,
				RoutingRule: &FederatedLogsRuleInput{
					Expression: `attributes["service.name"] == "integration-test"`,
				},
			},
		},
	}
	updateResp, err := client.FederatedLogsUpdateSetupWithContext(ctx, setupID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, setupID, updateResp.Setup.ID)
	require.Equal(t, updatedDescription, updateResp.Setup.Description)
}

func TestIntegrationFederatedLogs_Partition(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetFleetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newFleetIntegrationTestClient(t)
	ctx := context.Background()

	fleetID, cleanupFleet := createTestFleet(t, ctx)
	t.Cleanup(cleanupFleet)

	scopeOrgID := getFleetTestOrganizationID()
	connectionID, cleanupConn := createTestAwsConnection(t, client, testAccountID, scopeOrgID)
	t.Cleanup(cleanupConn)

	// Need a parent setup to host the partition.
	setupID := createTestSetup(t, client, ctx, connectionID, fleetID)

	partitionName := mock.GenerateRandomName(0) + "-partition"
	createInput := FederatedLogsCreatePartitionInput{
		Name:        partitionName,
		Description: "Created by integration test",
		Storage: FederatedLogsPartitionStorageInput{
			Table:           "nr_test_secondary_partition",
			DataLocationUri: "s3://nr-test-fed-logs-bucket/nr_test_secondary_partition",
		},
		DataRetentionPolicy: &FederatedLogsRetentionPolicyInput{
			Duration: 7,
			Unit:     FederatedLogsRetentionUnitTypes.DAYS,
		},
	}

	createResp, err := client.FederatedLogsCreatePartitionWithContext(ctx, createInput, setupID)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotEmpty(t, createResp.Partition.ID)
	require.Equal(t, partitionName, createResp.Partition.Name)
	require.False(t, createResp.Partition.IsDefault)

	partitionID := createResp.Partition.ID

	t.Cleanup(func() {
		deleteInput := FederatedLogsUpdatePartitionInput{
			LifecycleStatus: &FederatedLogsLifecycleStatusInput{
				Status: FederatedLogsLifecycleStateTypes.DELETING,
			},
		}
		if _, err := client.FederatedLogsUpdatePartitionWithContext(ctx, partitionID, deleteInput); err != nil {
			t.Logf("cleanup: failed to soft-delete partition %s: %v", partitionID, err)
		}
	})

	getResp, err := client.GetPartitionWithContext(ctx, partitionID)
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, partitionID, getResp.ID)
	require.Equal(t, partitionName, getResp.Name)

	// In-place update — mutable fields per FederatedLogsUpdatePartitionInput.
	active := false
	updateInput := FederatedLogsUpdatePartitionInput{
		Description: "updated by integration test",
		Active:      &active,
		DataRetentionPolicy: &FederatedLogsRetentionPolicyInput{
			Duration: 14,
			Unit:     FederatedLogsRetentionUnitTypes.DAYS,
		},
	}
	updateResp, err := client.FederatedLogsUpdatePartitionWithContext(ctx, partitionID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, partitionID, updateResp.Partition.ID)
}

func createTestAwsConnection(t *testing.T, client Federatedlogs, accountID int, scopeOrgID string) (string, func()) {
	t.Helper()
	ctx := context.Background()

	enabled := true
	input := EntityManagementAwsConnectionEntityCreateInput{
		Name:        mock.GenerateRandomName(0) + "-conn",
		Description: "Created by integration test",
		Enabled:     &enabled,
		Region:      "us-east-1",
		Credential: EntityManagementAwsCredentialsCreateInput{
			AssumeRole: EntityManagementAwsAssumeRoleConfigCreateInput{
				RoleArn: "arn:aws:iam::123456789012:role/nr-test-role",
			},
		},
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ORGANIZATION,
			ID:   scopeOrgID,
		},
	}

	resp, err := client.EntityManagementCreateAwsConnectionWithContext(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Entity.ID)

	connectionID := resp.Entity.ID
	cleanup := func() {
		// Need to fetch the current metadata.version for optimistic concurrency.
		getResp, err := client.GetEntityWithContext(ctx, connectionID)
		if err != nil || getResp == nil || *getResp == nil {
			t.Logf("cleanup: failed to read connection %s for delete: %v", connectionID, err)
			return
		}
		entity, ok := (*getResp).(*EntityManagementAwsConnectionEntity)
		if !ok {
			t.Logf("cleanup: unexpected entity type %T for connection %s", *getResp, connectionID)
			return
		}
		if _, err := client.EntityManagementDeleteWithContext(ctx, connectionID, entity.Metadata.Version); err != nil {
			t.Logf("cleanup: failed to delete connection %s: %v", connectionID, err)
		}
	}

	return connectionID, cleanup
}

// createTestSetup mints a FederatedLogsSetup using the supplied AWS connection
// for both ingest and query slots. Registers a t.Cleanup that soft-deletes the
// setup at end-of-test.
func createTestSetup(t *testing.T, client Federatedlogs, ctx context.Context, connectionID string, fleetID string) string {
	t.Helper()

	setupName := mock.GenerateRandomName(0) + "-setup"
	input := FederatedLogsCreateSetupInput{
		Name: setupName,
		Storage: FederatedLogsSetupStorageInput{
			DataLocationBucket:     "nr-test-fed-logs-bucket",
			Database:               "nr_test_fed_logs_db",
			DataIngestConnectionId: connectionID,
			QueryConnectionId:      connectionID,
			CloudProviderConfiguration: FederatedLogsCloudProviderConfigurationInput{
				Provider: FederatedLogsCloudProviderTypes.AWS,
				Region:   "us-east-1",
			},
		},
		DefaultPartition: FederatedLogsDefaultPartitionInput{
			Storage: FederatedLogsPartitionStorageInput{
				Table:           "nr_test_default_partition",
				DataLocationUri: "s3://nr-test-fed-logs-bucket/nr_test_default_partition",
			},
		},
		Forwarder: &FederatedLogsForwarderInput{
			Type: FederatedLogsForwarderTypeTypes.PIPELINE_CONTROL,
			PipelineControl: &FederatedLogsPipelineControlConfigurationInput{
				FleetId: fleetID,
			},
		},
	}

	resp, err := client.FederatedLogsCreateSetupWithContext(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Setup.ID)

	setupID := resp.Setup.ID
	t.Cleanup(func() {
		deleteInput := FederatedLogsUpdateSetupInput{
			LifecycleStatus: &FederatedLogsLifecycleStatusInput{
				Status: FederatedLogsLifecycleStateTypes.DELETING,
			},
		}
		if _, err := client.FederatedLogsUpdateSetupWithContext(ctx, setupID, deleteInput); err != nil {
			t.Logf("cleanup: failed to soft-delete setup %s: %v", setupID, err)
		}
	})

	return setupID
}

// createTestFleet provisions a Fleet entity for the lifetime of one test, using
// fleet-test credentials. Returns the fleet GUID and a cleanup func that
// deletes the fleet.
func createTestFleet(t *testing.T, ctx context.Context) (string, func()) {
	t.Helper()

	fleetClient := fleetcontrol.New(mock.NewFleetIntegrationTestConfig(t))

	fleetName := fmt.Sprintf("fed-logs-test-fleet-%d", time.Now().Unix())
	input := fleetcontrol.FleetControlFleetEntityCreateInput{
		Name:              fleetName,
		Description:       "Fleet created by federatedlogs integration test",
		ManagedEntityType: fleetcontrol.FleetControlManagedEntityTypeTypes.HOST,
		Product:           "Infrastructure",
		Scope: fleetcontrol.FleetControlScopedReferenceInput{
			ID:   getFleetTestOrganizationID(),
			Type: fleetcontrol.FleetControlEntityScopeTypes.ORGANIZATION,
		},
		OperatingSystem: &fleetcontrol.FleetControlOperatingSystemCreateInput{
			Type: fleetcontrol.FleetControlOperatingSystemTypeTypes.LINUX,
		},
	}

	resp, err := fleetClient.FleetControlCreateFleetWithContext(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Entity.ID)

	fleetID := resp.Entity.ID
	cleanup := func() {
		if _, err := fleetClient.FleetControlDeleteFleet(fleetID); err != nil {
			t.Logf("cleanup: failed to delete fleet %s: %v", fleetID, err)
		}
	}

	return fleetID, cleanup
}
