package nrdb

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_query() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		log.Fatal("example requires NEW_RELIC_ACCONT_ID to be set, got error: ", err)
	}

	// Initialize the client.
	client := New(cfg)

	// Execute a NRQL query to retrieve the average duration of transactions for
	// the "Example application" app.
	query := NRQL("SELECT average(duration) FROM Transaction TIMESERIES WHERE appName = 'Example application'")

	resp, err := client.Query(accountID, query)
	if err != nil {
		log.Fatal("error running NerdGraph query: ", err)
	}

	var durations []float64
	for _, r := range resp.Results {
		durations = append(durations, r["average.duration"].(float64))
	}

	// Output the raw time series values for transaction duration.
	fmt.Printf("durations: %v\n", durations)
}
