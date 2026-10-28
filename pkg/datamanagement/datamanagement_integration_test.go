//go:build integration
// +build integration

package datamanagement

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// testCardinalityLimitName is the well-known name of the dimensional metric
// per-metric cardinality limit that every New Relic account exposes.
const testCardinalityLimitName = "Dimensional Metric per-metric cardinality ingested per day"

// TestIntegrationDataManagement_GetLimits verifies that GetLimits returns a
// non-empty list containing the expected cardinality limit entry with correct
// metadata fields.
func TestIntegrationDataManagement_GetLimits(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	limits, err := client.GetLimits(accountID)
	require.NoError(t, err)
	require.NotNil(t, limits)
	require.NotEmpty(t, *limits, "expected at least one limit for the account")

	found := findCardinalityLimit(*limits)
	require.NotNil(t, found, "expected %q to be present in account limits", testCardinalityLimitName)
	require.Greater(t, found.Value, 0, "limit value should be positive")
	require.Equal(t, DataManagementCategoryTypes.INGEST, found.Category, "category should be INGEST")
	require.Equal(t, DataManagementUnitTypes.COUNT, found.Unit, "unit should be COUNT")
	require.NotEmpty(t, found.TimeInterval, "timeInterval should be populated")
	require.NotEmpty(t, found.Description, "description should be populated")
}

// TestIntegrationDataManagement_CreateAccountLimit_Default verifies that an
// account-wide cardinality limit override can be set and read back, and that a
// subsequent call with a new value updates it (pure upsert semantics).
func TestIntegrationDataManagement_CreateAccountLimit_Default(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	originalValue, err := fetchCurrentLimitValue(client, accountID)
	require.NoError(t, err, "should be able to read current limit before test")

	defer func() {
		_, _ = client.DataManagementCreateAccountLimit(accountID, DataManagementAccountLimitInput{
			Limit:          DataManagementLimitLookupInput{Name: testCardinalityLimitName},
			OverrideValue:  originalValue,
			OverrideReason: "integration test cleanup: restore original account-wide limit",
		})
	}()

	newValue := originalValue + 50000

	// SET account-wide override.
	resp, err := client.DataManagementCreateAccountLimit(accountID, DataManagementAccountLimitInput{
		Limit:          DataManagementLimitLookupInput{Name: testCardinalityLimitName},
		OverrideValue:  newValue,
		OverrideReason: "integration test: set account-wide cardinality limit",
	})
	require.NoError(t, err, "create account limit should succeed")
	require.NotNil(t, resp)
	require.Equal(t, testCardinalityLimitName, resp.Name)
	require.Equal(t, newValue, resp.Value)
	require.Equal(t, DataManagementCategoryTypes.INGEST, resp.Category)
	require.Equal(t, DataManagementUnitTypes.COUNT, resp.Unit)

	// READ BACK to confirm persistence.
	readBack, err := fetchCurrentLimitValue(client, accountID)
	require.NoError(t, err)
	require.Equal(t, newValue, readBack, "read-back value should match what was just set")

	// UPDATE to a new value (same mutation — upsert semantics).
	updatedValue := newValue + 10000
	resp2, err := client.DataManagementCreateAccountLimit(accountID, DataManagementAccountLimitInput{
		Limit:          DataManagementLimitLookupInput{Name: testCardinalityLimitName},
		OverrideValue:  updatedValue,
		OverrideReason: "integration test: update account-wide cardinality limit",
	})
	require.NoError(t, err, "update should succeed")
	require.Equal(t, updatedValue, resp2.Value, "response should reflect updated value")

	readAfterUpdate, err := fetchCurrentLimitValue(client, accountID)
	require.NoError(t, err)
	require.Equal(t, updatedValue, readAfterUpdate, "read-back after update should match")
}

// TestIntegrationDataManagement_CreateAccountLimit_PerMetric verifies that a
// per-metric cardinality override (with a Qualifier) can be submitted without
// error. The per-metric value cannot be read back via the limits API, so this
// test validates the mutation response fields only.
func TestIntegrationDataManagement_CreateAccountLimit_PerMetric(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	rand := mock.RandSeq(5)
	metricName := "test.cardinality.integration." + rand

	originalValue, err := fetchCurrentLimitValue(client, accountID)
	require.NoError(t, err)

	defer func() {
		// Reset per-metric override back to account-wide default.
		_, _ = client.DataManagementCreateAccountLimit(accountID, DataManagementAccountLimitInput{
			Limit:          DataManagementLimitLookupInput{Name: testCardinalityLimitName},
			OverrideValue:  originalValue,
			OverrideReason: "integration test cleanup: reset per-metric override for " + metricName,
			Qualifier:      metricName,
		})
	}()

	perMetricValue := 75000

	resp, err := client.DataManagementCreateAccountLimit(accountID, DataManagementAccountLimitInput{
		Limit:          DataManagementLimitLookupInput{Name: testCardinalityLimitName},
		OverrideValue:  perMetricValue,
		OverrideReason: "integration test: per-metric cardinality override for " + metricName,
		Qualifier:      metricName,
	})
	require.NoError(t, err, "per-metric override should succeed")
	require.NotNil(t, resp)
	require.Equal(t, testCardinalityLimitName, resp.Name)
	require.Equal(t, perMetricValue, resp.Value)
	require.Equal(t, DataManagementCategoryTypes.INGEST, resp.Category)
	require.Equal(t, DataManagementUnitTypes.COUNT, resp.Unit)
}

// TestIntegrationDataManagement_CreateAccountLimit_InvalidLimitName verifies
// that submitting an unrecognised limit name returns an error from the API.
func TestIntegrationDataManagement_CreateAccountLimit_InvalidLimitName(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	_, err = client.DataManagementCreateAccountLimit(accountID, DataManagementAccountLimitInput{
		Limit:          DataManagementLimitLookupInput{Name: "nonexistent.limit.name.xyz"},
		OverrideValue:  100000,
		OverrideReason: "integration test: invalid limit name",
	})
	require.Error(t, err, "submitting an unknown limit name should return an error")
}

// findCardinalityLimit returns the first limit matching testCardinalityLimitName, or nil.
func findCardinalityLimit(limits []DataManagementAccountLimit) *DataManagementAccountLimit {
	for i := range limits {
		if limits[i].Name == testCardinalityLimitName {
			return &limits[i]
		}
	}
	return nil
}

// fetchCurrentLimitValue reads the current account-wide cardinality limit value.
func fetchCurrentLimitValue(client Datamanagement, accountID int) (int, error) {
	limits, err := client.GetLimits(accountID)
	if err != nil {
		return 0, err
	}
	l := findCardinalityLimit(*limits)
	if l == nil {
		return 0, fmt.Errorf("cardinality limit %q not found for account %d", testCardinalityLimitName, accountID)
	}
	return l.Value, nil
}

func newIntegrationTestClient(t *testing.T) Datamanagement {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
