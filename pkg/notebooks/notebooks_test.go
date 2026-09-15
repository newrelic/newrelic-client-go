//go:build unit || integration
// +build unit integration

package notebooks

import (
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Notebooks {
	cfg := mock.NewIntegrationTestConfig(t)
	return New(cfg)
}

var testOrganizationID = os.Getenv("INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID")
