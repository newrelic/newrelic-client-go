// +build unit

package apm

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testApplicationInstanceLinks = ApplicationInstanceLinks{
		Application:     204261410,
		ApplicationHost: 204260579,
	}

	testApplicationInstance = ApplicationInstance{
		ID:              317218023,
		ApplicationName: "Billing Service",
		Host:            "host",
		Port:            80,
		Language:        "python",
		HealthStatus:    "unknown",
		Summary:         testApplicationSummary,
		Links:           testApplicationInstanceLinks,
	}

	testApplicationInstancesJson = `{
		"id": 317218023,
		"name": "Billing Service",
		"language": "python",
		"health_status": "unknown",
		"application_summary": {
			"response_time": 5.91,
			"throughput": 1,
			"error_rate": 0,
			"apdex_score": 1,
			"instance_count": 15,
		},
		"links": {
			"application": 204261410,
			"application_host": 204260579
		}
	}`
)

func TestListApplicationInstancesWithParams(t *testing.T) {
	t.Parallel()
	expectedHost := "appHost"
	expectedIDs := "123,456"

	apm := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		host := values.Get("filter[hostname]")
		ids := values.Get("filter[ids]")

		assert.Equal(t, expectedHost, host)
		assert.Equal(t, expectedIDs, ids)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"application_instances":[]}`))

		assert.NoError(t, err)
	}))

	params := ListApplicationInstancesParams{
		Hostname: expectedHost,
		IDs:      []int{123, 456},
	}

	_, err := apm.ListApplicationInstances(testApplication.ID, &params)

	assert.NoError(t, err)
}
