// +build integration

package events

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_batch() {
	// Initialize the client configuration.  An Insights insert key is required
	// to communicate with the backend API.
	cfg := config.New()
	cfg.InsightsInsertKey = os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")

	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	accountID, err := strconv.Atoi(envAccountID)
	if err != nil {
		log.Fatal("environment variable NEW_RELIC_ACCOUNT_ID required")
	}

	// Initialize the client.
	client := New(cfg)

	// Start batch mode
	if err := client.BatchMode(context.Background(), accountID); err != nil {
		log.Fatal("error starting batch mode:", err)
	}

	event := struct {
		EventType string  `json:"eventType"`
		Amount    float64 `json:"amount"`
	}{
		EventType: "Purchase",
		Amount:    123.45,
	}

	// Queueu a custom event.
	if err := client.EnqueueEvent(context.Background(), event); err != nil {
		log.Fatal("error posting custom event:", err)
	}

	// Force flush events
	if err := client.Flush(); err != nil {
		log.Fatal("error flushing event queue:", err)
	}
}
