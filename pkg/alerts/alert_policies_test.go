// +build unit

package alerts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func NewTestAlerts(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func TestGetAlertPolicy(t *testing.T) {
	t.Parallel()
	alerts := NewTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`
		{
			"policies": [
				{
					"id": 579506,
					"incident_preference": "PER_POLICY",
					"name": "test-alert-policy-1",
					"created_at": 1575438237690,
					"updated_at": 1575438237690
				},
				{
					"id": 579509,
					"incident_preference": "PER_POLICY",
					"name": "test-alert-policy-2",
					"created_at": 1575438240632,
					"updated_at": 1575438240632
				},
				{
					"id": 579510,
					"incident_preference": "PER_POLICY",
					"name": "alert",
					"created_at": 1575438240631,
					"updated_at": 1575438240631
				},
				{
					"id": 579511,
					"incident_preference": "PER_POLICY",
					"name": "alert",
					"created_at": 1575438240633,
					"updated_at": 1575438240633
				}
			]
		}
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	// GetAlertPolicy returns a pointer *AlertPolicy
	expected := &AlertPolicy{
		ID:                 579506,
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
		CreatedAt:          1575438237690,
		UpdatedAt:          1575438237690,
	}

	actual, err := alerts.GetAlertPolicy(579506)

	if err != nil {
		t.Fatalf("GetAlertPolicy error: %s", err)
	}

	if actual == nil {
		t.Fatalf("GetAlertPolicy result is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("GetAlertPolicy result differs from expected: %s", diff)
	}
}

func TestListAlertPolicies(t *testing.T) {
	t.Parallel()
	alerts := NewTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`
		{
			"policies": [
				{
					"id": 579506,
					"incident_preference": "PER_POLICY",
					"name": "test-alert-policy-1",
					"created_at": 1575438237690,
					"updated_at": 1575438237690
				},
				{
					"id": 579509,
					"incident_preference": "PER_POLICY",
					"name": "test-alert-policy-2",
					"created_at": 1575438240632,
					"updated_at": 1575438240632
				}
			]
		}
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	expected := []AlertPolicy{
		{
			ID:                 579506,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-1",
			CreatedAt:          1575438237690,
			UpdatedAt:          1575438237690,
		},
		{
			ID:                 579509,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-2",
			CreatedAt:          1575438240632,
			UpdatedAt:          1575438240632,
		},
	}

	actual, err := alerts.ListAlertPolicies(nil)

	if err != nil {
		t.Fatalf("ListAlertPolicies error: %s", err)
	}

	if actual == nil {
		t.Fatalf("ListAlertPolicies result is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("ListAlertPolicies result differs from expected: %s", diff)
	}
}

func TestListAlertPoliciesWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "test-alert-policy-1"

	alerts := NewTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		if name != expectedName {
			t.Errorf(`expected name filter "%s", recieved: "%s"`, expectedName, name)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`
		{
			"policies": [
				{
					"id": 579506,
					"incident_preference": "PER_POLICY",
					"name": "test-alert-policy-1",
					"created_at": 1575438237690,
					"updated_at": 1575438237690
				}
			]
		}
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	expected := []AlertPolicy{
		{
			ID:                 579506,
			IncidentPreference: "PER_POLICY",
			Name:               "test-alert-policy-1",
			CreatedAt:          1575438237690,
			UpdatedAt:          1575438237690,
		},
	}

	params := ListAlertPoliciesParams{
		Name: &expectedName,
	}

	actual, err := alerts.ListAlertPolicies(&params)

	if err != nil {
		t.Fatalf("ListApplications error: %s", err)
	}

	if actual == nil {
		t.Fatalf("ListApplications response is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("ListApplications response differs from expected: %s", diff)
	}
}

func TestCreateAlertPolicy(t *testing.T) {
	t.Parallel()
	alerts := NewTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`
		{
			"policy": {
				"id": 123,
				"incident_preference": "PER_POLICY",
				"name": "test-alert-policy-1",
				"created_at": 1575438237690,
				"updated_at": 1575438237690
			}
		}
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	policy := AlertPolicy{
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
	}

	expected := &AlertPolicy{
		ID:                 123,
		IncidentPreference: "PER_POLICY",
		Name:               "test-alert-policy-1",
		CreatedAt:          1575438237690,
		UpdatedAt:          1575438237690,
	}

	actual, err := alerts.CreateAlertPolicy(policy)

	if err != nil {
		t.Fatalf("CreateAlertPolicy error: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateAlertPolicy result is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("CreateAlertPolicy result differs from expected: %s", diff)
	}
}
