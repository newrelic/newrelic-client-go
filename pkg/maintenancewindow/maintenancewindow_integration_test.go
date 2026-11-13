//go:build integration
// +build integration

package maintenancewindow

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicelevel"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
)

func newIntegrationTestClient(t *testing.T) Maintenancewindow {
	tc := testhelpers.NewIntegrationTestConfig(t)
	return New(tc)
}

// createTestServiceLevel creates a service level indicator to act as the affected entity of a
// maintenance window, and returns its entity GUID. It is deleted when the test completes.
func createTestServiceLevel(t *testing.T, accountID int) string {
	tc := testhelpers.NewIntegrationTestConfig(t)
	slClient := servicelevel.New(tc)
	wlClient := workloads.New(tc)

	// Create a temporary workload in the test account to own the service level, so the test
	// works in any environment without a hardcoded entity GUID.
	workloadResp, err := wlClient.WorkloadCreate(accountID, workloads.WorkloadCreateInput{
		Name: fmt.Sprintf("integration-test-mw-workload-%s", testhelpers.RandSeq(5)),
		EntitySearchQueries: []workloads.WorkloadEntitySearchQueryInput{
			{Query: fmt.Sprintf("accountId = %d", accountID)},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, workloadResp)

	t.Cleanup(func() {
		_, _ = wlClient.WorkloadDelete(workloadResp.GUID) // best-effort cleanup
	})

	fmt.Println("waiting 20 seconds for workload to be indexed before attaching a service level to it...")
	time.Sleep(20 * time.Second)

	createResp, err := slClient.ServiceLevelCreate(common.EntityGUID(workloadResp.GUID), servicelevel.ServiceLevelIndicatorCreateInput{
		Name:        "integration-test-mw-sli",
		Description: "Service level used by maintenance window integration test",
		Events: servicelevel.ServiceLevelEventsCreateInput{
			AccountID: accountID,
			ValidEvents: &servicelevel.ServiceLevelEventsQueryCreateInput{
				From: servicelevel.NRQL("Transaction"),
			},
			BadEvents: &servicelevel.ServiceLevelEventsQueryCreateInput{
				From:  servicelevel.NRQL("Transaction"),
				Where: servicelevel.NRQL("duration > 2"),
			},
		},
		Objectives: []servicelevel.ServiceLevelObjectiveCreateInput{
			{
				Name:        "integration-test-mw-sli-objective",
				Description: "Objective description",
				Target:      99.9,
				TimeWindow: servicelevel.ServiceLevelObjectiveTimeWindowCreateInput{
					Rolling: servicelevel.ServiceLevelObjectiveRollingTimeWindowCreateInput{
						Count: 7,
						Unit:  servicelevel.ServiceLevelObjectiveRollingTimeWindowUnit(servicelevel.ServiceLevelObjectiveRollingTimeWindowUnitTypes.DAY),
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp)

	t.Cleanup(func() {
		_, _ = slClient.ServiceLevelDelete(createResp.GUID, createResp.ID) // best-effort cleanup
	})

	fmt.Println("waiting 20 seconds for entity to be indexed before using it in a maintenance window...")
	time.Sleep(20 * time.Second)

	return string(createResp.GUID)
}

func newMaintenanceWindowCreateInput(accountID int, affectedEntities []string, rrule string) MaintenanceWindowConfigCreateInput {
	return MaintenanceWindowConfigCreateInput{
		Name:               "integration-test-maintenance-window",
		Description:        "Maintenance window created by integration test",
		AffectedEntityType: MaintenanceWindowAffectedEntityTypeTypes.SERVICE_LEVEL,
		AffectedEntities:   affectedEntities,
		Duration:           "PT1H",
		Rrule:              rrule,
		StartTime:          NaiveDateTime(time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05")),
		Timezone:           "America/Los_Angeles",
		Scope: MaintenanceWindowScopedReferenceInput{
			ID:   strconv.Itoa(accountID),
			Type: MaintenanceWindowEntityScopeTypes.ACCOUNT,
		},
	}
}

func TestMaintenanceWindow_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	slGUID := createTestServiceLevel(t, testAccountID)

	// Create
	createInput := newMaintenanceWindowCreateInput(testAccountID, []string{slGUID}, "")
	createResp, err := client.MaintenanceWindowCreate(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotEmpty(t, createResp.ID)
	require.Equal(t, createInput.Name, createResp.Name)
	require.Equal(t, createInput.Description, createResp.Description)
	require.Equal(t, createInput.Duration, createResp.Duration)
	require.Equal(t, createInput.Timezone, createResp.Timezone)
	require.Equal(t, MaintenanceWindowAffectedEntityTypeTypes.SERVICE_LEVEL, createResp.AffectedEntityType)
	require.Equal(t, []string{slGUID}, createResp.AffectedEntities)
	require.Equal(t, createInput.Scope.ID, createResp.Scope.ID)
	require.Equal(t, MaintenanceWindowEntityScopeTypes.ACCOUNT, createResp.Scope.Type)

	deleted := false
	defer func() {
		if deleted {
			return
		}
		_, _ = client.MaintenanceWindowDelete(createResp.ID) // best-effort cleanup
	}()

	fmt.Println("waiting 30 seconds for maintenance window to be indexed before validating its creation...")
	time.Sleep(30 * time.Second)

	// Get by ID
	getResp, err := client.GetListByIds([]string{createResp.ID})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Len(t, getResp.MaintenanceWindows, 1)
	require.Equal(t, createResp.ID, getResp.MaintenanceWindows[0].ID)
	require.Equal(t, createInput.Name, getResp.MaintenanceWindows[0].Name)

	// Get by affected entity ID
	listResp, err := client.GetListByAffectedEntityId(slGUID)
	require.NoError(t, err)
	require.NotNil(t, listResp)
	found := false
	for _, mw := range listResp.MaintenanceWindows {
		if mw.ID == createResp.ID {
			found = true
			break
		}
	}
	require.True(t, found, "expected maintenance window %s to be listed for affected entity %s", createResp.ID, slGUID)

	// Update
	updateInput := MaintenanceWindowConfigUpdateInput{
		Name:        "integration-test-maintenance-window-updated",
		Description: "Maintenance window updated by integration test",
		Duration:    "PT2H",
	}
	updateResp, err := client.MaintenanceWindowUpdate(createResp.ID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, createResp.ID, updateResp.ID)
	require.Equal(t, updateInput.Name, updateResp.Name)
	require.Equal(t, updateInput.Description, updateResp.Description)
	require.Equal(t, updateInput.Duration, updateResp.Duration)

	// Delete
	deleteResp, err := client.MaintenanceWindowDelete(createResp.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
	require.Equal(t, createResp.ID, deleteResp.ID)
	deleted = true
}

func TestMaintenanceWindow_Recurring(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	slGUID := createTestServiceLevel(t, testAccountID)

	// Create
	createInput := newMaintenanceWindowCreateInput(testAccountID, []string{slGUID}, "FREQ=WEEKLY;BYDAY=SA")
	createResp, err := client.MaintenanceWindowCreate(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.Equal(t, createInput.Rrule, createResp.Rrule)

	deleted := false
	defer func() {
		if deleted {
			return
		}
		_, _ = client.MaintenanceWindowDelete(createResp.ID) // best-effort cleanup
	}()

	fmt.Println("waiting 5 seconds for maintenance window to be indexed before updating it...")
	time.Sleep(5 * time.Second)

	// Update — change the recurrence rule
	updateInput := MaintenanceWindowConfigUpdateInput{
		Rrule: "FREQ=DAILY",
	}
	updateResp, err := client.MaintenanceWindowUpdate(createResp.ID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, updateInput.Rrule, updateResp.Rrule)

	// Delete
	deleteResp, err := client.MaintenanceWindowDelete(createResp.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
	deleted = true
}

func TestMaintenanceWindow_NonExistentID(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	getResp, err := client.GetListByIds([]string{"00000000-0000-0000-0000-000000000000"})
	if err == nil {
		require.NotNil(t, getResp)
		require.Empty(t, getResp.MaintenanceWindows)
	}
}
