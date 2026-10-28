package federatedlogs

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type Federatedlogs struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

func New(config config.Config) Federatedlogs {
	client := http.NewClient(config)

	pkg := Federatedlogs{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
