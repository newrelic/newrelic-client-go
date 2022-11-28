//go:build unit || integration
// +build unit integration

package logconfigurations

import (
	"log"
	"os"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func LogConfigurationsParsingRule_Basic() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	accountID, enverr := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if enverr != nil {
		log.Fatal("environment variable NEW_RELIC_ACCOUNT_ID required")
	}

	// Create a new Parsing Rule.
	createInput := LogConfigurationsParsingRuleConfiguration{
		Attribute:   "attribute",
		Description: "Brief desc",
		Enabled:     true,
		Grok:        "sampleattribute=%{NUMBER:test:int}",
		Lucene:      "logtype:linux_messages",
		NRQL:        "SELECT * FROM Log WHERE logtype = 'linux_messages'",
	}

	created, err := client.LogConfigurationsCreateParsingRule(accountID, createInput)
	if err != nil {
		log.Fatal("error creating Parsing Rule: ", err)
	}

	//Update an existing Parsing Rule.
	updated, err := client.LogConfigurationsUpdateParsingRule(accountID, created.Rule.ID, LogConfigurationsParsingRuleConfiguration{
		Attribute:   "attribute",
		Description: "update",
		Enabled:     true,
		Grok:        "sampleattribute=%{NUMBER:test:int}",
		Lucene:      "logtype:linux_messages",
		NRQL:        "SELECT * FROM Log WHERE logtype = 'linux_messages'",
	})

	if err != nil {
		log.Fatal("error updating Parsing Rule: ", err)
	}

	// Delete an Parsing Rule.
	deleteInput := updated.Rule.ID
	_, err = client.LogConfigurationsDeleteParsingRule(accountID, deleteInput)
	if err != nil {
		log.Fatal("error deleting Parsing Rule: ", err)
	}
}
