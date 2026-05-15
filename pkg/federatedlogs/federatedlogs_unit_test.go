//go:build unit
// +build unit

package federatedlogs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testSetupID            = "ZmVkLWxvZ3Mtc2V0dXAtMTIzNDU2Nzg5MA"
	testPartitionID        = "ZmVkLWxvZ3MtcGFydGl0aW9uLTEyMzQ1Njc4OTA"
	testAwsConnectionID    = "ZmVkLWxvZ3MtYXdzLWNvbm4tMTIzNDU2Nzg5MA"
	testIngestConnectionID = "aW5nZXN0LWNvbm4tMTIzNDU2Nzg5MA"
	testQueryConnectionID  = "cXVlcnktY29ubi0xMjM0NTY3ODkwMA"
	testFleetID            = "ZmxlZXQtMTIzNDU2Nzg5MA"

	testCreateSetupResponseJSON = `
	{
		"data": {
			"federatedLogsCreateSetup": {
				"setup": {
					"id": "ZmVkLWxvZ3Mtc2V0dXAtMTIzNDU2Nzg5MA",
					"name": "Test Setup",
					"description": "Test federated logs setup",
					"active": true,
					"defaultPartitionId": "ZmVkLWxvZ3MtcGFydGl0aW9uLTEyMzQ1Njc4OTA",
					"storage": {
						"dataLocationBucket": "test-bucket",
						"database": "test_db",
						"dataIngestConnectionId": "aW5nZXN0LWNvbm4tMTIzNDU2Nzg5MA",
						"queryConnectionId": "cXVlcnktY29ubi0xMjM0NTY3ODkwMA",
						"cloudProviderConfiguration": {
							"provider": "AWS",
							"region": "us-east-1"
						}
					},
					"lifecycleStatus": {
						"status": "RESOURCE_CREATION_COMPLETE",
						"lastUpdatedAt": "2026-05-15T10:00:00Z"
					},
					"createdAt": "2026-05-15T10:00:00Z",
					"updatedAt": "2026-05-15T10:00:00Z"
				}
			}
		}
	}`

	testUpdateSetupResponseJSON = `
	{
		"data": {
			"federatedLogsUpdateSetup": {
				"setup": {
					"id": "ZmVkLWxvZ3Mtc2V0dXAtMTIzNDU2Nzg5MA",
					"name": "Test Setup Updated",
					"description": "Test federated logs setup - updated",
					"active": false,
					"storage": {
						"dataLocationBucket": "test-bucket",
						"database": "test_db",
						"dataIngestConnectionId": "aW5nZXN0LWNvbm4tMTIzNDU2Nzg5MA",
						"queryConnectionId": "cXVlcnktY29ubi0xMjM0NTY3ODkwMA",
						"cloudProviderConfiguration": {
							"provider": "AWS",
							"region": "us-east-1"
						}
					},
					"lifecycleStatus": {
						"status": "COMPLETE",
						"lastUpdatedAt": "2026-05-15T11:00:00Z"
					},
					"createdAt": "2026-05-15T10:00:00Z",
					"updatedAt": "2026-05-15T11:00:00Z"
				}
			}
		}
	}`

	testCreatePartitionResponseJSON = `
	{
		"data": {
			"federatedLogsCreatePartition": {
				"partition": {
					"id": "ZmVkLWxvZ3MtcGFydGl0aW9uLTEyMzQ1Njc4OTA",
					"name": "Test Partition",
					"description": "Test federated logs partition",
					"active": true,
					"isDefault": false,
					"setup": {
						"id": "ZmVkLWxvZ3Mtc2V0dXAtMTIzNDU2Nzg5MA",
						"name": "Test Setup",
						"active": true,
						"storage": {
							"dataLocationBucket": "test-bucket",
							"database": "test_db",
							"dataIngestConnectionId": "aW5nZXN0LWNvbm4tMTIzNDU2Nzg5MA",
							"queryConnectionId": "cXVlcnktY29ubi0xMjM0NTY3ODkwMA",
							"cloudProviderConfiguration": {"provider": "AWS", "region": "us-east-1"}
						},
						"lifecycleStatus": {"status": "COMPLETE", "lastUpdatedAt": "2026-05-15T10:00:00Z"},
						"createdAt": "2026-05-15T10:00:00Z",
						"updatedAt": "2026-05-15T10:00:00Z"
					},
					"storage": {
						"table": "log_transactions",
						"dataLocationUri": "s3://test-bucket/log_transactions"
					},
					"lifecycleStatus": {
						"status": "RESOURCE_CREATION_COMPLETE",
						"lastUpdatedAt": "2026-05-15T10:05:00Z"
					},
					"createdAt": "2026-05-15T10:05:00Z",
					"updatedAt": "2026-05-15T10:05:00Z"
				}
			}
		}
	}`

	testUpdatePartitionResponseJSON = `
	{
		"data": {
			"federatedLogsUpdatePartition": {
				"partition": {
					"id": "ZmVkLWxvZ3MtcGFydGl0aW9uLTEyMzQ1Njc4OTA",
					"name": "Test Partition Updated",
					"description": "Test federated logs partition - updated",
					"active": false,
					"isDefault": false,
					"setup": {
						"id": "ZmVkLWxvZ3Mtc2V0dXAtMTIzNDU2Nzg5MA",
						"name": "Test Setup",
						"active": true,
						"storage": {
							"dataLocationBucket": "test-bucket",
							"database": "test_db",
							"dataIngestConnectionId": "aW5nZXN0LWNvbm4tMTIzNDU2Nzg5MA",
							"queryConnectionId": "cXVlcnktY29ubi0xMjM0NTY3ODkwMA",
							"cloudProviderConfiguration": {"provider": "AWS", "region": "us-east-1"}
						},
						"lifecycleStatus": {"status": "COMPLETE", "lastUpdatedAt": "2026-05-15T10:00:00Z"},
						"createdAt": "2026-05-15T10:00:00Z",
						"updatedAt": "2026-05-15T10:00:00Z"
					},
					"storage": {
						"table": "log_transactions",
						"dataLocationUri": "s3://test-bucket/log_transactions"
					},
					"lifecycleStatus": {
						"status": "COMPLETE",
						"lastUpdatedAt": "2026-05-15T11:05:00Z"
					},
					"createdAt": "2026-05-15T10:05:00Z",
					"updatedAt": "2026-05-15T11:05:00Z"
				}
			}
		}
	}`

	testGetSetupResponseJSON = `
	{
		"data": {
			"actor": {
				"federatedLogs": {
					"setup": {
						"id": "ZmVkLWxvZ3Mtc2V0dXAtMTIzNDU2Nzg5MA",
						"name": "Test Setup",
						"description": "Test federated logs setup",
						"active": true,
						"defaultPartitionId": "ZmVkLWxvZ3MtcGFydGl0aW9uLTEyMzQ1Njc4OTA",
						"storage": {
							"dataLocationBucket": "test-bucket",
							"database": "test_db",
							"dataIngestConnectionId": "aW5nZXN0LWNvbm4tMTIzNDU2Nzg5MA",
							"queryConnectionId": "cXVlcnktY29ubi0xMjM0NTY3ODkwMA",
							"cloudProviderConfiguration": {"provider": "AWS", "region": "us-east-1"}
						},
						"lifecycleStatus": {
							"status": "COMPLETE",
							"lastUpdatedAt": "2026-05-15T11:00:00Z"
						},
						"createdAt": "2026-05-15T10:00:00Z",
						"updatedAt": "2026-05-15T11:00:00Z"
					}
				}
			}
		}
	}`

	testGetPartitionResponseJSON = `
	{
		"data": {
			"actor": {
				"federatedLogs": {
					"partition": {
						"id": "ZmVkLWxvZ3MtcGFydGl0aW9uLTEyMzQ1Njc4OTA",
						"name": "Test Partition",
						"description": "Test federated logs partition",
						"active": true,
						"isDefault": false,
						"setup": {
							"id": "ZmVkLWxvZ3Mtc2V0dXAtMTIzNDU2Nzg5MA",
							"name": "Test Setup",
							"active": true,
							"storage": {
								"dataLocationBucket": "test-bucket",
								"database": "test_db",
								"dataIngestConnectionId": "aW5nZXN0LWNvbm4tMTIzNDU2Nzg5MA",
								"queryConnectionId": "cXVlcnktY29ubi0xMjM0NTY3ODkwMA",
								"cloudProviderConfiguration": {"provider": "AWS", "region": "us-east-1"}
							},
							"lifecycleStatus": {"status": "COMPLETE", "lastUpdatedAt": "2026-05-15T10:00:00Z"},
							"createdAt": "2026-05-15T10:00:00Z",
							"updatedAt": "2026-05-15T10:00:00Z"
						},
						"storage": {
							"table": "log_transactions",
							"dataLocationUri": "s3://test-bucket/log_transactions"
						},
						"lifecycleStatus": {
							"status": "COMPLETE",
							"lastUpdatedAt": "2026-05-15T11:05:00Z"
						},
						"createdAt": "2026-05-15T10:05:00Z",
						"updatedAt": "2026-05-15T11:05:00Z"
					}
				}
			}
		}
	}`

	testCreateAwsConnectionResponseJSON = `
	{
		"data": {
			"entityManagementCreateAwsConnection": {
				"entity": {
					"__typename": "EntityManagementAwsConnectionEntity",
					"id": "ZmVkLWxvZ3MtYXdzLWNvbm4tMTIzNDU2Nzg5MA",
					"name": "Test AWS Connection",
					"description": "Test AWS Connection - New Relic Go Client",
					"enabled": true,
					"externalId": "ext-123",
					"region": "us-east-1",
					"credential": {
						"assumeRole": {
							"roleArn": "arn:aws:iam::123456789012:role/nr-test-role",
							"externalId": "ext-123"
						}
					},
					"metadata": {"version": 1},
					"scope": {"id": "12345", "type": "ACCOUNT"}
				}
			}
		}
	}`

	testDeleteResponseJSON = `
	{
		"data": {
			"entityManagementDelete": {
				"id": "ZmVkLWxvZ3MtYXdzLWNvbm4tMTIzNDU2Nzg5MA"
			}
		}
	}`

	testGetEntityAwsConnectionResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entity": {
						"__typename": "EntityManagementAwsConnectionEntity",
						"id": "ZmVkLWxvZ3MtYXdzLWNvbm4tMTIzNDU2Nzg5MA",
						"name": "Test AWS Connection",
						"description": "Test AWS Connection - New Relic Go Client",
						"enabled": true,
						"externalId": "ext-123",
						"region": "us-east-1",
						"credential": {
							"assumeRole": {
								"roleArn": "arn:aws:iam::123456789012:role/nr-test-role",
								"externalId": "ext-123"
							}
						},
						"metadata": {"version": 1},
						"scope": {"id": "12345", "type": "ACCOUNT"}
					}
				}
			}
		}
	}`

	testGetEntitySearchResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entitySearch": {
						"entities": [
							{
								"__typename": "EntityManagementAwsConnectionEntity",
								"id": "ZmVkLWxvZ3MtYXdzLWNvbm4tMTIzNDU2Nzg5MA",
								"name": "Test AWS Connection",
								"region": "us-east-1",
								"enabled": true
							}
						],
						"nextCursor": ""
					}
				}
			}
		}
	}`
)

func TestUnitFederatedLogs_CreateSetup(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testCreateSetupResponseJSON, http.StatusOK)

	input := FederatedLogsCreateSetupInput{
		Name:        "Test Setup",
		Description: "Test federated logs setup",
		Storage: FederatedLogsSetupStorageInput{
			DataLocationBucket:     "test-bucket",
			Database:               "test_db",
			DataIngestConnectionId: testIngestConnectionID,
			QueryConnectionId:      testQueryConnectionID,
			CloudProviderConfiguration: FederatedLogsCloudProviderConfigurationInput{
				Provider: FederatedLogsCloudProviderTypes.AWS,
				Region:   "us-east-1",
			},
		},
		DefaultPartition: FederatedLogsDefaultPartitionInput{
			Storage: FederatedLogsPartitionStorageInput{
				Table:           "log_transactions",
				DataLocationUri: "s3://test-bucket/log_transactions",
			},
		},
	}

	result, err := client.FederatedLogsCreateSetup(input)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testSetupID, result.Setup.ID)
	require.Equal(t, "Test Setup", result.Setup.Name)
	require.True(t, result.Setup.Active)
	require.Equal(t, testPartitionID, result.Setup.DefaultPartitionId)
	require.Equal(t, "test-bucket", result.Setup.Storage.DataLocationBucket)
	require.Equal(t, FederatedLogsCloudProviderTypes.AWS, result.Setup.Storage.CloudProviderConfiguration.Provider)
}

func TestUnitFederatedLogs_UpdateSetup(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testUpdateSetupResponseJSON, http.StatusOK)

	updateInput := FederatedLogsUpdateSetupInput{
		Name:        "Test Setup Updated",
		Description: "Test federated logs setup - updated",
	}

	result, err := client.FederatedLogsUpdateSetup(testSetupID, updateInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testSetupID, result.Setup.ID)
	require.Equal(t, "Test Setup Updated", result.Setup.Name)
	require.False(t, result.Setup.Active)
	require.Equal(t, FederatedLogsLifecycleStateSetupTypes.COMPLETE, result.Setup.LifecycleStatus.Status)
}

func TestUnitFederatedLogs_CreatePartition(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testCreatePartitionResponseJSON, http.StatusOK)

	input := FederatedLogsCreatePartitionInput{
		Name:        "Test Partition",
		Description: "Test federated logs partition",
		Storage: FederatedLogsPartitionStorageInput{
			Table:           "log_transactions",
			DataLocationUri: "s3://test-bucket/log_transactions",
		},
	}

	result, err := client.FederatedLogsCreatePartition(input, testSetupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testPartitionID, result.Partition.ID)
	require.Equal(t, "Test Partition", result.Partition.Name)
	require.False(t, result.Partition.IsDefault)
	require.Equal(t, "log_transactions", result.Partition.Storage.Table)
	require.Equal(t, testSetupID, result.Partition.Setup.ID)
}

func TestUnitFederatedLogs_UpdatePartition(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testUpdatePartitionResponseJSON, http.StatusOK)

	updateInput := FederatedLogsUpdatePartitionInput{
		Name:        "Test Partition Updated",
		Description: "Test federated logs partition - updated",
	}

	result, err := client.FederatedLogsUpdatePartition(testPartitionID, updateInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testPartitionID, result.Partition.ID)
	require.Equal(t, "Test Partition Updated", result.Partition.Name)
	require.False(t, result.Partition.Active)
	require.Equal(t, FederatedLogsLifecycleStatePartitionTypes.COMPLETE, result.Partition.LifecycleStatus.Status)
}

func TestUnitFederatedLogs_GetSetup(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testGetSetupResponseJSON, http.StatusOK)

	result, err := client.GetSetup(testSetupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testSetupID, result.ID)
	require.Equal(t, "Test Setup", result.Name)
	require.Equal(t, testPartitionID, result.DefaultPartitionId)
}

func TestUnitFederatedLogs_GetPartition(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testGetPartitionResponseJSON, http.StatusOK)

	result, err := client.GetPartition(testPartitionID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testPartitionID, result.ID)
	require.Equal(t, "Test Partition", result.Name)
	require.Equal(t, testSetupID, result.Setup.ID)
}

func TestUnitEntityManagement_CreateAwsConnection(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testCreateAwsConnectionResponseJSON, http.StatusOK)

	input := EntityManagementAwsConnectionEntityCreateInput{
		Name:        "Test AWS Connection",
		Description: "Test AWS Connection - New Relic Go Client",
		Enabled:     true,
		ExternalId:  "ext-123",
		Region:      "us-east-1",
		Credential: EntityManagementAwsCredentialsCreateInput{
			AssumeRole: EntityManagementAwsAssumeRoleConfigCreateInput{
				RoleArn: "arn:aws:iam::123456789012:role/nr-test-role",
			},
		},
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   "12345",
		},
	}

	result, err := client.EntityManagementCreateAwsConnection(input)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testAwsConnectionID, result.Entity.ID)
	require.Equal(t, "Test AWS Connection", result.Entity.Name)
	require.True(t, result.Entity.Enabled)
}

func TestUnitEntityManagement_Delete(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testDeleteResponseJSON, http.StatusOK)

	result, err := client.EntityManagementDelete(testAwsConnectionID, 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testAwsConnectionID, result.ID)
}

func TestUnitEntityManagement_GetEntity_AwsConnection(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testGetEntityAwsConnectionResponseJSON, http.StatusOK)

	result, err := client.GetEntity(testAwsConnectionID)

	require.NoError(t, err)
	require.NotNil(t, result)

	awsEntity, ok := (*result).(*EntityManagementAwsConnectionEntity)
	require.True(t, ok, "Fetched entity was not an EntityManagementAwsConnectionEntity")
	require.Equal(t, testAwsConnectionID, awsEntity.ID)
	require.Equal(t, "Test AWS Connection", awsEntity.Name)
	require.Equal(t, "us-east-1", awsEntity.Region)
}

func TestUnitEntityManagement_GetEntitySearch(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, testGetEntitySearchResponseJSON, http.StatusOK)

	result, err := client.GetEntitySearch("", "")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.Entities)
	awsEntity, ok := result.Entities[0].(*EntityManagementAwsConnectionEntity)
	require.True(t, ok, "Search result entity was not an EntityManagementAwsConnectionEntity")
	require.Equal(t, testAwsConnectionID, awsEntity.ID)
}

func TestUnitFederatedLogs_GetSetup_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors": [{"message": "Not Found"}]}`, http.StatusNotFound)

	result, err := client.GetSetup("non-existent-id")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestUnitFederatedLogs_CreateSetup_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors": [{"message": "Internal Server Error"}]}`, http.StatusInternalServerError)

	result, err := client.FederatedLogsCreateSetup(FederatedLogsCreateSetupInput{Name: "x"})

	require.Error(t, err)
	require.Nil(t, result)
}
