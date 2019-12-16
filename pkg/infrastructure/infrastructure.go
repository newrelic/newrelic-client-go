package infrastructure

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var baseURLs = map[config.RegionType]string{
	config.Region.US:      "https://infra-api.newrelic.com/v2",
	config.Region.EU:      "https://infra-api.eu.newrelic.com/v2",
	config.Region.Staging: "https://staging-infra-api.newrelic.com/v2",
}

// Infrastructure is used to communicate with the New Relic Infrastructure product.
type Infrastructure struct {
	client http.NewRelicClient
}

// New is used to create a new Infrastructure client instance.
func New(config config.Config) Infrastructure {
	if config.BaseURL == "" {
		config.BaseURL = baseURLs[config.Region]
	}

	c := http.NewClient(config).
		SetError(&ErrorResponse{}).
		SetPager(&LinkBodyPager{})

	pkg := Infrastructure{
		client: c,
	}

	return pkg
}
