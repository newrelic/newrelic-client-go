// Package synthetics provides a programmatic API for interacting with the New Relic Synthetics product.
package synthetics

import (
	"net/http"
	"strings"

	nrhttp "github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/logging"
)

// Synthetics is used to communicate with the New Relic Synthetics product.
type Synthetics struct {
	client nrhttp.Client
	config config.Config
	logger logging.Logger
	pager  nrhttp.Pager
}

// ErrorResponse represents an error response from New Relic Synthetics.
type ErrorResponse struct {
	Message            string        `json:"error,omitempty"`
	Messages           []ErrorDetail `json:"errors,omitempty"`
	ServerErrorMessage string        `json:"message,omitempty"`
}

// ErrorDetail represents an single error from New Relic Synthetics.
type ErrorDetail struct {
	Message string `json:"error,omitempty"`
}

// Error surfaces an error message from the New Relic Synthetics error response.
func (e *ErrorResponse) Error() string {
	if e.ServerErrorMessage != "" {
		return e.ServerErrorMessage
	}

	if e.Message != "" {
		return e.Message
	}

	if len(e.Messages) > 0 {
		messages := []string{}
		for _, m := range e.Messages {
			messages = append(messages, m.Message)
		}
		return strings.Join(messages, ", ")
	}

	return ""
}

// New creates a new instance of ErrorResponse.
func (e *ErrorResponse) New() nrhttp.ErrorResponse {
	return &ErrorResponse{}
}

func (e *ErrorResponse) IsNotFound() bool {
	return false
}

func (e *ErrorResponse) IsRetryableError() bool {
	return false
}

// IsUnauthorized checks the response for a 401 Unauthorize HTTP status code.
func (e *ErrorResponse) IsUnauthorized(resp *http.Response) bool {
	return resp.StatusCode == http.StatusUnauthorized
}

// New is used to create a new Synthetics client instance.
func New(config config.Config) Synthetics {
	client := nrhttp.NewClient(config)
	client.SetAuthStrategy(&nrhttp.PersonalAPIKeyCapableV2Authorizer{})
	client.SetErrorValue(&ErrorResponse{})

	pkg := Synthetics{
		client: client,
		config: config,
		logger: config.GetLogger(),
		pager:  &nrhttp.LinkHeaderPager{},
	}

	return pkg
}
