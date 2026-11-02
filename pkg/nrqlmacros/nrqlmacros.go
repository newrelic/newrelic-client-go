package nrqlmacros

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

// Nrqlmacros exposes CRUD for NRQL Macros via the New Relic Entity Management API.
// Macros are reusable NRQL expression fragments that can be referenced in queries
// using the $ prefix (e.g. $myMacro). They are scoped to an organization.
type Nrqlmacros struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

// New returns a Nrqlmacros client backed by the given config.
func New(config config.Config) Nrqlmacros {
	client := http.NewClient(config)
	pkg := Nrqlmacros{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
