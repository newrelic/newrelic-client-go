//go:build unit
// +build unit

package nrqlmacros

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testOrgID   = "fb33fea3-4d7e-4736-9701-acb59a634fdf"
	testMacroID = "NDgyOTY3M3xFTlRJVFlNQU5BR0VNRU5UfE5SUVFMX01BQ1JPfDEyMzQ1Njc4OQ"

	testCreateMacroResponseJSON = `
	{
		"data": {
			"entityManagementCreateNrqlMacro": {
				"entity": {
					"id": "NDgyOTY3M3xFTlRJVFlNQU5BR0VNRU5UfE5SUVFMX01BQ1JPfDEyMzQ1Njc4OQ",
					"name": "filterProdTraffic",
					"expression": "environment = 'production' AND appName NOT LIKE '%test%'",
					"type": "NRQL_MACRO",
					"numParameters": 0,
					"scope": {
						"id": "fb33fea3-4d7e-4736-9701-acb59a634fdf",
						"type": "ORGANIZATION"
					},
					"metadata": {
						"version": 1,
						"createdAt": 1730000000000,
						"updatedAt": 1730000000000
					},
					"tags": []
				}
			}
		}
	}`

	testUpdateMacroResponseJSON = `
	{
		"data": {
			"entityManagementUpdateNrqlMacro": {
				"entity": {
					"id": "NDgyOTY3M3xFTlRJVFlNQU5BR0VNRU5UfE5SUVFMX01BQ1JPfDEyMzQ1Njc4OQ",
					"name": "filterProdTraffic",
					"expression": "environment = 'staging' AND appName NOT LIKE '%test%'",
					"type": "NRQL_MACRO",
					"numParameters": 0,
					"scope": {
						"id": "fb33fea3-4d7e-4736-9701-acb59a634fdf",
						"type": "ORGANIZATION"
					},
					"metadata": {
						"version": 2,
						"createdAt": 1730000000000,
						"updatedAt": 1730000001000
					},
					"tags": []
				}
			}
		}
	}`

	testGetMacroResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entity": {
						"__typename": "EntityManagementNrqlMacroEntity",
						"id": "NDgyOTY3M3xFTlRJVFlNQU5BR0VNRU5UfE5SUVFMX01BQ1JPfDEyMzQ1Njc4OQ",
						"name": "filterProdTraffic",
						"expression": "environment = 'production' AND appName NOT LIKE '%test%'",
						"type": "NRQL_MACRO",
						"numParameters": 0,
						"scope": {
							"id": "fb33fea3-4d7e-4736-9701-acb59a634fdf",
							"type": "ORGANIZATION"
						},
						"metadata": {
							"version": 1,
							"createdAt": 1730000000000,
							"updatedAt": 1730000000000
						},
						"tags": []
					}
				}
			}
		}
	}`

	testDeleteMacroResponseJSON = `
	{
		"data": {
			"entityManagementDelete": {
				"id": "NDgyOTY3M3xFTlRJVFlNQU5BR0VNRU5UfE5SUVFMX01BQ1JPfDEyMzQ1Njc4OQ"
			}
		}
	}`

	testSearchMacroResponseJSON = `
	{
		"data": {
			"actor": {
				"entityManagement": {
					"entitySearch": {
						"entities": [
							{
								"__typename": "EntityManagementNrqlMacroEntity",
								"id": "NDgyOTY3M3xFTlRJVFlNQU5BR0VNRU5UfE5SUVFMX01BQ1JPfDEyMzQ1Njc4OQ",
								"name": "filterProdTraffic",
								"expression": "environment = 'production' AND appName NOT LIKE '%test%'",
								"type": "NRQL_MACRO",
								"numParameters": 0,
								"scope": {
									"id": "fb33fea3-4d7e-4736-9701-acb59a634fdf",
									"type": "ORGANIZATION"
								},
								"metadata": {
									"version": 1
								},
								"tags": []
							}
						],
						"nextCursor": null
					}
				}
			}
		}
	}`
)

func TestUnitNrqlmacros_CreateMacro(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testCreateMacroResponseJSON, http.StatusOK)

	createInput := EntityManagementNRQLMacroEntityCreateInput{
		Name:       "filterProdTraffic",
		Expression: "environment = 'production' AND appName NOT LIKE '%test%'",
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ORGANIZATION,
			ID:   testOrgID,
		},
	}

	result, err := client.EntityManagementCreateNRQLMacro(createInput)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testMacroID, result.Entity.ID)
	require.Equal(t, "filterProdTraffic", result.Entity.Name)
	require.Equal(t, createInput.Expression, result.Entity.Expression)
}

func TestUnitNrqlmacros_UpdateMacro(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testUpdateMacroResponseJSON, http.StatusOK)

	updateInput := EntityManagementNRQLMacroEntityUpdateInput{
		Expression: "environment = 'staging' AND appName NOT LIKE '%test%'",
	}

	result, err := client.EntityManagementUpdateNRQLMacro(testMacroID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testMacroID, result.Entity.ID)
	require.Equal(t, updateInput.Expression, result.Entity.Expression)
	require.Equal(t, 2, result.Entity.Metadata.Version)
}

func TestUnitNrqlmacros_GetEntity(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testGetMacroResponseJSON, http.StatusOK)

	result, err := client.GetEntity(testMacroID)
	require.NoError(t, err)
	require.NotNil(t, result)

	macroEntity, ok := (*result).(*EntityManagementNRQLMacroEntity)
	require.True(t, ok, "expected *EntityManagementNRQLMacroEntity, got %T", *result)
	require.Equal(t, testMacroID, macroEntity.ID)
	require.Equal(t, "filterProdTraffic", macroEntity.Name)
}

func TestUnitNrqlmacros_DeleteMacro(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testDeleteMacroResponseJSON, http.StatusOK)

	result, err := client.EntityManagementDelete(testMacroID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testMacroID, result.ID)
}

func TestUnitNrqlmacros_DeleteMacroWithContext(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testDeleteMacroResponseJSON, http.StatusOK)

	result, err := client.EntityManagementDeleteWithContext(context.Background(), testMacroID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testMacroID, result.ID)
}

func TestUnitNrqlmacros_GetEntitySearch(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, testSearchMacroResponseJSON, http.StatusOK)

	result, err := client.GetEntitySearch("", "type = 'NRQL_MACRO'")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Entities, 1)

	macroEntity, ok := result.Entities[0].(*EntityManagementNRQLMacroEntity)
	require.True(t, ok, "expected *EntityManagementNRQLMacroEntity")
	require.Equal(t, testMacroID, macroEntity.ID)
}

func TestUnitNrqlmacros_GetEntity_Error(t *testing.T) {
	t.Parallel()
	client := newMockClient(t, `{"errors": [{"message": "Not Found"}]}`, http.StatusNotFound)

	result, err := client.GetEntity("non-existent-id")
	require.Error(t, err)
	require.Nil(t, result)
}
