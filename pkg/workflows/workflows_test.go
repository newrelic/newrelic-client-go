//go:build integration
// +build integration

package workflows

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/newrelic"
	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Workflows {
	cfg := mock.NewIntegrationTestConfig(t)
	client := New(cfg)

	return client
}

func newrelicIntegrationTestClient(t *testing.T) newrelic.NewRelic {
	tc := mock.NewIntegrationTestConfig(t)

	client, _ := newrelic.New(
		newrelic.ConfigPersonalAPIKey(tc.PersonalAPIKey),
		newrelic.ConfigInsightsInsertKey(tc.InsightsInsertKey),
		newrelic.ConfigUserAgent(tc.UserAgent),
	)

	return *client
}
