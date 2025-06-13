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

	// For simple query, we expect otherResult and totalResult to be single objects
	_, isSingleOtherResult := res.OtherResult.(NRDBResult)
	_, isSingleTotalResult := res.TotalResult.(NRDBResult)

	assert.True(t, isSingleOtherResult, "Expected otherResult to be a single NRDBResult")
	assert.True(t, isSingleTotalResult, "Expected totalResult to be a single NRDBResult")

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

	// For FACET query, we expect otherResult and totalResult to be single objects
	singleOtherResult, isSingleOtherResult := res.OtherResult.(NRDBResult)
	_, isSingleTotalResult := res.TotalResult.(NRDBResult)

	assert.True(t, isSingleOtherResult, "Expected otherResult to be a single NRDBResult")
	require.IsType(t, NRDBResult{}, singleOtherResult, "Expected singleOtherResult to be a map[string]interface{}")

	assert.True(t, isSingleTotalResult, "Expected totalResult to be a single NRDBResult")

	// Results should contain multiple items due to FACET
	require.GreaterOrEqual(t, len(res.Results), 1)
	require.GreaterOrEqual(t, len(res.Results), 1)
	if count, ok := singleOtherResult["count"]; ok {
		require.IsType(t, float64(0), count, "Expected 'count' to be a numerical value")
	}
}

// TestIntegrationPerformNRQLQueryTimeseries tests a query with TIMESERIES
func TestIntegrationPerformNRQLQueryTimeseries(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction TIMESERIES 1 hour SINCE 1 day ago"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For TIMESERIES query without FACET, we expect otherResult and totalResult to be single objects
	_, isSingleOtherResult := res.OtherResult.(NRDBResult)
	_, isSingleTotalResult := res.TotalResult.(NRDBResult)

	assert.True(t, isSingleOtherResult, "Expected otherResult to be a single NRDBResult")
	assert.True(t, isSingleTotalResult, "Expected totalResult to be a single NRDBResult")

	// Results should contain multiple items for each time window
	require.GreaterOrEqual(t, len(res.Results), 1)
}

// TestIntegrationPerformNRQLQueryFacetTimeseries tests a query with both FACET and TIMESERIES
func TestIntegrationPerformNRQLQueryFacetTimeseries(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago LIMIT 3"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For FACET + TIMESERIES queries, we expect otherResult and totalResult to be arrays
	otherResults, isArrayOtherResult := res.OtherResult.([]NRDBResult)
	totalResults, isArrayTotalResult := res.TotalResult.([]NRDBResult)

	assert.True(t, isArrayOtherResult, "Expected otherResult to be an array of NRDBResult")
	require.IsType(t, []NRDBResult{}, otherResults, "Expected otherResults to be a []map[string]interface{}")
	assert.True(t, isArrayTotalResult, "Expected totalResult to be an array of NRDBResult")
	require.IsType(t, []NRDBResult{}, totalResults, "Expected totalResults to be a []map[string]interface{}")

	// If they are arrays, they should have multiple elements
	if isArrayOtherResult {
		require.GreaterOrEqual(t, len(otherResults), 1, "Expected otherResults to have at least one element")
	}

	if isArrayTotalResult {
		require.GreaterOrEqual(t, len(totalResults), 1, "Expected totalResults to have at least one element")
	}

	// Results should contain multiple items
	require.GreaterOrEqual(t, len(res.Results), 1)
}

// TestIntegrationPerformNRQLQueryMultipleFacets tests a query with multiple FACETs
func TestIntegrationPerformNRQLQueryMultipleFacets(t *testing.T) {
	t.Parallel()

	query := "SELECT count(*) FROM Transaction FACET appName, name SINCE 1 day ago LIMIT 5"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For multiple FACET query, we expect otherResult and totalResult to be single objects
	_, isSingleOtherResult := res.OtherResult.(NRDBResult)
	_, isSingleTotalResult := res.TotalResult.(NRDBResult)

	assert.True(t, isSingleOtherResult, "Expected otherResult to be a single NRDBResult")
	assert.True(t, isSingleTotalResult, "Expected totalResult to be a single NRDBResult")

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

	// For TIMESERIES with COMPARE WITH, we expect otherResult and totalResult to be single objects
	// and currentResults and previousResults to be populated
	_, isSingleOtherResult := res.OtherResult.(NRDBResult)
	_, isSingleTotalResult := res.TotalResult.(NRDBResult)

	assert.True(t, isSingleOtherResult, "Expected otherResult to be a single NRDBResult")
	assert.True(t, isSingleTotalResult, "Expected totalResult to be a single NRDBResult")

	// Both currentResults and previousResults should be populated
	require.GreaterOrEqual(t, len(res.CurrentResults), 1)
	require.GreaterOrEqual(t, len(res.PreviousResults), 1)
}

// TestIntegrationPerformNRQLQueryFacetTimeseriesDifferentQuery tests a query with both FACET and TIMESERIES using a different NRQL query
func TestIntegrationPerformNRQLQueryFacetTimeseriesDifferentQuery(t *testing.T) {
	t.Parallel()

	query := "SELECT average(duration) FROM Transaction FACET appName TIMESERIES 30 minutes SINCE 1 day ago LIMIT 3"
	client := newNRDBIntegrationTestClient(t)

	res, err := client.PerformNRQLQuery(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	// For FACET + TIMESERIES queries, we expect otherResult and totalResult to be arrays
	otherResults, isArrayOtherResult := res.OtherResult.([]NRDBResult)
	totalResults, isArrayTotalResult := res.TotalResult.([]NRDBResult)

	assert.True(t, isArrayOtherResult, "Expected otherResult to be an array of NRDBResult")
	require.IsType(t, []NRDBResult{}, otherResults, "Expected otherResults to be a []map[string]interface{}")
	assert.True(t, isArrayTotalResult, "Expected totalResult to be an array of NRDBResult")
	require.IsType(t, []NRDBResult{}, totalResults, "Expected totalResults to be a []map[string]interface{}")

	// If they are arrays, they should have multiple elements
	if isArrayOtherResult {
		require.GreaterOrEqual(t, len(otherResults), 1, "Expected otherResults to have at least one element")
	}

	if isArrayTotalResult {
		require.GreaterOrEqual(t, len(totalResults), 1, "Expected totalResults to have at least one element")
	}

	// Results should contain multiple items
	require.GreaterOrEqual(t, len(res.Results), 1)
}
