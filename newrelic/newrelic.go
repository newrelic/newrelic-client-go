package newrelic

import (
	"errors"
	"net/http"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/plugins"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

// NewRelic is a collection of New Relic APIs.
type NewRelic struct {
	Alerts     alerts.Alerts
	APM        apm.APM
	Dashboards dashboards.Dashboards
	Entities   entities.Entities
	Plugins    plugins.Plugins
	Synthetics synthetics.Synthetics
}

// New returns a collection of New Relic APIs.
func New(opts ...ConfigOption) (*NewRelic, error) {
	config := config.Config{}

	// Loop through config options
	for _, fn := range opts {
		if nil != fn {
			if err := fn(&config); err != nil {
				return nil, err
			}
		}
	}

	if config.APIKey == "" && config.PersonalAPIKey == "" {
		return nil, errors.New("use of ConfigAPIKey and/or ConfigPersonalAPIKey is required")
	}

	nr := &NewRelic{
		Alerts:     alerts.New(config),
		APM:        apm.New(config),
		Dashboards: dashboards.New(config),
		Entities:   entities.New(config),
		Plugins:    plugins.New(config),
		Synthetics: synthetics.New(config),
	}

	return nr, nil
}

// ConfigOption configures the Config when provided to NewApplication.
// https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys
type ConfigOption func(*config.Config) error

// ConfigAPIKey sets the New Relic Admin API key this client will use.
// One of ConfigAPIKey or ConfigPersonalAPIKey must be used to create a client.
// https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys
func ConfigAPIKey(apiKey string) ConfigOption {
	return func(cfg *config.Config) error {
		cfg.APIKey = apiKey
		return nil
	}
}

// ConfigPersonalAPIKey sets the New Relic Personal API key this client will use.
// One of ConfigAPIKey or ConfigPersonalAPIKey must be used to create a client.
// https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys
func ConfigPersonalAPIKey(personalAPIKey string) ConfigOption {
	return func(cfg *config.Config) error {
		cfg.PersonalAPIKey = personalAPIKey
		return nil
	}
}

// ConfigRegion sets the New Relic Region this client will use.
func ConfigRegion(region string) ConfigOption {
	return func(cfg *config.Config) error {
		cfg.Region = region
		return nil
	}
}

// ConfigHTTPTimeout sets the timeout for HTTP requests.
func ConfigHTTPTimeout(t time.Duration) ConfigOption {
	return func(cfg *config.Config) error {
		var timeout = &t
		cfg.Timeout = timeout
		return nil
	}
}

// ConfigHTTPTransport sets the HTTP Transporter.
func ConfigHTTPTransport(transport *http.RoundTripper) ConfigOption {
	return func(cfg *config.Config) error {
		if transport != nil {
			cfg.HTTPTransport = transport
			return nil
		}

		return errors.New("HTTP Transport can not be nil")
	}
}

// ConfigUserAgent sets the HTTP UserAgent for API requests.
func ConfigUserAgent(ua string) ConfigOption {
	return func(cfg *config.Config) error {
		if ua != "" {
			cfg.UserAgent = ua
			return nil
		}

		return errors.New("user-agent can not be empty")
	}
}

// ConfigBaseURL sets the base URL used to make requests to the REST API V2.
func ConfigBaseURL(url string) ConfigOption {
	return func(cfg *config.Config) error {
		if url != "" {
			cfg.BaseURL = url
			return nil
		}

		return errors.New("base URL can not be empty")
	}
}

// ConfigInfrastructureBaseURL sets the base URL used to make requests to the Infrastructure API.
func ConfigInfrastructureBaseURL(url string) ConfigOption {
	return func(cfg *config.Config) error {
		if url != "" {
			cfg.InfrastructureBaseURL = url
			return nil
		}

		return errors.New("Infrastructure base URL can not be empty")
	}
}

// ConfigSyntheticsBaseURL sets the base URL used to make requests to the Synthetics API.
func ConfigSyntheticsBaseURL(url string) ConfigOption {
	return func(cfg *config.Config) error {
		if url != "" {
			cfg.SyntheticsBaseURL = url
			return nil
		}

		return errors.New("Synthetics base URL can not be empty")
	}
}

// ConfigNerdGraphBaseURL sets the base URL used to make requests to the NerdGraph API.
func ConfigNerdGraphBaseURL(url string) ConfigOption {
	return func(cfg *config.Config) error {
		if url != "" {
			cfg.NerdGraphBaseURL = url
			return nil
		}

		return errors.New("NerdGraph base URL can not be empty")
	}
}
