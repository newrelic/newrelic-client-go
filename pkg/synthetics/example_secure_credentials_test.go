package synthetics

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_secureCredentials() {
	// Initialize the client configuration.  An Admin API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.AdminAPIKey = os.Getenv("NEW_RELIC_ADMIN_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Get the Synthetics secure credentials for this account.
	credentials, err := client.GetSecureCredentials()
	if err != nil {
		log.Fatal("error getting Synthetics secure credentials: ", err)
	}

	// Get a single Synthetics secure credential belonging to this account.
	credential, err := client.GetSecureCredential(credentials[0].Key)
	if err != nil {
		log.Fatal("error getting Synthetics secure credential: ", err)
	}

	// Add a secure credential for use with Synthetics.
	credential, err = client.AddSecureCredential("key", "value", "description")
	if err != nil {
		log.Fatal("error adding Synthetics secure credential: ", err)
	}

	// Delete a Synthetics secure credential.
	err = client.DeleteSecureCredential(credential.Key)
	if err != nil {
		log.Fatal("error deleting Synthetics secure credential: ", err)
	}
}
