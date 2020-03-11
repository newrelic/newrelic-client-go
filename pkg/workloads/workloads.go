package workloads

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Workloads is used to communicate with the New Relic Workloads product.
type Workloads struct {
	client *http.GraphQLClient
	logger logging.Logger
}

// New returns a new client for interacting with New Relic One workloads.
func New(config config.Config) Workloads {
	return Workloads{
		client: http.NewGraphQLClient(config),
		logger: config.GetLogger(),
	}
}

// BaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var BaseURLs = region.NerdGraphBaseURLs
