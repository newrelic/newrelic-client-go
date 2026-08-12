//go:build unit
// +build unit

package pathpoint

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testGetFlowResponseJSON = `{
	"data": {
		"actor": {
			"account": {
				"pathPoint": {
					"flow": {
						"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
						"name": "Supply Chain Flow",
						"id": "wl-123",
						"healthRollup": "AUTOMATIC_ROLL_UP",
						"healthStatus": "OPERATIONAL",
						"refreshInterval": "FIVE_MINUTES",
						"excludedKpis": [],
						"kpis": [],
						"stages": {
							"items": [],
							"totalCount": 0
						},
						"version": 1700000000000
					}
				}
			}
		}
	}
}`

var testGetFlowWithStagesResponseJSON = `{
	"data": {
		"actor": {
			"account": {
				"pathPoint": {
					"flow": {
						"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
						"name": "Supply Chain Flow",
						"id": "wl-456",
						"healthRollup": "AUTOMATIC_ROLL_UP",
						"healthStatus": "DEGRADED",
						"excludedKpis": [],
						"kpis": [],
						"stages": {
							"items": [
								{
									"id": "stage-001",
									"name": "Checkout",
									"healthRollup": "AUTOMATIC_ROLL_UP",
									"healthStatus": "DEGRADED",
									"link": "https://example.com/checkout",
									"stageKpis": [],
									"levels": {
										"items": [
											{
												"id": "level-001",
												"healthStatus": "DEGRADED",
												"steps": {
													"items": [
														{
															"id": "step-001",
															"name": "Payment",
															"healthStatus": "DISRUPTED",
															"isExcluded": false,
															"link": "https://example.com/payment",
															"config": {
																"healthRollup": "WORST_STATUS_WINS",
																"thresholdType": "FIXED",
																"thresholdValue": 1
															},
															"entitySearchQuery": {
																"isExcluded": false,
																"query": "domain = 'APM' AND type = 'APPLICATION'"
															},
															"signals": [
																{
																	"guid": "signal-guid-abc",
																	"name": "Payment Service",
																	"type": "ENTITY",
																	"isExcluded": false
																}
															]
														}
													],
													"totalCount": 1
												}
											}
										],
										"totalCount": 1
									},
									"related": {
										"source": false,
										"target": true
									}
								}
							],
							"totalCount": 1
						},
						"version": 1700000000000
					}
				}
			}
		}
	}
}`

var testCreateFlowResponseJSON = `{
	"data": {
		"pathPointCreate": {
			"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
			"name": "Supply Chain Flow",
			"id": "wl-123",
			"healthRollup": "AUTOMATIC_ROLL_UP",
			"healthStatus": "",
			"refreshInterval": "FIVE_MINUTES",
			"excludedKpis": [],
			"kpis": [],
			"stages": {
				"items": [],
				"totalCount": 0
			},
			"version": 1700000000000
		}
	}
}`

var testUpdateFlowResponseJSON = `{
	"data": {
		"pathPointUpdate": {
			"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
			"name": "Updated Pathpoint Flow",
			"id": "wl-123",
			"healthRollup": "ALERT_CONDITIONS",
			"healthStatus": "",
			"refreshInterval": "TEN_MINUTES",
			"excludedKpis": [],
			"kpis": [],
			"stages": {
				"items": [],
				"totalCount": 0
			},
			"version": 1700000001000
		}
	}
}`

var testDeleteFlowResponseJSON = `{
	"data": {
		"pathPointDelete": {
			"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
			"name": "Supply Chain Flow"
		}
	}
}`

var testGetFlowUnknownStatusResponseJSON = `{
	"data": {
		"actor": {
			"account": {
				"pathPoint": {
					"flow": {
						"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
						"name": "Supply Chain Flow",
						"id": "wl-123",
						"healthRollup": "AUTOMATIC_ROLL_UP",
						"healthStatus": "UNKNOWN",
						"refreshInterval": "FIVE_MINUTES",
						"excludedKpis": ["kpi-a", "kpi-b"],
						"kpis": [],
						"stages": {
							"items": [],
							"totalCount": 0
						},
						"message": "Step 'Broken Step' could not be created.",
						"version": 1700000000000
					}
				}
			}
		}
	}
}`

var testGetFlowWithAlertSignalResponseJSON = `{
	"data": {
		"actor": {
			"account": {
				"pathPoint": {
					"flow": {
						"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
						"name": "Supply Chain Flow",
						"id": "wl-789",
						"healthRollup": "AUTOMATIC_ROLL_UP",
						"healthStatus": "DISRUPTED",
						"excludedKpis": [],
						"kpis": [],
						"stages": {
							"items": [
								{
									"id": "stage-002",
									"name": "Fulfillment",
									"healthRollup": "ALERT_CONDITIONS",
									"healthStatus": "DISRUPTED",
									"link": "",
									"stageKpis": [],
									"levels": {
										"items": [
											{
												"id": "level-002",
												"healthStatus": "DISRUPTED",
												"steps": {
													"items": [
														{
															"id": "step-002",
															"name": "Alert Step",
															"healthStatus": "DISRUPTED",
															"isExcluded": true,
															"link": "",
															"config": {
																"healthRollup": "BEST_STATUS_WINS",
																"thresholdType": "PERCENTAGE",
																"thresholdValue": 80
															},
															"entitySearchQuery": {
																"isExcluded": true,
																"query": "domain = 'INFRA'"
															},
															"signals": [
																{
																	"guid": "alert-signal-guid",
																	"name": "High Error Rate",
																	"type": "ALERT",
																	"isExcluded": true
																}
															]
														}
													],
													"nextCursor": "cursor-abc",
													"totalCount": 5
												}
											}
										],
										"nextCursor": "level-cursor",
										"totalCount": 3
									},
									"related": {
										"source": true,
										"target": false
									}
								}
							],
							"nextCursor": "stage-cursor",
							"totalCount": 2
						},
						"version": 1700000000000
					}
				}
			}
		}
	}
}`

var testCreateFlowWithPartialFailureResponseJSON = `{
	"data": {
		"pathPointCreate": {
			"guid": "MTE5NTA0MDN8TkdFUHxGTE9XfDAxOWZkZDM1LTJmZWQtN2Y0My05OWJjLWYyNzM1YmZiN2Y0NA",
			"name": "Supply Chain Flow",
			"id": "wl-123",
			"healthRollup": "AUTOMATIC_ROLL_UP",
			"healthStatus": "",
			"refreshInterval": "FIVE_MINUTES",
			"excludedKpis": ["flow-kpi-1"],
			"message": "KPI 'flow-kpi-1' could not be created.",
			"kpis": [],
			"stages": {
				"items": [],
				"totalCount": 0
			},
			"version": 1700000000000
		}
	}
}`

func TestUnitGetFlow(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testGetFlowResponseJSON, http.StatusOK)

	actual, err := client.GetFlow(testPathpointAccountID, testFlowGUID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testFlowGUID, actual.GUID)
	assert.Equal(t, testFlowName, actual.Name)
	assert.Equal(t, PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP, actual.HealthRollup)
	assert.Equal(t, PathPointStatusValueTypes.OPERATIONAL, actual.HealthStatus)
	assert.Equal(t, PathPointRefreshIntervalTypes.FIVE_MINUTES, actual.RefreshInterval)
}

func TestUnitGetFlow_WithStages(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testGetFlowWithStagesResponseJSON, http.StatusOK)

	actual, err := client.GetFlow(testPathpointAccountID, testFlowGUID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testFlowGUID, actual.GUID)
	assert.Equal(t, PathPointStatusValueTypes.DEGRADED, actual.HealthStatus)

	require.Equal(t, 1, actual.Stages.TotalCount)
	require.Len(t, actual.Stages.Items, 1)

	stage := actual.Stages.Items[0]
	assert.Equal(t, "stage-001", stage.ID)
	assert.Equal(t, "Checkout", stage.Name)
	assert.Equal(t, PathPointStageHealthRollupTypes.AUTOMATIC_ROLL_UP, stage.HealthRollup)
	assert.Equal(t, PathPointStatusValueTypes.DEGRADED, stage.HealthStatus)

	require.Equal(t, 1, stage.Levels.TotalCount)
	require.Len(t, stage.Levels.Items, 1)

	level := stage.Levels.Items[0]
	assert.Equal(t, "level-001", level.ID)
	require.Equal(t, 1, level.Steps.TotalCount)
	require.Len(t, level.Steps.Items, 1)

	step := level.Steps.Items[0]
	assert.Equal(t, "step-001", step.ID)
	assert.Equal(t, "Payment", step.Name)
	assert.Equal(t, PathPointStatusValueTypes.DISRUPTED, step.HealthStatus)
	assert.Equal(t, PathPointStepHealthRollupTypes.WORST_STATUS_WINS, step.Config.HealthRollup)
	assert.Equal(t, PathPointThresholdTypeTypes.FIXED, step.Config.ThresholdType)
	assert.Equal(t, 1, step.Config.ThresholdValue)

	require.Len(t, step.Signals, 1)
	assert.Equal(t, EntityGUID("signal-guid-abc"), step.Signals[0].GUID)
	assert.Equal(t, PathPointSignalTypeTypes.ENTITY, step.Signals[0].Type)
}

func TestUnitPathPointCreate(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testCreateFlowResponseJSON, http.StatusOK)

	input := PathPointFlowInput{
		Name:            testFlowName,
		HealthRollup:    PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP,
		RefreshInterval: PathPointRefreshIntervalTypes.FIVE_MINUTES,
	}
	scope := PathPointScopeInput{
		ID:   testPathpointAccountID,
		Type: PathPointScopeTypeTypes.ACCOUNT,
	}

	actual, err := client.PathPointCreate(input, scope)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testFlowGUID, actual.GUID)
	assert.Equal(t, testFlowName, actual.Name)
	assert.Equal(t, PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP, actual.HealthRollup)
	assert.Equal(t, PathPointRefreshIntervalTypes.FIVE_MINUTES, actual.RefreshInterval)
	assert.Equal(t, "wl-123", actual.ID)
}

func TestUnitPathPointUpdate(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testUpdateFlowResponseJSON, http.StatusOK)

	update := PathPointFlowUpdateInput{
		Name:            "Updated Pathpoint Flow",
		HealthRollup:    PathPointFlowHealthRollupTypes.ALERT_CONDITIONS,
		RefreshInterval: PathPointRefreshIntervalTypes.TEN_MINUTES,
	}

	actual, err := client.PathPointUpdate(testFlowGUID, update)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testFlowGUID, actual.GUID)
	assert.Equal(t, "Updated Pathpoint Flow", actual.Name)
	assert.Equal(t, PathPointFlowHealthRollupTypes.ALERT_CONDITIONS, actual.HealthRollup)
	assert.Equal(t, PathPointRefreshIntervalTypes.TEN_MINUTES, actual.RefreshInterval)
}

func TestUnitPathPointDelete(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testDeleteFlowResponseJSON, http.StatusOK)

	actual, err := client.PathPointDelete(testFlowGUID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, testFlowGUID, actual.GUID)
	assert.Equal(t, testFlowName, actual.Name)
}

func TestUnitGetFlow_UnknownStatus_ExcludedKpis_Message(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testGetFlowUnknownStatusResponseJSON, http.StatusOK)

	actual, err := client.GetFlow(testPathpointAccountID, testFlowGUID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, PathPointStatusValueTypes.UNKNOWN, actual.HealthStatus)
	assert.Equal(t, []string{"kpi-a", "kpi-b"}, actual.ExcludedKpis)
	assert.Equal(t, "Step 'Broken Step' could not be created.", actual.Message)
}

func TestUnitGetFlow_AlertSignal_IsExcluded_Pagination(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testGetFlowWithAlertSignalResponseJSON, http.StatusOK)

	actual, err := client.GetFlow(testPathpointAccountID, testFlowGUID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, PathPointStatusValueTypes.DISRUPTED, actual.HealthStatus)

	require.Equal(t, 2, actual.Stages.TotalCount)
	assert.Equal(t, "stage-cursor", actual.Stages.NextCursor)
	require.Len(t, actual.Stages.Items, 1)

	stage := actual.Stages.Items[0]
	assert.Equal(t, PathPointStageHealthRollupTypes.ALERT_CONDITIONS, stage.HealthRollup)
	assert.True(t, stage.Related.Source)
	assert.False(t, stage.Related.Target)

	require.Equal(t, 3, stage.Levels.TotalCount)
	assert.Equal(t, "level-cursor", stage.Levels.NextCursor)

	level := stage.Levels.Items[0]
	require.Equal(t, 5, level.Steps.TotalCount)
	assert.Equal(t, "cursor-abc", level.Steps.NextCursor)
	require.Len(t, level.Steps.Items, 1)

	step := level.Steps.Items[0]
	assert.True(t, step.IsExcluded)
	assert.Equal(t, PathPointStepHealthRollupTypes.BEST_STATUS_WINS, step.Config.HealthRollup)
	assert.Equal(t, PathPointThresholdTypeTypes.PERCENTAGE, step.Config.ThresholdType)
	assert.Equal(t, 80, step.Config.ThresholdValue)
	assert.True(t, step.EntitySearchQuery.IsExcluded)

	require.Len(t, step.Signals, 1)
	signal := step.Signals[0]
	assert.Equal(t, EntityGUID("alert-signal-guid"), signal.GUID)
	assert.Equal(t, PathPointSignalTypeTypes.ALERT, signal.Type)
	assert.True(t, signal.IsExcluded)
}

func TestUnitGetFlow_Error(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, `{"errors": [{"message": "flow not found"}]}`, http.StatusNotFound)

	result, err := client.GetFlow(testPathpointAccountID, "nonexistent-guid")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestUnitPathPointCreate_Error(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, `{"errors": [{"message": "Internal Server Error"}]}`, http.StatusInternalServerError)

	input := PathPointFlowInput{Name: testFlowName}
	scope := PathPointScopeInput{ID: testPathpointAccountID, Type: PathPointScopeTypeTypes.ACCOUNT}

	result, err := client.PathPointCreate(input, scope)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestUnitPathPointCreate_PartialFailure_ExcludedKpis(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testCreateFlowWithPartialFailureResponseJSON, http.StatusOK)

	input := PathPointFlowInput{Name: testFlowName}
	scope := PathPointScopeInput{ID: testPathpointAccountID, Type: PathPointScopeTypeTypes.ACCOUNT}

	actual, err := client.PathPointCreate(input, scope)

	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, []string{"flow-kpi-1"}, actual.ExcludedKpis)
	assert.Equal(t, "KPI 'flow-kpi-1' could not be created.", actual.Message)
}

func TestUnitPathPointUpdate_Error(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, `{"errors": [{"message": "flow not found"}]}`, http.StatusNotFound)

	update := PathPointFlowUpdateInput{Name: "Updated Flow"}

	result, err := client.PathPointUpdate("nonexistent-guid", update)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestUnitPathPointDelete_Error(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, `{"errors": [{"message": "flow not found"}]}`, http.StatusNotFound)

	result, err := client.PathPointDelete("nonexistent-guid")

	require.Error(t, err)
	require.Nil(t, result)
}
