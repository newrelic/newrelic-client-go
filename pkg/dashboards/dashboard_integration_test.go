// +build integration

package dashboards

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDashboards(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEW_RELIC_ADMIN_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires NEW_RELIC_ADMIN_API_KEY to be set")
	}

	dashboards := New(config.Config{
		AdminAPIKey: apiKey,
		LogLevel:    "debug",
	})

	d := Dashboard{
		Metadata: DashboardMetadata{
			Version: 1,
		},
		Title:      "newrelic-client-go-test",
		Visibility: VisibilityTypes.Owner,
		Editable:   EditableTypes.Owner,
	}

	// Test: Create
	created, err := dashboards.CreateDashboard(d)

	require.NoError(t, err)
	require.NotNil(t, created)

	// Test: List
	params := ListDashboardsParams{
		Title: "newrelic-client-go",
	}
	multiple, err := dashboards.ListDashboards(&params)

	require.NoError(t, err)
	require.NotNil(t, multiple)

	// Test: Get
	single, err := dashboards.GetDashboard(multiple[0].ID)

	require.NoError(t, err)
	require.NotNil(t, single)

	// Test: Update
	single.Title = "updated"
	updated, err := dashboards.UpdateDashboard(*single)

	require.NoError(t, err)
	require.NotNil(t, updated)

	// Test: Delete
	deleted, err := dashboards.DeleteDashboard(updated.ID)

	require.NoError(t, err)
	require.NotNil(t, deleted)
}
