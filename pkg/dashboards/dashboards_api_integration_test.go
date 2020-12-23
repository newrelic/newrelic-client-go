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
							Markdown: DashboardMarkdownWidgetConfigurationInput{
								Text: "Test Text widget **markdown**",
							},
						},
					},
				},
			},
		},
	}

	// Test: Create
	result, err := client.DashboardCreate(mock.TestAccountID, dashboardInput)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, len(result.Errors))
	require.NotNil(t, result.EntityResult.GUID)

	dashGUID := result.EntityResult.GUID

	// Test: Delete
	delRes, err := client.DashboardDelete(dashGUID)
	require.NoError(t, err)
	require.NotNil(t, delRes)
	assert.Equal(t, 0, len(delRes.Errors))
	assert.Equal(t, DashboardDeleteResultStatusTypes.SUCCESS, delRes.Status)
}
