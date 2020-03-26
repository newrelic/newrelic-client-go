// Package config provides cross-cutting configuration support for the newrelic-client-go project.
package config

import (
	"net/http"
	"time"

	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/version"
	"github.com/newrelic/newrelic-client-go/pkg/region"
)

// Config contains all the configuration data for the API Client.
type Config struct {
	// PersonalAPIKey to authenticate API requests
	// see: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#personal-api-key
	PersonalAPIKey string

	// AdminAPIKey to authenticate API requests
	// Note this will be deprecated in the future!
	// see: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#admin
	AdminAPIKey string

	// Region of the New Relic platform to use
	Region *region.Region

	// Timeout is the client timeout for HTTP requests.
	Timeout *time.Duration

	// HTTPTransport allows customization of the client's underlying transport.
	HTTPTransport http.RoundTripper

	// UserAgent updates the default user agent string used by the client.
	UserAgent string

	// BaseURL updates the default base URL used by the client during requests to
	// the V2 REST API.
	BaseURL string

	// SyntheticsBaseURL updates the default base URL used by the client during
	// requests to the Synthetics API.
	SyntheticsBaseURL string

	// InfrastructureBaseURL updates the default base URL used by the client during
	// requests to the Infrastructure API.
	InfrastructureBaseURL string

	// NerdGraph updates the default base URL used by the client during requests
	// to the NerdGraph API.
	NerdGraphBaseURL string

	// ServiceName is for New Relic internal use only.
	ServiceName string

	// LogLevel can be one of the following values:
	// "panic", "fatal", "error", "warn", "info", "debug", "trace"
	LogLevel string

	// LogJSON toggles formatting of log entries in JSON format.
	LogJSON bool

	// Logger allows customization of the client's underlying logger.
	Logger logging.Logger
}

// New creates a default configuration and returns it
func New() Config {
	regCopy := *region.Default

	return Config{
		Region:    &regCopy,
		UserAgent: "newrelic/newrelic-client-go",
		LogLevel:  "info",
	}
}

// GetLogger returns a logger instance based on the config values.
func (c *Config) GetLogger() logging.Logger {
	if c.Logger != nil {
		return c.Logger
	}

	l := logging.NewStructuredLogger().
		SetDefaultFields(map[string]string{"newrelic-client-go": version.Version}).
		LogJSON(c.LogJSON).
		SetLogLevel(c.LogLevel)

	c.Logger = l
	return l
}
