package apm

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func Example_keyTransaction() {
	// Initialize the client configuration. A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Search the key transactions for the current account by name.
	listParams := &ListKeyTransactionsParams{
		Name: "Example key transaction",
	}

	transactions, err := client.ListKeyTransactions(listParams)
	if err != nil {
		log.Fatal("error listing key transactions:", err)
	}

	// Get a key transaction by ID.  This example assumes that at least one key
	// transaction has been returned by the list endpoint, but in practice it is
	// possible that an empty slice is returned.
	transaction, err := client.GetKeyTransaction(transactions[0].ID)
	if err != nil {
		log.Fatal("error getting key transaction:", err)
	}

	// Output the key transaction's health status.
	fmt.Printf("Key transaction status: %s\n", transaction.HealthStatus)
}
