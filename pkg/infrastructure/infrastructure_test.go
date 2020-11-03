// +build unit

package infrastructure

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Infrastructure {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}
