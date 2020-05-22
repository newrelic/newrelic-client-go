package edge

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

var (
	testAccountIDStr      = os.Getenv("NEW_RELIC_MASTER_ACCOUNT_ID")
	testTraceObserverName = "testTraceObserver"
	testProviderRegion    = EdgeProviderRegionTypes.AWS_US_EAST_1
)

func TestIntegrationTraceObserver(t *testing.T) {
	if testAccountIDStr == "" {
		t.Skip("a master account is required to run trace observer integration tests")
		return
	}

	testAccountID, err := strconv.Atoi(testAccountIDStr)
	if err != nil {
		t.Fatal(err)
	}

	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.CreateTraceObserver(testAccountID, testTraceObserverName, testProviderRegion)

	require.NoError(t, err)
	require.NotNil(t, created)

	// Test: List
	traceObservers, err := client.ListTraceObservers(testAccountID)
	require.NoError(t, err)
	require.Greater(t, len(traceObservers), 0)

	// Test: Delete
	deleted, err := client.DeleteTraceObserver(testAccountID, created.ID)

	require.NoError(t, err)
	require.NotNil(t, deleted)
}

func newIntegrationTestClient(t *testing.T) Edge {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
