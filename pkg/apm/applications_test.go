package apm

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	client "github.com/newrelic/newrelic-client-go/internal"
)

func NewTestAPM(handler http.Handler) APM {
	ts := httptest.NewServer(handler)

	c := New(client.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func TestListApplications(t *testing.T) {
	apm := NewTestAPM(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`
		{
			"applications": [
				{
					"id": 204261410,
					"name": "Billing Service",
					"language": "python",
					"health_status": "unknown",
					"reporting": true,
					"last_reported_at": "2019-12-11T19:09:10+00:00",
					"application_summary": {
						"response_time": 5.91,
						"throughput": 1,
						"error_rate": 0,
						"apdex_target": 0.5,
						"apdex_score": 1,
						"host_count": 1,
						"instance_count": 15
					},
					"end_user_summary": {
						"response_time": 3.8,
						"throughput": 1660,
						"apdex_target": 2.5,
						"apdex_score": 0.78
					},
					"settings": {
						"app_apdex_threshold": 0.5,
						"end_user_apdex_threshold": 7,
						"enable_real_user_monitoring": true,
						"use_server_side_config": false
					},
					"links": {
						"application_instances": [
							204261411
						],
						"servers": [],
						"application_hosts": [
							204260579
						]
					}
				}
			]
		}
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	applicationSummary := ApplicationSummary{
		ResponseTime:            5.91,
		Throughput:              1,
		ErrorRate:               0,
		ApdexTarget:             0.5,
		ApdexScore:              1,
		HostCount:               1,
		InstanceCount:           1,
		ConcurrentInstanceCount: 15,
	}

	applicationEndUserSummary := ApplicationEndUserSummary{
		ResponseTime: 3.8,
		Throughput:   1660,
		ApdexTarget:  2.5,
		ApdexScore:   0.78,
	}

	applicationSettings := ApplicationSettings{
		AppApdexThreshold:        0.5,
		EndUserApdexThreshold:    7,
		EnableRealUserMonitoring: true,
		UseServerSideConfig:      false,
	}

	applicationLinks := ApplicationLinks{
		ServerIDs:     []int{},
		HostIDs:       []int{204260579},
		InstanceIDs:   []int{204261411},
		AlertPolicyID: 1234,
	}

	expected := []Application{
		{
			ID:             204261410,
			Name:           "Billing Service",
			Language:       "python",
			HealthStatus:   "unknown",
			Reporting:      true,
			LastReportedAt: "2019-12-11T19:09:10+00:00",
			Summary:        applicationSummary,
			EndUserSummary: applicationEndUserSummary,
			Settings:       applicationSettings,
			Links:          applicationLinks,
		},
	}

	actual, err := apm.ListApplications()

	if err != nil {
		t.Log(err)
		t.Fatal("ListApplications error")
	}

	if actual == nil {
		t.Log(err)
		t.Fatal("ListApplications error")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("applications not parsed correctly: %s", diff)
	}
}
