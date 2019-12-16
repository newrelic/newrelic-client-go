// +build integration

package infrastructure

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestIntegrationListAlertConditions(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey: apiKey,
	})

	c, err := api.ListAlertConditions(586577)

	if err != nil {
		t.Fatalf("ListAlertConditions error: %s", err)
	}

	if len(c) != 2 {
		t.Fatalf("expected 2 conditions, received %d", len(c))
	}
}
