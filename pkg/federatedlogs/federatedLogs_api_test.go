//go:build unit
// +build unit

package federatedlogs

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testEntityID      = "test-entity-id-123"
	testSetupEntityID = "test-setup-entity-id-456"
	testPartitionID   = "test-partition-id-789"
	testAccountScopeID = "12345678"
	testEntityVersion  = 1

	testCreateAwsConnectionResponseJSON = `{
		"entityManagementCreateAwsConnection": {
			"entity": {
				"id": "test-entity-id-123",
				"name": "test-aws-connection",
				"type": "AWSCONNECTIONENTITY",
				"description": "Test AWS connection",
				"enabled": true,
				"region": "us-east-1",
				"credential": {
					"assumeRole": {
						"externalId": "test-external-id",
						"roleArn": "arn:aws:iam::123456789012:role/test-role"
					}
				},
				"scope": {
					"id": "12345678",
					"type": "ACCOUNT"
				},
				"tags": [],
				"settings": [],
				"metadata": {
					"version": 1
				}
			}
		}
	}`

	testCreateFederatedLogPartitionResponseJSON = `{
		"entityManagementCreateFederatedLogPartition": {
			"entity": {
				"id": "test-partition-id-789",
				"name": "test-partition",
				"type": "FEDERATEDLOGPARTITIONENTITY",
				"description": "Test log partition",
				"dataLocationUri": "s3://test-bucket/logs/partition",
				"isDefault": false,
				"partitionTable": "test_partition_table",
				"status": "ACTIVE",
				"scope": {
					"id": "12345678",
					"type": "ACCOUNT"
				},
				"tags": [],
				"metadata": {
					"version": 1
				}
			}
		}
	}`

	testCreateFederatedLogSetupResponseJSON = `{
		"entityManagementCreateFederatedLogSetup": {
			"entity": {
				"id": "test-setup-entity-id-456",
				"name": "test-federated-log-setup",
				"type": "FEDERATEDLOGSETUPENTITY",
				"description": "Test federated log setup",
				"cloudProvider": "AWS",
				"cloudProviderRegion": "us-east-1",
				"dataLocationBucket": "test-log-bucket",
				"partitionDatabase": "test_partition_db",
				"status": "ACTIVE",
				"scope": {
					"id": "12345678",
					"type": "ACCOUNT"
				},
				"tags": [],
				"metadata": {
					"version": 1
				}
			}
		}
	}`

	testUpdateFederatedLogPartitionResponseJSON = `{
		"entityManagementUpdateFederatedLogPartition": {
			"entity": {
				"id": "test-partition-id-789",
				"name": "updated-partition",
				"type": "FEDERATEDLOGPARTITIONENTITY",
				"description": "Updated log partition",
				"dataLocationUri": "s3://test-bucket/logs/updated-partition",
				"isDefault": true,
				"partitionTable": "updated_partition_table",
				"status": "INACTIVE",
				"scope": {
					"id": "12345678",
					"type": "ACCOUNT"
				},
				"tags": [],
				"metadata": {
					"version": 2
				}
			}
		}
	}`

	testUpdateFederatedLogSetupResponseJSON = `{
		"entityManagementUpdateFederatedLogSetup": {
			"entity": {
				"id": "test-setup-entity-id-456",
				"name": "updated-federated-log-setup",
				"type": "FEDERATEDLOGSETUPENTITY",
				"description": "Updated federated log setup",
				"cloudProvider": "AWS",
				"cloudProviderRegion": "us-west-2",
				"dataLocationBucket": "updated-log-bucket",
				"partitionDatabase": "updated_partition_db",
				"status": "INACTIVE",
				"scope": {
					"id": "12345678",
					"type": "ACCOUNT"
				},
				"tags": [],
				"metadata": {
					"version": 2
				}
			}
		}
	}`

	testGetFederatedLogSetupEntityResponseJSON = `{
		"actor": {
			"entityManagement": {
				"entity": {
					"__typename": "EntityManagementFederatedLogSetupEntity",
					"id": "test-setup-entity-id-456",
					"name": "test-federated-log-setup",
					"type": "FEDERATEDLOGSETUPENTITY",
					"description": "Test federated log setup",
					"cloudProvider": "AWS",
					"cloudProviderRegion": "us-east-1",
					"dataLocationBucket": "test-log-bucket",
					"partitionDatabase": "test_partition_db",
					"status": "ACTIVE",
					"scope": {
						"id": "12345678",
						"type": "ACCOUNT"
					},
					"tags": [],
					"metadata": {
						"version": 1
					}
				}
			}
		}
	}`

	testGetFederatedLogPartitionEntityResponseJSON = `{
		"actor": {
			"entityManagement": {
				"entity": {
					"id": "test-partition-id-789",
					"name": "test-partition",
					"type": "FEDERATEDLOGPARTITIONENTITY",
					"description": "Test log partition",
					"dataLocationUri": "s3://test-bucket/logs/partition",
					"isDefault": false,
					"partitionTable": "test_partition_table",
					"status": "ACTIVE",
					"scope": {
						"id": "12345678",
						"type": "ACCOUNT"
					},
					"tags": [],
					"metadata": {
						"version": 1
					}
				}
			}
		}
	}`

	testGetEntitySearchResponseJSON = `{
		"actor": {
			"entityManagement": {
				"entitySearch": {
					"entities": [
						{
							"__typename": "EntityManagementFederatedLogSetupEntity",
							"id": "test-setup-entity-id-456",
							"name": "test-federated-log-setup",
							"type": "FEDERATEDLOGSETUPENTITY",
							"cloudProvider": "AWS",
							"cloudProviderRegion": "us-east-1",
							"dataLocationBucket": "test-log-bucket",
							"partitionDatabase": "test_partition_db",
							"status": "ACTIVE",
							"scope": {
								"id": "12345678",
								"type": "ACCOUNT"
							},
							"tags": [],
							"metadata": {
								"version": 1
							}
						}
					],
					"nextCursor": ""
				}
			}
		}
	}`
)

func TestEntityManagementCreateAwsConnection(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateAwsConnectionResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusCreated)

	input := EntityManagementAwsConnectionEntityCreateInput{
		Name:        "test-aws-connection",
		Description: "Test AWS connection",
		Enabled:     true,
		Region:      "us-east-1",
		Credential: EntityManagementAwsCredentialsCreateInput{
			AssumeRole: EntityManagementAwsAssumeRoleConfigCreateInput{
				RoleArn: "arn:aws:iam::123456789012:role/test-role",
			},
		},
		Scope: EntityManagementScopedReferenceInput{
			ID:   testAccountScopeID,
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
		},
	}

	actual, err := client.EntityManagementCreateAwsConnection(input)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testEntityID, actual.Entity.ID)
	assert.Equal(t, "test-aws-connection", actual.Entity.Name)
	assert.Equal(t, "us-east-1", actual.Entity.Region)
	assert.Equal(t, true, actual.Entity.Enabled)
	assert.Equal(t, EntityManagementDynamicString("test-external-id"), actual.Entity.Credential.AssumeRole.ExternalId)
	assert.Equal(t, "arn:aws:iam::123456789012:role/test-role", actual.Entity.Credential.AssumeRole.RoleArn)
	assert.Equal(t, testAccountScopeID, actual.Entity.Scope.ID)
	assert.Equal(t, EntityManagementEntityScopeTypes.ACCOUNT, actual.Entity.Scope.Type)
}

func TestEntityManagementCreateAwsConnection_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"unauthorized"}]}`, http.StatusUnauthorized)

	input := EntityManagementAwsConnectionEntityCreateInput{
		Name: "test-aws-connection",
	}

	actual, err := client.EntityManagementCreateAwsConnection(input)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestEntityManagementCreateFederatedLogPartition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateFederatedLogPartitionResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusCreated)

	input := EntityManagementFederatedLogPartitionEntityCreateInput{
		Name:            "test-partition",
		Description:     "Test log partition",
		DataLocationUri: "s3://test-bucket/logs/partition",
		IsDefault:       false,
		PartitionTable:  "test_partition_table",
		SetupId:         testSetupEntityID,
		Status:          EntityManagementLogPartitionStatusTypes.ACTIVE,
		Scope: EntityManagementScopedReferenceInput{
			ID:   testAccountScopeID,
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
		},
	}

	actual, err := client.EntityManagementCreateFederatedLogPartition(input)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testPartitionID, actual.Entity.ID)
	assert.Equal(t, "test-partition", actual.Entity.Name)
	assert.Equal(t, "s3://test-bucket/logs/partition", actual.Entity.DataLocationUri)
	assert.Equal(t, false, actual.Entity.IsDefault)
	assert.Equal(t, "test_partition_table", actual.Entity.PartitionTable)
	assert.Equal(t, EntityManagementLogPartitionStatusTypes.ACTIVE, actual.Entity.Status)
	assert.Equal(t, testAccountScopeID, actual.Entity.Scope.ID)
}

func TestEntityManagementCreateFederatedLogPartition_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"internal server error"}]}`, http.StatusInternalServerError)

	input := EntityManagementFederatedLogPartitionEntityCreateInput{
		Name:    "test-partition",
		SetupId: testSetupEntityID,
		Status:  EntityManagementLogPartitionStatusTypes.ACTIVE,
	}

	actual, err := client.EntityManagementCreateFederatedLogPartition(input)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestEntityManagementCreateFederatedLogSetup(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateFederatedLogSetupResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusCreated)

	input := EntityManagementFederatedLogSetupEntityCreateInput{
		Name:                "test-federated-log-setup",
		Description:         "Test federated log setup",
		CloudProvider:       EntityManagementCloudProviderTypes.AWS,
		CloudProviderRegion: "us-east-1",
		DataLocationBucket:  "test-log-bucket",
		PartitionDatabase:   "test_partition_db",
	}

	actual, err := client.EntityManagementCreateFederatedLogSetup(input)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testSetupEntityID, actual.Entity.ID)
	assert.Equal(t, "test-federated-log-setup", actual.Entity.Name)
	assert.Equal(t, EntityManagementCloudProviderTypes.AWS, actual.Entity.CloudProvider)
	assert.Equal(t, "us-east-1", actual.Entity.CloudProviderRegion)
	assert.Equal(t, "test-log-bucket", actual.Entity.DataLocationBucket)
	assert.Equal(t, "test_partition_db", actual.Entity.PartitionDatabase)
	assert.Equal(t, EntityManagementFederatedLogSetupStatusTypes.ACTIVE, actual.Entity.Status)
}

func TestEntityManagementCreateFederatedLogSetup_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"bad request"}]}`, http.StatusBadRequest)

	input := EntityManagementFederatedLogSetupEntityCreateInput{
		Name:          "test-federated-log-setup",
		CloudProvider: EntityManagementCloudProviderTypes.AWS,
	}

	actual, err := client.EntityManagementCreateFederatedLogSetup(input)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestEntityManagementUpdateFederatedLogPartition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateFederatedLogPartitionResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusOK)

	input := EntityManagementFederatedLogPartitionEntityUpdateInput{
		Name:            "updated-partition",
		Description:     "Updated log partition",
		DataLocationUri: "s3://test-bucket/logs/updated-partition",
		IsDefault:       true,
		PartitionTable:  "updated_partition_table",
		Status:          EntityManagementLogPartitionStatusTypes.INACTIVE,
	}

	actual, err := client.EntityManagementUpdateFederatedLogPartition(input, testPartitionID, testEntityVersion)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testPartitionID, actual.Entity.ID)
	assert.Equal(t, "updated-partition", actual.Entity.Name)
	assert.Equal(t, "s3://test-bucket/logs/updated-partition", actual.Entity.DataLocationUri)
	assert.Equal(t, true, actual.Entity.IsDefault)
	assert.Equal(t, "updated_partition_table", actual.Entity.PartitionTable)
	assert.Equal(t, EntityManagementLogPartitionStatusTypes.INACTIVE, actual.Entity.Status)
	assert.Equal(t, testAccountScopeID, actual.Entity.Scope.ID)
}

func TestEntityManagementUpdateFederatedLogPartition_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"not found"}]}`, http.StatusNotFound)

	input := EntityManagementFederatedLogPartitionEntityUpdateInput{
		Name: "updated-partition",
	}

	actual, err := client.EntityManagementUpdateFederatedLogPartition(input, "nonexistent-id", 1)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestEntityManagementUpdateFederatedLogSetup(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateFederatedLogSetupResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusOK)

	input := EntityManagementFederatedLogSetupEntityUpdateInput{
		Name:                "updated-federated-log-setup",
		Description:         "Updated federated log setup",
		CloudProvider:       EntityManagementCloudProviderTypes.AWS,
		CloudProviderRegion: "us-west-2",
		DataLocationBucket:  "updated-log-bucket",
		PartitionDatabase:   "updated_partition_db",
		Status:              EntityManagementFederatedLogSetupStatusTypes.INACTIVE,
	}

	actual, err := client.EntityManagementUpdateFederatedLogSetup(input, testSetupEntityID, testEntityVersion)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testSetupEntityID, actual.Entity.ID)
	assert.Equal(t, "updated-federated-log-setup", actual.Entity.Name)
	assert.Equal(t, "us-west-2", actual.Entity.CloudProviderRegion)
	assert.Equal(t, "updated-log-bucket", actual.Entity.DataLocationBucket)
	assert.Equal(t, "updated_partition_db", actual.Entity.PartitionDatabase)
	assert.Equal(t, EntityManagementFederatedLogSetupStatusTypes.INACTIVE, actual.Entity.Status)
}

func TestEntityManagementUpdateFederatedLogSetup_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"conflict"}]}`, http.StatusConflict)

	input := EntityManagementFederatedLogSetupEntityUpdateInput{
		Name: "updated-federated-log-setup",
	}

	actual, err := client.EntityManagementUpdateFederatedLogSetup(input, "nonexistent-id", 1)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestGetFederatedLogSetupEntityWithContext(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetFederatedLogSetupEntityResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetFederatedLogSetupEntityWithContext(t.Context(), testSetupEntityID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testSetupEntityID, actual.ID)
	assert.Equal(t, "test-federated-log-setup", actual.Name)
	assert.Equal(t, EntityManagementCloudProviderTypes.AWS, actual.CloudProvider)
	assert.Equal(t, "us-east-1", actual.CloudProviderRegion)
	assert.Equal(t, "test-log-bucket", actual.DataLocationBucket)
	assert.Equal(t, "test_partition_db", actual.PartitionDatabase)
	assert.Equal(t, EntityManagementFederatedLogSetupStatusTypes.ACTIVE, actual.Status)
	assert.Equal(t, testAccountScopeID, actual.Scope.ID)
}

func TestGetFederatedLogSetupEntityWithContext_NotFound(t *testing.T) {
	t.Parallel()
	respJSON := `{ "data": { "actor": { "entityManagement": { "entity": null } } } }`
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetFederatedLogSetupEntityWithContext(t.Context(), "nonexistent-id")

	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestGetFederatedLogSetupEntityWithContext_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"internal error"}]}`, http.StatusInternalServerError)

	actual, err := client.GetFederatedLogSetupEntityWithContext(t.Context(), testSetupEntityID)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestGetFederatedLogPartitionEntityWithContext(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetFederatedLogPartitionEntityResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetFederatedLogPartitionEntityWithContext(t.Context(), testPartitionID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testPartitionID, actual.ID)
	assert.Equal(t, "test-partition", actual.Name)
	assert.Equal(t, "s3://test-bucket/logs/partition", actual.DataLocationUri)
	assert.Equal(t, false, actual.IsDefault)
	assert.Equal(t, "test_partition_table", actual.PartitionTable)
	assert.Equal(t, EntityManagementLogPartitionStatusTypes.ACTIVE, actual.Status)
	assert.Equal(t, testAccountScopeID, actual.Scope.ID)
}

func TestGetFederatedLogPartitionEntityWithContext_NotFound(t *testing.T) {
	t.Parallel()
	respJSON := `{ "data": { "actor": { "entityManagement": { "entity": null } } } }`
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetFederatedLogPartitionEntityWithContext(t.Context(), "nonexistent-id")

	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestGetFederatedLogPartitionEntityWithContext_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"internal error"}]}`, http.StatusInternalServerError)

	actual, err := client.GetFederatedLogPartitionEntityWithContext(t.Context(), testPartitionID)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestGetEntitySearch(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetEntitySearchResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetEntitySearch("", `type = 'FEDERATEDLOGSETUPENTITY'`)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Len(t, actual.Entities, 1)
	assert.Equal(t, "", actual.NextCursor)

	setup, ok := actual.Entities[0].(*EntityManagementFederatedLogSetupEntity)
	require.True(t, ok, "expected entity to be *EntityManagementFederatedLogSetupEntity")
	assert.Equal(t, testSetupEntityID, setup.ID)
	assert.Equal(t, "test-federated-log-setup", setup.Name)
}

func TestGetEntitySearch_EmptyResult(t *testing.T) {
	t.Parallel()
	respJSON := `{ "data": { "actor": { "entityManagement": { "entitySearch": { "entities": [], "nextCursor": "" } } } } }`
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetEntitySearch("", `type = 'FEDERATEDLOGSETUPENTITY'`)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Len(t, actual.Entities, 0)
}

func TestGetEntitySearch_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"query error"}]}`, http.StatusBadRequest)

	actual, err := client.GetEntitySearch("", "invalid query")

	assert.Error(t, err)
	assert.Nil(t, actual)
}

var (
	testGetEntityResponseJSON = `{
		"actor": {
			"entityManagement": {
				"entity": {
					"__typename": "EntityManagementFederatedLogSetupEntity",
					"id": "test-setup-entity-id-456",
					"name": "test-federated-log-setup",
					"type": "FEDERATEDLOGSETUPENTITY",
					"scope": {
						"id": "12345678",
						"type": "ACCOUNT"
					},
					"tags": [],
					"metadata": {
						"version": 1
					}
				}
			}
		}
	}`

	testGetAwsConnectionEntityResponseJSON = `{
		"actor": {
			"entityManagement": {
				"entity": {
					"id": "test-entity-id-123",
					"name": "test-aws-connection",
					"type": "AWSCONNECTIONENTITY",
					"description": "Test AWS connection",
					"enabled": true,
					"region": "us-east-1",
					"credential": {
						"assumeRole": {
							"externalId": "test-external-id",
							"roleArn": "arn:aws:iam::123456789012:role/test-role"
						}
					},
					"scope": {
						"id": "12345678",
						"type": "ACCOUNT"
					},
					"tags": [],
					"settings": [],
					"metadata": {
						"version": 1
					}
				}
			}
		}
	}`
)

func TestGetEntity(t *testing.T) {
	t.Parallel()
	// EntityManagementEntityInterface has no custom unmarshaler on ActorStitchedFields,
	// so the entity field can only be populated as a nil interface via this path.
	// Concrete entity retrieval uses GetFederatedLogSetupEntityWithContext / GetFederatedLogPartitionEntityWithContext.
	respJSON := `{ "data": { "actor": { "entityManagement": { "entity": null } } } }`
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetEntity(testSetupEntityID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Nil(t, *actual)
}

func TestGetEntityWithContext(t *testing.T) {
	t.Parallel()
	respJSON := `{ "data": { "actor": { "entityManagement": { "entity": null } } } }`
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetEntityWithContext(context.Background(), testSetupEntityID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Nil(t, *actual)
}

func TestGetEntity_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"not found"}]}`, http.StatusNotFound)

	actual, err := client.GetEntity("nonexistent-id")

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestGetAwsConnectionEntityWithContext(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetAwsConnectionEntityResponseJSON)
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetAwsConnectionEntityWithContext(context.Background(), testEntityID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testEntityID, actual.ID)
	assert.Equal(t, "test-aws-connection", actual.Name)
	assert.Equal(t, "us-east-1", actual.Region)
	assert.Equal(t, true, actual.Enabled)
	assert.Equal(t, "arn:aws:iam::123456789012:role/test-role", actual.Credential.AssumeRole.RoleArn)
	assert.Equal(t, EntityManagementDynamicString("test-external-id"), actual.Credential.AssumeRole.ExternalId)
	assert.Equal(t, testAccountScopeID, actual.Scope.ID)
}

func TestGetAwsConnectionEntityWithContext_NotFound(t *testing.T) {
	t.Parallel()
	respJSON := `{ "data": { "actor": { "entityManagement": { "entity": null } } } }`
	client := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := client.GetAwsConnectionEntityWithContext(context.Background(), "nonexistent-id")

	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestGetAwsConnectionEntityWithContext_Error(t *testing.T) {
	t.Parallel()
	client := newMockResponse(t, `{"errors":[{"message":"internal error"}]}`, http.StatusInternalServerError)

	actual, err := client.GetAwsConnectionEntityWithContext(context.Background(), testEntityID)

	assert.Error(t, err)
	assert.Nil(t, actual)
}
