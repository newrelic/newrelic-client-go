// +build integration

package dashboards

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
)

func TestIntegrationDashboards(t *testing.T) {
	t.Parallel()

	tc := mock.NewIntegrationTestConfig(t)

	dashboards := New(tc)

	d := Dashboard{
		Metadata: DashboardMetadata{
			Version: 1,
		},
		Title:           "newrelic-client-go-test",
		Visibility:      VisibilityTypes.Owner,
		Editable:        EditableTypes.Owner,
		GridColumnCount: 3,
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
