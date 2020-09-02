package apm

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_application_instances() {
	// Initialize the client configuration.  A Personal API key
	// is required to communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	applicationID := 12345
	instanceID := 12345678

	params := ListApplicationInstancesParams{
		Hostname: "hostname",
		IDs:      []int{instanceID},
	}

	// List an application's instances.
	instances, err := client.ListApplicationInstances(applicationID, &params)
	if err != nil {
		log.Fatal("error listing application instances:", err)
	}

	// Output the application instance count.
	fmt.Printf("Instance count: %d", len(instances))

	// Get a specific application instance.
	instance, err := client.GetApplicationInstance(applicationID, instanceID)
	if err != nil {
		log.Fatal("error deleting application:", err)
	}

	// Output the application instance's host and port.
	fmt.Printf("Host: %s, Port: %d\n", instance.Host, instance.Port)
}
