//go:build unit
// +build unit

package aimonitoring

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	testDatasetID = "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY"
	testVersionID = "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDAxYTFlNTc5LTFlMmUtN2ZkMi1hYzQ0LTE0NmI2Y2Y0Mzc1YQ"

	testDatasetSearchResponseJSON = `{
		"data": {
			"actor": {
				"entityManagement": {
					"entitySearch": {
						"entities": [
							{
								"__typename": "EntityManagementAimDatasetEntity",
								"id": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY",
								"name": "PII redaction test",
								"type": "AIM_DATASET",
								"description": "redaction fixtures",
								"latestVersion": 2,
								"latestVersionId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDAxYTFlNTc5LTFlMmUtN2ZkMi1hYzQ0LTE0NmI2Y2Y0Mzc1YQ",
								"mode": "TEXT",
								"variables": ["input"],
								"versionCount": 2,
								"tags": [{"key": "team", "values": ["aim"]}],
								"metadata": {"createdAt": 1794851282478, "updatedAt": 1794851282675, "version": 1}
							}
						],
						"nextCursor": "next-page-cursor"
					}
				}
			}
		}
	}`

	testVersionSearchResponseJSON = `{
		"data": {
			"actor": {
				"entityManagement": {
					"entitySearch": {
						"entities": [
							{
								"__typename": "EntityManagementAimDatasetVersionEntity",
								"id": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDAxYTFlNTc5LTFlMmUtN2ZkMi1hYzQ0LTE0NmI2Y2Y0Mzc1YQ",
								"name": "v1",
								"type": "AIM_DATASET_VERSION",
								"aimDataset": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY",
								"version": 1,
								"message": "Initial version",
								"recordCount": 9000,
								"metadata": {"createdAt": 1794851282478, "updatedAt": 1794851282675, "version": 1}
							}
						],
						"nextCursor": null
					}
				}
			}
		}
	}`

	testDatasetDeleteResponseJSON = `{
		"data": {
			"aimDatasetDelete": {
				"datasetId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY"
			}
		}
	}`

	testDatasetDeleteVersionResponseJSON = `{
		"data": {
			"aimDatasetDeleteVersion": {
				"versionId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDAxYTFlNTc5LTFlMmUtN2ZkMi1hYzQ0LTE0NmI2Y2Y0Mzc1YQ"
			}
		}
	}`
)

type graphQLRequestBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// newMockResponseCapturingRequest returns a client backed by a mock server that
// serves mockJSONResponse and records the body of the last request it received.
func newMockResponseCapturingRequest(t *testing.T, mockJSONResponse string, statusCode int) (Aimonitoring, *graphQLRequestBody) {
	var lastRequest graphQLRequestBody

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, &lastRequest))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err = w.Write([]byte(mockJSONResponse))
		require.NoError(t, err)
	}))
	t.Cleanup(ts.Close)

	return New(mock.NewTestConfig(t, ts)), &lastRequest
}

func TestUnitGetEntitySearch_Datasets(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testDatasetSearchResponseJSON, http.StatusOK)

	result, err := client.GetEntitySearch(nil, "type = 'AIM_DATASET'")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "type = 'AIM_DATASET'", lastRequest.Variables["query"])
	assert.Nil(t, lastRequest.Variables["cursor"])

	assert.Equal(t, "next-page-cursor", result.NextCursor)
	require.Len(t, result.Entities, 1)

	dataset, ok := result.Entities[0].(*EntityManagementAimDatasetEntity)
	require.True(t, ok, "expected *EntityManagementAimDatasetEntity, got %T", result.Entities[0])
	assert.Equal(t, testDatasetID, dataset.ID)
	assert.Equal(t, "PII redaction test", dataset.Name)
	assert.Equal(t, "redaction fixtures", dataset.Description)
	assert.Equal(t, 2, dataset.LatestVersion)
	assert.Equal(t, testVersionID, dataset.LatestVersionId)
	assert.Equal(t, EntityManagementAimDatasetModeTypes.TEXT, dataset.Mode)
	assert.Equal(t, []string{"input"}, dataset.Variables)
	assert.Equal(t, []EntityManagementTag{{Key: "team", Values: []string{"aim"}}}, dataset.Tags)
}

func TestUnitGetEntitySearch_DatasetVersions(t *testing.T) {
	t.Parallel()

	client, _ := newMockResponseCapturingRequest(t, testVersionSearchResponseJSON, http.StatusOK)

	result, err := client.GetEntitySearch(nil, "type = 'AIM_DATASET_VERSION'")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Empty(t, result.NextCursor)
	require.Len(t, result.Entities, 1)

	version, ok := result.Entities[0].(*EntityManagementAimDatasetVersionEntity)
	require.True(t, ok, "expected *EntityManagementAimDatasetVersionEntity, got %T", result.Entities[0])
	assert.Equal(t, testVersionID, version.ID)
	assert.Equal(t, "v1", version.Name)
	assert.Equal(t, testDatasetID, version.AimDataset)
	assert.Equal(t, 1, version.Version)
	assert.Equal(t, "Initial version", version.Message)
	assert.Equal(t, 9000, version.RecordCount)
}

func TestUnitGetEntitySearch_SendsCursorVariable(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testDatasetSearchResponseJSON, http.StatusOK)

	cursor := "current-page-cursor"
	_, err := client.GetEntitySearch(&cursor, "type = 'AIM_DATASET'")
	require.NoError(t, err)

	assert.Contains(t, lastRequest.Query, "$cursor: String")
	assert.Contains(t, lastRequest.Query, "cursor: $cursor")
	assert.Equal(t, cursor, lastRequest.Variables["cursor"])
}

func TestUnitAimDatasetDelete(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testDatasetDeleteResponseJSON, http.StatusOK)

	result, err := client.AimDatasetDelete(AimDatasetDeleteInput{DatasetId: testDatasetID})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, map[string]interface{}{"datasetId": testDatasetID}, lastRequest.Variables["dataset"])
	assert.Equal(t, testDatasetID, result.DatasetId)
}

func TestUnitAimDatasetDeleteVersion(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testDatasetDeleteVersionResponseJSON, http.StatusOK)

	result, err := client.AimDatasetDeleteVersion(AimDatasetDeleteVersionInput{
		DatasetId: testDatasetID,
		VersionId: testVersionID,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, map[string]interface{}{"datasetId": testDatasetID, "versionId": testVersionID}, lastRequest.Variables["version"])
	assert.Equal(t, testVersionID, result.VersionId)
}

var (
	testCreateResponseJSON = `{
		"data": {
			"aimDatasetCreate": {
				"datasetId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY",
				"versionId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDAxYTFlNTc5LTFlMmUtN2ZkMi1hYzQ0LTE0NmI2Y2Y0Mzc1YQ"
			}
		}
	}`

	testCreateFromCsvResponseJSON = `{
		"data": {
			"aimDatasetCreateFromCsv": {
				"datasetId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY",
				"versionId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDAxYTFlNTc5LTFlMmUtN2ZkMi1hYzQ0LTE0NmI2Y2Y0Mzc1YQ"
			}
		}
	}`

	testUpdateResponseJSON = `{
		"data": {
			"aimDatasetUpdate": {
				"datasetId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY"
			}
		}
	}`

	testUpsertRecordsResponseJSON = `{
		"data": {
			"aimDatasetUpsertRecords": {
				"datasetId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwwMWExZTU3OS0xZDhhLTdkNTktYmM2OC0wMjhiOTQ0YzMwNDY",
				"versionId": "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDAxYTFlNTc5LTFlMmUtN2ZkMi1hYzQ0LTE0NmI2Y2Y0Mzc1YQ"
			}
		}
	}`

	testRecordsResponseJSON = `{
		"data": {
			"actor": {
				"account": {
					"aimDataset": {
						"records": {
							"records": [
								{"id": 1, "input": {"prompt": "hello"}, "expected": "hi", "context": "greeting", "additionalFields": {"lang": "en"}},
								{"id": 2, "input": [{"role": "user", "content": "hello"}]}
							],
							"nextCursor": "next-page-cursor",
							"totalCount": 9000
						}
					}
				}
			}
		}
	}`
)

func TestUnitAimDatasetCreate(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testCreateResponseJSON, http.StatusOK)

	result, err := client.AimDatasetCreate(AimDatasetCreateInput{
		Name:  "PII redaction test",
		Mode:  AimDatasetModeTypes.TEXT,
		Scope: AimDatasetScopeInput{ID: "11859098", Type: AimDatasetScopeTypeTypes.ACCOUNT},
		Records: []AimDatasetRecordInput{
			{Input: map[string]interface{}{"prompt": "hello"}, Expected: "hi"},
		},
		Tags: []AimDatasetTagInput{{Key: "team", Values: []string{"aim"}}},
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, map[string]interface{}{
		"name":  "PII redaction test",
		"mode":  "TEXT",
		"scope": map[string]interface{}{"id": "11859098", "type": "ACCOUNT"},
		"records": []interface{}{
			map[string]interface{}{"input": map[string]interface{}{"prompt": "hello"}, "expected": "hi"},
		},
		"tags": []interface{}{
			map[string]interface{}{"key": "team", "values": []interface{}{"aim"}},
		},
	}, lastRequest.Variables["dataset"])
	assert.Equal(t, testDatasetID, result.DatasetId)
	assert.Equal(t, testVersionID, result.VersionId)
}

func TestUnitAimDatasetCreateFromCsv(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testCreateFromCsvResponseJSON, http.StatusOK)

	result, err := client.AimDatasetCreateFromCsv(AimDatasetCreateFromCsvInput{
		Name:  "PII redaction test",
		Csv:   "prompt,answer\nhello,hi\n",
		Scope: AimDatasetScopeInput{ID: "11859098", Type: AimDatasetScopeTypeTypes.ACCOUNT},
		Mapping: AimDatasetCsvMappingInput{
			Mode:        AimDatasetCsvModeTypes.TEXT,
			InputKeys:   []string{"prompt"},
			ExpectedKey: "answer",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, map[string]interface{}{
		"name":  "PII redaction test",
		"csv":   "prompt,answer\nhello,hi\n",
		"scope": map[string]interface{}{"id": "11859098", "type": "ACCOUNT"},
		"mapping": map[string]interface{}{
			"mode":                "TEXT",
			"inputKeys":           []interface{}{"prompt"},
			"expectedKey":         "answer",
			"additionalFieldKeys": nil,
		},
	}, lastRequest.Variables["dataset"])
	assert.Equal(t, testDatasetID, result.DatasetId)
	assert.Equal(t, testVersionID, result.VersionId)
}

func TestUnitAimDatasetUpdate(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testUpdateResponseJSON, http.StatusOK)

	result, err := client.AimDatasetUpdate(AimDatasetUpdateInput{
		DatasetId:   testDatasetID,
		Description: "updated description",
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, map[string]interface{}{
		"datasetId":   testDatasetID,
		"description": "updated description",
	}, lastRequest.Variables["dataset"])
	assert.Equal(t, testDatasetID, result.DatasetId)
}

func TestUnitAimDatasetUpsertRecords_ReplaceOmitsInsert(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testUpsertRecordsResponseJSON, http.StatusOK)

	result, err := client.AimDatasetUpsertRecords(AimDatasetUpsertRecordsInput{
		DatasetId: testDatasetID,
		Records:   []AimDatasetRecordInput{{Input: "hello"}},
		Message:   "replace all records",
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, map[string]interface{}{
		"datasetId": testDatasetID,
		"records":   []interface{}{map[string]interface{}{"input": "hello"}},
		"message":   "replace all records",
	}, lastRequest.Variables["records"])
	assert.Equal(t, testDatasetID, result.DatasetId)
	assert.Equal(t, testVersionID, result.VersionId)
}

func TestUnitAimDatasetUpsertRecords_InsertAfterZeroPrepends(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testUpsertRecordsResponseJSON, http.StatusOK)

	afterID := 0
	_, err := client.AimDatasetUpsertRecords(AimDatasetUpsertRecordsInput{
		DatasetId: testDatasetID,
		Insert: &AimDatasetInsertRecordInput{
			AfterId: &afterID,
			Records: []AimDatasetRecordInput{{Input: "hello"}},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, map[string]interface{}{
		"datasetId": testDatasetID,
		"insert": map[string]interface{}{
			"afterId": float64(0),
			"records": []interface{}{map[string]interface{}{"input": "hello"}},
		},
	}, lastRequest.Variables["records"])
}

func TestUnitAimDatasetUpsertRecords_InsertWithoutAfterIdAppends(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testUpsertRecordsResponseJSON, http.StatusOK)

	_, err := client.AimDatasetUpsertRecords(AimDatasetUpsertRecordsInput{
		DatasetId: testDatasetID,
		Insert: &AimDatasetInsertRecordInput{
			Records: []AimDatasetRecordInput{{Input: "hello"}},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, map[string]interface{}{
		"datasetId": testDatasetID,
		"insert": map[string]interface{}{
			"records": []interface{}{map[string]interface{}{"input": "hello"}},
		},
	}, lastRequest.Variables["records"])
}

func TestUnitAimDatasetUpsertRecords_DeleteRecordIds(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testUpsertRecordsResponseJSON, http.StatusOK)

	_, err := client.AimDatasetUpsertRecords(AimDatasetUpsertRecordsInput{
		DatasetId:       testDatasetID,
		DeleteRecordIds: []int{3, 5},
	})
	require.NoError(t, err)

	assert.Equal(t, map[string]interface{}{
		"datasetId":       testDatasetID,
		"deleteRecordIds": []interface{}{float64(3), float64(5)},
	}, lastRequest.Variables["records"])
}

func TestUnitGetRecords(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testRecordsResponseJSON, http.StatusOK)

	result, err := client.GetRecords(11859098, testDatasetID, testVersionID, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, map[string]interface{}{
		"accountID":    float64(11859098),
		"aimDatasetId": testDatasetID,
		"versionId":    testVersionID,
		"cursor":       nil,
		"recordIds":    nil,
	}, lastRequest.Variables)

	assert.Equal(t, "next-page-cursor", result.NextCursor)
	assert.Equal(t, 9000, result.TotalCount)
	require.Len(t, result.Records, 2)

	assert.Equal(t, 1, result.Records[0].ID)
	assert.Equal(t, map[string]interface{}{"prompt": "hello"}, result.Records[0].Input)
	assert.Equal(t, "hi", result.Records[0].Expected)
	assert.Equal(t, "greeting", result.Records[0].Context)
	assert.Equal(t, map[string]interface{}{"lang": "en"}, result.Records[0].AdditionalFields)

	assert.Equal(t, 2, result.Records[1].ID)
	assert.Equal(t, []interface{}{map[string]interface{}{"role": "user", "content": "hello"}}, result.Records[1].Input)
}

func TestUnitGetRecords_SendsCursorAndRecordIds(t *testing.T) {
	t.Parallel()

	client, lastRequest := newMockResponseCapturingRequest(t, testRecordsResponseJSON, http.StatusOK)

	cursor := "current-page-cursor"
	_, err := client.GetRecords(11859098, testDatasetID, testVersionID, &cursor, []int{1, 2})
	require.NoError(t, err)

	assert.Contains(t, lastRequest.Query, "account(id: $accountID)")
	assert.Contains(t, lastRequest.Query, "cursor: $cursor")
	assert.Equal(t, cursor, lastRequest.Variables["cursor"])
	assert.Equal(t, []interface{}{float64(1), float64(2)}, lastRequest.Variables["recordIds"])
}

func TestUnitAimonitoring_ReturnsNerdGraphError(t *testing.T) {
	t.Parallel()

	errorResponseJSON := `{
		"errors": [
			{"message": "Not authorized to access dataset"}
		]
	}`

	afterID := 0
	cursor := "current-page-cursor"

	tests := map[string]func(client Aimonitoring) (interface{}, error){
		"AimDatasetCreate": func(client Aimonitoring) (interface{}, error) {
			return client.AimDatasetCreate(AimDatasetCreateInput{Name: "PII redaction test"})
		},
		"AimDatasetCreateFromCsv": func(client Aimonitoring) (interface{}, error) {
			return client.AimDatasetCreateFromCsv(AimDatasetCreateFromCsvInput{Name: "PII redaction test"})
		},
		"AimDatasetUpdate": func(client Aimonitoring) (interface{}, error) {
			return client.AimDatasetUpdate(AimDatasetUpdateInput{DatasetId: testDatasetID})
		},
		"AimDatasetUpsertRecords": func(client Aimonitoring) (interface{}, error) {
			return client.AimDatasetUpsertRecords(AimDatasetUpsertRecordsInput{
				DatasetId: testDatasetID,
				Insert:    &AimDatasetInsertRecordInput{AfterId: &afterID},
			})
		},
		"AimDatasetDelete": func(client Aimonitoring) (interface{}, error) {
			return client.AimDatasetDelete(AimDatasetDeleteInput{DatasetId: testDatasetID})
		},
		"AimDatasetDeleteVersion": func(client Aimonitoring) (interface{}, error) {
			return client.AimDatasetDeleteVersion(AimDatasetDeleteVersionInput{DatasetId: testDatasetID, VersionId: testVersionID})
		},
		"GetEntitySearch": func(client Aimonitoring) (interface{}, error) {
			return client.GetEntitySearch(&cursor, "type = 'AIM_DATASET'")
		},
		"GetRecords": func(client Aimonitoring) (interface{}, error) {
			return client.GetRecords(11859098, testDatasetID, testVersionID, &cursor, []int{1})
		},
	}

	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, _ := newMockResponseCapturingRequest(t, errorResponseJSON, http.StatusOK)

			result, err := call(client)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "Not authorized to access dataset")
			assert.Nil(t, result)
		})
	}
}
