//go:build integration
// +build integration

package dashboards

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationDashboard_Billboard(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	dashboardName := "newrelic-client-go test-dashboard-" + mock.RandSeq(5)

	// Create RawConfiguration with ChartStyles for line widget
	lineRawConfig := RawConfiguration{
		NRQLQueries: []DashboardWidgetNRQLQueryInput{
			{
				AccountID: testAccountID,
				Query:     "FROM Metric SELECT average(duration) TIMESERIES",
			},
		},
		ChartStyles: &DashboardWidgetChartStyles{
			LineInterpolation: DashboardLineInterpolationTypes.SMOOTH,
			Gradient: &DashboardWidgetChartStylesGradient{
				Enabled: true,
			},
		},
	}
	lineRawConfigJSON, _ := json.Marshal(lineRawConfig)

	dashboardInput := DashboardInput{
		Description: "Test description",
		Name:        dashboardName,
		Permissions: entities.DashboardPermissionsTypes.PRIVATE,
		Pages: []DashboardPageInput{
			{
				Description: "Test page description",
				Name:        "Test Page",
				Widgets: []DashboardWidgetInput{
					{
						Title: "Test Billboard Widget",
						Configuration: DashboardWidgetConfigurationInput{
							Billboard: &DashboardBillboardWidgetConfigurationInput{
								NRQLQueries: []DashboardWidgetNRQLQueryInput{
									{
										AccountID: testAccountID,
										Query:     "FROM Metric SELECT 1",
									},
								},
								// Thresholds: []DashboardBillboardWidgetThresholdInput
							},
						},
					},
					{
						Title: "Test Line Widget with ChartStyles",
						Configuration: DashboardWidgetConfigurationInput{
							Line: &DashboardLineWidgetConfigurationInput{
								NRQLQueries: []DashboardWidgetNRQLQueryInput{
									{
										AccountID: testAccountID,
										Query:     "FROM Metric SELECT average(duration) TIMESERIES",
									},
								},
							},
						},
						RawConfiguration: lineRawConfigJSON,
					},
				},
			},
		},
	}

	// Test: DashboardCreate
	result, err := client.DashboardCreate(testAccountID, dashboardInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, len(result.Errors))
	require.NotNil(t, result.EntityResult.GUID)

	dashGUID := result.EntityResult.GUID

	// Test: GetDashboardEntity
	dash, err := client.GetDashboardEntity(dashGUID)
	require.NoError(t, err)
	require.NotNil(t, dash)

	assert.Equal(t, dashGUID, dash.GUID)
	assert.Equal(t, dashboardInput.Description, dash.Description)
	assert.Equal(t, dashboardInput.Name, dash.Name)
	assert.Equal(t, entities.DashboardEntityPermissions(string(dashboardInput.Permissions)), dash.Permissions)

	// Input and Pages are different types so we can not easily compare them...
	assert.Equal(t, len(dashboardInput.Pages), len(dash.Pages))
	require.Equal(t, 1, len(dash.Pages))
	require.Equal(t, 2, len(dash.Pages[0].Widgets))

	// Check billboard widget
	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Title, dash.Pages[0].Widgets[0].Title)
	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries[0].Query, dash.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries[0].Query)
	assert.Nil(t, dash.Pages[0].Widgets[0].Configuration.Billboard.Thresholds)
	assert.Greater(t, len(dash.Pages[0].Widgets[0].RawConfiguration), 1)

	// Check line widget with ChartStyles
	assert.Equal(t, dashboardInput.Pages[0].Widgets[1].Title, dash.Pages[0].Widgets[1].Title)
	assert.Equal(t, dashboardInput.Pages[0].Widgets[1].Configuration.Line.NRQLQueries[0].Query, dash.Pages[0].Widgets[1].Configuration.Line.NRQLQueries[0].Query)

	// Verify ChartStyles in RawConfiguration
	var lineWidgetRawConfig RawConfiguration
	err = json.Unmarshal(dash.Pages[0].Widgets[1].RawConfiguration, &lineWidgetRawConfig)
	require.NoError(t, err)
	require.NotNil(t, lineWidgetRawConfig.ChartStyles)
	assert.Equal(t, DashboardLineInterpolationTypes.SMOOTH, lineWidgetRawConfig.ChartStyles.LineInterpolation)
	require.NotNil(t, lineWidgetRawConfig.ChartStyles.Gradient)
	assert.True(t, lineWidgetRawConfig.ChartStyles.Gradient.Enabled)

	testWarningThreshold := 10.0

	// Test: DashboardUpdate - Update ChartStyles
	updatedLineRawConfig := RawConfiguration{
		NRQLQueries: []DashboardWidgetNRQLQueryInput{
			{
				AccountID: testAccountID,
				Query:     "FROM Metric SELECT average(duration) TIMESERIES",
			},
		},
		ChartStyles: &DashboardWidgetChartStyles{
			LineInterpolation: DashboardLineInterpolationTypes.LINEAR,
			Gradient: &DashboardWidgetChartStylesGradient{
				Enabled: false,
			},
		},
	}
	updatedLineRawConfigJSON, _ := json.Marshal(updatedLineRawConfig)

	updatedDashboard := DashboardInput{
		Name:        dash.Name,
		Permissions: result.EntityResult.Permissions,
		Pages: []DashboardPageInput{
			{
				Name: dash.Pages[0].Name,
				Widgets: []DashboardWidgetInput{
					{
						Title: "Test BillboardText Widget",
						Configuration: DashboardWidgetConfigurationInput{
							Billboard: &DashboardBillboardWidgetConfigurationInput{
								NRQLQueries: []DashboardWidgetNRQLQueryInput{
									{
										AccountID: testAccountID,
										Query:     "FROM Metric SELECT 1",
									},
								},
								Thresholds: []DashboardBillboardWidgetThresholdInput{
									{
										AlertSeverity: entities.DashboardAlertSeverityTypes.WARNING,
										Value:         &testWarningThreshold,
									},
								},
							},
						},
					},
					{
						Title: "Test Line Widget with Updated ChartStyles",
						Configuration: DashboardWidgetConfigurationInput{
							Line: &DashboardLineWidgetConfigurationInput{
								NRQLQueries: []DashboardWidgetNRQLQueryInput{
									{
										AccountID: testAccountID,
										Query:     "FROM Metric SELECT average(duration) TIMESERIES",
									},
								},
							},
						},
						RawConfiguration: updatedLineRawConfigJSON,
					},
				},
			},
		},
	}

	upDash, err := client.DashboardUpdate(updatedDashboard, dashGUID)
	require.NoError(t, err)
	require.NotNil(t, upDash)

	require.Equal(t, 1, len(upDash.EntityResult.Pages))
	require.Equal(t, 2, len(upDash.EntityResult.Pages[0].Widgets))

	// Check updated billboard widget
	assert.Equal(t, updatedDashboard.Pages[0].Widgets[0].Title, upDash.EntityResult.Pages[0].Widgets[0].Title)
	require.Equal(t, 1, len(upDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries))
	assert.Equal(t, updatedDashboard.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries[0].Query, upDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries[0].Query)
	require.Equal(t, 1, len(upDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.Thresholds))
	assert.Equal(t, updatedDashboard.Pages[0].Widgets[0].Configuration.Billboard.Thresholds[0].AlertSeverity, upDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.Thresholds[0].AlertSeverity)
	assert.Equal(t, *updatedDashboard.Pages[0].Widgets[0].Configuration.Billboard.Thresholds[0].Value, upDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.Thresholds[0].Value)

	// Check updated line widget with ChartStyles
	assert.Equal(t, updatedDashboard.Pages[0].Widgets[1].Title, upDash.EntityResult.Pages[0].Widgets[1].Title)

	// Verify updated ChartStyles in RawConfiguration
	var updatedLineWidgetRawConfig RawConfiguration
	err = json.Unmarshal(upDash.EntityResult.Pages[0].Widgets[1].RawConfiguration, &updatedLineWidgetRawConfig)
	require.NoError(t, err)
	require.NotNil(t, updatedLineWidgetRawConfig.ChartStyles)
	assert.Equal(t, DashboardLineInterpolationTypes.LINEAR, updatedLineWidgetRawConfig.ChartStyles.LineInterpolation)
	require.NotNil(t, updatedLineWidgetRawConfig.ChartStyles.Gradient)
	assert.False(t, updatedLineWidgetRawConfig.ChartStyles.Gradient.Enabled)

	// Test removal of threshold (set back to initial input)
	removeThresholdDash, err := client.DashboardUpdate(dashboardInput, dashGUID)
	require.NoError(t, err)
	require.NotNil(t, removeThresholdDash)

	require.Equal(t, 1, len(removeThresholdDash.EntityResult.Pages))
	require.Equal(t, 2, len(removeThresholdDash.EntityResult.Pages[0].Widgets))

	// Check billboard widget (threshold removed)
	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Title, removeThresholdDash.EntityResult.Pages[0].Widgets[0].Title)
	require.Equal(t, 1, len(removeThresholdDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries))
	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries[0].Query, removeThresholdDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.NRQLQueries[0].Query)
	assert.Nil(t, removeThresholdDash.EntityResult.Pages[0].Widgets[0].Configuration.Billboard.Thresholds)

	// Verify line widget still has original ChartStyles
	assert.Equal(t, dashboardInput.Pages[0].Widgets[1].Title, removeThresholdDash.EntityResult.Pages[0].Widgets[1].Title)
	var finalLineWidgetRawConfig RawConfiguration
	err = json.Unmarshal(removeThresholdDash.EntityResult.Pages[0].Widgets[1].RawConfiguration, &finalLineWidgetRawConfig)
	require.NoError(t, err)
	require.NotNil(t, finalLineWidgetRawConfig.ChartStyles)
	assert.Equal(t, DashboardLineInterpolationTypes.SMOOTH, finalLineWidgetRawConfig.ChartStyles.LineInterpolation)

	// Test: DashboardDelete
	delRes, err := client.DashboardDelete(dashGUID)
	require.NoError(t, err)
	require.NotNil(t, delRes)
	assert.Equal(t, 0, len(delRes.Errors))
	assert.Equal(t, DashboardDeleteResultStatusTypes.SUCCESS, delRes.Status)
}

// TestIntegrationDashboard_EmptyPage tests creating a dashboard with a page comprising no widgets
func TestIntegrationDashboard_EmptyPage(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	dashboardName := "newrelic-client-go-test-dashboard-empty-pages" + mock.RandSeq(5)
	dashboardInput := DashboardInput{
		Description: "newrelic-client-go-test-dashboard-description",
		Name:        dashboardName,
		Permissions: entities.DashboardPermissionsTypes.PUBLIC_READ_WRITE,
		Pages: []DashboardPageInput{{
			Name:    "Test Page",
			Widgets: []DashboardWidgetInput{},
		}},
	}

	// Test: Create Dashboard
	result, err := client.DashboardCreate(testAccountID, dashboardInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, len(result.Errors))
	require.NotNil(t, result.EntityResult.GUID)

	dashGUID := result.EntityResult.GUID

	// Test: Get Dashboard
	dash, err := client.GetDashboardEntity(dashGUID)
	require.NoError(t, err)
	require.NotNil(t, dash)

	assert.Equal(t, dashGUID, dash.GUID)
	assert.Equal(t, dashboardInput.Description, dash.Description)
	assert.Equal(t, dashboardInput.Name, dash.Name)
	assert.Equal(t, entities.DashboardEntityPermissions(string(dashboardInput.Permissions)), dash.Permissions)

	// Test: Update Dashboard
	updatedDashboard := DashboardInput{
		Name:        dash.Name,
		Permissions: result.EntityResult.Permissions,
		Pages: []DashboardPageInput{
			{
				Name: dash.Pages[0].Name,
				Widgets: []DashboardWidgetInput{
					{
						Title: "Test BillboardText Widget",
						Configuration: DashboardWidgetConfigurationInput{
							Billboard: &DashboardBillboardWidgetConfigurationInput{
								NRQLQueries: []DashboardWidgetNRQLQueryInput{
									{
										AccountID: testAccountID,
										Query:     "FROM Metric SELECT 1",
									},
								},
							},
						},
					},
				},
			},
			{
				Name:    "Test Page Two",
				Widgets: []DashboardWidgetInput{},
			},
		},
	}

	upDash, err := client.DashboardUpdate(updatedDashboard, dashGUID)
	require.NoError(t, err)
	require.NotNil(t, upDash)

	//// Test: Delete Dashboard
	delRes, err := client.DashboardDelete(dashGUID)
	require.NoError(t, err)
	require.NotNil(t, delRes)
	assert.Equal(t, 0, len(delRes.Errors))
	assert.Equal(t, DashboardDeleteResultStatusTypes.SUCCESS, delRes.Status)
}
