package alerts

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/infrastructure"
)

// Alerts is used to communicate with New Relic Alerts.
type Alerts struct {
	client      http.NewRelicClient
	infraClient http.NewRelicClient
	logger      logging.Logger
	pager       http.Pager
}

// New is used to create a new Alerts client instance.
func New(config config.Config) Alerts {
	infraConfig := config

	if infraConfig.InfrastructureBaseURL == "" {
		infraConfig.InfrastructureBaseURL = region.Parse(config.Region).BaseURL()
	}

	infraConfig.BaseURL = infraConfig.InfrastructureBaseURL

	infraClient := http.NewClient(infraConfig)
	infraClient.SetErrorValue(&infrastructure.ErrorResponse{})

	client := http.NewClient(config)
	client.AuthStrategy = &http.PersonalAPIKeyCapableV2Authorizer{}

	pkg := Alerts{
		client:      client,
		infraClient: infraClient,
		logger:      config.GetLogger(),
		pager:       &http.LinkHeaderPager{},
	}

	return pkg
}
