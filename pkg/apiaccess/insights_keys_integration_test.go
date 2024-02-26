package apiaccess

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

func TestIntegrationAPIAccess_InsightsInsertKeys(t *testing.T) {
	t.Parallel()
	t.Skip("For manual use only.")

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test: Create
	createResult, err := client.CreateInsightsInsertKey(testAccountID)
	require.NoError(t, err)
	require.NotZero(t, createResult)

	// Test: List
	listResult, err := client.ListInsightsInsertKeys(testAccountID)
	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	getResult, err := client.GetInsightsInsertKey(testAccountID, listResult[0].ID)
	require.NoError(t, err)
	require.NotZero(t, getResult)

	// Test: Delete
	updateResult, err := client.DeleteInsightsInsertKey(testAccountID, getResult.ID)
	require.NoError(t, err)
	require.NotNil(t, updateResult)

}

func TestIntegrationAPIAccess_InsightsQueryKeys(t *testing.T) {
	t.Parallel()
	t.Skip("For manual use only.")

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test: Create
	createResult, err := client.CreateInsightsQueryKey(testAccountID)
	require.NoError(t, err)
	require.NotZero(t, createResult)

	// Test: List
	listResult, err := client.ListInsightsQueryKeys(testAccountID)
	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	getResult, err := client.GetInsightsQueryKey(testAccountID, listResult[0].ID)
	require.NoError(t, err)
	require.NotZero(t, getResult)

	// Test: Delete
	deleteResult, err := client.DeleteInsightsQueryKey(testAccountID, getResult.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

// nolint: unused
func newIntegrationTestClient(t *testing.T) APIAccess {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
