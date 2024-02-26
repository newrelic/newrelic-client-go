//go:build unit || integration
// +build unit integration

package logconfigurations

import (
	"log"
	"os"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v3/pkg/config"
)

func LogConfigurationsObfuscationRule_Basic() {
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

	// Create a new rule.
	createInput := LogConfigurationsCreateObfuscationRuleInput{
		Description: "Some Thing",
		Name:        "ruleName",
		Filter:      "SELECT * FROM Log WHERE container_name = 'noise'",
		Enabled:     true,
		Actions: []LogConfigurationsCreateObfuscationActionInput{
			{
				Attributes:   []string{"awsAccountId"},
				ExpressionId: "1376",
				Method:       "MASK",
			},
		},
	}

	created, err := client.LogConfigurationsCreateObfuscationRule(accountID, createInput)
	if err != nil {
		log.Fatal("error creating obfuscation rule: ", err)
	}

	// Update an existing obfuscation expression.
	ruleName := "Example updated rule"
	updateInput := LogConfigurationsUpdateObfuscationRuleInput{
		Name:    ruleName,
		Actions: []LogConfigurationsUpdateObfuscationActionInput{
			//{
			//	Attributes:   []string{"awsAccountId"},
			//	ExpressionId: "1376",
			//	Method:       "MASK",
			//},
		},
	}

	updated, err := client.LogConfigurationsUpdateObfuscationRule(created.CreatedBy.ID, updateInput)
	if err != nil {
		log.Fatal("error updating obfuscation rule: ", err)
	}

	// Delete an existing events to metrics rule.
	deleteInput := updated.ID
	_, err = client.LogConfigurationsDeleteObfuscationRule(accountID, deleteInput)
	if err != nil {
		log.Fatal("error deleting obfuscation rule: ", err)
	}
}
