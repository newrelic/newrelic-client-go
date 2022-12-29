//go:build integration
// +build integration

package servicelevel

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestServiceLevel_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// GUID of Dummy App
	guid := common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")

	createInput := ServiceLevelIndicatorCreateInput{
		Name:        "integration-test-sli",
		Description: "integration test service level",
		Events: ServiceLevelEventsCreateInput{
			AccountID: testAccountID,
			ValidEvents: &ServiceLevelEventsQueryCreateInput{
				From: NRQL("Transaction"),
			},
			BadEvents: &ServiceLevelEventsQueryCreateInput{
				From: NRQL("Transaction"),
			},
		},
		Objectives: []ServiceLevelObjectiveCreateInput{
			{
				Name:        "intgration-test-sli-objective",
				Description: "testing",
				Target:      1.1,
				TimeWindow: ServiceLevelObjectiveTimeWindowCreateInput{
					Rolling: ServiceLevelObjectiveRollingTimeWindowCreateInput{
						Count: 1,
						Unit:  ServiceLevelObjectiveRollingTimeWindowUnit(ServiceLevelObjectiveRollingTimeWindowUnitTypes.DAY),
					},
				},
			},
		},
	}

	// Create
	createResp, err := client.ServiceLevelCreate(guid, createInput)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	fmt.Println("waiting 5 seconds for entity to be indexed before validating its creation...")
	time.Sleep(5 * time.Second)

	// Get
	getResp, err := client.GetIndicators(createResp.GUID)
	require.NoError(t, err)
	require.NotNil(t, getResp)

	// Update
	updateInput := ServiceLevelIndicatorUpdateInput{
		Description: "integration test service level updated",
	}
	updateResp, err := client.ServiceLevelUpdate(createResp.GUID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updateResp)

	// Delete secure credential
	deleteResp, err := client.ServiceLevelDelete(createResp.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
}

func TestServiceLevel_GoodOrBadEventsRequiredError(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// GUID of Dummy App
	guid := common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")

	// This input is missing some required fields to create a service level (i.e. GoodEvents, BadEvents)
	createInput := ServiceLevelIndicatorCreateInput{
		Name:        "integration-test-sli",
		Description: "integration test service level",
		Events: ServiceLevelEventsCreateInput{
			AccountID: testAccountID,
			ValidEvents: &ServiceLevelEventsQueryCreateInput{
				From: NRQL("Transaction"),
			},
		},
		Objectives: []ServiceLevelObjectiveCreateInput{
			{
				Name:        "intgration-test-sli-objective",
				Description: "testing",
				Target:      1.1,
				TimeWindow: ServiceLevelObjectiveTimeWindowCreateInput{
					Rolling: ServiceLevelObjectiveRollingTimeWindowCreateInput{
						Count: 1,
						Unit:  ServiceLevelObjectiveRollingTimeWindowUnit(ServiceLevelObjectiveRollingTimeWindowUnitTypes.DAY),
					},
				},
			},
		},
	}

	// Create
	createResp, err := client.ServiceLevelCreate(guid, createInput)
	require.Nil(t, createResp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Defining a new SLI requires a good or bad events query")
}

func newIntegrationTestClient(t *testing.T) Servicelevel {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
