package pathpoint

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type Pathpoint struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

func New(config config.Config) Pathpoint {
	client := http.NewClient(config)
	pkg := Pathpoint{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
