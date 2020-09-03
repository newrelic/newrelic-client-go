// +build unit

package apm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetApplicationInstance(t *testing.T) {
	t.Parallel()

	testApplicationID := 204261410
	testApplicationInstanceID := 317218023
	testApplicationHostID := 204260579

	testApplicationInstanceSummary := ApplicationSummary{
		ResponseTime:  5.91,
		Throughput:    1,
		ErrorRate:     0,
		ApdexScore:    1,
		InstanceCount: 1,
	}

	testApplicationInstanceEndUserSummary := ApplicationEndUserSummary{
		ResponseTime: 5.91,
		Throughput:   1,
		ApdexTarget:  0,
		ApdexScore:   1,
	}

	testApplicationInstanceLinks := ApplicationInstanceLinks{
		Application:     testApplicationID,
		ApplicationHost: testApplicationHostID,
	}

	testApplicationInstance := ApplicationInstance{
		ID:              testApplicationInstanceID,
		ApplicationName: "Billing Service",
		Host:            "host",
		Port:            80,
		Language:        "python",
		HealthStatus:    "unknown",
		Summary:         testApplicationInstanceSummary,
		EndUserSummary:  testApplicationInstanceEndUserSummary,
		Links:           testApplicationInstanceLinks,
	}

	testApplicationInstanceJson := fmt.Sprintf(`{
		"id": %d,
		"application_name": "Billing Service",
		"language": "python",
		"port": 80,
		"host": "host",
		"health_status": "unknown",
		"application_summary": {
			"response_time": 5.91,
			"throughput": 1,
			"error_rate": 0,
			"apdex_score": 1,
			"instance_count": 1
		},
		"end_user_summary": {
			"response_time": 5.91,
			"throughput": 1,
			"apdex_score": 1
		  },
		"links": {
			"application": %d,
			"application_host": %d
		}
	}`, testApplicationInstanceID, testApplicationID, testApplicationHostID)

	responseJSON := fmt.Sprintf(`{ "application_instance": %s}`, testApplicationInstanceJson)
	apm := newMockResponse(t, responseJSON, http.StatusOK)

	actual, err := apm.GetApplicationInstance(testApplication.ID, testApplicationInstance.ID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testApplicationInstance, actual)
}
