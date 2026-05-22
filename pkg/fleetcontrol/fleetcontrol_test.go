//go:build unit || integration
// +build unit integration

package fleetcontrol

import (
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Fleetcontrol {
	tc := mock.NewFleetIntegrationTestConfig(t)
	return New(tc)
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Fleetcontrol {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func getTestOrganizationID() string {
	if id := os.Getenv("NEW_RELIC_FLEET_TEST_ORGANIZATION_ID"); id != "" {
		return id
	}
	return "b961cf81-d62b-4359-8822-7b1d6dadd374"
}

var testOrganizationID = getTestOrganizationID()
