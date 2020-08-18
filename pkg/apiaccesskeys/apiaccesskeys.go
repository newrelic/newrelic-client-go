// Package apikeys provides a programmatic API for interacting with New Relic API keys
package apiaccesskeys

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// APIKeys is used to communicate with the New Relic APIKeys product.
type APIAccessKeys struct {
	client http.Client
	logger logging.Logger
}

// New returns a new client for interacting with New Relic One entities.
func New(config config.Config) APIAccessKeys {
	return APIAccessKeys{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}
