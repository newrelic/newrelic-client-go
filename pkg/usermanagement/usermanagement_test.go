package usermanagement

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Usermanagement {
	cfg := mock.NewIntegrationTestConfig(t)
	client := New(cfg)

	return client
}
