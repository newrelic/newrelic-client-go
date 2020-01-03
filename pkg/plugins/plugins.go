package plugins

import (
	"strconv"
	"strings"

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
	return Plugins{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}
}

// ListPluginsParams represents a set of filters to be
// used when querying New Relic key transactions.
type ListPluginsParams struct {
	GUID string
	IDs  []int
}

// ListPlugins returns a list of Plugins associated with an account.
func (plugins *Plugins) ListPlugins(params *ListPluginsParams) ([]*Plugin, error) {
	response := pluginsResponse{}
	results := []*Plugin{}
	paramsMap := buildListPluginsParamsMap(params)
	nextURL := "/plugins.json"

	for nextURL != "" {
		resp, err := plugins.client.Get(nextURL, &paramsMap, &response)

		if err != nil {
			return nil, err
		}

		results = append(results, response.Plugins...)

		paging := plugins.pager.Parse(resp)
		nextURL = paging.Next
	}

	return results, nil
}

func buildListPluginsParamsMap(params *ListPluginsParams) map[string]string {
	paramsMap := map[string]string{}

	if params == nil {
		return paramsMap
	}

	if params.GUID != "" {
		paramsMap["filter[guid]"] = params.GUID
	}

	if params.IDs != nil && len(params.IDs) > 0 {
		paramsMap["filter[ids]"] = intArrayToString(params.IDs)
	}

	return paramsMap
}

func intArrayToString(integers []int) string {
	sArray := []string{}

	for _, n := range integers {
		sArray = append(sArray, strconv.Itoa(n))
	}

	return strings.Join(sArray, ",")
}

type pluginsResponse struct {
	Plugins []*Plugin `json:"plugins,omitempty"`
}
