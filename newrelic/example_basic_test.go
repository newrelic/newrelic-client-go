package newrelic

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apm"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/plugins"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
)

func Example_basic() {
	// Initialize the client.
	client, err := New(
		ConfigPersonalAPIKey(os.Getenv("NEW_RELIC_API_KEY")),
		ConfigRegion("US"),
		ConfigLogLevel("DEBUG"),
	)
	if err != nil {
		log.Fatal("error initializing client:", err)
	}

	// Interact with the New Relic Alerts product.
	policies, err := client.Alerts.ListPolicies(&alerts.ListPoliciesParams{
		Name: "Example policy",
	})
	if err != nil {
		log.Fatal("error listing alert policies:", err)
	}

	fmt.Printf("Policies: %v+\n", policies)

	// Interact with the New Relic APM product.
	apps, err := client.APM.ListApplications(&apm.ListApplicationsParams{
		Name: "Example application",
	})
	if err != nil {
		log.Fatal("error listing APM applications:", err)
	}

	fmt.Printf("Applications: %v+\n", apps)

	// Interact with New Relic One entities.
	entities, err := client.Entities.GetEntitySearch(
		entities.EntitySearchOptions{},
		"",
		entities.EntitySearchQueryBuilder{
			Name: "Example entity",
		},
		[]entities.EntitySearchSortCriteria{},
	)
	if err != nil {
		log.Fatal("error listing entities:", err)
	}

	fmt.Printf("Entities: %v+\n", entities)

	// Interact with the New Relic One NerdGraph API.
	query := "{ actor { user { email } } }"

	resp, err := client.NerdGraph.Query(query, nil)
	if err != nil {
		log.Fatal("error executing query:", err)
	}

	fmt.Printf("Query response: %v+\n", resp)

	// Interact with the New Relic Plugins product.
	plugins, err := client.Plugins.ListPlugins(&plugins.ListPluginsParams{
		IDs: []int{1234, 5678},
	})
	if err != nil {
		log.Fatal("error listing plugins:", err)
	}

	fmt.Printf("Plugins: %v+\n", plugins)

	// Interact with the New Relic Synthetics product.
	monitors, err := client.Synthetics.ListMonitors()
	if err != nil {
		log.Fatal("error listing monitors:", err)
	}

	fmt.Printf("Synthetics monitors: %v+\n", monitors)

	// Interact with the New Relic One Workloads product.
	workload, err := client.Workloads.WorkloadCreate(12345678, workloads.WorkloadCreateInput{Name: "test-workload"})
	if err != nil {
		log.Fatal("error creating workload:", err)
	}

	fmt.Printf("Workload created: %v+\n", workload)
}
