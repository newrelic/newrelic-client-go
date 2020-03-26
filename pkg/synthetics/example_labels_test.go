package synthetics

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_labels() {
	// Initialize the client configuration.  An Admin API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.AdminAPIKey = os.Getenv("NEW_RELIC_ADMIN_API_KEY")

	// Initialize the client.
	client := New(cfg)

	monitorID := "fe3c7f8e-9200-4f4b-8643-13e156f22a2a"

	// Get the labels for a given Synthetics monitor.
	labels, err := client.GetMonitorLabels(monitorID)
	if err != nil {
		log.Fatal("error updating Synthetics monitor: ", err)
	}

	// Output the label keys and values.
	for _, l := range labels {
		fmt.Printf("monitor ID: %s, key: %s, value: %s", monitorID, l.Type, l.Value)
	}

	// Add a label to a given Synthetics monitor.
	err = client.AddMonitorLabel(monitorID, "key", "value")
	if err != nil {
		log.Fatal("error adding Synthetics label: ", err)
	}

	// Delete a label from a Synthetics monitor.
	err = client.DeleteMonitorLabel(monitorID, "key", "value")
	if err != nil {
		log.Fatal("error deleting Synthetics label: ", err)
	}
}
