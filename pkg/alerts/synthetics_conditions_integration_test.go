//go:build integration
// +build integration

package alerts

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationSyntheticsConditions(t *testing.T) {
	t.Parallel()

	tc := mock.NewIntegrationTestConfig(t)

	var (
		testRandStr                      = mock.RandSeq(5)
		testIntegrationSyntheticsMonitor = synthetics.Monitor{
			Name:         fmt.Sprintf("test-synthetics-alert-conditions-monitor-%s", testRandStr),
			Type:         synthetics.MonitorTypes.Ping,
			Frequency:    15,
			URI:          "https://google.com",
			Locations:    []string{"AWS_US_EAST_1"},
			Status:       synthetics.MonitorStatus.Enabled,
			SLAThreshold: 7,
			APIVersion:   "LATEST",
		}
		testIntegrationPolicy = Policy{
			Name:               fmt.Sprintf("test-synthetics-alert-conditions-policy-%s", testRandStr),
			IncidentPreference: "PER_POLICY",
		}
		testIntegrationSyntheticsCondition = SyntheticsCondition{
			Name: fmt.Sprintf("test-synthetics-alert-condition-%s", testRandStr),
		}
	)

	alerts := newIntegrationTestClient(t)
	synth := synthetics.New(tc)

	// Setup
	monitor, err := synth.CreateMonitor(testIntegrationSyntheticsMonitor)

	require.NoError(t, err)

	policy, err := alerts.CreatePolicy(testIntegrationPolicy)

	require.NoError(t, err)

	// Deferred Teardown
	defer func() {
		// Teardown
		_, err = alerts.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("Error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}

		err = synth.DeleteMonitor(monitor.ID)
		if err != nil {
			t.Logf("Error cleaning up synthetics monitor %s (%s): %s",
				monitor.ID, testIntegrationSyntheticsMonitor.Name, err)
		}
	}()

	// Test: Create
	testIntegrationSyntheticsCondition.MonitorID = monitor.ID

	// addition of a 10 second sleep before a Synthetics Condition is created,
	// since the API is currently flakey and is undergoing rerouting of traffic, causing the monitor to
	// not be found during creation, intermittently
	time.Sleep(10 * time.Second)
	created, err := alerts.CreateSyntheticsCondition(policy.ID, testIntegrationSyntheticsCondition)

	require.NoError(t, err)
	require.NotZero(t, created)

	// Test: List
	conditions, err := alerts.ListSyntheticsConditions(policy.ID)

	require.NoError(t, err)
	require.NotZero(t, conditions)

	// Test: Get
	condition, err := alerts.GetSyntheticsCondition(policy.ID, conditions[0].ID)

	require.NoError(t, err)
	require.NotZero(t, condition)

	// Test: Update
	condition.Name = fmt.Sprintf("test-synthetics-alert-condition-updated-%s", testRandStr)
	updated, err := alerts.UpdateSyntheticsCondition(*condition)

	require.NoError(t, err)
	require.NotZero(t, updated)

	// Test: Delete
	deleted, err := alerts.DeleteSyntheticsCondition(updated.ID)

	require.NoError(t, err)
	require.NotZero(t, deleted)
}
