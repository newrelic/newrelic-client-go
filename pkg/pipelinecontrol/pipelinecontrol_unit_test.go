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
			"EntityManagementCreateFederatedLogSetup": {
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
