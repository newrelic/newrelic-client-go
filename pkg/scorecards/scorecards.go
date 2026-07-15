// Package scorecards manages New Relic entity-management resources that share
// the NGEP (NerdGraph Entity Platform) interface — Teams, Scorecards,
// ScorecardRules, and Collections. See NGEP_ANALYSIS.md for the relationship
// model and API semantics that shaped this package.
package scorecards

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

// Scorecards exposes CRUD for Teams, Scorecards, ScorecardRules, and
// Collections against the New Relic Entity Management API.
type Scorecards struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

// New returns a Scorecards client backed by the given config.
func New(config config.Config) Scorecards {
	client := http.NewClient(config)
	pkg := Scorecards{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
