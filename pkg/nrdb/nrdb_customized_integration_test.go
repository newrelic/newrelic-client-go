//go:build integration
// +build integration

package nrdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var accountID = 3957524

// TestIntegrationPerformNRQLQuerySimple tests a simple query with no FACET or TIMESERIES
func TestIntegrationPerformNRQLQuerySimple(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction SINCE 1 day ago"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For simple query, we expect otherResult and totalResult to be arrays with a single element
	require.NotNil(t, res.OtherResult, "OtherResult should not be nil")
	require.NotNil(t, res.TotalResult, "TotalResult should not be nil")

	// Arrays should be non-nil, length may be 0 if no otherResult
	assert.IsType(t, NRDBMultiResultCustomized{}, res.OtherResult, "OtherResult should always be a NRDBMultiResultCustomized")
	assert.IsType(t, NRDBMultiResultCustomized{}, res.TotalResult, "TotalResult should always be a NRDBMultiResultCustomized")

	// Results should contain at least one item
	require.GreaterOrEqual(t, len(res.Results), 1)
}

// TestIntegrationPerformNRQLQueryFacet tests a query with FACET
func TestIntegrationPerformNRQLQueryFacet(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction FACET appName SINCE 1 day ago LIMIT 5"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For FACET query, we expect otherResult and totalResult to be arrays with single elements
	require.NotNil(t, res.OtherResult)
	require.IsType(t, NRDBMultiResultCustomized{}, res.OtherResult)

	if len(res.OtherResult) > 0 {
		// Check that we can directly access first element since we know it's an array
		if count, ok := res.OtherResult[0]["count"]; ok {
			require.IsType(t, float64(0), count, "Expected 'count' to be a numerical value")
		}
	}

	// Results should contain multiple items due to FACET
	require.GreaterOrEqual(t, len(res.Results), 1)
}

// TestIntegrationPerformNRQLQueryTimeseries tests a query with TIMESERIES
func TestIntegrationPerformNRQLQueryTimeseries(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction TIMESERIES 1 hour SINCE 1 day ago"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For TIMESERIES query without FACET, we expect otherResult and totalResult to be arrays with single elements
	assert.IsType(t, NRDBMultiResultCustomized{}, res.OtherResult)
	assert.IsType(t, NRDBMultiResultCustomized{}, res.TotalResult)

	// Results should contain multiple items for each time window
	require.GreaterOrEqual(t, len(res.Results), 1)
}

// TestIntegrationPerformNRQLQueryFacetTimeseries tests a query with both FACET and TIMESERIES
func TestIntegrationPerformNRQLQueryFacetTimeseries(t *testing.T) {
	t.Parallel()

	facetTimeseriesQueries := []string{
		"SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago LIMIT 3",
		"SELECT average(duration) FROM Transaction FACET appName TIMESERIES 30 minutes SINCE 1 day ago LIMIT 3",
		"SELECT count(*) as 'count', average(duration) as 'avgDuration', max(duration) as 'maxDuration' FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago LIMIT 3",
	}

	client := newNRDBIntegrationTestClient(t)

	for _, query := range facetTimeseriesQueries {
		t.Run(query, func(t *testing.T) {
			res, err := client.PerformNRQLQuery(accountID, NRQL(query))

			require.NoError(t, err)
			require.NotNil(t, res)

			// For FACET + TIMESERIES queries, we expect otherResult and totalResult to be arrays with multiple elements
			assert.IsType(t, NRDBMultiResultCustomized{}, res.OtherResult)
			assert.IsType(t, NRDBMultiResultCustomized{}, res.TotalResult)

			// Arrays should have elements
			require.GreaterOrEqual(t, len(res.OtherResult), 1, "Expected otherResults to have at least one element")
			require.GreaterOrEqual(t, len(res.TotalResult), 1, "Expected totalResults to have at least one element")

			// Results should contain multiple items
			require.GreaterOrEqual(t, len(res.Results), 1)
		})
	}
}

// TestIntegrationPerformNRQLQueryMultipleFacets tests a query with multiple FACETs
func TestIntegrationPerformNRQLQueryMultipleFacets(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction FACET appName, name SINCE 1 day ago LIMIT 5"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For multiple FACET query, we expect otherResult and totalResult to be arrays (with single elements)
	assert.IsType(t, NRDBMultiResultCustomized{}, res.OtherResult, "OtherResult should always be a NRDBMultiResultCustomized")
	assert.IsType(t, NRDBMultiResultCustomized{}, res.TotalResult, "TotalResult should always be a NRDBMultiResultCustomized")

	// Since this is a FACET query without TIMESERIES, likely just one element in each array
	if len(res.OtherResult) > 0 {
		// We can now directly access the first element
		require.IsType(t, NRDBResult{}, res.OtherResult[0], "Array element should be NRDBResult type")
	}

	// Results should contain multiple items due to multiple FACETs
	require.GreaterOrEqual(t, len(res.Results), 1)
}

// TestIntegrationPerformNRQLQueryTimeseriesCompare tests a query with TIMESERIES and COMPARE WITH
func TestIntegrationPerformNRQLQueryTimeseriesCompare(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction TIMESERIES 1 hour COMPARE WITH 1 week ago SINCE 1 day ago"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For TIMESERIES with COMPARE WITH, we expect otherResult and totalResult to be arrays (with single elements)
	// and currentResults and previousResults to be populated
	assert.IsType(t, NRDBMultiResultCustomized{}, res.OtherResult, "OtherResult should always be a NRDBMultiResultCustomized")
	assert.IsType(t, NRDBMultiResultCustomized{}, res.TotalResult, "TotalResult should always be a NRDBMultiResultCustomized")

	// Both currentResults and previousResults should be populated
	require.GreaterOrEqual(t, len(res.CurrentResults), 1)
	require.GreaterOrEqual(t, len(res.PreviousResults), 1)
}
