package entities

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_entity() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.Config{
		PersonalAPIKey: os.Getenv("NEW_RELIC_API_KEY"),
	}

	// Initialize the client.
	client := New(cfg)

	// Search the current account for entities by name and type.
	searchParams := SearchEntitiesParams{
		Name: "Example entity",
		Type: EntityTypes.Application,
	}

	entities, err := client.SearchEntities(searchParams)
	if err != nil {
		log.Fatal("error searching entities:", err)
	}

	// Get several entities by GUID.
	var entityGuids []string
	for _, e := range entities {
		entityGuids = append(entityGuids, e.GUID)
	}

	entities, err = client.GetEntities(entityGuids)
	if err != nil {
		log.Fatal("error getting entities:", err)
	}

	// Get an entity by GUID.
	entity, err := client.GetEntity(entities[0].GUID)
	if err != nil {
		log.Fatal("error getting entity:", err)
	}

	// Output the entity's URL.
	fmt.Printf("Entity name: %s, URL: %s\n", entity.Name, entity.Permalink)
}
