package accounts

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_accounts() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// List the accounts this user is authorized to view.
	params := ListAccountsParams{
		Scope: &RegionScopeTypes.GLOBAL,
	}

	accounts, err := client.ListAccounts(params)
	if err != nil {
		log.Fatal("error retrieving accounts:", err)
	}

	log.Printf("accounts count: %d", len(accounts))
}
