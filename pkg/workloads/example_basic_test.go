//go:build integration
// +build integration

package workloads

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func Example_basic() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	accountID := 12345678
	entityGUID := common.EntityGUID("MjUwODy1OXxOUjF8V097S0xPQ3R8ODcz")

	// Create a new workload.
	createInput := WorkloadCreateInput{
		Name:        "Example workload",
		EntityGUIDs: []common.EntityGUID{entityGUID},
		EntitySearchQueries: []WorkloadEntitySearchQueryInput{
			{
				Query: fmt.Sprintf("(accountId IN ('%d')) AND (((name like 'Example application' or id = 'Example application' or domainId = 'Example application')))", accountID),
			},
		},
	}

	workload, err := client.WorkloadCreate(accountID, createInput)
	if err != nil {
		log.Fatal("error creating workload: ", err)
	}

	// Duplicate an existing workload.
	duplicate, err := client.WorkloadDuplicate(accountID, workload.GUID, WorkloadDuplicateInput{
		Name: workload.Name + "-duplicate",
	})
	if err != nil {
		log.Fatal("error duplicating workload: ", err)
	}

	// Update an existing workload.
	workloadName := "Example updated workload"
	updateInput := WorkloadUpdateInput{
		Name: workloadName,
	}

	updated, err := client.WorkloadUpdate(duplicate.GUID, updateInput)
	if err != nil {
		log.Fatal("error updating workload: ", err)
	}

	// Delete an existing workload.
	_, err = client.WorkloadDelete(updated.GUID)
	if err != nil {
		log.Fatal("error deleting workload: ", err)
	}
}

func Example_intelligentWorkload() {
	// Initialize the client configuration. A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	accountID := 12345678
	entityGUID := common.EntityGUID("MjUwODy1OXxOUjF8V097S0xPQ3R8ODcz")

	// Create a new intelligent workload using dynamic flows.
	// New Relic auto-discovers related entities via Transaction 360 distributed tracing.
	// If set alongside entityGuids or entitySearchQueries, dynamicFlows takes precedence
	// and others will be ignored.
	createInput := WorkloadCreateInput{
		Name: "Example intelligent workload",
		ScopeAccounts: &WorkloadScopeAccountsInput{
			AccountIDs: []int{accountID},
		},
		DynamicFlows: []WorkloadDynamicFlowInput{
			{
				EntityGUID:      entityGUID,
				TransactionName: "WebTransaction/Action/index",
			},
		},
		StatusConfig: &WorkloadStatusConfigInput{
			AlertPolicy: &WorkloadAlertPolicyInput{
				Enabled: true,
			},
		},
	}

	workload, err := client.WorkloadCreate(accountID, createInput)
	if err != nil {
		log.Fatal("error creating intelligent workload: ", err)
	}

	// Update an existing intelligent workload.
	updated, err := client.WorkloadUpdate(workload.GUID, WorkloadUpdateInput{
		Name: "Example intelligent workload - updated",
		DynamicFlows: []WorkloadUpdateDynamicFlowInput{
			{
				EntityGUID:      entityGUID,
				TransactionName: "WebTransaction/Action/index",
			},
		},
		StatusConfig: &WorkloadUpdateStatusConfigInput{
			AlertPolicy: &WorkloadUpdateAlertPolicyInput{
				Enabled: false,
			},
		},
	})
	if err != nil {
		log.Fatal("error updating intelligent workload: ", err)
	}

	// Delete an existing intelligent workload.
	_, err = client.WorkloadDelete(updated.GUID)
	if err != nil {
		log.Fatal("error deleting intelligent workload: ", err)
	}
}
