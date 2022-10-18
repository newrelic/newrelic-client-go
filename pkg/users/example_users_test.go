package users

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func Example_accounts() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	user, err := client.GetUser()
	if err != nil {
		log.Fatal("error retrieving user:", err)
	}

	log.Printf("User name:  %s", user.Name)
	log.Printf("User email: %s", user.Email)
	log.Printf("User ID:    %d", user.ID)
}
