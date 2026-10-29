//go:build integration
// +build integration

package servicelevel

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newCompositeComponentIndicatorCreateInput(accountID int, name string) ServiceLevelIndicatorCreateInput {
	return ServiceLevelIndicatorCreateInput{
		Name:        name,
		Description: "Composite integration test component",
		Events: ServiceLevelEventsCreateInput{
			AccountID: accountID,
			ValidEvents: &ServiceLevelEventsQueryCreateInput{
				From: NRQL("Transaction"),
			},
			BadEvents: &ServiceLevelEventsQueryCreateInput{
				From:  NRQL("Transaction"),
				Where: NRQL("duration > 2"),
			},
		},
		Objectives: []ServiceLevelObjectiveCreateInput{
			{
				Name:   name + "-objective",
				Target: 99.9,
				TimeWindow: ServiceLevelObjectiveTimeWindowCreateInput{
					Rolling: ServiceLevelObjectiveRollingTimeWindowCreateInput{
						Count: 7,
						Unit:  ServiceLevelObjectiveRollingTimeWindowUnit(ServiceLevelObjectiveRollingTimeWindowUnitTypes.DAY),
					},
				},
			},
		},
	}
}

func TestServiceLevelComposite_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	workloadGUID := common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew)

	// Create two component SLIs — composite requires at least two.
	component1, err := client.ServiceLevelCreate(workloadGUID, newCompositeComponentIndicatorCreateInput(testAccountID, "composite-integration-test-component-1"))
	require.NoError(t, err)
	require.NotNil(t, component1)

	component2, err := client.ServiceLevelCreate(workloadGUID, newCompositeComponentIndicatorCreateInput(testAccountID, "composite-integration-test-component-2"))
	require.NoError(t, err)
	require.NotNil(t, component2)

	componentsDeleted := false
	defer func() {
		if componentsDeleted {
			return
		}
		_, _ = client.ServiceLevelDelete(component1.GUID, component1.ID) // best-effort cleanup
		_, _ = client.ServiceLevelDelete(component2.GUID, component2.ID)
	}()

	createInput := ServiceLevelCompositeCreateInput{
		AccountID:   testAccountID,
		Name:        "composite-integration-test",
		Description: "Composite integration test",
		Components: []ServiceLevelComponentCreateInput{
			{GUID: component1.GUID},
			{GUID: component2.GUID},
		},
		Objectives: []ServiceLevelCompositeObjectiveCreateInput{
			{
				Name:   "composite-integration-test-objective",
				Target: 99,
				TimeWindow: ServiceLevelObjectiveTimeWindowCreateInput{
					Rolling: ServiceLevelObjectiveRollingTimeWindowCreateInput{
						Count: 7,
						Unit:  ServiceLevelObjectiveRollingTimeWindowUnit(ServiceLevelObjectiveRollingTimeWindowUnitTypes.DAY),
					},
				},
			},
		},
	}

	// Create
	createResp, err := client.ServiceLevelCompositeCreate(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.Equal(t, "composite-integration-test", createResp.Name)
	require.Len(t, createResp.Components, 2)
	require.Equal(t, ServiceLevelWeightTypeTypes.AVERAGE, createResp.WeightType) // no weights given -> server defaults to equal distribution

	compositeDeleted := false
	defer func() {
		if compositeDeleted {
			return
		}
		_, _ = client.ServiceLevelCompositeDelete(createResp.GUID) // best-effort cleanup
	}()

	fmt.Println("waiting 5 seconds for entity to be indexed before validating its creation...")
	time.Sleep(5 * time.Second)

	// Get
	getResp, err := client.GetComposites([]common.EntityGUID{createResp.GUID})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Len(t, *getResp, 1)

	// Update
	updateInput := ServiceLevelCompositeUpdateInput{
		Description: "composite integration test updated",
	}
	updateResp, err := client.ServiceLevelCompositeUpdate(updateInput, createResp.GUID)
	require.NoError(t, err)
	require.NotNil(t, updateResp)

	// Delete composite first — components may not be deletable while referenced.
	deleteResp, err := client.ServiceLevelCompositeDelete(createResp.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
	compositeDeleted = true

	// Delete the two component SLIs now that nothing references them.
	_, err = client.ServiceLevelDelete(component1.GUID, component1.ID)
	require.NoError(t, err)
	_, err = client.ServiceLevelDelete(component2.GUID, component2.ID)
	require.NoError(t, err)
	componentsDeleted = true
}

func TestServiceLevelComposite_MinimumComponentsError(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	workloadGUID := common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew)

	// A single real component — not zero — deliberately, since this struct's
	// `Components` field is `json:"components,omitempty"`, so an empty or nil
	// slice gets dropped from the request entirely and fails GraphQL's non-null
	// constraint on the argument before ever reaching the actual "at least 2
	// components" business check. One component is below the minimum (2) but
	// still a non-empty slice, so it serializes correctly and isolates the test
	// to the count rule specifically rather than entity-resolution or GraphQL
	// null-argument errors.
	onlyComponent, err := client.ServiceLevelCreate(workloadGUID, newCompositeComponentIndicatorCreateInput(testAccountID, "composite-integration-test-min-components-component"))
	require.NoError(t, err)
	require.NotNil(t, onlyComponent)
	defer func() {
		_, _ = client.ServiceLevelDelete(onlyComponent.GUID, onlyComponent.ID) // best-effort cleanup
	}()

	createInput := ServiceLevelCompositeCreateInput{
		AccountID: testAccountID,
		Name:      "composite-integration-test-min-components",
		Components: []ServiceLevelComponentCreateInput{
			{GUID: onlyComponent.GUID},
		},
		Objectives: []ServiceLevelCompositeObjectiveCreateInput{
			{
				Target: 99,
				TimeWindow: ServiceLevelObjectiveTimeWindowCreateInput{
					Rolling: ServiceLevelObjectiveRollingTimeWindowCreateInput{
						Count: 7,
						Unit:  ServiceLevelObjectiveRollingTimeWindowUnit(ServiceLevelObjectiveRollingTimeWindowUnitTypes.DAY),
					},
				},
			},
		},
	}

	createResp, err := client.ServiceLevelCompositeCreate(createInput)
	require.Nil(t, createResp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must include at least two component SLIs")
}
