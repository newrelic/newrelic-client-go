package eventstometrics

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type Eventstometrics struct {
	client http.Client
	logger logging.Logger
}

func New(config config.Config) Eventstometrics {
	return Eventstometrics{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}
