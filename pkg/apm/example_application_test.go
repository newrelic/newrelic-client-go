package apm

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v3/pkg/config"
)

func Example_application() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Search the applications for the current account by name.
	listParams := &ListApplicationsParams{
		Name: "Example application",
	}

	apps, err := client.ListApplications(listParams)
	if err != nil {
		log.Fatal("error listing applications:", err)
	}

	// Get an application by ID.  This example assumes that at least one application
	// has been returned by the list endpoint, but in practice it is possible
	// that an empty slice is returned.
	app, err := client.GetApplication(apps[0].ID)
	if err != nil {
		log.Fatal("error getting application:", err)
	}

	// Update an application's settings.  The following example updates the
	// application's Apdex threshold.
	updateParams := UpdateApplicationParams{
		Name: app.Name,
		Settings: ApplicationSettings{
			AppApdexThreshold: 0.6,
		},
	}

	app, err = client.UpdateApplication(app.ID, updateParams)
	if err != nil {
		log.Fatal("error updating application settings:", err)
	}

	// Delete an application that is no longer reporting data.
	if !app.Reporting {
		_, err = client.DeleteApplication(app.ID)
		if err != nil {
			log.Fatal("error deleting application:", err)
		}
	}
}
