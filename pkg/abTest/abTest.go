// Manually added
package abtest

import (
	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

type AbTest struct {
	client http.Client
	logger logging.Logger
}

// New is used to create a new Ab Test.
func New(config config.Config) AbTest {
	client := http.NewClient(config)

	pkg := AbTest{
		client: client,
		logger: config.GetLogger(),
	}

	return pkg
}
