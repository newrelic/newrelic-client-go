package edge

import (
	"log"
	"os"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v3/pkg/config"
)

func Example_trace_observer() {
	accountIDStr := os.Getenv("ACCOUNT_ID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Fatal("error parsing account ID")
	}

	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Create a new trace observer.
	traceObserver, err := client.CreateTraceObserver(accountID, "myObserver", EdgeProviderRegionTypes.AWS_US_EAST_1)
	if err != nil {
		log.Fatal("error creating trace observer:", err)
	}

	// List the existing trace observers.
	traceObservers, err := client.ListTraceObservers(accountID)
	if err != nil {
		log.Fatal("error creating trace observer:", err)
	}

	log.Printf("trace observer count: %d", len(traceObservers))

	// Delete an existing trace observer.
	_, err = client.DeleteTraceObserver(accountID, traceObserver.ID)
	if err != nil {
		log.Fatal("error deleting trace observer:", err)
	}
}
