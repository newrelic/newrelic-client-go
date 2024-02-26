//go:build unit || integration
// +build unit integration

package logconfigurations

import (
	"log"
	"os"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v3/pkg/config"
)

func LogConfigurationsObfuscationExpression_Basic() {
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

	// Create a new expression.
	createInput := LogConfigurationsCreateObfuscationExpressionInput{
		Description: "Brief desc",
		Name:        "some thing",
		Regex:       "(^http.*)",
	}

	created, err := client.LogConfigurationsCreateObfuscationExpression(accountID, createInput)
	if err != nil {
		log.Fatal("error creating Obfuscation expression: ", err)
	}

	// Update an existing obfuscation expression.
	expressionName := "Example updated workload"
	updateInput := LogConfigurationsUpdateObfuscationExpressionInput{
		Name: expressionName,
	}

	updated, err := client.LogConfigurationsUpdateObfuscationExpression(created.CreatedBy.ID, updateInput)
	if err != nil {
		log.Fatal("error updating obfuscation expression: ", err)
	}

	// Delete an existing events to metrics rule.
	deleteInput := updated.ID
	_, err = client.LogConfigurationsDeleteObfuscationExpression(accountID, deleteInput)
	if err != nil {
		log.Fatal("error deleting obfuscation expression: ", err)
	}
}
