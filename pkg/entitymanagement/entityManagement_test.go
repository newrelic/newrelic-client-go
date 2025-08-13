//go:build unit
// +build unit

package entitymanagement

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	testAccountID = "12345"
	testRuleID    = "test-rule-id"
	testVersion   = 1

	testCreateResponseJSON = `
	{
		"data": {
			"entityManagementCreatePipelineCloudRule": {
				"entity": {
					"id": "test-rule-id",
					"metadata": {
						"version": 1
					},
					"name": "Test Rule",
					"description": "A test rule",
					"nrql": "SELECT * FROM Log"
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
						"id": "test-rule-id",
						"metadata": {
							"version": 1
						},
						"name": "Test Rule",
						"description": "A test rule",
						"nrql": "SELECT * FROM Log",
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
				"id": "test-rule-id"
			}
		}
	}`
)

func newMockClient(t *testing.T, mockJSONResponse string, statusCode int) Entitymanagement {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

func TestUnitEntityManagement_CreatePipelineCloudRule(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateResponseJSON, http.StatusOK)

	createInput := EntityManagementPipelineCloudRuleEntityCreateInput{
		Name:        "Test Rule",
		Description: "A test rule",
		NRQL:        "SELECT * FROM Log",
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

	result, err := client.EntityManagementDelete(testRuleID, testVersion)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testRuleID, result.ID)
}

func TestUnitEntityManagement_DeleteEntityWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testDeleteResponseJSON, http.StatusOK)

	result, err := client.EntityManagementDeleteWithContext(context.Background(), testRuleID, testVersion)

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
