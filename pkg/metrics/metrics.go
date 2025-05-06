package metrics

import (
	"errors"
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

// Metrics is used to send metric data to the New Relic Metrics API
type Metrics struct {
	client http.Client
	config config.Config
	logger logging.Logger
}

// New is used to create a new Metrics client instance.
func New(cfg config.Config) Metrics {
	cfg.Compression = config.Compression.Gzip
	client := http.NewClient(cfg)
	client.SetAuthStrategy(&http.LicenseKeyAuthorizer{})

	pkg := Metrics{
		client: client,
		config: cfg,
		logger: cfg.GetLogger(),
	}

	return pkg
}

// CreateMetricEntry reports a metric entry to New Relic.
// It's up to the caller to send a valid Metric API payload, no checking done here
func (m *Metrics) CreateMetricEntry(metricEntry any) error {
	if metricEntry == nil {
		return errors.New("metrics: CreateMetricEntry: metricEntry is nil, nothing to do")
	}
	_, err := m.client.Post(m.config.Region().MetricsURL(), nil, metricEntry, nil)

	// If no error is returned then the call succeeded
	if err != nil {
		m.logger.Error("metrics: Error: CreateMetricEntry: %s", err.Error())
		return err
	}

	return nil
}
