package entityrelationship

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"testing"
)

type Entityrelationship struct {
	client http.Client
	logger logging.Logger
	config config.Config
}

// nolint
func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Entityrelationship {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

func New(config config.Config) Entityrelationship {
	client := http.NewClient(config)
	pkg := Entityrelationship{
		client: client,
		logger: config.GetLogger(),
		config: config,
	}
	return pkg
}
