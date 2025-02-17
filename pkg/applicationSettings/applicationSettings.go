package applicationsettings

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"testing"
)

type ApplicationSettings struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

func New(config config.Config) ApplicationSettings {
	client := http.NewClient(config)
	pkg := ApplicationSettings{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}

func newIntegrationTestClient(t *testing.T) ApplicationSettings {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}

func newMockResponseApm(t *testing.T, mockJSONResponse string, statusCode int) ApplicationSettings {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

// nolint
func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) entities.Entities {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return entities.New(tc)
}
