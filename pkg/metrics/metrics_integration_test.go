//go:build integration
// +build integration

// Requires NEW_RELIC_LICENSE_KEY env var

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	nr "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

type testDatum struct {
	Datum interface{}
	err   error
}

var testData = []testDatum{
	{
		Datum: struct {
			Name       string            `json:"name"`
			Type       string            `json:"type"`
			Value      int               `json:"value"`
			Attributes map[string]string `json:"attributes"`
		}{
			Name:  "service.errors.all",
			Type:  "gauge",
			Value: 9,
			Attributes: map[string]string{
				"service.response.statuscode": "400",
			},
		},
		err: nil,
	},
}

func TestIntegrationMetrics(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	for _, test := range testData {
		err := client.CreateMetricsEntry(test.Datum)
		if test.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, test.err, err)
		}
	}
}

func newIntegrationTestClient(t *testing.T) Metrics {
	tc := nr.NewIntegrationTestConfig(t)

	return New(tc)
}
