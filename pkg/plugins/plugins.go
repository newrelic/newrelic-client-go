package plugins

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Plugins is used to communicate with the New Relic Plugins product.
type Plugins struct {
	client http.NewRelicClient
	pager  http.Pager
}

// New is used to create a new Plugins client instance.
func New(config config.Config) Plugins {
	pkg := Plugins{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}

// ListPlugins returns a list of Plugins associated with an account.
func (plugins *Plugins) ListPlugins() ([]*Plugin, error) {
	response := pluginsResponse{}
	results := []*Plugin{}
	nextURL := "/plugins.json"

	for nextURL != "" {
		resp, err := plugins.client.Get(nextURL, nil, &response)

		if err != nil {
			return nil, err
		}

		results = append(results, response.Plugins...)

		paging := plugins.pager.Parse(resp)
		nextURL = paging.Next
	}

	return results, nil
}

type pluginsResponse struct {
	Plugins []*Plugin `json:"plugins,omitempty"`
}
