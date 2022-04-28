package syntheticsv2

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/logging"
)

type SyntheticsV2 struct {
	client http.Client
	config config.Config
	logger logging.Logger
	pager  http.Pager
}

func New(config config.Config) SyntheticsV2 {

	client := http.NewClient(config)
	client.SetAuthStrategy(&http.PersonalAPIKeyCapableV2Authorizer{})

	pkg := SyntheticsV2{
		client: client,
		config: config,
		logger: config.GetLogger(),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}
