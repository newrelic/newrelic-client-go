package infrastructure

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/newrelic/newrelic-client-go/internal/http"
)

// LinkBodyPager represents a pagination implementation that parses the pagination context from the response body.
type LinkBodyPager struct{}

// Parse is used to parse a pagination context from an HTTP response.
func (l *LinkBodyPager) Parse(res *resty.Response) (*http.Paging, error) {
	paging := http.Paging{}

	body := res.Body()
	if len(body) == 0 {
		return &paging, nil
	}

	linksResponse := struct {
		Links struct {
			Next string `json:"next"`
		} `json:"links"`
	}{}

	err := json.Unmarshal(body, &linksResponse)
	if err != nil {
		return nil, err
	}

	if linksResponse.Links.Next != "" {
		paging.Next = linksResponse.Links.Next
	}

	return &paging, nil
}
