//go:build unit
// +build unit

package federatedlogs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testAwsConnectionID = "ZmVkLWxvZ3MtYXdzLWNvbm4tMTIzNDU2Nzg5MA"

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
