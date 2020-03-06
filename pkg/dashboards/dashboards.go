package dashboards

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Dashboards is used to communicate with the New Relic Dashboards product.
type Dashboards struct {
	client http.NewRelicClient
	logger logging.Logger
	pager  http.Pager
}

// New is used to create a new Dashboards client instance.
func New(config config.Config) Dashboards {
	pkg := Dashboards{
		client: http.NewClient(config),
		logger: config.GetLogger(),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}
