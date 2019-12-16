package infrastructure

// ErrorResponse represents an error response from New Relic Infrastructure.
type ErrorResponse struct {
	Errors []*ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail represents the details of an error response from New Relic Infrastructure.
type ErrorDetail struct {
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// Error surfaces an error message from the Infrastructure error response.
func (e *ErrorResponse) Error() string {
	if e != nil && len(e.Errors) > 0 && e.Errors[0].Detail != "" {
		return e.Errors[0].Detail
	}
	return "Unknown error"
}
