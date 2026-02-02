//go:build unit || integration
// +build unit integration

package fleetcontrol

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Fleetcontrol {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Fleetcontrol {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

var (
	testOrganizationID = "fb33fea3-4d7e-4736-9701-acb59a634fdf"
	// testOrganizationID = "b961cf81-d62b-4359-8822-7b1d6dadd374"
)
