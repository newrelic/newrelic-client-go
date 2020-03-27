package nerdgraph

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_query() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Execute a NRQL query to retrieve the average duration of transactions for
	// the "Example application" app.
	query := `
	query($accountId: Int!, $nrqlQuery: Nrql!) {
		actor {
			account(id: $accountId) {
				nrql(query: $nrqlQuery, timeout: 5) {
					results
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"accountId": 12345678,
		"nrqlQuery": "SELECT average(duration) FROM Transaction TIMESERIES WHERE appName = 'Example application'",
	}

	resp, err := client.Query(query, variables)
	if err != nil {
		log.Fatal("error running NerdGraph query: ", err)
	}

	queryResp := resp.(QueryResponse)
	actor := queryResp.Actor.(map[string]interface{})
	account := actor["account"].(map[string]interface{})
	nrql := account["nrql"].(map[string]interface{})
	results := nrql["results"].([]interface{})

	var durations []float64
	for _, r := range results {
		data := r.(map[string]interface{})
		durations = append(durations, data["average.duration"].(float64))
	}

	// Output the raw time series values for transaction duration.
	fmt.Printf("durations: %v\n", durations)
}
