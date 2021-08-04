// Experimental.  For NR internal use only.
package servicelevel

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/logging"
)

type Servicelevel struct {
	client http.Client
	logger logging.Logger
}

func New(config config.Config) Servicelevel {
	return Servicelevel{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}
