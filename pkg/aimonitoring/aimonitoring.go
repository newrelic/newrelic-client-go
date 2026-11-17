package aimonitoring

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

// Aimonitoring exposes AI Monitoring datasets and dataset versions via the
// New Relic Entity Management API.
type Aimonitoring struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

// New returns an Aimonitoring client backed by the given config.
func New(config config.Config) Aimonitoring {
	client := http.NewClient(config)
	pkg := Aimonitoring{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
