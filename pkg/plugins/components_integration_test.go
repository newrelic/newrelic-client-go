// +build integration

package plugins

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestIntegrationComponents(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEW_RELIC_ADMIN_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires NEW_RELIC_ADMIN_API_KEY to be set")
	}

	api := New(config.Config{
		AdminAPIKey: apiKey,
		LogLevel:    "debug",
	})

	a, err := api.ListComponents(nil)

	require.NoError(t, err)
	require.NotNil(t, a)

	if len(a) < 1 {
		t.Skipf("Skipping `GetComponent` integration test due to zero plugins found")
	}

	c, err := api.GetComponent(a[0].ID)

	require.NoError(t, err)
	require.NotNil(t, c)

	m, err := api.ListComponentMetrics(c.ID, nil)

	require.NoError(t, err)
	require.NotNil(t, m)

	if len(m) < 1 {
		t.Skipf("Skipping `GetComponentMetricData` integration test due to zero plugin metrics found")
	}
	params := GetComponentMetricDataParams{
		Names: []string{m[0].Name},
	}
	_, err = api.GetComponentMetricData(a[0].ID, &params)

	require.NoError(t, err)
}
