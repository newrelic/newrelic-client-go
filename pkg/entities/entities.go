// Package entities provides a programmatic API for interacting with New Relic One entities.
package entities

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

// Entities is used to communicate with the New Relic Entities product.
type Entities struct {
	client http.Client
	logger logging.Logger
}

// New returns a new client for interacting with New Relic One entities.
func New(config config.Config) Entities {
	return Entities{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}

func BuildEntitySearchQuery(name string, domain string, entityType string, tags []map[string]string) string {
	params := map[string]string{
		"name":   name,
		"domain": domain,
		"type":   entityType,
	}

	count := 0
	query := ""
	for k, v := range params {
		if v == "" {
			continue
		}

		if count == 0 {
			query = fmt.Sprintf("%s = '%s'", k, v)
		} else {
			query = fmt.Sprintf("%s AND %s = '%s'", query, k, v)
		}
		count++
	}

	if len(tags) > 0 {
		query = fmt.Sprintf("%s AND %s", query, BuildTagsQueryFragment(tags))
	}

	return query
}

func BuildTagsQueryFragment(tags []map[string]string) string {
	var query string

	for i, t := range tags {
		var q string
		if i > 0 {
			q = fmt.Sprintf(" AND tags.`%s` = '%s'", t["key"], t["value"])
		} else {
			q = fmt.Sprintf("tags.`%s` = '%s'", t["key"], t["value"])
		}

		query = fmt.Sprintf("%s%s", query, q)
	}

	return query
}
