// Manually added types
package abtest

// AbTestGetVariationResponse - A response returning a VariationKey and a potential array of abTestError
type AbTestGetVariationResponse struct {
	// An array of errors
	Errors []abTestError `json:"errors,omitempty"`
	// Variation key denoted the inclusion status of an accountId for a given experiment
	VariationKey string `json:"variationKey,omitempty"`
}

// abTestError - An error Message combined with its Type of error
type abTestError struct {
	// A description of the error
	Message string `json:"message"`
	// The kind of error which occurred
	Type string `json:"type"`
}
