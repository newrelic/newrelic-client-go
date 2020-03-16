package entities

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func Example_tags() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.Config{
		PersonalAPIKey: os.Getenv("NEW_RELIC_API_KEY"),
	}

	// Initialize the client.
	client := New(cfg)

	// Search the current account for entities by tag.
	searchParams := SearchEntitiesParams{
		Tags: &TagValue{
			Key:   "exampleKey",
			Value: "exampleValue",
		},
	}

	entities, err := client.SearchEntities(searchParams)
	if err != nil {
		log.Fatal("error searching entities:", err)
	}

	// List the tags associated with a given entity.  This example assumes that
	// at least one entity has been returned by the search endpoint, but in
	// practice it is possible that an empty slice is returned.
	entityGUID := entities[0].GUID
	tags, err := client.ListTags(entityGUID)
	if err != nil {
		log.Fatal("error listing tags:", err)
	}

	// Output all tags and their values.
	for _, t := range tags {
		fmt.Printf("Key: %s, Values: %v\n", t.Key, t.Values)
	}

	// Add tags to a given entity.
	addTags := []Tag{
		{
			Key: "environment",
			Values: []string{
				"production",
			},
		},
		{
			Key: "teams",
			Values: []string{
				"ops",
				"product-development",
			},
		},
	}

	err = client.AddTags(entityGUID, addTags)
	if err != nil {
		log.Fatal("error adding tags to entity:", err)
	}

	// Delete tag values from a given entity.
	// This example deletes the "ops" value from the "teams" tag.
	tagValuesToDelete := []TagValue{
		{
			Key:   "teams",
			Value: "ops",
		},
	}

	err = client.DeleteTagValues(entityGUID, tagValuesToDelete)
	if err != nil {
		log.Fatal("error deleting tag values from entity:", err)
	}

	// Delete tags from a given entity.
	// This example delete the "teams" tag and all its values from the entity.
	err = client.DeleteTags(entityGUID, []string{"teams"})
	if err != nil {
		log.Fatal("error deleting tags from entity:", err)
	}

	// Replace all existing tags for a given entity with the given set.
	datacenterTag := Tag{
		Key: "datacenter",
		Values: []string{
			"east",
		},
	}

	replaceTags := []Tag{datacenterTag}

	err = client.ReplaceTags(entityGUID, replaceTags)
	if err != nil {
		log.Fatal("error replacing tags for entity:", err)
	}
}
