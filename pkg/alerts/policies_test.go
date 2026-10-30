//go:build unit
// +build unit

package alerts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/internal/serialization"
)

var (
	testTimestampStringMs = "1575438237690"
	testTimestamp         = serialization.EpochTime(time.Unix(0, 1575438237690*int64(time.Millisecond)).UTC())

	testPoliciesResponseJSON = `{
		"policies": [
			{
				"id": 579506,
				"incident_preference": "PER_POLICY",
				"name": "test-alert-policy-1",
				"created_at": ` + testTimestampStringMs + `,
				"updated_at": ` + testTimestampStringMs + `
			},
			{
				"id": 579509,
				"incident_preference": "PER_POLICY",
				"name": "test-alert-policy-2",
				"created_at": ` + testTimestampStringMs + `,
				"updated_at": ` + testTimestampStringMs + `
			}
		]
	}`

	testPolicyResponseJSON = `{
		"policy": {
			"id": 579506,
			"incident_preference": "PER_POLICY",
			"name": "test-alert-policy-1",
			"created_at": ` + testTimestampStringMs + `,
			"updated_at": ` + testTimestampStringMs + `
		}
	}`

	testPolicyResponseUpdatedJSON = `{
		"policy": {
			"id": 579506,
			"incident_preference": "PER_CONDITION",
			"name": "test-alert-policy-updated",
			"created_at": ` + testTimestampStringMs + `,
			"updated_at": ` + testTimestampStringMs + `
		}
	}`
)

func TestGetPolicy(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testPoliciesResponseJSON, http.StatusOK)

	expected := &Policy{
		ID:                 579506,
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               "test-alert-policy-1",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.GetPolicy(579506)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestListPolicies(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testPoliciesResponseJSON, http.StatusOK)

	expected := []Policy{
		{
			ID:                 579506,
			IncidentPreference: IncidentPreferenceTypes.PerPolicy,
			Name:               "test-alert-policy-1",
			CreatedAt:          &testTimestamp,
			UpdatedAt:          &testTimestamp,
		},
		{
			ID:                 579509,
			IncidentPreference: IncidentPreferenceTypes.PerPolicy,
			Name:               "test-alert-policy-2",
			CreatedAt:          &testTimestamp,
			UpdatedAt:          &testTimestamp,
		},
	}

	actual, err := alerts.ListPolicies(nil)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestListPoliciesWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "does-not-exist"

	alerts := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		require.Equal(t, expectedName, name)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{ "policies": [] }`))

		require.NoError(t, err)
	}))

	params := ListPoliciesParams{
		Name: expectedName,
	}

	expectedCount := 0

	actual, err := alerts.ListPolicies(&params)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expectedCount, len(actual))
}

func TestCreatePolicy(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testPolicyResponseJSON, http.StatusOK)

	policy := Policy{
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               "test-alert-policy-1",
	}

	expected := &Policy{
		ID:                 579506,
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               "test-alert-policy-1",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.CreatePolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestUpdatePolicy(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testPolicyResponseUpdatedJSON, http.StatusOK)

	policy := Policy{
		ID:                 579506,
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               "test-alert-policy-1",
	}

	expected := &Policy{
		ID:                 579506,
		IncidentPreference: IncidentPreferenceTypes.PerCondition,
		Name:               "test-alert-policy-updated",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.UpdatePolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestDeletePolicy(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testPolicyResponseJSON, http.StatusOK)

	expected := &Policy{
		ID:                 579506,
		IncidentPreference: IncidentPreferenceTypes.PerPolicy,
		Name:               "test-alert-policy-1",
		CreatedAt:          &testTimestamp,
		UpdatedAt:          &testTimestamp,
	}

	actual, err := alerts.DeletePolicy(579506)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

// NerdGraph Policy Tests (with entityGuid support)

var (
	testAlertsPolicyQueryResponseJSON = `{
		"data": {
			"actor": {
				"account": {
					"alerts": {
						"policy": {
							"id": "123456",
							"name": "test-alert-policy-1",
							"incidentPreference": "PER_POLICY",
							"accountId": 123456,
							"entityGuid": "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2"
						}
					}
				}
			}
		}
	}`

	testAlertsPolicySearchResponseJSON = `{
		"data": {
			"actor": {
				"account": {
					"alerts": {
						"policiesSearch": {
							"nextCursor": null,
							"totalCount": 2,
							"policies": [
								{
									"id": "123456",
									"name": "test-alert-policy-1",
									"incidentPreference": "PER_POLICY",
									"accountId": 123456,
									"entityGuid": "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2"
								},
								{
									"id": "123457",
									"name": "test-alert-policy-2",
									"incidentPreference": "PER_CONDITION",
									"accountId": 123456,
									"entityGuid": "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU3"
								}
							]
						}
					}
				}
			}
		}
	}`

	testAlertsPolicyCreateResponseJSON = `{
		"data": {
			"alertsPolicyCreate": {
				"id": "123456",
				"name": "test-alert-policy-1",
				"incidentPreference": "PER_POLICY",
				"accountId": 123456,
				"entityGuid": "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2"
			}
		}
	}`

	testAlertsPolicyUpdateResponseJSON = `{
		"data": {
			"alertsPolicyUpdate": {
				"id": "123456",
				"name": "test-alert-policy-updated",
				"incidentPreference": "PER_CONDITION",
				"accountId": 123456,
				"entityGuid": "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2"
			}
		}
	}`
)

func TestQueryPolicy(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testAlertsPolicyQueryResponseJSON, http.StatusOK)

	expected := &AlertsPolicy{
		ID:                 "123456",
		Name:               "test-alert-policy-1",
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		AccountID:          123456,
		EntityGuid:         "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2",
	}

	actual, err := alerts.QueryPolicy(123456, "123456")

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.IncidentPreference, actual.IncidentPreference)
	require.Equal(t, expected.AccountID, actual.AccountID)
	require.Equal(t, expected.EntityGuid, actual.EntityGuid)
}

func TestQueryPolicySearch(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testAlertsPolicySearchResponseJSON, http.StatusOK)

	actual, err := alerts.QueryPolicySearch(123456, AlertsPoliciesSearchCriteriaInput{})

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Len(t, actual, 2)

	// Check first policy
	require.Equal(t, "123456", actual[0].ID)
	require.Equal(t, "test-alert-policy-1", actual[0].Name)
	require.Equal(t, "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2", actual[0].EntityGuid)

	// Check second policy
	require.Equal(t, "123457", actual[1].ID)
	require.Equal(t, "test-alert-policy-2", actual[1].Name)
	require.Equal(t, "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU3", actual[1].EntityGuid)
}

func TestCreatePolicyMutation(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testAlertsPolicyCreateResponseJSON, http.StatusOK)

	policy := AlertsPolicyInput{
		Name:               "test-alert-policy-1",
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
	}

	expected := &AlertsPolicy{
		ID:                 "123456",
		Name:               "test-alert-policy-1",
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_POLICY,
		AccountID:          123456,
		EntityGuid:         "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2",
	}

	actual, err := alerts.CreatePolicyMutation(123456, policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.IncidentPreference, actual.IncidentPreference)
	require.Equal(t, expected.AccountID, actual.AccountID)
	require.Equal(t, expected.EntityGuid, actual.EntityGuid)
}

func TestUpdatePolicyMutation(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testAlertsPolicyUpdateResponseJSON, http.StatusOK)

	policy := AlertsPolicyUpdateInput{
		Name:               "test-alert-policy-updated",
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_CONDITION,
	}

	expected := &AlertsPolicy{
		ID:                 "123456",
		Name:               "test-alert-policy-updated",
		IncidentPreference: AlertsIncidentPreferenceTypes.PER_CONDITION,
		AccountID:          123456,
		EntityGuid:         "MTIzNDU2fEFMRVJUU3xQT0xJQ1l8MTIzNDU2",
	}

	actual, err := alerts.UpdatePolicyMutation(123456, "123456", policy)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.IncidentPreference, actual.IncidentPreference)
	require.Equal(t, expected.AccountID, actual.AccountID)
	require.Equal(t, expected.EntityGuid, actual.EntityGuid)
}

func TestQueryPolicySearchByName(t *testing.T) {
	t.Parallel()
	expectedName := "test-alert-policy-1"
	respJSON := fmt.Sprintf(`{"data":{"actor":{"account":{"alerts":{"policiesSearch":{"nextCursor":null,"totalCount":1,"policies":[{"id":"123456","name":%q,"incidentPreference":"PER_POLICY","accountId":123456}]}}}}}}`, expectedName)

	alerts := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Variables struct {
				SearchCriteria struct {
					Name     string `json:"name"`
					NameLike string `json:"nameLike"`
				} `json:"criteria"`
			} `json:"variables"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, expectedName, body.Variables.SearchCriteria.Name)
		require.Empty(t, body.Variables.SearchCriteria.NameLike)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respJSON))
	}))

	actual, err := alerts.QueryPolicySearch(123456, AlertsPoliciesSearchCriteriaInput{Name: expectedName})

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Len(t, actual, 1)
	require.Equal(t, expectedName, actual[0].Name)
}

func TestQueryPolicySearchByNameLike(t *testing.T) {
	t.Parallel()
	expectedNameLike := "test-alert"
	respJSON := `{"data":{"actor":{"account":{"alerts":{"policiesSearch":{"nextCursor":null,"totalCount":2,"policies":[{"id":"123456","name":"test-alert-policy-1","incidentPreference":"PER_POLICY","accountId":123456},{"id":"123457","name":"test-alert-policy-2","incidentPreference":"PER_POLICY","accountId":123456}]}}}}}}`

	alerts := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Variables struct {
				SearchCriteria struct {
					Name     string `json:"name"`
					NameLike string `json:"nameLike"`
				} `json:"criteria"`
			} `json:"variables"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, expectedNameLike, body.Variables.SearchCriteria.NameLike)
		require.Empty(t, body.Variables.SearchCriteria.Name)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respJSON))
	}))

	actual, err := alerts.QueryPolicySearch(123456, AlertsPoliciesSearchCriteriaInput{NameLike: expectedNameLike})

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Len(t, actual, 2)
	require.Equal(t, "test-alert-policy-1", actual[0].Name)
	require.Equal(t, "test-alert-policy-2", actual[1].Name)
}

func TestQueryPolicySearchByIDs(t *testing.T) {
	t.Parallel()
	expectedIDs := []string{"123456", "123457"}
	respJSON := `{"data":{"actor":{"account":{"alerts":{"policiesSearch":{"nextCursor":null,"totalCount":2,"policies":[{"id":"123456","name":"test-alert-policy-1","incidentPreference":"PER_POLICY","accountId":123456},{"id":"123457","name":"test-alert-policy-2","incidentPreference":"PER_POLICY","accountId":123456}]}}}}}}`

	alerts := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Variables struct {
				Criteria struct {
					IDs []string `json:"ids"`
				} `json:"criteria"`
			} `json:"variables"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, expectedIDs, body.Variables.Criteria.IDs)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respJSON))
	}))

	actual, err := alerts.QueryPolicySearch(123456, AlertsPoliciesSearchCriteriaInput{IDs: expectedIDs})

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Len(t, actual, 2)
	require.Equal(t, "123456", actual[0].ID)
	require.Equal(t, "123457", actual[1].ID)
}

// TestQueryPolicySearchPaginatesWithCriteria verifies that when NerdGraph returns a
// nextCursor the criteria variable is carried on every subsequent page request —
// not just the first. This is the regression test for the "searchCriteria" vs
// "criteria" variable name mismatch: a wrong key would silently pass on page 1
// but a server-side filter would still be absent on all pages.
func TestQueryPolicySearchPaginatesWithCriteria(t *testing.T) {
	t.Parallel()
	expectedName := "test-alert-policy-1"
	cursor := "next-page-cursor"
	requestCount := 0

	page1JSON := fmt.Sprintf(`{"data":{"actor":{"account":{"alerts":{"policiesSearch":{"nextCursor":%q,"totalCount":2,"policies":[{"id":"123456","name":%q,"incidentPreference":"PER_POLICY","accountId":123456}]}}}}}}`, cursor, expectedName)
	page2JSON := fmt.Sprintf(`{"data":{"actor":{"account":{"alerts":{"policiesSearch":{"nextCursor":null,"totalCount":2,"policies":[{"id":"123457","name":%q,"incidentPreference":"PER_POLICY","accountId":123456}]}}}}}}`, expectedName)

	alerts := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Variables struct {
				Criteria struct {
					Name string `json:"name"`
				} `json:"criteria"`
				Cursor *string `json:"cursor"`
			} `json:"variables"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		// criteria.name must be present on every page request
		require.Equal(t, expectedName, body.Variables.Criteria.Name, "criteria.name must be sent on request %d", requestCount+1)

		w.Header().Set("Content-Type", "application/json")
		requestCount++
		if requestCount == 1 {
			_, _ = w.Write([]byte(page1JSON))
		} else {
			_, _ = w.Write([]byte(page2JSON))
		}
	}))

	actual, err := alerts.QueryPolicySearch(123456, AlertsPoliciesSearchCriteriaInput{Name: expectedName})

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Len(t, actual, 2, "both pages should be aggregated")
	require.Equal(t, 2, requestCount, "should have made exactly 2 requests (one per page)")
	require.Equal(t, "123456", actual[0].ID)
	require.Equal(t, "123457", actual[1].ID)
}
