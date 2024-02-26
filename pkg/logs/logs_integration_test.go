//go:build integration
// +build integration

// Requires NEW_RELIC_LICENSE_KEY envvar (APM License Key)

package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	nr "github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

type testDatum struct {
	Datum interface{}
	err   error
}

var testData = []testDatum{
	{
		Datum: struct {
			Message string `json:"message"`
		}{
			Message: "INFO: simple message test",
		},
		err: nil,
	},
}

func TestIntegrationLogs(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	for _, test := range testData {
		err := client.CreateLogEntry(test.Datum)
		if test.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, test.err, err)
		}
	}
}

func newIntegrationTestClient(t *testing.T) Logs {
	tc := nr.NewIntegrationTestConfig(t)

	return New(tc)
}
