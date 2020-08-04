package logs

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Events is used to send custom events to NRDB.
type Logs struct {
	client http.Client
	config config.Config
	logger logging.Logger
}

// New is used to create a new Logs client instance.
func New(cfg config.Config) Logs {
	cfg.Compression = config.Compression.Gzip

	client := http.NewClient(cfg)
	client.SetAuthStrategy(&http.InsightsInsertKeyAuthorizer{})

	pkg := Logs{
		client: client,
		config: cfg,
		logger: cfg.GetLogger(),
	}

	return pkg
}
