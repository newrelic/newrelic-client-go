package entities

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

func Example_entity() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Search the current account for entities by name and type.
	queryBuilder := EntitySearchQueryBuilder{
		Name: "Example entity",
		Type: EntitySearchQueryBuilderTypeTypes.APPLICATION,
	}

	entitySearch, err := client.GetEntitySearch(
		EntitySearchOptions{},
		"",
		queryBuilder,
		[]EntitySearchSortCriteria{},
	)
	if err != nil {
		log.Fatal("error searching entities:", err)
	}

	// Get several entities by GUID.
	var entityGuids []common.EntityGUID
	for _, x := range entitySearch.Results.Entities {
		e := x.(*GenericEntityOutline)
		entityGuids = append(entityGuids, e.GUID)
	}

	entities, err := client.GetEntities(entityGuids)
	if err != nil {
		log.Fatal("error getting entities:", err)
	}
	fmt.Printf("GetEntities returned %d entities", len((*entities)))

	// Get an entity by GUID.
	entity, err := client.GetEntity(entityGuids[0])
	if err != nil {
		log.Fatal("error getting entity:", err)
	}

	// Output the entity's URL.
	fmt.Printf("Entity name: %s, URL: %s\n", (*entity).(*GenericEntity).Name, (*entity).(*GenericEntity).Permalink)
}
