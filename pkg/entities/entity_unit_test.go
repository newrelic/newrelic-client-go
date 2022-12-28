//go:build unit
// +build unit

package entities

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGetEntity_SyntheticsMonitorEntity(t *testing.T) {
	t.Parallel()

	entityClient := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"data": {
				"actor": {
					"entity": {
						"__typename": "SyntheticMonitorEntity",
						"account": {
							"id": 2520528,
							"name": "Account 2520528",
							"reportingEventTypes": [
								"SystemSample",
								"Transaction",
								"TransactionError",
								"TransactionTrace",
								"WorkloadStatus",
								"flexStatusSample",
								"pixieHealthCheck"
							]
						},
						"accountId": 2520528,
						"alertSeverity": "CRITICAL",
						"domain": "SYNTH",
						"entityType": "SYNTHETIC_MONITOR_ENTITY",
						"goldenMetrics": {
							"metrics": [
								{
									"name": "medianDurationS",
									"title": "Median duration (s)"
								},
								{
									"name": "failures",
									"title": "Failures"
								}
							]
						},
						"goldenTags": {
							"tags": [
								{
									"key": "monitorStatus"
								},
								{
									"key": "monitorType"
								},
								{
									"key": "period"
								},
								{
									"key": "publicLocation"
								},
								{
									"key": "privateLocation"
								}
							]
						},
						"indexedAt": 1672657388166,
						"monitorSummary": {
							"locationsFailing": 1,
							"locationsRunning": 1,
							"status": "ENABLED",
							"successRate": 0.0
						},
						"monitorType": "SIMPLE",
						"name": "cyber-ng-asset-ng-stage-eu-west-1",
						"period": 1,
						"recentAlertViolations": [
							{
								"agentUrl": "https://synthetics.newrelic.com/accounts/2520528/monitors/786c74ee-aa35-4ca4-8945-72f163e3a8f2/results/09b27348-180f-40c0-b92f-e3eef2e70577",
								"alertSeverity": "CRITICAL",
								"closedAt": null,
								"label": "Monitor failed for location stage_eks_minion on 'cyber-ng-asset-ng-stage-eu-west-1'",
								"level": "3",
								"openedAt": 1672657387181,
								"violationId": 3799829796,
								"violationUrl": "https://alerts.newrelic.com/accounts/2520528/incidents/1040599543/violations?id=3799829796"
							},
							{
								"agentUrl": "https://synthetics.newrelic.com/accounts/2520528/monitors/786c74ee-aa35-4ca4-8945-72f163e3a8f2/results/b3d80d9b-6b04-4209-9fe4-debf2191aa8c",
								"alertSeverity": "CRITICAL",
								"closedAt": 1672473018774,
								"label": "Monitor failed for location stage_eks_minion on 'cyber-ng-asset-ng-stage-eu-west-1'",
								"level": "3",
								"openedAt": 1672213802788,
								"violationId": 3780902468,
								"violationUrl": "https://alerts.newrelic.com/accounts/2520528/incidents/1033979345/violations?id=3780902468"
							},
							{
								"agentUrl": "https://synthetics.newrelic.com/accounts/2520528/monitors/786c74ee-aa35-4ca4-8945-72f163e3a8f2/results/7e684d0e-b058-4b91-b205-e4df77cc7cc6",
								"alertSeverity": "CRITICAL",
								"closedAt": 1671880128781,
								"label": "Monitor failed for location stage_eks_minion on 'cyber-ng-asset-ng-stage-eu-west-1'",
								"level": "3",
								"openedAt": 1671620905292,
								"violationId": 3708808439,
								"violationUrl": "https://alerts.newrelic.com/accounts/2520528/incidents/1024826483/violations?id=3708808439"
							}
						],
						"relatedEntities": {
							"nextCursor": null,
							"results": [
								{
									"__typename": "EntityRelationshipDetectedEdge",
									"createdAt": 1672657388437,
									"type": "CONNECTS_TO"
								}
							]
						},
						"relationships": null,
						"reporting": true,
						"serviceLevel": {
							"indicators": null
						},
						"tags": [
							{
								"key": "account",
								"values": [
									"Account 2520528"
								]
							},
							{
								"key": "accountId",
								"values": [
									"2520528"
								]
							},
							{
								"key": "apdexTarget",
								"values": [
									"7.0"
								]
							},
							{
								"key": "environment",
								"values": [
									"ng-stage-eu-west-1"
								]
							},
							{
								"key": "monitorStatus",
								"values": [
									"Enabled"
								]
							},
							{
								"key": "monitorType",
								"values": [
									"Ping"
								]
							},
							{
								"key": "period",
								"values": [
									"1"
								]
							},
							{
								"key": "platform",
								"values": [
									"cyberng"
								]
							},
							{
								"key": "privateLocation",
								"values": [
									"stage-eks-minion"
								]
							},
							{
								"key": "redirectIsFailure",
								"values": [
									"false"
								]
							},
							{
								"key": "service",
								"values": [
									"cyber-ng-asset-ng-stage-eu-west-1"
								]
							},
							{
								"key": "shouldBypassHeadRequest",
								"values": [
									"false"
								]
							},
							{
								"key": "trustedAccountId",
								"values": [
									"2520528"
								]
							},
							{
								"key": "useTlsValidation",
								"values": [
									"true"
								]
							}
						],
						"tagsWithMetadata": [
							{
								"key": "account",
								"values": [
									{
										"mutable": false,
										"value": "Account 2520528"
									}
								]
							},
							{
								"key": "accountId",
								"values": [
									{
										"mutable": false,
										"value": "2520528"
									}
								]
							},
							{
								"key": "apdexTarget",
								"values": [
									{
										"mutable": false,
										"value": "7.0"
									}
								]
							},
							{
								"key": "environment",
								"values": [
									{
										"mutable": true,
										"value": "ng-stage-eu-west-1"
									}
								]
							},
							{
								"key": "monitorStatus",
								"values": [
									{
										"mutable": false,
										"value": "Enabled"
									}
								]
							},
							{
								"key": "monitorType",
								"values": [
									{
										"mutable": false,
										"value": "Ping"
									}
								]
							},
							{
								"key": "period",
								"values": [
									{
										"mutable": false,
										"value": "1"
									}
								]
							},
							{
								"key": "platform",
								"values": [
									{
										"mutable": true,
										"value": "cyberng"
									}
								]
							},
							{
								"key": "privateLocation",
								"values": [
									{
										"mutable": false,
										"value": "stage-eks-minion"
									}
								]
							},
							{
								"key": "redirectIsFailure",
								"values": [
									{
										"mutable": false,
										"value": "false"
									}
								]
							},
							{
								"key": "service",
								"values": [
									{
										"mutable": true,
										"value": "cyber-ng-asset-ng-stage-eu-west-1"
									}
								]
							},
							{
								"key": "shouldBypassHeadRequest",
								"values": [
									{
										"mutable": false,
										"value": "false"
									}
								]
							},
							{
								"key": "trustedAccountId",
								"values": [
									{
										"mutable": false,
										"value": "2520528"
									}
								]
							},
							{
								"key": "useTlsValidation",
								"values": [
									{
										"mutable": false,
										"value": "true"
									}
								]
							}
						],
						"type": "MONITOR"
					}
				}
			},
			"errors": [
				{
					"locations": [
						{
							"column": 3,
							"line": 646
						}
					],
					"message": "This field is deprecated! Please use relatedEntities instead.",
					"path": [
						"actor",
						"entity",
						"relationships"
					]
				}
			]
		}`))

		assert.NoError(t, err)
	}))

	result, err := entityClient.GetEntity(common.EntityGUID("MjUyMDUyOHxTWU5USHxNT05JVE9SfDViZjE3ODgwLTNhYjQtNGEyNC04ODFiLTI2YmU0OTA4NDk2Yg"))

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// nolint
func newTestClient(t *testing.T, handler http.Handler) Entities {
	ts := httptest.NewServer(handler)
	tc := mock.NewTestConfig(t, ts)

	c := New(tc)

	return c
}

// nolint
func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Entities {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}
