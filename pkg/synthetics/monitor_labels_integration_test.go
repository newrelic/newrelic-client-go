// +build integration

package synthetics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

var (
	testIntegrationLabeledMonitor = Monitor{
		Type:         MonitorTypes.APITest,
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

func TestIntegrationGetMonitorLabels(t *testing.T) {
	t.Parallel()

	tc := mock.NewIntegrationTestConfig(t)

	synthetics := New(tc)

	// Setup
	rand := mock.RandSeq(5)
	testIntegrationLabeledMonitor.Name = fmt.Sprintf("test-synthetics-monitor-%s", rand)
	monitor, err := synthetics.CreateMonitor(testIntegrationLabeledMonitor)
	require.NoError(t, err)

	// Test: Get
	_, err = synthetics.GetMonitorLabels(monitor.ID)
	require.NoError(t, err)

	// Deferred teardown
	defer func() {
		err = synthetics.DeleteMonitor(monitor.ID)

		if err != nil {
			t.Logf("error cleaning up monitor %s (%s): %s", monitor.ID, monitor.Name, err)
		}
	}()
}
