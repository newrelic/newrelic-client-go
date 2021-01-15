// +build integration

package dashboards

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/entities"
	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Dashboards {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}

func TestIntegrationDashboard_Basic(t *testing.T) {
	t.Parallel()

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
	result, err := client.DashboardCreate(mock.TestAccountID, dashboardInput)

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
