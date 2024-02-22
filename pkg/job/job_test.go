//go:build unit || integration
// +build unit integration

package job

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Job {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Job {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

var (
	contractId = "aaaaaaaa-8428-4762-acdf-858bb6fd6db2"
	customerId = "CC-0000000000"

	org1Id   = "aaaaaaaa-3f0b-4be7-8de8-319c23bdf9e8"
	org1Name = "Org 1"

	jobId = "aaaaaaaa-0d73-4db6-a53c-8d56feee6f0f"
)
