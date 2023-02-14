package agentapplication

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

// AgentApplication is used to communicate with the New Relic's GraphQL API, NerdGraph.
type AgentApplication struct {
	client http.Client
	logger logging.Logger
}

// New returns a new GraphQL client for interacting with New Relic's GraphQL API, NerdGraph.
func New(config config.Config) AgentApplication {
	return AgentApplication{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}
