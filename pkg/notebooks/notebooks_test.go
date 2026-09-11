//go:build unit || integration
// +build unit integration

package notebooks

import (
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// newIntegrationTestClient builds a Notebooks client bound to the Fleet test
// account, which is the only account in our current test infrastructure with
// the Notebooks entitlement enabled. Uses NEW_RELIC_FLEET_TEST_API_KEY.
func newIntegrationTestClient(t *testing.T) Notebooks {
	cfg := mock.NewFleetIntegrationTestConfig(t)
	return New(cfg)
}

// getTestOrganizationID returns the organization UUID for Blob Storage API
// calls. Falls back to the known fleet-test org UUID when the env var is not
// set - keeping the tests runnable from a fresh checkout without extra
// configuration.
func getTestOrganizationID() string {
	if id := os.Getenv("NEW_RELIC_FLEET_TEST_ORGANIZATION_ID"); id != "" {
		return id
	}
	return "b961cf81-d62b-4359-8822-7b1d6dadd374"
}

var testOrganizationID = getTestOrganizationID()
