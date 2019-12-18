// +build unit

package apm

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func NewTestAPM(handler http.Handler) APM {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func TestListApplications(t *testing.T) {
	t.Parallel()
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
						"instance_count": 15,
						"concurrent_instance_count": 1
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
						],
						"alert_policy": 1234
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
		InstanceCount:           15,
		ConcurrentInstanceCount: 1,
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

	actual, err := apm.ListApplications(nil)

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

func TestListApplicationsWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "appName"
	expectedHost := "appHost"
	expectedLanguage := "appLanguage"
	expectedIDs := "123,456"

	apm := NewTestAPM(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		if name != expectedName {
			t.Errorf(`expected name filter "%s", recieved: "%s"`, expectedName, name)
		}

		host := values.Get("filter[host]")
		if host != expectedHost {
			t.Errorf(`expected host filter "%s", recieved: "%s"`, expectedHost, host)
		}

		ids := values.Get("filter[ids]")
		if ids != expectedIDs {
			t.Errorf(`expected ids filter "%s", recieved: "%s"`, expectedIDs, ids)
		}

		language := values.Get("filter[language]")
		if language != expectedLanguage {
			t.Errorf(`expected language filter "%s", recieved: "%s"`, expectedLanguage, language)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"applications":[]}`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	params := ListApplicationsParams{
		Name:     &expectedName,
		Host:     &expectedHost,
		IDs:      []int{123, 456},
		Language: &expectedLanguage,
	}

	_, err := apm.ListApplications(&params)

	if err != nil {
		t.Fatalf("ListApplications error: %s", err)
	}
}
