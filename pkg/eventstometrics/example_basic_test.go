//go:build integration
// +build integration

package eventstometrics

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

	// List the events to metrics rules for a given account.
	rules, err := client.ListRules(accountID)
	if err != nil {
		log.Fatal("error listing events to metrics rules: ", err)
	}

	// Get a specific events to metrics rule by ID.
	rule, err := client.GetRule(accountID, rules[0].ID)
	if err != nil {
		log.Fatal("error getting events to metric rule: ", err)
	}

	log.Printf("Rule name: %s", rule.Name)

	// Create a new events to metrics rule.
	createInput := []EventsToMetricsCreateRuleInput{
		{
			AccountID:   accountID,
			Name:        "Example rule",
			Description: "Example description",
			NRQL:        "SELECT uniqueCount(account_id) AS `Transaction.account_id` FROM Transaction FACET appName, name",
		},
	}

	rules, err = client.CreateRules(createInput)
	if err != nil {
		log.Fatal("error creating events to metrics rules: ", err)
	}

	// Update an existing events to metrics rule.
	updateInput := []EventsToMetricsUpdateRuleInput{
		{
			AccountID: accountID,
			RuleId:    rules[0].ID,
			Enabled:   false,
		},
	}

	rules, err = client.UpdateRules(updateInput)
	if err != nil {
		log.Fatal("error updating events to metrics rules: ", err)
	}

	// Delete an existing events to metrics rule.
	deleteInput := []EventsToMetricsDeleteRuleInput{
		{
			AccountID: accountID,
			RuleId:    rules[0].ID,
		},
	}
	_, err = client.DeleteRules(deleteInput)
	if err != nil {
		log.Fatal("error deleting event to meetrics rules: ", err)
	}
}
