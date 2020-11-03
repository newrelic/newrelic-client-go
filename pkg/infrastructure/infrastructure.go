// Package infrastructure provides metadata about the underlying Infrastructure API.
package infrastructure

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Infrastructure is used to communicate with the New Relic Infrastructure product.
type Infrastructure struct {
	client http.Client
}

// New is used to create a new APM client instance.
func New(config config.Config) Infrastructure {
	client := http.NewClient(config)
	client.SetAuthStrategy(&http.PersonalAPIKeyCapableV2Authorizer{})

	return Infrastructure{
		client: client,
	}
}
