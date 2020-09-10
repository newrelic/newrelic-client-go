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
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Search the current account for entities by name and type.
	searchParams := SearchEntitiesParams{
		Name: "Example entity",
		Type: EntityTypeTypes.APM_APPLICATION_ENTITY,
	}

	entities, err := client.SearchEntities(searchParams)
	if err != nil {
		log.Fatal("error searching entities:", err)
	}

	// Get several entities by GUID.
	var entityGuids []string
	for _, e := range entities {
		entityGuids = append(entityGuids, (*e).(*Entity).GUID)
	}

	entities, err = client.GetEntities(entityGuids)
	if err != nil {
		log.Fatal("error getting entities:", err)
	}

	// Get an entity by GUID.
	entity, err := client.GetEntity((*entities[0]).(*Entity).GUID)
	if err != nil {
		log.Fatal("error getting entity:", err)
	}

	// Output the entity's URL.
	fmt.Printf("Entity name: %s, URL: %s\n", (*entity).(*Entity).Name, (*entity).(*Entity).Permalink)
}
