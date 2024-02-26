//go:build integration
// +build integration

package servicelevel

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v3/pkg/common"
	"github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Servicelevel {
	tc := testhelpers.NewIntegrationTestConfig(t)
	return New(tc)
}

func newServiceLevelIndicatorCreateInput(eventsInput ServiceLevelEventsCreateInput) ServiceLevelIndicatorCreateInput {
	return ServiceLevelIndicatorCreateInput{
		Name:        "integration-test-sli",
		Description: "Service level description",
		Events:      eventsInput,
		Objectives: []ServiceLevelObjectiveCreateInput{
			{
				Name:        "intgration-test-sli-objective",
				Description: "Objective description",
				Target:      99.9,
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

func TestServiceLevel_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// GUID of Dummy App
	guid := common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)

	eventsInput := ServiceLevelEventsCreateInput{
		AccountID: testAccountID,
		ValidEvents: &ServiceLevelEventsQueryCreateInput{
			From:  NRQL("Transaction"),
			Where: NRQL("appName = 'foo'"),
		},
		BadEvents: &ServiceLevelEventsQueryCreateInput{
			From:  NRQL("Transaction"),
			Where: NRQL("appName = 'foo' AND duration > 2"),
		},
	}
	createInput := newServiceLevelIndicatorCreateInput(eventsInput)

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

func TestServiceLevel_CDF(t *testing.T) {
	t.Skip("Skipping due to new schema fields being behind a feature flag until the feature is generally available to all customers.")
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// GUID of Dummy App
	guid := common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)

	eventsInput := ServiceLevelEventsCreateInput{
		AccountID: testAccountID,
		ValidEvents: &ServiceLevelEventsQueryCreateInput{
			Select: &ServiceLevelEventsQuerySelectCreateInput{
				Attribute: "some.distributed.attribute",
				Function:  ServiceLevelEventsQuerySelectFunction(ServiceLevelEventsQuerySelectFunctionTypes.GET_FIELD),
			},
			From:  NRQL("Metric"),
			Where: NRQL("appName = 'foo'"),
		},
		BadEvents: &ServiceLevelEventsQueryCreateInput{
			Select: &ServiceLevelEventsQuerySelectCreateInput{
				Attribute: "some.distributed.attribute",
				Function:  ServiceLevelEventsQuerySelectFunction(ServiceLevelEventsQuerySelectFunctionTypes.GET_CDF_COUNT),
				Threshold: 2.5,
			},
			From:  NRQL("Metric"),
			Where: NRQL("appName = 'foo'"),
		},
	}
	createInput := newServiceLevelIndicatorCreateInput(eventsInput)

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

	// Delete secure credential
	deleteResp, err := client.ServiceLevelDelete(createResp.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
}

func TestServiceLevel_GoodOrBadEventsRequiredError(t *testing.T) {
	t.Parallel()

	testAccountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// GUID of Dummy App
	guid := common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)

	// This input is missing some required fields to create a service level (i.e. GoodEvents, BadEvents)
	eventsInput := ServiceLevelEventsCreateInput{
		AccountID: testAccountID,
		ValidEvents: &ServiceLevelEventsQueryCreateInput{
			From: NRQL("Transaction"),
		},
	}
	createInput := newServiceLevelIndicatorCreateInput(eventsInput)

	// Create
	createResp, err := client.ServiceLevelCreate(guid, createInput)
	require.Nil(t, createResp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Defining a new SLI requires a good or bad events query")
}
