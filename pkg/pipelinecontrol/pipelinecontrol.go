package pipelinecontrol

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

type Pipelinecontrol struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

func New(config config.Config) Pipelinecontrol {
	client := http.NewClient(config)
	pkg := Pipelinecontrol{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}

func newMockClient(t *testing.T, mockJSONResponse string, statusCode int) Pipelinecontrol {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)
	return New(tc)
}

func newIntegrationTestClient(t *testing.T) Pipelinecontrol {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
