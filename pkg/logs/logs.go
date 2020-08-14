package logs

import (
	"errors"
	"fmt"
	"os"

	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Logs is used to send log data to the New Relic Log API
type Logs struct {
	client http.Client
	config config.Config
	logger logging.Logger
}

// New is used to create a new Logs client instance.
func New(cfg config.Config) Logs {
	cfg.Compression = config.Compression.Gzip

	client := http.NewClient(cfg)
	client.SetAuthStrategy(&http.LicenseKeyAuthorizer{})

	pkg := Logs{
		client: client,
		config: cfg,
		logger: cfg.GetLogger(),
	}

	return pkg
}

// CreateLogEntry reports a log entry to New Relic.
// It's up to the caller to send a valid Log API payload, no checking done here
func (l *Logs) CreateLogEntry(logEntry interface{}) error {
	if logEntry == nil {
		return errors.New("logs: CreateLogEntry: logEntry is nil, nothing to do")
	}

	rsp, err := l.client.Post(l.config.Region().LogsURL(), nil, logEntry, nil)

	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "CreateLogEntry: rsp: %+v\n", rsp)
	if !(rsp.StatusCode < 299) {
		return fmt.Errorf("failed creating log entry: %v", rsp)
	}

	return nil
}
