//go:build unit
// +build unit

package synthetics

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockSyntheticsStartAutomatedTestResponse = `{
	"data": {
		"syntheticsStartAutomatedTest": {
			"batchId": "` + mockBatchID + `"
			}
		}
	}`

	mockAutomatedTestResultsResponse = `{
	   "data": {
		  "actor": {
			 "account": {
				"synthetics": {
				   "automatedTestResult": {
					  "config": {
						 "batchName": "sample-batch",
						 "branch": "sample-branch",
						 "commit": "sample-commit",
						 "deepLink": "sample-deepLink",
						 "platform": "sample-platform",
						 "repository": "sample-repository"
					  },
					  "status": "FAILED",
					  "tests": [
						 {
							"automatedTestMonitorConfig": {
							   "isBlocking": true,
							   "overrides": {
								  "domain": [
									 {
										"domain": "sample-domain",
										"override": "sample-override"
									 }
								  ],
								  "location": "sample-location",
								  "secureCredential": [
									 {
										"key": "sample-key",
										"overrideKey": "sample-override-key"
									 }
								  ],
								  "startingUrl": "sample-starting-url.com"
							   }
							},
							"batchId": "` + mockBatchID + `",
							"duration": 0,
							"error": "An unexpected error occurred",
							"id": "55bb193c-f0d0-4eed-b0de-45590949c99e",
							"location": "AWS_US_WEST_2",
							"locationLabel": "Portland, OR, USA",
							"monitorGuid": "` + monitorOneGUID + `",
							"monitorId": "` + monitorOneID + `",
							"monitorName": "Test Monitor One",
							"result": "FAILED",
							"resultsUrl": "sample-results-url.com",
							"type": "SCRIPT_BROWSER",
							"typeLabel": "Scripted Browser"
						 },
						 {
							"automatedTestMonitorConfig": {
							   "isBlocking": false,
							   "overrides": null
							},
							"batchId": "` + mockBatchID + `",
							"duration": 0,
							"error": "",
							"id": "6ce400b2-cba2-43cf-b139-1267d26c522f",
							"location": "AWS_US_WEST_2",
							"locationLabel": "Portland, OR, USA",
							"monitorGuid": "` + monitorTwoGUID + `",
							"monitorId": "` + monitorTwoID + `",
							"monitorName": "Test Monitor Two",
							"result": "SUCCESS",
							"resultsUrl": "sample-results-url.com",
							"type": "SCRIPT_BROWSER",
							"typeLabel": "Scripted Browser"
						 },
						 {
							"automatedTestMonitorConfig": {
							   "isBlocking": false,
							   "overrides": null
							},
							"batchId": "` + mockBatchID + `",
							"duration": 0,
							"error": "Unknown Error",
							"id": "9000dd3e-eaa6-42fb-8cda-aad0b54147c8",
							"location": "AWS_US_WEST_2",
							"locationLabel": "Portland, OR, USA",
							"monitorGuid": "` + monitorThreeGUID + `",
							"monitorId": "` + monitorThreeID + `",
							"monitorName": "Test Monitor Three",
							"result": "FAILED",
							"resultsUrl": "sample-results-url.com",
							"type": "SCRIPT_BROWSER",
							"typeLabel": "Scripted Browser"
						 }
					  ]
				   }
				}
			 }
		  }
	   }
	}`

	monitorOneGUID   = "McUyMDUyOHxTWU5USHxNT05JVE9SfGE0MTk2MzRhLWRhNzgtNGQzNC1hN2NmLTBmMDE1MjQyAjA2AA"
	monitorTwoGUID   = "McUyMDUyOHxTWU5USHxNT05JVE9SfDA1YjVkZDQ4LTJlNzQtNDhjMi05OTQ0LTdkYTU2YTc3BmByBB"
	monitorThreeGUID = "McUyMDUyOHxTWU5USHxNT05JVE9SfDljMmUwMzAxLTBiNWMtNGZmMy05M2JhLTE4ODg5MGVkCzCzCC"

	monitorOneID   = "3ffb2a4f-459d-49c2-85df-726e86a3d773"
	monitorTwoID   = "12b5c6a8-331f-4609-926e-ca6adb8fb6af"
	monitorThreeID = "942e0948-0ab1-4a19-8177-3288259e4cc6"

	mockBatchID = "d0622283-87b7-44a1-9622-0d66971b0c1d"
	// random seven digit number
	mockAccountID = rand.Intn(999999) + 1000000

	configInput = SyntheticsAutomatedTestConfigInput{
		BatchName:  "sample-batch",
		Branch:     "sample-branch",
		Commit:     "sample-commit",
		DeepLink:   "sample-deepLink",
		Platform:   "sample-platform",
		Repository: "sample-repository",
	}

	testsInput = []SyntheticsAutomatedTestMonitorInput{
		{
			Config: SyntheticsAutomatedTestMonitorConfigInput{
				IsBlocking: true,
				Overrides: &SyntheticsAutomatedTestOverridesInput{
					Domain: SyntheticsScriptDomainOverrideInput{
						Domain:   "sample-domain",
						Override: "sample-override",
					},
					SecureCredential: SyntheticsSecureCredentialOverrideInput{
						Key:         "sample-key",
						OverrideKey: "sample-override-key",
					},
					StartingURL: "sample-starting-url.com",
					Location:    "sample-location",
				},
			},
			MonitorGUID: EntityGUID(monitorOneGUID),
		},
		{
			Config:      SyntheticsAutomatedTestMonitorConfigInput{},
			MonitorGUID: EntityGUID(monitorTwoGUID),
		},
		{
			Config:      SyntheticsAutomatedTestMonitorConfigInput{},
			MonitorGUID: EntityGUID(monitorThreeGUID),
		},
	}

	automatedTestResults = SyntheticsAutomatedTestResult{
		Config: SyntheticsAutomatedTestConfig{
			BatchName:  "sample-batch",
			Branch:     "sample-branch",
			Commit:     "sample-commit",
			DeepLink:   "sample-deepLink",
			Platform:   "sample-platform",
			Repository: "sample-repository",
		},
		Status: "FAILED",
		Tests: []SyntheticsAutomatedTestJobResult{
			// the first monitor has incorrect values (which are dissmilar to those specified in the mock response)
			{
				AutomatedTestMonitorConfig: SyntheticsAutomatedTestMonitorConfig{
					IsBlocking: true,
					Overrides:  nil,
				},
				BatchId:       mockBatchID,
				Duration:      0,
				Error:         "An unexpected error occurred",
				ID:            "55bb193c-f0d0-4eed-b0de-45590949c99e",
				Location:      "AWS_US_WEST_2",
				LocationLabel: "Columbus, OH, USA",
				MonitorGUID:   EntityGUID(monitorOneGUID),
				MonitorId:     monitorOneID,
				MonitorName:   "Test Monitor One",
				Result:        "PASSED",
				ResultsURL:    "sample-results-url.com",
				Type:          "SCRIPT_BROWSER",
				TypeLabel:     "Scripted Browser",
			},
			{
				AutomatedTestMonitorConfig: SyntheticsAutomatedTestMonitorConfig{
					IsBlocking: false,
					Overrides:  nil,
				},
				BatchId:       mockBatchID,
				Duration:      0,
				Error:         "",
				ID:            "6ce400b2-cba2-43cf-b139-1267d26c522f",
				Location:      "AWS_US_WEST_2",
				LocationLabel: "Portland, OR, USA",
				MonitorGUID:   EntityGUID(monitorTwoGUID),
				MonitorId:     monitorTwoID,
				MonitorName:   "Test Monitor Two",
				Result:        "SUCCESS",
				ResultsURL:    "sample-results-url.com",
				Type:          "SCRIPT_BROWSER",
				TypeLabel:     "Scripted Browser",
			},
			{
				AutomatedTestMonitorConfig: SyntheticsAutomatedTestMonitorConfig{
					IsBlocking: false,
					Overrides:  nil,
				},
				BatchId:       mockBatchID,
				Duration:      0,
				Error:         "Unknown Error",
				ID:            "9000dd3e-eaa6-42fb-8cda-aad0b54147c8",
				Location:      "AWS_US_WEST_2",
				LocationLabel: "Portland, OR, USA",
				MonitorGUID:   EntityGUID(monitorThreeGUID),
				MonitorId:     monitorThreeID,
				MonitorName:   "Test Monitor Three",
				Result:        "FAILED",
				ResultsURL:    "sample-results-url.com",
				Type:          "SCRIPT_BROWSER",
				TypeLabel:     "Scripted Browser",
			},
		},
	}
)

func TestUnit_SyntheticsStartAutomatedTest(t *testing.T) {
	t.Parallel()
	t.Skipf(
		`Temporarily skipping tests associated with the Synthetics Automated Tests feature, ` +
			`given the API is currently unstable and endpoint access is not configured to all accounts at the moment.
	`)

	synthetics := newMockResponse(t, mockSyntheticsStartAutomatedTestResponse, http.StatusOK)

	expected := SyntheticsAutomatedTestStartResult{BatchId: mockBatchID}
	actual, err := synthetics.SyntheticsStartAutomatedTest(configInput, testsInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestUnit_AutomatedTestResults_Failure(t *testing.T) {
	t.Parallel()
	t.Skipf(
		`Temporarily skipping tests associated with the Synthetics Automated Tests feature, ` +
			`given the API is currently unstable and endpoint access is not configured to all accounts at the moment.
	`)

	synthetics := newMockResponse(t, mockAutomatedTestResultsResponse, http.StatusOK)

	expected := automatedTestResults
	actual, err := synthetics.GetAutomatedTestResult(mockAccountID, mockBatchID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	// The expected and actual responses are expected to be dissimilar, given that
	// the first test in the list of tests comprises incorrect values of attributes
	assert.NotEqual(t, expected, *actual)
}

func TestUnit_AutomatedTestResults(t *testing.T) {
	t.Parallel()
	t.Skipf(
		`Temporarily skipping tests associated with the Synthetics Automated Tests feature, ` +
			`given the API is currently unstable and endpoint access is not configured to all accounts at the moment.
	`)
	synthetics := newMockResponse(t, mockAutomatedTestResultsResponse, http.StatusOK)

	expected := automatedTestResults

	// Since the test TestUnit_AutomatedTestResults_Failure is seen to fail due to dissimilarities
	// between the expected and actual responses, the following lines of code correct the
	// erroneous fields accordingly.

	expected.Tests[0].AutomatedTestMonitorConfig.Overrides = &SyntheticsAutomatedTestOverrides{
		Domain: []SyntheticsScriptDomainOverride{
			{
				Domain:   "sample-domain",
				Override: "sample-override",
			},
		},
		SecureCredential: []SyntheticsSecureCredentialOverride{
			{
				Key:         "sample-key",
				OverrideKey: "sample-override-key",
			},
		},
		StartingURL: "sample-starting-url.com",
		Location:    "sample-location",
	}
	expected.Tests[0].LocationLabel = "Portland, OR, USA"
	expected.Tests[0].Result = SyntheticsJobStatusTypes.FAILED

	actual, err := synthetics.GetAutomatedTestResult(mockAccountID, mockBatchID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}
