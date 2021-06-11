package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data interface{} `json:"data"`
}

// GraphQLError represents a single error.
type GraphQLError struct {
	Message    string   `json:"message,omitempty"`
	Path       []string `json:"path,omitempty"`
	Extensions struct {
		ErrorClass string `json:"errorClass,omitempty"`
		ErrorCode  string `json:"error_code,omitempty"`
	} `json:"extensions,omitempty"`
	DownstreamResponse []GraphQLDownstreamResponse `json:"downstreamResponse,omitempty"`
}

// GraphQLDownstreamResponse represents an error's downstream response.
type GraphQLDownstreamResponse struct {
	Extensions struct {
		Code             string `json:"code,omitempty"`
		ValidationErrors []struct {
			Name   string `json:"name,omitempty"`
			Reason string `json:"reason,omitempty"`
		} `json:"validationErrors,omitempty"`
	} `json:"extensions,omitempty"`
	Message string `json:"message,omitempty"`
}

// GraphQLErrorResponse represents a default error response body.
type GraphQLErrorResponse struct {
	Errors []GraphQLError `json:"errors"`
}

func (r *GraphQLErrorResponse) Error() string {
	if len(r.Errors) > 0 {
		messages := []string{}
		for _, e := range r.Errors {

			if e.Message != "" {
				messages = append(messages, e.Message)
			}

			if e.DownstreamResponse != nil {
				f, _ := json.Marshal(e.DownstreamResponse)
				messages = append(messages, string(f))
			}
		}
		return strings.Join(messages, ", ")
	}

	return ""
}

// IsNotFound determines if the error is due to a missing resource.
func (r *GraphQLErrorResponse) IsNotFound() bool {
	return false
}

// IsRetryableError determines if the error is due to a server timeout, or another error that we might want to retry.
func (r *GraphQLErrorResponse) IsRetryableError() bool {

	if len(r.Errors) == 0 {
		return false
	}

	for _, err := range r.Errors {
		errorClass := err.Extensions.ErrorClass
		if errorClass == "TIMEOUT" || errorClass == "INTERNAL_SERVER_ERROR" || errorClass == "SERVER_ERROR" {
			return true
		}

		for _, downstreamErr := range err.DownstreamResponse {
			code := downstreamErr.Extensions.Code
			if code == "INTERNAL_SERVER_ERROR" || code == "SERVER_ERROR" {
				return true
			}
		}
	}

	return false
}

// IsUnauthorized checks a NerdGraph response for a 401 Unauthorize HTTP status code,
// then falls back to check the nested extensions error_code field for `BAD_API_KEY`.
func (r *GraphQLErrorResponse) IsUnauthorized(resp *http.Response) bool {
	if len(r.Errors) == 0 {
		return false
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return true
	}

	// Handle invalid or missing API key
	for _, err := range r.Errors {
		if err.Extensions.ErrorCode == "BAD_API_KEY" {
			return true
		}
	}

	return false
}

// New creates a new instance of GraphQLErrorRepsonse.
func (r *GraphQLErrorResponse) New() ErrorResponse {
	return &GraphQLErrorResponse{}
}
