//go:build integration
// +build integration

package dashboards

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Dashboards {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}

func TestIntegrationDashboard_Nil(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: GetDashboardEntity
	dash, err := client.GetDashboardEntity(`bad-guid`)
	require.NotNil(t, err)
	require.Error(t, err)
	assert.Nil(t, dash)
}

func TestIntegrationDashboard_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	dashboardName := "newrelic-client-go test-dashboard-" + mock.RandSeq(5)
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
						Title: "Test Text Widget",
						Configuration: DashboardWidgetConfigurationInput{
							Markdown: &DashboardMarkdownWidgetConfigurationInput{
								Text: "Test Text widget **markdown**",
							},
						},
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
	assert.Equal(t, dashboardInput.Permissions, dash.Permissions)

	// Input and Pages are different types so we can not easily compare them...
	assert.Equal(t, len(dashboardInput.Pages), len(dash.Pages))
	require.Equal(t, 1, len(dash.Pages))
	require.Equal(t, 1, len(dash.Pages[0].Widgets))

	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Title, dash.Pages[0].Widgets[0].Title)
	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Configuration.Markdown.Text, dash.Pages[0].Widgets[0].Configuration.Markdown.Text)
	assert.Greater(t, len(dash.Pages[0].Widgets[0].RawConfiguration), 1)

	// Test: DashboardUpdate
	updatedDashboard := DashboardInput{
		Name:        dash.Name,
		Permissions: dash.Permissions,
		Pages: []DashboardPageInput{
			{
				Name: dash.Pages[0].Name,
				Widgets: []DashboardWidgetInput{
					{
						// Even though the config isn't changing, we have to pass it. 2021-01-11 JT
						Configuration: dashboardInput.Pages[0].Widgets[0].Configuration,
						ID:            dash.Pages[0].Widgets[0].ID,
						Title:         "Updated Title",
					},
				},
			},
		},
	}

	upDash, err := client.DashboardUpdate(updatedDashboard, dashGUID)
	require.NoError(t, err)
	require.NotNil(t, upDash)

	require.Equal(t, 1, len(upDash.EntityResult.Pages))
	require.Equal(t, 1, len(upDash.EntityResult.Pages[0].Widgets))
	assert.Equal(t, updatedDashboard.Pages[0].Widgets[0].Title, upDash.EntityResult.Pages[0].Widgets[0].Title)

	// Test: DashboardDelete
	delRes, err := client.DashboardDelete(dashGUID)
	require.NoError(t, err)
	require.NotNil(t, delRes)
	assert.Equal(t, 0, len(delRes.Errors))
	assert.Equal(t, DashboardDeleteResultStatusTypes.SUCCESS, delRes.Status)
}

func TestIntegrationDashboard_LinkedEntities(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	dashboardAName := "newrelic-client-go test-dashboard-" + mock.RandSeq(5)
	dashboardAInput := DashboardInput{
		Description: "Test description",
		Name:        dashboardAName,
		Permissions: entities.DashboardPermissionsTypes.PRIVATE,
		Pages: []DashboardPageInput{
			{
				Description: "Test page description",
				Name:        "Test Page",
				Widgets: []DashboardWidgetInput{
					{
						Title: "Test Text Widget",
						Configuration: DashboardWidgetConfigurationInput{
							Markdown: &DashboardMarkdownWidgetConfigurationInput{
								Text: "Test Text widget **markdown**",
							},
						},
					},
				},
			},
		},
	}

	// Create a dashboard to reference in linked entity GUIDs
	resultDashA, err := client.DashboardCreate(testAccountID, dashboardAInput)
	require.NoError(t, err)
	require.NotNil(t, resultDashA)
	defer client.DashboardDelete(resultDashA.EntityResult.GUID) // Clean up dashboard A

	dashboardBName := "newrelic-client-go test-dashboard-" + mock.RandSeq(5)
	dashboardBInput := DashboardInput{
		Description: "Testing dashboard widget with linked entities",
		Name:        dashboardBName,
		Permissions: entities.DashboardPermissionsTypes.PRIVATE,
		Pages: []DashboardPageInput{
			{
				Description: "Test page description",
				Name:        "Test Page",
				Widgets: []DashboardWidgetInput{
					{
						Title: "Widget with linked entities",
						Configuration: DashboardWidgetConfigurationInput{
							Bar: &DashboardBarWidgetConfigurationInput{
								NRQLQueries: []DashboardWidgetNRQLQueryInput{
									{
										AccountID: testAccountID,
										Query:     "FROM Transaction SELECT average(duration) FACET appName",
									},
								},
							},
						},
						LinkedEntityGUIDs: []common.EntityGUID{
							common.EntityGUID(resultDashA.EntityResult.Pages[0].GUID),
						},
					},
				},
			},
		},
	}

	// Test: Create dashboard with a widget that includes `linkedEntityGuids`
	resultDashB, err := client.DashboardCreate(testAccountID, dashboardBInput)
	require.NoError(t, err)
	require.NotNil(t, resultDashB)
	defer client.DashboardDelete(resultDashB.EntityResult.GUID) // Clean up dashboard B

	assert.Equal(t, 0, len(resultDashB.Errors))
	assert.NotNil(t, resultDashB.EntityResult.GUID)

	// Test: GetDashboardEntity
	dashB, err := client.GetDashboardEntity(resultDashB.EntityResult.GUID)
	require.NoError(t, err)
	require.NotNil(t, dashB)
	assert.Greater(t, len(dashB.Pages[0].Widgets[0].LinkedEntities), 0)
}

func TestIntegrationDashboard_InvalidInput(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	dashboardName := "newrelic-client-go test-dashboard-" + mock.RandSeq(5)
	dashboardInput := DashboardInput{
		Description: "Testing dashboard widget with linked entities",
		Name:        dashboardName,
		Permissions: entities.DashboardPermissionsTypes.PRIVATE,
		Pages: []DashboardPageInput{
			{
				Description: "Test page description",
				Name:        "Test Page",
				Widgets: []DashboardWidgetInput{
					{
						Title: "Widget with bad NRQL",
						Configuration: DashboardWidgetConfigurationInput{
							Bar: &DashboardBarWidgetConfigurationInput{
								NRQLQueries: []DashboardWidgetNRQLQueryInput{
									{
										AccountID: testAccountID,
										Query:     "This is bad NRQL input",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Test: DashboardCreate
	dash, err := client.DashboardCreate(testAccountID, dashboardInput)

	require.Nil(t, dash)
	require.Error(t, err)
}

func TestIntegrationDashboard_Variables(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	dashboardName := "newrelic-client-go test-dashboard-" + mock.RandSeq(5)
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
						Title: "Test Text Widget",
						Configuration: DashboardWidgetConfigurationInput{
							Markdown: &DashboardMarkdownWidgetConfigurationInput{
								Text: "Test Text widget **markdown**",
							},
						},
					},
				},
			},
		},
		Variables: []DashboardVariableInput{
			{
				DefaultValues: &[]DashboardVariableDefaultItemInput{
					{
						Value: DashboardVariableDefaultValueInput{
							String: "default value",
						},
					},
				},
				IsMultiSelection: true,
				Items: []DashboardVariableEnumItemInput{
					{
						Title: "VARIABLE",
						Value: "Variable",
					},
				},
				Name: "variable",
				Options: DashboardVariableOptionsInput{
					IgnoreTimeRange: true,
				},
				NRQLQuery: &DashboardVariableNRQLQueryInput{
					AccountIDs: []int{testAccountID},
					Query:      nrdb.NRQL("SELECT average(duration) FROM Transaction TIMESERIES WHERE appName = 'Dummy App'"),
				},
				ReplacementStrategy: DashboardVariableReplacementStrategyTypes.DEFAULT,
				Title:               "variable title",
				Type:                DashboardVariableTypeTypes.NRQL,
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
	assert.Equal(t, dashboardInput.Permissions, dash.Permissions)
	assert.Equal(t, dashboardInput.Variables[0].Options.IgnoreTimeRange, dash.Variables[0].Options.IgnoreTimeRange)

	// Input and Pages are different types so we can not easily compare them...
	assert.Equal(t, len(dashboardInput.Pages), len(dash.Pages))
	require.Equal(t, 1, len(dash.Pages))
	require.Equal(t, 1, len(dash.Pages[0].Widgets))
	require.Equal(t, 1, len(dash.Variables))

	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Title, dash.Pages[0].Widgets[0].Title)
	assert.Equal(t, dashboardInput.Pages[0].Widgets[0].Configuration.Markdown.Text, dash.Pages[0].Widgets[0].Configuration.Markdown.Text)
	assert.Greater(t, len(dash.Pages[0].Widgets[0].RawConfiguration), 1)

	// Test: DashboardUpdate
	updatedDashboard := DashboardInput{
		Name:        dash.Name,
		Permissions: dash.Permissions,
		Pages: []DashboardPageInput{
			{
				Name: dash.Pages[0].Name,
				Widgets: []DashboardWidgetInput{
					{
						// Even though the config isn't changing, we have to pass it. 2021-01-11 JT
						Configuration: dashboardInput.Pages[0].Widgets[0].Configuration,
						ID:            dash.Pages[0].Widgets[0].ID,
						Title:         "Updated Title",
					},
				},
			},
		},
		Variables: []DashboardVariableInput{
			{
				DefaultValues: &[]DashboardVariableDefaultItemInput{
					{
						Value: DashboardVariableDefaultValueInput{
							String: "default value-updated",
						},
					},
				},
				IsMultiSelection: true,
				Items: []DashboardVariableEnumItemInput{
					{
						Title: "VARIABLE",
						Value: "Variable",
					},
				},
				Name: "variableUpdated",
				Options: DashboardVariableOptionsInput{
					IgnoreTimeRange: true,
				},
				NRQLQuery: &DashboardVariableNRQLQueryInput{
					AccountIDs: []int{testAccountID},
					Query:      nrdb.NRQL("SELECT average(duration) FROM Transaction TIMESERIES WHERE appName = 'Dummy App'"),
				},
				ReplacementStrategy: DashboardVariableReplacementStrategyTypes.DEFAULT,
				Title:               "variable title",
				Type:                DashboardVariableTypeTypes.NRQL,
			},
		},
	}

	upDash, err := client.DashboardUpdate(updatedDashboard, dashGUID)
	require.NoError(t, err)
	require.NotNil(t, upDash)

	require.Equal(t, 1, len(upDash.EntityResult.Pages))
	require.Equal(t, 1, len(upDash.EntityResult.Pages[0].Widgets))
	assert.Equal(t, updatedDashboard.Pages[0].Widgets[0].Title, upDash.EntityResult.Pages[0].Widgets[0].Title)

	// Test: DashboardDelete
	delRes, err := client.DashboardDelete(dashGUID)
	require.NoError(t, err)
	require.NotNil(t, delRes)
	assert.Equal(t, 0, len(delRes.Errors))
	assert.Equal(t, DashboardDeleteResultStatusTypes.SUCCESS, delRes.Status)
}
