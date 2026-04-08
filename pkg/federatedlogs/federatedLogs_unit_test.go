//go:build unit
// +build unit

package federatedlogs

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testAccountID            = "12345"
	testSetupID              = "MXxGRURFUkFURURfTE9HX1NFVFVQX0VOVElUWXwxMjM0NXw1Njc4OQ"
	testPartitionID          = "MXxGRURFUkFURURfTE9HX1BBUlRJVElPTl9FTlRJVFl8MTIzNDV8OTg3NjU"
	testAwsConnectionID      = "MXxBV1NfQ09OTkVDVElPTl9FTlRJVFl8MTIzNDV8MTExMTE"
	testDataProcessingConnID = "MXxEQVRBX1BST0NFU1NJTkdfQ09OTkVDVElPTnwxMjM0NXwyMjIyMg"
	testQueryConnectionID    = "MXxRVUVSWV9DT05ORUNUSU9OfDEyMzQ1fDMzMzMz"

	testCreateAwsConnectionResponseJSON = `
	{
		"data": {
			"entityManagementCreateAwsConnection": {
				"entity": {
					"id": "MXxBV1NfQ09OTkVDVElPTl9FTlRJVFl8MTIzNDV8MTExMTE",
					"name": "Test AWS Connection",
					"description": "Test AWS Connection for Federated Logs",
					"enabled": true,
					"externalId": "arn:aws:iam::123456789012:role/NewRelicRole",
					"region": "us-east-1",
					"type": "AwsConnectionEntity",
					"metadata": {
						"version": 1
					}
				}
			}
		}
	}`

	testCreateFederatedLogSetupResponseJSON = `
	{
		"data": {
			"entityManagementCreateFederatedLogSetup": {
				"entity": {
					"id": "MXxGRURFUkFURURfTE9HX1NFVFVQX0VOVElUWXwxMjM0NXw1Njc4OQ",
					"name": "Test Federated Log Setup",
					"description": "Test Federated Log Setup for S3",
					"cloudProvider": "AWS",
					"cloudProviderRegion": "us-east-1",
					"dataLocationBucket": "s3://my-log-bucket",
					"metadata": {
						"version": 1
					}
				}
			}
		}
	}`

	testCreateFederatedLogPartitionResponseJSON = `
	{
		"data": {
			"entityManagementCreateFederatedLogPartition": {
				"entity": {
					"id": "MXxGRURFUkFURURfTE9HX1BBUlRJVElPTl9FTlRJVFl8MTIzNDV8OTg3NjU",
					"name": "Test Log Partition",
					"description": "Test Federated Log Partition",
					"dataLocationURI": "s3://my-log-bucket/partition-path",
					"isDefault": false,
					"partitionTable": "logs_partition_table",
					"status": "ACTIVE",
					"metadata": {
						"version": 1
					}
				}
			}
		}
	}`

	testUpdateFederatedLogSetupResponseJSON = `
	{
		"data": {
			"entityManagementUpdateFederatedLogSetup": {
				"entity": {
					"id": "MXxGRURFUkFURURfTE9HX1NFVFVQX0VOVElUWXwxMjM0NXw1Njc4OQ",
					"name": "Test Federated Log Setup Updated",
					"description": "Updated Test Federated Log Setup",
					"cloudProvider": "AWS",
					"cloudProviderRegion": "us-west-2",
					"dataLocationBucket": "s3://my-updated-log-bucket",
					"metadata": {
						"version": 2
					}
				}
			}
		}
	}`

	testUpdateFederatedLogPartitionResponseJSON = `
	{
		"data": {
			"entityManagementUpdateFederatedLogPartition": {
				"entity": {
					"id": "MXxGRURFUkFURURfTE9HX1BBUlRJVElPTl9FTlRJVFl8MTIzNDV8OTg3NjU",
					"name": "Test Log Partition Updated",
					"description": "Updated Test Federated Log Partition",
					"isDefault": true,
					"partitionTable": "logs_partition_table_v2",
					"metadata": {
						"version": 2
					}
				}
			}
		}
	}`
)

func TestUnitEntityManagement_CreateAwsConnection(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateAwsConnectionResponseJSON, http.StatusOK)

	createInput := EntityManagementAwsConnectionEntityCreateInput{
		Name:        "Test AWS Connection",
		Description: "Test AWS Connection for Federated Logs",
		Enabled:     true,
		ExternalId:  "arn:aws:iam::123456789012:role/NewRelicRole",
		Region:      "us-east-1",
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreateAwsConnection(createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testAwsConnectionID, result.Entity.ID)
	require.Equal(t, "Test AWS Connection", result.Entity.Name)
	require.Equal(t, true, result.Entity.Enabled)
	require.Equal(t, "us-east-1", result.Entity.Region)
}

func TestUnitEntityManagement_CreateAwsConnectionWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateAwsConnectionResponseJSON, http.StatusOK)

	createInput := EntityManagementAwsConnectionEntityCreateInput{
		Name:        "Test AWS Connection",
		Description: "Test AWS Connection for Federated Logs",
		Enabled:     true,
		ExternalId:  "arn:aws:iam::123456789012:role/NewRelicRole",
		Region:      "us-east-1",
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	ctx := context.Background()
	result, err := client.EntityManagementCreateAwsConnectionWithContext(ctx, createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testAwsConnectionID, result.Entity.ID)
	require.Equal(t, "Test AWS Connection", result.Entity.Name)
}

func TestUnitEntityManagement_CreateFederatedLogSetup(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogSetupResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogSetupEntityCreateInput{
		Name:                       "Test Federated Log Setup",
		Description:                "Test Federated Log Setup for S3",
		CloudProvider:              EntityManagementCloudProviderTypes.AWS,
		CloudProviderRegion:        "us-east-1",
		DataLocationBucket:         "s3://my-log-bucket",
		DataProcessingConnectionId: testDataProcessingConnID,
		PartitionDatabase:          "federated_logs_db",
		QueryConnectionId:          testQueryConnectionID,
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreateFederatedLogSetup(createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testSetupID, result.Entity.ID)
	require.Equal(t, "Test Federated Log Setup", result.Entity.Name)
	require.Equal(t, "AWS", string(result.Entity.CloudProvider))
	require.Equal(t, "s3://my-log-bucket", result.Entity.DataLocationBucket)
}

func TestUnitEntityManagement_CreateFederatedLogSetupWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogSetupResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogSetupEntityCreateInput{
		Name:                       "Test Federated Log Setup",
		Description:                "Test Federated Log Setup for S3",
		CloudProvider:              EntityManagementCloudProviderTypes.AWS,
		CloudProviderRegion:        "us-east-1",
		DataLocationBucket:         "s3://my-log-bucket",
		DataProcessingConnectionId: testDataProcessingConnID,
		PartitionDatabase:          "federated_logs_db",
		QueryConnectionId:          testQueryConnectionID,
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	ctx := context.Background()
	result, err := client.EntityManagementCreateFederatedLogSetupWithContext(ctx, createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testSetupID, result.Entity.ID)
	require.Equal(t, "Test Federated Log Setup", result.Entity.Name)
}

func TestUnitEntityManagement_CreateFederatedLogPartition(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogPartitionResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogPartitionEntityCreateInput{
		Name:            "Test Log Partition",
		Description:     "Test Federated Log Partition",
		DataLocationUri: "s3://my-log-bucket/partition-path",
		IsDefault:       false,
		PartitionTable:  "logs_partition_table",
		SetupId:         testSetupID,
		Status:          EntityManagementLogPartitionStatusTypes.ACTIVE,
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
		RetentionPolicy: EntityManagementRetentionPolicyCreateInput{
			Duration: 30,
			Unit:     "DAYS",
		},
	}

	result, err := client.EntityManagementCreateFederatedLogPartition(createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testPartitionID, result.Entity.ID)
	require.Equal(t, "Test Log Partition", result.Entity.Name)
	require.Equal(t, false, result.Entity.IsDefault)
	require.Equal(t, "logs_partition_table", result.Entity.PartitionTable)
}

func TestUnitEntityManagement_CreateFederatedLogPartitionWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogPartitionResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogPartitionEntityCreateInput{
		Name:            "Test Log Partition",
		Description:     "Test Federated Log Partition",
		DataLocationUri: "s3://my-log-bucket/partition-path",
		IsDefault:       false,
		PartitionTable:  "logs_partition_table",
		SetupId:         testSetupID,
		Status:          EntityManagementLogPartitionStatusTypes.ACTIVE,
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	ctx := context.Background()
	result, err := client.EntityManagementCreateFederatedLogPartitionWithContext(ctx, createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testPartitionID, result.Entity.ID)
	require.Equal(t, "Test Log Partition", result.Entity.Name)
}

func TestUnitEntityManagement_UpdateFederatedLogSetup(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testUpdateFederatedLogSetupResponseJSON, http.StatusOK)

	updateInput := EntityManagementFederatedLogSetupEntityUpdateInput{
		Name:                "Test Federated Log Setup Updated",
		Description:         "Updated Test Federated Log Setup",
		CloudProviderRegion: "us-west-2",
		DataLocationBucket:  "s3://my-updated-log-bucket",
	}

	result, err := client.EntityManagementUpdateFederatedLogSetup(
		updateInput,
		testSetupID,
		2,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testSetupID, result.Entity.ID)
	require.Equal(t, "Test Federated Log Setup Updated", result.Entity.Name)
	require.Equal(t, "us-west-2", result.Entity.CloudProviderRegion)
	require.Equal(t, "s3://my-updated-log-bucket", result.Entity.DataLocationBucket)
}

func TestUnitEntityManagement_UpdateFederatedLogSetupWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testUpdateFederatedLogSetupResponseJSON, http.StatusOK)

	updateInput := EntityManagementFederatedLogSetupEntityUpdateInput{
		Name:                "Test Federated Log Setup Updated",
		Description:         "Updated Test Federated Log Setup",
		CloudProviderRegion: "us-west-2",
		DataLocationBucket:  "s3://my-updated-log-bucket",
	}

	ctx := context.Background()
	result, err := client.EntityManagementUpdateFederatedLogSetupWithContext(
		ctx,
		updateInput,
		testSetupID,
		2,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testSetupID, result.Entity.ID)
	require.Equal(t, "Test Federated Log Setup Updated", result.Entity.Name)
}

func TestUnitEntityManagement_UpdateFederatedLogPartition(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testUpdateFederatedLogPartitionResponseJSON, http.StatusOK)

	updateInput := EntityManagementFederatedLogPartitionEntityUpdateInput{
		Name:            "Test Log Partition Updated",
		Description:     "Updated Test Federated Log Partition",
		IsDefault:       true,
		PartitionTable:  "logs_partition_table_v2",
		RetentionPolicy: EntityManagementRetentionPolicyUpdateInput{
			Duration: 60,
			Unit:     "DAYS",
		},
	}

	result, err := client.EntityManagementUpdateFederatedLogPartition(
		updateInput,
		testPartitionID,
		2,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testPartitionID, result.Entity.ID)
	require.Equal(t, "Test Log Partition Updated", result.Entity.Name)
	require.Equal(t, true, result.Entity.IsDefault)
	require.Equal(t, "logs_partition_table_v2", result.Entity.PartitionTable)
}

func TestUnitEntityManagement_UpdateFederatedLogPartitionWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testUpdateFederatedLogPartitionResponseJSON, http.StatusOK)

	updateInput := EntityManagementFederatedLogPartitionEntityUpdateInput{
		Name:           "Test Log Partition Updated",
		Description:    "Updated Test Federated Log Partition",
		IsDefault:      true,
		PartitionTable: "logs_partition_table_v2",
	}

	ctx := context.Background()
	result, err := client.EntityManagementUpdateFederatedLogPartitionWithContext(
		ctx,
		updateInput,
		testPartitionID,
		2,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testPartitionID, result.Entity.ID)
	require.Equal(t, "Test Log Partition Updated", result.Entity.Name)
}
