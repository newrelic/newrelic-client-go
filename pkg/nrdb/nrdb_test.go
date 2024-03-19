//go:build unit || integration
// +build unit integration

package nrdb

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Nrdb {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func newNRDBIntegrationTestClient(t *testing.T) Nrdb {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
