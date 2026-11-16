// Experimental.  For NR internal use only.
package maintenancewindow

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type Maintenancewindow struct {
	client http.Client
	logger logging.Logger
}

func New(config config.Config) Maintenancewindow {
	return Maintenancewindow{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}
