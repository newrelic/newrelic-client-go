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
	testNRDBQuerySingleResponse = `
{
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
	testNRDBQueryArrayResponse = `
{
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

	// Response with null values for otherResult and totalResult
	testNRDBQueryNullResponse = `
{
  "data": {
    "actor": {
      "account": {
        "nrql": {
          "currentResults": null,
          "metadata": {
            "eventTypes": ["Transaction"],
            "facets": null,
            "messages": [],
            "timeWindow": {
              "begin": 1709442593550,
              "compareWith": null,
              "end": 1709446193550,
              "since": "1 day AGO",
              "until": "NOW"
            }
          },
          "nrql": "SELECT count(*) FROM Transaction SINCE 1 day ago",
          "otherResult": null,
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
              "count": 1542
            }
          ],
          "totalResult": null
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
	accountID := 3957524
	query := "SELECT count(*) FROM Transaction FACET appName SINCE 1 day ago LIMIT 5"

	// Execute the function being tested
	result, err := nrdbObject.PerformNRQLQuery(accountID, NRQL(query))

	// Verify basic results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, NRQL(query), result.NRQL)
	assert.Len(t, result.Results, 2)

	// Verify otherResult is an array with one element (converted from single object)
	assert.IsType(t, NRDBMultiResultCustomized{}, result.OtherResult)
	require.Len(t, result.OtherResult, 1)
	assert.Equal(t, float64(42), result.OtherResult[0]["count"])

	// Verify totalResult is an array with one element (converted from single object)
	assert.IsType(t, NRDBMultiResultCustomized{}, result.TotalResult)
	require.Len(t, result.TotalResult, 1)
	assert.Equal(t, float64(1542), result.TotalResult[0]["count"])

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
	accountID := 3957524
	query := "SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago LIMIT 3"

	// Execute the function being tested
	result, err := nrdbObject.PerformNRQLQuery(accountID, NRQL(query))

	// Verify basic results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, NRQL(query), result.NRQL)
	assert.Len(t, result.Results, 2)

	// Verify otherResult is an array (preserved as-is)
	assert.IsType(t, NRDBMultiResultCustomized{}, result.OtherResult)
	require.Len(t, result.OtherResult, 2)
	assert.Equal(t, float64(42), result.OtherResult[0]["count"])
	assert.Equal(t, float64(36), result.OtherResult[1]["count"])

	// Verify totalResult is an array (preserved as-is)
	assert.IsType(t, NRDBMultiResultCustomized{}, result.TotalResult)
	require.Len(t, result.TotalResult, 2)
	assert.Equal(t, float64(142), result.TotalResult[0]["count"])
	assert.Equal(t, float64(136), result.TotalResult[1]["count"])

	// Verify results content
	assert.Equal(t, "App1", result.Results[0]["appName"])
	assert.Equal(t, float64(100), result.Results[0]["count"])
	assert.Equal(t, "App2", result.Results[1]["appName"])
	assert.Equal(t, float64(50), result.Results[1]["count"])
}

func TestUnitPerformNRQLQueryNullFields(t *testing.T) {
	t.Parallel()
	nrdbObject := newMockResponse(t, testNRDBQueryNullResponse, http.StatusOK)

	// Using a constant account ID for testing
	accountID := 3957524
	query := "SELECT count(*) FROM Transaction SINCE 1 day ago"

	// Execute the function being tested
	result, err := nrdbObject.PerformNRQLQuery(accountID, NRQL(query))

	// Verify basic results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, NRQL(query), result.NRQL)
	assert.Len(t, result.Results, 1)

	// Verify otherResult and totalResult are empty arrays (not nil)
	assert.IsType(t, NRDBMultiResultCustomized{}, result.OtherResult)
	assert.Empty(t, result.OtherResult)
	assert.IsType(t, NRDBMultiResultCustomized{}, result.TotalResult)
	assert.Empty(t, result.TotalResult)

	// Verify the result still works
	assert.Equal(t, float64(1542), result.Results[0]["count"])
}
