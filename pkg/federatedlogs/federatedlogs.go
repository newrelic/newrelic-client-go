package federatedlogs

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type FederatedLogs struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

func New(config config.Config) FederatedLogs {
	client := http.NewClient(config)

	pkg := FederatedLogs{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
