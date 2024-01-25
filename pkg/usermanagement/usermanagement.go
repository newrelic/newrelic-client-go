package usermanagement

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type UserManagement struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

func New(config config.Config) UserManagement {
	client := http.NewClient(config)

	pkg := UserManagement{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
