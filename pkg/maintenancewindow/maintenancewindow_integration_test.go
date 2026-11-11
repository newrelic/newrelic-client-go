//go:build integration
// +build integration

package maintenancewindow

import (
	"fmt"
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

func TestMaintenanceWindow_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	slClient := servicelevel.New(testhelpers.NewIntegrationTestConfig(t))

	// A maintenance window requires at least one real affected entity, so a
	// genuine SLI is created here (attached to the dummy test app) purely to
	// obtain a real GUID for affected_entities.
	appGUID := common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew)
	sliResp, err := slClient.ServiceLevelCreate(appGUID, servicelevel.ServiceLevelIndicatorCreateInput{
		Name:        "integration-test-mw-sli",
		Description: "Created for a maintenance window integration test",
		Events: servicelevel.ServiceLevelEventsCreateInput{
			AccountID: testAccountID,
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
				Name:   "integration-test-mw-sli-objective",
				Target: 99.9,
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
	require.NotNil(t, sliResp)

	sliDeleted := false
	defer func() {
		if sliDeleted {
			return
		}
		_, _ = slClient.ServiceLevelDelete(sliResp.GUID, sliResp.ID) // best-effort cleanup
	}()

	createInput := MaintenanceWindowConfigCreateInput{
		Name:        "integration-test-maintenance-window",
		Description: "Created by an integration test",
		Scope: MaintenanceWindowScopedReferenceInput{
			ID:   fmt.Sprintf("%d", testAccountID),
			Type: MaintenanceWindowEntityScopeTypes.ACCOUNT,
		},
		StartTime:          NaiveDateTime(time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05")),
		Duration:           workloads.Duration("PT1H"),
		Timezone:           "UTC",
		AffectedEntityType: MaintenanceWindowAffectedEntityTypeTypes.SERVICE_LEVEL,
		AffectedEntities:   []string{string(sliResp.GUID)},
	}

	// Create
	createResp, err := client.MaintenanceWindowCreate(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.Equal(t, createInput.Name, createResp.Name)
	require.Equal(t, MaintenanceWindowEntityScopeTypes.ACCOUNT, createResp.Scope.Type)
	require.Equal(t, []string{string(sliResp.GUID)}, createResp.AffectedEntities)

	deleted := false
	defer func() {
		if deleted {
			return
		}
		_, _ = client.MaintenanceWindowDelete(createResp.ID) // best-effort cleanup
	}()

	// Get by ID
	listResp, err := client.GetListByIds([]string{createResp.ID})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.MaintenanceWindows, 1)
	require.Equal(t, createResp.ID, listResp.MaintenanceWindows[0].ID)

	// Update
	updateInput := MaintenanceWindowConfigUpdateInput{
		Description: "Updated by an integration test",
	}
	updateResp, err := client.MaintenanceWindowUpdate(createResp.ID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, updateInput.Description, updateResp.Description)

	// Delete
	deleteResp, err := client.MaintenanceWindowDelete(createResp.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
	deleted = true

	// Clean up the SLI created for this test.
	_, err = slClient.ServiceLevelDelete(sliResp.GUID, sliResp.ID)
	require.NoError(t, err)
	sliDeleted = true
}
