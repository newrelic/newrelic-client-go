package apm

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_labels() {
	// Initialize the client configuration. A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Get an APM application by ID.
	app, err := client.GetApplication(12345678)
	if err != nil {
		log.Fatal("error getting application:", err)
	}

	// List the existing labels for this account.
	labels, err := client.ListLabels()
	if err != nil {
		log.Fatal("error listing labels:", err)
	}

	// Output the concatenated label key and associated application IDs for each.
	for _, l := range labels {
		fmt.Printf("Label key: %s, Application IDs: %v\n", l.Key, l.Links.Applications)
	}

	// Add a label to the application that describes its data center's location.
	label := Label{
		Category: "Datacenter",
		Name:     "East",
		Links: LabelLinks{
			Applications: []int{app.ID},
		},
	}

	l, err := client.CreateLabel(label)
	if err != nil {
		log.Fatal("error creating label:", err)
	}

	// Delete a label from all linked applications and servers.
	_, err = client.DeleteLabel(l.Key)
	if err != nil {
		log.Fatal("error deleting label:", err)
	}
}
