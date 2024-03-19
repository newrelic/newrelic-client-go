//go:build unit
// +build unit

package nrdb

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testNRQLQuery         = "SELECT * from Metric where entity.guid ='MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTUzNDQ4MjAy' and endTimestamp = 1709446129592"
	testNRDBQueryResponse = `{
  "data": {
    "actor": {
      "account": {
        "nrql": {
          "currentResults": null,
          "embeddedChartUrl": "https://chart-embed.service.newrelic.com/charts/0ab47fb4-b9e9-4d31-a3bc-539a37a313c2",
          "eventDefinitions": [
            {
              "attributes": [
                {
                  "definition": "The start time for the metric, in Unix time.",
                  "documentationUrl": "https://docs.newrelic.com/attribute-dictionary?attribute=timestamp&event=Metric",
                  "label": "timestamp",
                  "name": "timestamp"
                },
                {
                  "definition": "The source of this data. For example: metricAPI.",
                  "documentationUrl": "https://docs.newrelic.com/attribute-dictionary?attribute=newrelic.source&event=Metric",
                  "label": "newrelic.source",
                  "name": "newrelic.source"
                },
                {
                  "definition": "Name of the metric.",
                  "documentationUrl": "https://docs.newrelic.com/attribute-dictionary?attribute=metricName&event=Metric",
                  "label": "metricName",
                  "name": "metricName"
                },
                {
                  "definition": "The length of the time window.",
                  "documentationUrl": "https://docs.newrelic.com/attribute-dictionary?attribute=interval.ms&event=Metric",
                  "label": "interval.ms",
                  "name": "interval.ms"
                },
                {
                  "definition": "The end of the time range associated with the metric, in Unix time, in milliseconds. This is calculated by adding the metric interval to the timestamp of the metric (timestamp + interval.ms).",
                  "documentationUrl": "https://docs.newrelic.com/attribute-dictionary?attribute=endTimestamp&event=Metric",
                  "label": "endTimestamp",
                  "name": "endTimestamp"
                }
              ],
              "definition": "Represents a metric data point (a measurement over a range of time, or a sample at a specific point in time) with multiple attributes attached, which allow for in-depth analysis and querying. This metric data comes from our Metric API, our Telemetry SDKs, network performance monitoring, and some of our open-source exporters/integrations.",
              "label": "Metric",
              "name": "Metric"
            }
          ],
          "metadata": {
            "eventTypes": [
              "Metric"
            ],
            "facets": null,
            "messages": [
              "We tried to transform your Metric query into a query that is compatible with the previous Infrastructure metric format, but there are some expressions that can’t be transformed. More info here[https://docs.newrelic.com/docs/query-your-data/nrql-new-relic-query-language/nrql-query-tutorials/query-infrastructure-dimensional-metrics-nrql#known-limitations]."
            ],
            "timeWindow": {
              "begin": 1709442593550,
              "compareWith": null,
              "end": 1709446193550,
              "since": "60 MINUTES AGO",
              "until": "NOW"
            }
          },
          "nrql": "` + testNRQLQuery + `",
          "otherResult": null,
          "previousResults": null,
          "queryProgress": {
            "completed": true,
            "queryId": null,
            "resultExpiration": null,
            "retryAfter": null,
            "retryDeadline": null
          },
          "rawResponse": {
            "metadata": {
              "accounts": [
                3806526
              ],
              "beginTime": "2024-03-03T05:09:53Z",
              "beginTimeMillis": 1709442593550,
              "contents": [
                {
                  "function": "events",
                  "limit": 100,
                  "order": {
                    "column": "timestamp",
                    "descending": true
                  }
                }
              ],
              "endTime": "2024-03-03T06:09:53Z",
              "endTimeMillis": 1709446193550,
              "eventType": "Metric",
              "eventTypes": [
                "Metric"
              ],
              "guid": "82b231a0-982f-8372-0b9b-96110c7b812f",
              "messages": [
                "We tried to transform your Metric query into a query that is compatible with the previous Infrastructure metric format, but there are some expressions that can’t be transformed. More info here[https://docs.newrelic.com/docs/query-your-data/nrql-new-relic-query-language/nrql-query-tutorials/query-infrastructure-dimensional-metrics-nrql#known-limitations]."
              ],
              "openEnded": true,
              "rawCompareWith": "",
              "rawSince": "60 MINUTES AGO",
              "rawUntil": "NOW",
              "routerGuid": "82b231a0-982f-8372-0b9b-96110c7b812f",
              "timeAggregations": [
                "raw metrics"
              ],
              "warn": [
                "We tried to transform your Metric query into a query that is compatible with the previous Infrastructure metric format, but there are some expressions that can’t be transformed. More info here[https://docs.newrelic.com/docs/query-your-data/nrql-new-relic-query-language/nrql-query-tutorials/query-infrastructure-dimensional-metrics-nrql#known-limitations]."
              ]
            },
            "performanceStats": {
              "exceedsRetentionWindow": false,
              "inspectedCount": 32420,
              "matchCount": 1,
              "omittedCount": 0,
              "wallClockTime": 34
            },
            "results": [
              {
                "events": [
                  {
                    "appId": 553448202,
                    "appName": "Dummy App Pro Max",
                    "endTimestamp": 1709446129592,
                    "entity.guid": "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTUzNDQ4MjAy",
                    "metricName": "newrelic.internal.usage",
                    "newrelic.internal.usage": {
                      "count": 49500,
                      "type": "count"
                    },
                    "timestamp": 1709446121841,
                    "type": "timeslice",
                    "usage.agent.language": "python",
                    "usage.agent.type": "apm",
                    "usage.metric": "Metrics",
                    "usage.newrelic.source": "agent"
                  }
                ]
              }
            ]
          },
          "results": [
            {
              "appId": 553448202,
              "appName": "Dummy App Pro Max",
              "endTimestamp": 1709446129592,
              "entity.guid": "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTUzNDQ4MjAy",
              "metricName": "newrelic.internal.usage",
              "newrelic.internal.usage": {
                "count": 49500,
                "type": "count"
              },
              "timestamp": 1709446121841,
              "type": "timeslice",
              "usage.agent.language": "python",
              "usage.agent.type": "apm",
              "usage.metric": "Metrics",
              "usage.newrelic.source": "agent"
            }
          ],
          "staticChartUrl": "https://chart-image.service.newrelic.com/image/5bef0ad8-2c65-4d51-bd3d-b35ecec9cb25",
          "suggestedFacets": [],
          "totalResult": null
        }
      }
    }
  }
}`
)

func TestUnitNRDBQuery(t *testing.T) {
	t.Parallel()
	nrdbObject := newMockResponse(t, testNRDBQueryResponse, http.StatusCreated)

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		t.Skipf("test requires NEW_RELIC_ACOUNT_ID")
	}

	actual, err := nrdbObject.Query(
		accountID,
		NRQL(testNRQLQuery),
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func TestUnitNRDBQueryWithAdditionalOptions(t *testing.T) {
	t.Parallel()
	nrdbObject := newMockResponse(t, testNRDBQueryResponse, http.StatusCreated)

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		t.Skipf("test requires NEW_RELIC_ACOUNT_ID")
	}

	actual, err := nrdbObject.QueryWithAdditionalOptions(
		accountID,
		NRQL(testNRQLQuery),
		30,
		false,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
