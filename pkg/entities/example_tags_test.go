package entities

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/v3/pkg/config"
)

func Example_tags() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	// Search the current account for entities by tag.
	queryBuilder := EntitySearchQueryBuilder{
		Tags: []EntitySearchQueryBuilderTag{
			{
				Key:   "exampleKey",
				Value: "exampleValue",
			},
		},
	}

	entities, err := client.GetEntitySearch(
		EntitySearchOptions{},
		"",
		queryBuilder,
		[]EntitySearchSortCriteria{},
	)
	if err != nil {
		log.Fatal("error searching entities:", err)
	}

	// List the tags associated with a given entity.  This example assumes that
	// at least one entity has been returned by the search endpoint, but in
	// practice it is possible that an empty slice is returned.
	entityGUID := entities.Results.Entities[0].(*GenericEntityOutline).GUID
	tags, err := client.ListTags(entityGUID)
	if err != nil {
		log.Fatal("error listing tags:", err)
	}

	// Output all tags and their values.
	for _, t := range tags {
		fmt.Printf("Key: %s, Values: %v\n", t.Key, t.Values)
	}

	// Add tags to a given entity.
	addTags := []TaggingTagInput{
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

	res, err := client.TaggingAddTagsToEntity(entityGUID, addTags)
	if err != nil || len(res.Errors) > 0 {
		log.Fatal("error adding tags to entity:", err)
	}

	// Delete tag values from a given entity.
	// This example deletes the "ops" value from the "teams" tag.
	tagValuesToDelete := []TaggingTagValueInput{
		{
			Key:   "teams",
			Value: "ops",
		},
	}

	res, err = client.TaggingDeleteTagValuesFromEntity(entityGUID, tagValuesToDelete)
	if err != nil {
		log.Fatal("error deleting tag values from entity:", err)
	}
	if res != nil {
		for _, v := range res.Errors {
			log.Print("error deleting tags from entity: ", v)
		}
	}

	// Delete tags from a given entity.
	// This example delete the "teams" tag and all its values from the entity.
	res, err = client.TaggingDeleteTagFromEntity(entityGUID, []string{"teams"})
	if err != nil {
		log.Fatal("error deleting tags from entity:", err)
	}
	if res != nil {
		for _, v := range res.Errors {
			log.Print("error deleting tags from entity: ", v)
		}
	}

	// Replace all existing tags for a given entity with the given set.
	datacenterTag := []TaggingTagInput{
		{
			Key: "datacenter",
			Values: []string{
				"east",
			},
		},
	}

	res, err = client.TaggingReplaceTagsOnEntity(entityGUID, datacenterTag)
	if err != nil {
		log.Fatal("error replacing tags for entity:", err)
	}
	if res != nil {
		for _, v := range res.Errors {
			log.Print("error replacing tags for entity: ", v)
		}
	}
}
