//go:build unit
// +build unit

package alerts

import (
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
