//go:build unit
// +build unit

package pipelinecontrol

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testAccountID = "12345"
	testRuleID    = "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx"
	// testVersion   = 1

	testCreateResponseJSON = `
	{
		"data": {
			"entityManagementCreatePipelineCloudRule": {
				"entity": {
					"id": "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx",
					"metadata": {
						"version": 1
					},
					"name": "Test Rule",
					"description": "Test Pipeline Cloud Rule - New Relic Go Client",
					"nrql": "DELETE FROM Log WHERE (container_name = 'mario')"
				}
			}
		}
	}`

	testUpdateResponseJSON = `
	{
		"data": {
			"entityManagementUpdatePipelineCloudRule": {
				"entity": {
					"id": "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx",
					"metadata": {
						"version": 2
					},
					"name": "Test Rule Updated",
					"description": "Test Pipeline Cloud Rule - New Relic Go Client",
					"nrql": "DELETE FROM Log WHERE (container_name = 'mario')"
				}
			}
		}
	}`

	testGetResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entity": {
						"__typename": "EntityManagementPipelineCloudRuleEntity",
						"id": "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx",
						"metadata": {
							"version": 1
						},
						"name": "Test Rule",
						"description": "Test Pipeline Cloud Rule - New Relic Go Client",
						"nrql": "DELETE FROM Log WHERE (container_name = 'mario')",
						"scope": {
							"id": "12345",
							"type": "ACCOUNT"
						}
					}
				}
			}
		}
	}`

	testDeleteResponseJSON = `
	{
		"data": {
			"entityManagementDelete": {
				"id": "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx"
			}
		}
	}`
)

func TestUnitEntityManagement_CreatePipelineCloudRule(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateResponseJSON, http.StatusOK)

	createInput := EntityManagementPipelineCloudRuleEntityCreateInput{
		Name:        "Test Rule",
		Description: "Test Pipeline Cloud Rule - New Relic Go Client",
		NRQL:        "DELETE FROM Log WHERE (container_name = 'mario')",
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreatePipelineCloudRule(createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testRuleID, result.Entity.ID)
	require.Equal(t, "Test Rule", result.Entity.Name)
}

func TestUnitEntityManagement_UpdatePipelineCloudRule(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testUpdateResponseJSON, http.StatusOK)

	updateInput := EntityManagementPipelineCloudRuleEntityUpdateInput{
		Name:        "Test Rule Updated",
		Description: "Test Pipeline Cloud Rule - New Relic Go Client Updated",
		NRQL:        "DELETE FROM Log WHERE (container_name = 'shrimp')",
	}

	result, err := client.EntityManagementUpdatePipelineCloudRule(
		testRuleID,
		updateInput,
		//testVersion+1,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testRuleID, result.Entity.ID)
	require.Equal(t, "Test Rule Updated", result.Entity.Name)
}

func TestUnitEntityManagement_GetEntity(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testGetResponseJSON, http.StatusOK)

	result, err := client.GetEntity(testRuleID)

	require.NoError(t, err)
	require.NotNil(t, result)

	ruleEntity, ok := (*result).(*EntityManagementPipelineCloudRuleEntity)
	require.True(t, ok)
	require.Equal(t, testRuleID, ruleEntity.ID)
	require.Equal(t, "Test Rule", ruleEntity.Name)
}

func TestUnitEntityManagement_DeleteEntity(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testDeleteResponseJSON, http.StatusOK)

	result, err := client.EntityManagementDelete(
		testRuleID,
		//testVersion,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testRuleID, result.ID)
}

func TestUnitEntityManagement_DeleteEntityWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testDeleteResponseJSON, http.StatusOK)

	result, err := client.EntityManagementDeleteWithContext(
		context.Background(),
		testRuleID,
		//testVersion,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testRuleID, result.ID)
}

func TestUnitEntityManagement_GetEntity_Error(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, `{"errors": [{"message": "Not Found"}]}`, http.StatusNotFound)

	result, err := client.GetEntity("non-existent-id")

	require.Error(t, err)
	require.Nil(t, result)
}

var (
	testFederatedLogSetupID = "MXxYY0YYxYYYYYYYYY6YY0YYY1YYY1YYYYY9YYY6YYYxY3YxYyyyYY04YYY3YYy0YYYyYYYyYyYyYYYyYYYy"

	testCreateFederatedLogSetupResponseJSON = `
	{
		"data": {
			"entityManagementCreateFederatedLogSetup": {
				"entity": {
					"id": "MXxYY0YYxYYYYYYYYY6YY0YYY1YYY1YYYYY9YYY6YYYxY3YxYyyyYY04YYY3YYy0YYYyYYYyYyYyYYYyYYYy",
					"metadata": {
						"version": 1
					},
					"name": "Test Federated Log Setup",
					"description": "Test Federated Log Setup - New Relic Go Client",
					"cloudProvider": "AWS",
					"cloudProviderRegion": "us-east-1",
					"dataLocationBucket": "my-test-bucket",
					"nrAccountId": "12345",
					"nrRegion": "US",
					"status": "ACTIVE"
				}
			}
		}
	}`

	testCreateFederatedLogSetupErrorResponseJSON = `
	{
		"errors": [
			{
				"message": "Invalid input: missing required field"
			}
		]
	}`
)

func TestUnitEntityManagement_CreateFederatedLogSetup(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogSetupResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogSetupEntityCreateInput{
		Name:                       "Test Federated Log Setup",
		Description:                "Test Federated Log Setup - New Relic Go Client",
		CloudProvider:              EntityManagementCloudProviderTypes.AWS,
		CloudProviderRegion:        "us-east-1",
		DataLocationBucket:         "my-test-bucket",
		DataProcessingConnectionId: "connection-123",
		NrAccountId:                "12345",
		NrRegion:                   EntityManagementNrRegionTypes.US,
		QueryConnectionId:          "query-connection-456",
		Status:                     EntityManagementFederatedLogSetupStatusTypes.ACTIVE,
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreateFederatedLogSetup(createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testFederatedLogSetupID, result.Entity.ID)
	require.Equal(t, "Test Federated Log Setup", result.Entity.Name)
	require.Equal(t, EntityManagementCloudProviderTypes.AWS, result.Entity.CloudProvider)
	require.Equal(t, "us-east-1", result.Entity.CloudProviderRegion)
	require.Equal(t, EntityManagementNrRegionTypes.US, result.Entity.NrRegion)
	require.Equal(t, EntityManagementFederatedLogSetupStatusTypes.ACTIVE, result.Entity.Status)
}

func TestUnitEntityManagement_CreateFederatedLogSetupWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogSetupResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogSetupEntityCreateInput{
		Name:                       "Test Federated Log Setup",
		Description:                "Test Federated Log Setup - New Relic Go Client",
		CloudProvider:              EntityManagementCloudProviderTypes.AWS,
		CloudProviderRegion:        "us-east-1",
		DataLocationBucket:         "my-test-bucket",
		DataProcessingConnectionId: "connection-123",
		NrAccountId:                "12345",
		NrRegion:                   EntityManagementNrRegionTypes.US,
		QueryConnectionId:          "query-connection-456",
		Status:                     EntityManagementFederatedLogSetupStatusTypes.ACTIVE,
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreateFederatedLogSetupWithContext(context.Background(), createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testFederatedLogSetupID, result.Entity.ID)
	require.Equal(t, "Test Federated Log Setup", result.Entity.Name)
}

func TestUnitEntityManagement_CreateFederatedLogSetup_Error(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogSetupErrorResponseJSON, http.StatusBadRequest)

	createInput := EntityManagementFederatedLogSetupEntityCreateInput{
		Name:          "Test Federated Log Setup",
		CloudProvider: EntityManagementCloudProviderTypes.AWS,
	}

	result, err := client.EntityManagementCreateFederatedLogSetup(createInput)

	require.Error(t, err)
	require.Nil(t, result)
}

var (
	testAwsConnectionID = "MXxZZ0ZZxZZZZZZZZZ7ZZ0ZZZ1ZZZ1ZZZZZ0ZZZ7ZZZxZ4ZxZzzzZZ05ZZZ4ZZz0ZZZzZZZzZzZzZZZzZZZz"

	testCreateAwsConnectionResponseJSON = `
	{
		"data": {
			"EntityManagementCreateAwsConnection": {
				"entity": {
					"id": "MXxZZ0ZZxZZZZZZZZZ7ZZ0ZZZ1ZZZ1ZZZZZ0ZZZ7ZZZxZ4ZxZzzzZZ05ZZZ4ZZz0ZZZzZZZzZzZzZZZzZZZz",
					"metadata": {
						"version": 1
					},
					"name": "Test AWS Connection",
					"description": "Test AWS Connection - New Relic Go Client",
					"enabled": true,
					"region": "us-east-1",
					"externalId": "ext-123",
					"credential": {
						"assumeRole": {
							"roleArn": "arn:aws:iam::123456789012:role/test-role",
							"externalId": "ext-456"
						}
					},
					"scope": {
						"id": "12345",
						"type": "ACCOUNT"
					},
					"settings": [
						{
							"key": "setting1",
							"value": "value1"
						}
					],
					"tags": [],
					"type": "AWS_CONNECTION"
				}
			}
		}
	}`

	testCreateAwsConnectionErrorResponseJSON = `
	{
		"errors": [
			{
				"message": "Invalid input: missing required field"
			}
		]
	}`

	testGetEntitySearchResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entitySearch": {
						"entities": [
							{
								"__typename": "EntityManagementPipelineCloudRuleEntity",
								"id": "MXxXX0XXxXXXXXXXXX5XX0XXX1XXX1XXXXX8XXX5XXXxX2XxXxxxXX03XXX2XXx0XXXxXXXxXxXxXXXxXXXx",
								"metadata": {
									"version": 1
								},
								"name": "Test Rule",
								"description": "Test Pipeline Cloud Rule",
								"nrql": "DELETE FROM Log WHERE (container_name = 'mario')",
								"scope": {
									"id": "12345",
									"type": "ACCOUNT"
								},
								"tags": []
							}
						],
						"nextCursor": ""
					}
				}
			}
		}
	}`

	testGetEntitySearchEmptyResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entitySearch": {
						"entities": [],
						"nextCursor": ""
					}
				}
			}
		}
	}`

	testGetEntityFederatedLogSetupResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entity": {
						"__typename": "EntityManagementFederatedLogSetupEntity",
						"id": "MXxYY0YYxYYYYYYYYY6YY0YYY1YYY1YYYYY9YYY6YYYxY3YxYyyyYY04YYY3YYy0YYYyYYYyYyYyYYYyYYYy",
						"metadata": {
							"version": 1
						},
						"name": "Test Federated Log Setup",
						"description": "Test Federated Log Setup",
						"cloudProvider": "AWS",
						"cloudProviderRegion": "us-east-1",
						"dataLocationBucket": "my-test-bucket",
						"nrAccountId": "12345",
						"nrRegion": "US",
						"status": "ACTIVE",
						"scope": {
							"id": "12345",
							"type": "ACCOUNT"
						},
						"tags": []
					}
				}
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
						"id": "MXxZZ0ZZxZZZZZZZZZ7ZZ0ZZZ1ZZZ1ZZZZZ0ZZZ7ZZZxZ4ZxZzzzZZ05ZZZ4ZZz0ZZZzZZZzZzZzZZZzZZZz",
						"metadata": {
							"version": 1
						},
						"name": "Test AWS Connection",
						"description": "Test AWS Connection",
						"enabled": true,
						"region": "us-east-1",
						"scope": {
							"id": "12345",
							"type": "ACCOUNT"
						},
						"tags": []
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
		Description: "Test AWS Connection - New Relic Go Client",
		Enabled:     true,
		Region:      "us-east-1",
		ExternalId:  "ext-123",
		Credential: EntityManagementAwsCredentialsCreateInput{
			AssumeRole: EntityManagementAwsAssumeRoleConfigCreateInput{
				RoleArn: "arn:aws:iam::123456789012:role/test-role",
			},
		},
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
	require.Equal(t, "Test AWS Connection - New Relic Go Client", result.Entity.Description)
	require.True(t, result.Entity.Enabled)
	require.Equal(t, "us-east-1", result.Entity.Region)
}

func TestUnitEntityManagement_CreateAwsConnectionWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateAwsConnectionResponseJSON, http.StatusOK)

	createInput := EntityManagementAwsConnectionEntityCreateInput{
		Name:        "Test AWS Connection",
		Description: "Test AWS Connection - New Relic Go Client",
		Enabled:     true,
		Region:      "us-east-1",
		Credential: EntityManagementAwsCredentialsCreateInput{
			AssumeRole: EntityManagementAwsAssumeRoleConfigCreateInput{
				RoleArn: "arn:aws:iam::123456789012:role/test-role",
			},
		},
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreateAwsConnectionWithContext(context.Background(), createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testAwsConnectionID, result.Entity.ID)
	require.Equal(t, "Test AWS Connection", result.Entity.Name)
}

func TestUnitEntityManagement_CreateAwsConnection_Error(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateAwsConnectionErrorResponseJSON, http.StatusBadRequest)

	createInput := EntityManagementAwsConnectionEntityCreateInput{
		Name: "Test AWS Connection",
	}

	result, err := client.EntityManagementCreateAwsConnection(createInput)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestUnitEntityManagement_GetEntitySearch(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testGetEntitySearchResponseJSON, http.StatusOK)

	result, err := client.GetEntitySearch("", "type = 'PIPELINE_CLOUD_RULE'")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Entities, 1)

	ruleEntity, ok := result.Entities[0].(*EntityManagementPipelineCloudRuleEntity)
	require.True(t, ok)
	require.Equal(t, testRuleID, ruleEntity.ID)
}

func TestUnitEntityManagement_GetEntitySearchWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testGetEntitySearchResponseJSON, http.StatusOK)

	result, err := client.GetEntitySearchWithContext(context.Background(), "", "type = 'PIPELINE_CLOUD_RULE'")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Entities, 1)
}

func TestUnitEntityManagement_GetEntitySearch_Empty(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testGetEntitySearchEmptyResponseJSON, http.StatusOK)

	result, err := client.GetEntitySearch("", "type = 'NONEXISTENT'")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Entities, 0)
}

func TestUnitEntityManagement_GetEntitySearch_Error(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, `{"errors": [{"message": "Invalid query"}]}`, http.StatusBadRequest)

	result, err := client.GetEntitySearch("", "invalid query")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestUnitEntityManagement_GetEntity_FederatedLogSetup(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testGetEntityFederatedLogSetupResponseJSON, http.StatusOK)

	result, err := client.GetEntity(testFederatedLogSetupID)

	require.NoError(t, err)
	require.NotNil(t, result)

	setupEntity, ok := (*result).(*EntityManagementFederatedLogSetupEntity)
	require.True(t, ok, "Fetched entity was not of the expected FederatedLogSetupEntity type")
	require.Equal(t, testFederatedLogSetupID, setupEntity.ID)
	require.Equal(t, "Test Federated Log Setup", setupEntity.Name)
	require.Equal(t, EntityManagementCloudProviderTypes.AWS, setupEntity.CloudProvider)
}

func TestUnitEntityManagement_GetEntity_AwsConnection(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testGetEntityAwsConnectionResponseJSON, http.StatusOK)

	result, err := client.GetEntity(testAwsConnectionID)

	require.NoError(t, err)
	require.NotNil(t, result)

	connEntity, ok := (*result).(*EntityManagementAwsConnectionEntity)
	require.True(t, ok, "Fetched entity was not of the expected AwsConnectionEntity type")
	require.Equal(t, testAwsConnectionID, connEntity.ID)
	require.Equal(t, "Test AWS Connection", connEntity.Name)
	require.True(t, connEntity.Enabled)
}

var (
	testFederatedLogPartitionID = "MXxWW0WWxWWWWWWWWW8WW0WWW1WWW1WWWWW2WWW8WWWxW5WxWwwwWW06WWW5WWw0WWWwWWWwWwWwWWWwWWWw"

	testCreateFederatedLogPartitionResponseJSON = `
	{
		"data": {
			"EntityManagementCreateFederatedLogPartition": {
				"entity": {
					"id": "MXxWW0WWxWWWWWWWWW8WW0WWW1WWW1WWWWW2WWW8WWWxW5WxWwwwWW06WWW5WWw0WWWwWWWwWwWwWWWwWWWw",
					"metadata": {
						"version": 1
					},
					"name": "Test Partition",
					"description": "Test Federated Log Partition - New Relic Go Client",
					"dataLocationUri": "s3://my-bucket/logs/partition1",
					"isDefault": false,
					"nrAccountId": "12345",
					"partitionDatabase": "my_database",
					"partitionTable": "my_table",
					"status": "ACTIVE",
					"retentionPolicy": {
						"duration": 30,
						"unit": "DAYS"
					},
					"scope": {
						"id": "12345",
						"type": "ACCOUNT"
					},
					"tags": [],
					"type": "FEDERATED_LOG_PARTITION"
				}
			}
		}
	}`

	testCreateFederatedLogPartitionErrorResponseJSON = `
	{
		"errors": [
			{
				"message": "Invalid input: missing required field"
			}
		]
	}`
)

func TestUnitEntityManagement_CreateFederatedLogPartition(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogPartitionResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogPartitionEntityCreateInput{
		Name:              "Test Partition",
		Description:       "Test Federated Log Partition - New Relic Go Client",
		DataLocationUri:   "s3://my-bucket/logs/partition1",
		IsDefault:         false,
		NrAccountId:       "12345",
		PartitionDatabase: "my_database",
		PartitionTable:    "my_table",
		SetupId:           "setup-789",
		Status:            EntityManagementLogPartitionStatusTypes.ACTIVE,
		RetentionPolicy: &EntityManagementRetentionPolicyCreateInput{
			Duration: 30,
			Unit:     EntityManagementRetentionUnitTypes.DAYS,
		},
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreateFederatedLogPartition(createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testFederatedLogPartitionID, result.Entity.ID)
	require.Equal(t, "Test Partition", result.Entity.Name)
	require.Equal(t, "s3://my-bucket/logs/partition1", result.Entity.DataLocationUri)
	require.False(t, result.Entity.IsDefault)
	require.Equal(t, "12345", result.Entity.NrAccountId)
	require.Equal(t, "my_database", result.Entity.PartitionDatabase)
	require.Equal(t, "my_table", result.Entity.PartitionTable)
	require.Equal(t, EntityManagementLogPartitionStatusTypes.ACTIVE, result.Entity.Status)
	require.Equal(t, 30, result.Entity.RetentionPolicy.Duration)
	require.Equal(t, EntityManagementRetentionUnitTypes.DAYS, result.Entity.RetentionPolicy.Unit)
}

func TestUnitEntityManagement_CreateFederatedLogPartitionWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogPartitionResponseJSON, http.StatusOK)

	createInput := EntityManagementFederatedLogPartitionEntityCreateInput{
		Name:              "Test Partition",
		Description:       "Test Federated Log Partition - New Relic Go Client",
		DataLocationUri:   "s3://my-bucket/logs/partition1",
		IsDefault:         false,
		NrAccountId:       "12345",
		PartitionDatabase: "my_database",
		PartitionTable:    "my_table",
		SetupId:           "setup-789",
		Status:            EntityManagementLogPartitionStatusTypes.ACTIVE,
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   testAccountID,
		},
	}

	result, err := client.EntityManagementCreateFederatedLogPartitionWithContext(context.Background(), createInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testFederatedLogPartitionID, result.Entity.ID)
	require.Equal(t, "Test Partition", result.Entity.Name)
}

func TestUnitEntityManagement_CreateFederatedLogPartition_Error(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateFederatedLogPartitionErrorResponseJSON, http.StatusBadRequest)

	createInput := EntityManagementFederatedLogPartitionEntityCreateInput{
		Name: "Test Partition",
	}

	result, err := client.EntityManagementCreateFederatedLogPartition(createInput)

	require.Error(t, err)
	require.Nil(t, result)
}
