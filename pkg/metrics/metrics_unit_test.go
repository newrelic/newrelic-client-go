//go:build unit
// +build unit

package metrics

import (
	"net/http"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"

	"github.com/stretchr/testify/assert"
)

var (
	testMetrics = `{
		"name": "service.errors.all",
		"type": "gauge",
		"value": 9,
		"attributes": {
			"service.response.statuscode": "400"
		}
	}`
)

func getMockMetricsResponse(t *testing.T, mockJSONResponse string, statusCode int) Metrics {
	mockServer := mock.NewMockServer(t, mockJSONResponse, statusCode)
	testConfig := mock.NewTestConfig(t, mockServer)

	return New(testConfig)
}

func TestCreateMetricEntry(t *testing.T) {
	t.Parallel()

	// Arrange
	metricsClient := getMockMetricsResponse(t, ``, http.StatusAccepted)

	// Act
	err := metricsClient.CreateMetricEntry(testMetrics)

	// Assert successful metrics call
	assert.NoError(t, err)

}

func TestNilMetricEntry_FailsWithError(t *testing.T) {
	t.Parallel()

	// Arrange
	metricsClient := getMockMetricsResponse(t, ``, http.StatusAccepted)

	// Act
	err := metricsClient.CreateMetricEntry(nil)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Equal(t, "metrics: CreateMetricEntry: metricEntry is nil, nothing to do", err.Error())
}
