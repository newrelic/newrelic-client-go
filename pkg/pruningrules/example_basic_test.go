//go:build integration
// +build integration

package pruningrules

import (
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func Example_basic() {
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	client := New(cfg)

	accountID := 12345678

	// List all pruning rules for the account.
	listResult, err := client.GetList(accountID)
	if err != nil {
		log.Fatal("error listing pruning rules: ", err)
	}

	if len(listResult.Rules) > 0 {
		// Fetch a specific rule by ID.
		_, err = client.GetPruningRuleByID(accountID, listResult.Rules[0].ID)
		if err != nil {
			log.Fatal("error getting pruning rule by ID: ", err)
		}
	}

	// Create a pruning rule that drops the collector.name attribute from metric aggregates.
	createInput := []NRQLDropRulesCreateDropRuleInput{
		{
			Action:      NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES,
			Description: "Drop collector.name from my.service.latency metric aggregates",
			NRQL:        "SELECT collector.name FROM Metric WHERE metricName = 'my.service.latency'",
		},
	}

	created, err := client.NRQLDropRulesCreate(accountID, createInput)
	if err != nil {
		log.Fatal("error creating pruning rule: ", err)
	}

	// Clean up.
	_, err = client.NRQLDropRulesDelete(accountID, []string{created.Successes[0].ID})
	if err != nil {
		log.Fatal("error deleting pruning rule: ", err)
	}
}
