//go:build unit
// +build unit

package nrdb

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Response with otherResult and totalResult as single objects (non FACET + TIMESERIES)
	testNRDBQuerySingleResponse = `{
  "data": {
    "actor": {
      "account": {
        "nrql": {
          "currentResults": null,
          "metadata": {
            "eventTypes": ["Transaction"],
            "facets": ["appName"],
            "messages": [],
            "timeWindow": {
              "begin": 1709442593550,
              "compareWith": null,
              "end": 1709446193550,
              "since": "1 day AGO",
              "until": "NOW"
            }
          },
          "nrql": "SELECT count(*) FROM Transaction FACET appName SINCE 1 day ago LIMIT 5",
          "otherResult": {
            "count": 42
          },
          "previousResults": null,
          "queryProgress": {
            "completed": true,
            "queryId": null,
            "resultExpiration": null,
            "retryAfter": null,
            "retryDeadline": null
          },
          "results": [
            {
              "appName": "App1",
              "count": 1000
            },
            {
              "appName": "App2", 
              "count": 500
            }
          ],
          "totalResult": {
            "count": 1542
          }
        }
      }
    }
  }
}`

	// Response with otherResult and totalResult as arrays (FACET + TIMESERIES)
	testNRDBQueryArrayResponse = `{
  "data": {
    "actor": {
      "account": {
        "nrql": {
          "currentResults": null,
          "metadata": {
            "eventTypes": ["Transaction"],
            "facets": ["appName"],
            "messages": [],
            "timeWindow": {
              "begin": 1709442593550,
              "compareWith": null,
              "end": 1709446193550,
              "since": "1 day AGO",
              "until": "NOW"
            }
          },
          "nrql": "SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago LIMIT 3",
          "otherResult": [
            {
              "beginTimeSeconds": 1748910826,
              "endTimeSeconds": 1748914426,
              "count": 42
            },
            {
              "beginTimeSeconds": 1748914426,
              "endTimeSeconds": 1748918026,
              "count": 36
            }
          ],
          "previousResults": null,
          "queryProgress": {
            "completed": true,
            "queryId": null,
            "resultExpiration": null,
            "retryAfter": null,
            "retryDeadline": null
          },
          "results": [
            {
              "appName": "App1",
              "beginTimeSeconds": 1748910826,
              "count": 100,
              "endTimeSeconds": 1748914426
            },
            {
              "appName": "App2", 
              "beginTimeSeconds": 1748910826,
              "count": 50,
              "endTimeSeconds": 1748914426
            }
          ],
          "totalResult": [
            {
              "beginTimeSeconds": 1748910826,
              "endTimeSeconds": 1748914426,
              "count": 142
            },
            {
              "beginTimeSeconds": 1748914426,
              "endTimeSeconds": 1748918026,
              "count": 136
            }
          ]
        }
      }
    }
  }
}`
)

func TestUnitPerformNRQLQuerySingleOtherAndTotalResult(t *testing.T) {
	t.Parallel()
	nrdbObject := newMockResponse(t, testNRDBQuerySingleResponse, http.StatusOK)

	// Using a constant account ID for testing
	accountID := 1111111
	query := "SELECT count(*) FROM Transaction FACET appName SINCE 1 day ago LIMIT 5"

	// Execute the function being tested
	result, err := nrdbObject.PerformNRQLQuery(accountID, NRQL(query))

	// Verify basic results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, NRQL(query), result.NRQL)
	assert.Len(t, result.Results, 2)

	// Verify otherResult is a single object
	otherResult, isSingleOtherResult := result.OtherResult.(NRDBResult)
	assert.True(t, isSingleOtherResult, "Expected otherResult to be a single NRDBResult")
	assert.Equal(t, float64(42), otherResult["count"])

	// Verify totalResult is a single object
	totalResult, isSingleTotalResult := result.TotalResult.(NRDBResult)
	assert.True(t, isSingleTotalResult, "Expected totalResult to be a single NRDBResult")
	assert.Equal(t, float64(1542), totalResult["count"])

	// Verify results content
	assert.Equal(t, "App1", result.Results[0]["appName"])
	assert.Equal(t, float64(1000), result.Results[0]["count"])
	assert.Equal(t, "App2", result.Results[1]["appName"])
	assert.Equal(t, float64(500), result.Results[1]["count"])
}

func TestUnitPerformNRQLQueryArrayOtherAndTotalResult(t *testing.T) {
	t.Parallel()
	nrdbObject := newMockResponse(t, testNRDBQueryArrayResponse, http.StatusOK)

	// Using a constant account ID for testing
	accountID := 1111111
	query := "SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago LIMIT 3"

	// Execute the function being tested
	result, err := nrdbObject.PerformNRQLQuery(accountID, NRQL(query))

	// Verify basic results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, NRQL(query), result.NRQL)
	assert.Len(t, result.Results, 2)

	// Verify otherResult is an array
	otherResults, isArrayOtherResult := result.OtherResult.([]NRDBResult)
	assert.True(t, isArrayOtherResult, "Expected otherResult to be an array of NRDBResult")
	require.Len(t, otherResults, 2)
	assert.Equal(t, float64(42), otherResults[0]["count"])
	assert.Equal(t, float64(36), otherResults[1]["count"])

	// Verify totalResult is an array
	totalResults, isArrayTotalResult := result.TotalResult.([]NRDBResult)
	assert.True(t, isArrayTotalResult, "Expected totalResult to be an array of NRDBResult")
	require.Len(t, totalResults, 2)
	assert.Equal(t, float64(142), totalResults[0]["count"])
	assert.Equal(t, float64(136), totalResults[1]["count"])

	// Verify results content
	assert.Equal(t, "App1", result.Results[0]["appName"])
	assert.Equal(t, float64(100), result.Results[0]["count"])
	assert.Equal(t, "App2", result.Results[1]["appName"])
	assert.Equal(t, float64(50), result.Results[1]["count"])
}
