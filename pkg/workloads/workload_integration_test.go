// +build integration

package workloads

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIntegrationGetWorkload(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	actual, err := client.GetWorkload(2508259, 791)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

func TestIntegrationListWorkloads(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	actual, err := client.ListWorkloads(2508259)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

// nolint
func newIntegrationTestClient(t *testing.T) Workloads {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing for graphql requires your personal API key")
	}

	return New(config.Config{
		PersonalAPIKey: apiKey,
		UserAgent:      "newrelic/newrelic-client-go",
		LogLevel:       "debug",
	})
}
