//go:build integration
// +build integration

package nrqldroprules

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_basic() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	accountID := 12345678

	// List the NRQL drop rules for a given account.
	_, err := client.GetList(accountID)
	if err != nil {
		log.Fatal("error listing NRQL drop rules: ", err)
	}

	// Create a new events to metrics rule.
	createInput := []NRQLDropRulesCreateDropRuleInput{
		{
			Action:      NRQLDropRulesActionTypes.DROP_DATA,
			Description: "NRQL drop rule description",
			NRQL:        "SELECT * FROM Log WHERE container_name = 'noise'",
			Source:      "Logging",
		},
	}

	created, err := client.NRQLDropRulesCreate(accountID, createInput)
	if err != nil {
		log.Fatal("error creating NRQL drop rules: ", err)
	}

	// Delete an existing events to metrics rule.
	deleteInput := []string{created.Successes[0].ID}
	_, err = client.NRQLDropRulesDelete(accountID, deleteInput)
	if err != nil {
		log.Fatal("error deleting NRQL drop rules: ", err)
	}
}
