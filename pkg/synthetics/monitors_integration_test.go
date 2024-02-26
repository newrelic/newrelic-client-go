//go:build integration
// +build integration

package synthetics

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

var (
	testIntegrationMonitor = Monitor{
		Type:         MonitorTypes.Ping,
		Frequency:    15,
		URI:          "https://google.com",
		Locations:    []string{"AWS_US_EAST_1"},
		Status:       MonitorStatus.Disabled,
		SLAThreshold: 7,
		UserID:       0,
		APIVersion:   "LATEST",
		Options:      MonitorOptions{},
	}
)

func TestIntegrationMonitors(t *testing.T) {
	t.Skipf("Synthetics REST API is deprecated")

	tc := mock.NewIntegrationTestConfig(t)

	synthetics := New(tc)

	testIntegrationMonitor.Name = generateSyntheticsEntityNameForIntegrationTest("MONITOR", false)

	// Test: Create
	created, err := synthetics.CreateMonitor(testIntegrationMonitor)

	require.NoError(t, err)
	require.NotNil(t, created)

	// Test: List
	monitors, err := synthetics.ListMonitors()

	require.NoError(t, err)
	require.NotNil(t, monitors)
	require.Greater(t, len(monitors), 0)

	// Test: Get
	monitorID := created.ID
	monitor, err := synthetics.GetMonitor(monitorID)

	require.NoError(t, err)
	require.NotNil(t, *monitor)

	// Test: Update
	updatedName := generateSyntheticsEntityNameForIntegrationTest("MONITOR", true)
	monitor.Name = updatedName
	updated, err := synthetics.UpdateMonitor(*monitor)

	require.NoError(t, err)
	require.NotNil(t, *updated)

	monitor, err = synthetics.GetMonitor(monitorID)

	require.NoError(t, err)
	require.Equal(t, updatedName, monitor.Name)

	// Test: Delete
	err = synthetics.DeleteMonitor(monitorID)

	require.NoError(t, err)
}
