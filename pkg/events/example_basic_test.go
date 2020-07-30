// +build integration

package events

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	nr "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func Example_basic() {
	// Initialize the client configuration.  An Insights insert key is required
	// to communicate with the backend API.
	cfg := config.New()
	cfg.InsightsInsertKey = os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")

	// Initialize the client.
	client := New(cfg)

	event := struct {
		EventType string  `json:"eventType"`
		Amount    float64 `json:"amount"`
	}{
		EventType: "Purchase",
		Amount:    123.45,
	}

	// Post a custom event.
	if err := client.CreateEvent(nr.TestAccountID, event); err != nil {
		log.Fatal("error posting custom event:", err)
	}
}
