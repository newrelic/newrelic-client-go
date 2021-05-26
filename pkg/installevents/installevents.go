package installevents

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// NerdStorage is used to communicate with the New Relic Workloads product.
type InstallEvents struct {
	client http.Client
	logger logging.Logger
}

// New returns a new client for interacting with the New Relic One NerdStorage
// document store.
func New(config config.Config) InstallEvents {
	return InstallEvents{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}
