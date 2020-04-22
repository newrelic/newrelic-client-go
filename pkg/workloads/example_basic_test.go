// +build integration

package workloads

import (
	"fmt"
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
	entityGUID := "MjUwODy1OXxOUjF8V097S0xPQ3R8ODcz"

	// List the workloads for a given account.
	workloads, err := client.ListWorkloads(accountID)
	if err != nil {
		log.Fatal("error listing workloads: ", err)
	}

	// Get a specific workload by ID.  This example assumes that at least one
	// workload has been returned by the list endpoint, but in practice it is
	// possible that an empty slice is returned.
	workload, err := client.GetWorkload(accountID, workloads[0].GUID)
	if err != nil {
		log.Fatal("error getting workload: ", err)
	}

	// Create a new workload.
	createInput := CreateInput{
		Name:        "Example workload",
		EntityGUIDs: []string{entityGUID},
		EntitySearchQueries: []EntitySearchQueryInput{
			{
				Query: fmt.Sprintf("(accountId IN ('%d')) AND (((name like 'Example application' or id = 'Example application' or domainId = 'Example application')))", accountID),
			},
		},
	}

	workload, err = client.CreateWorkload(accountID, createInput)
	if err != nil {
		log.Fatal("error creating workload: ", err)
	}

	// Duplicate an existing workload.
	duplicate, err := client.DuplicateWorkload(accountID, workload.GUID, nil)
	if err != nil {
		log.Fatal("error duplicating workload: ", err)
	}

	// Update an existing workload.
	workloadName := "Example updated workload"
	updateInput := UpdateInput{
		Name:                workloadName,
		EntityGUIDs:         createInput.EntityGUIDs,
		EntitySearchQueries: createInput.EntitySearchQueries,
	}

	updated, err := client.UpdateWorkload(duplicate.GUID, updateInput)
	if err != nil {
		log.Fatal("error updating workload: ", err)
	}

	// Delete an existing workload.
	_, err = client.DeleteWorkload(updated.GUID)
	if err != nil {
		log.Fatal("error deleting workload: ", err)
	}
}
