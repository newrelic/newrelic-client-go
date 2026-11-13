//go:build unit || integration
// +build unit integration

package maintenancewindow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newMockClient(t *testing.T, mockJSONResponse string, statusCode int) Maintenancewindow {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	t.Cleanup(ts.Close)

	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

func TestMaintenanceWindowCreate_Mock(t *testing.T) {
	t.Parallel()

	mockResponse := `{
		"data": {
			"maintenanceWindowCreate": {
				"id": "12345",
				"name": "tf-test-mw",
				"description": "a test maintenance window",
				"duration": "PT1H",
				"rrule": "",
				"timezone": "America/Los_Angeles",
				"affectedEntityType": "SERVICE_LEVEL",
				"affectedEntities": ["guid-1", "guid-2"],
				"startTime": "2026-12-01T22:00:00",
				"scope": {"id": "1067061", "type": "ACCOUNT"},
				"metadata": {
					"createdAt": 1700000000000,
					"updatedAt": 1700000000000,
					"createdBy": {"id": "user-1"},
					"updatedBy": {"id": "user-1"}
				}
			}
		}
	}`

	client := newMockClient(t, mockResponse, 200)

	input := MaintenanceWindowConfigCreateInput{
		Name:               "tf-test-mw",
		Description:        "a test maintenance window",
		Scope:              MaintenanceWindowScopedReferenceInput{ID: "1067061", Type: MaintenanceWindowEntityScopeTypes.ACCOUNT},
		StartTime:          "2026-12-01T22:00:00",
		Duration:           "PT1H",
		Timezone:           "America/Los_Angeles",
		AffectedEntityType: MaintenanceWindowAffectedEntityTypeTypes.SERVICE_LEVEL,
		AffectedEntities:   []string{"guid-1", "guid-2"},
	}

	resp, err := client.MaintenanceWindowCreate(input)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "12345", resp.ID)
	assert.Equal(t, "tf-test-mw", resp.Name)
	assert.Equal(t, "a test maintenance window", resp.Description)
	assert.Equal(t, "PT1H", string(resp.Duration))
	assert.Equal(t, "America/Los_Angeles", resp.Timezone)
	assert.Equal(t, MaintenanceWindowAffectedEntityTypeTypes.SERVICE_LEVEL, resp.AffectedEntityType)
	assert.Equal(t, []string{"guid-1", "guid-2"}, resp.AffectedEntities)
	assert.Equal(t, "1067061", resp.Scope.ID)
	assert.Equal(t, MaintenanceWindowEntityScopeTypes.ACCOUNT, resp.Scope.Type)
	assert.Equal(t, NaiveDateTime("2026-12-01T22:00:00"), resp.StartTime)
}

func TestMaintenanceWindowUpdate_Mock(t *testing.T) {
	t.Parallel()

	mockResponse := `{
		"data": {
			"maintenanceWindowUpdate": {
				"id": "12345",
				"name": "tf-test-mw",
				"description": "updated description",
				"duration": "PT2H",
				"rrule": "",
				"timezone": "America/Los_Angeles",
				"affectedEntityType": "SERVICE_LEVEL",
				"affectedEntities": ["guid-1"],
				"startTime": "2026-12-01T22:00:00",
				"scope": {"id": "1067061", "type": "ACCOUNT"},
				"metadata": {
					"createdAt": 1700000000000,
					"updatedAt": 1700000001000,
					"createdBy": {"id": "user-1"},
					"updatedBy": {"id": "user-1"}
				}
			}
		}
	}`

	client := newMockClient(t, mockResponse, 200)

	input := MaintenanceWindowConfigUpdateInput{
		Description: "updated description",
		Duration:    "PT2H",
	}

	resp, err := client.MaintenanceWindowUpdate("12345", input)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "12345", resp.ID)
	assert.Equal(t, "updated description", resp.Description)
	assert.Equal(t, "PT2H", string(resp.Duration))
	assert.Equal(t, []string{"guid-1"}, resp.AffectedEntities)
}

func TestMaintenanceWindowDelete_Mock(t *testing.T) {
	t.Parallel()

	mockResponse := `{
		"data": {
			"maintenanceWindowDelete": {
				"id": "12345",
				"name": "tf-test-mw"
			}
		}
	}`

	client := newMockClient(t, mockResponse, 200)

	resp, err := client.MaintenanceWindowDelete("12345")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "12345", resp.ID)
}

func TestGetListByIds_Mock(t *testing.T) {
	t.Parallel()

	mockResponse := `{
		"data": {
			"actor": {
				"maintenanceWindow": {
					"listByIds": {
						"maintenanceWindows": [
							{
								"id": "12345",
								"name": "tf-test-mw",
								"description": "a test maintenance window",
								"duration": "PT1H",
								"rrule": "",
								"timezone": "America/Los_Angeles",
								"affectedEntityType": "SERVICE_LEVEL",
								"affectedEntities": ["guid-1", "guid-2"],
								"startTime": "2026-12-01T22:00:00",
								"scope": {"id": "1067061", "type": "ACCOUNT"},
								"metadata": {
									"createdAt": 1700000000000,
									"updatedAt": 1700000000000,
									"createdBy": {"id": "user-1"},
									"updatedBy": {"id": "user-1"}
								}
							}
						]
					}
				}
			}
		}
	}`

	client := newMockClient(t, mockResponse, 200)

	resp, err := client.GetListByIds([]string{"12345"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.MaintenanceWindows, 1)

	mw := resp.MaintenanceWindows[0]
	assert.Equal(t, "12345", mw.ID)
	assert.Equal(t, "tf-test-mw", mw.Name)
	assert.Equal(t, []string{"guid-1", "guid-2"}, mw.AffectedEntities)
}
