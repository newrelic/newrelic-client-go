// Package region describes the operational regions defined for New Relic
//
// Regions are geographical locations where the New Relic platform operates
// and this package provides an abstraction layer for handling them within
// the New Relic Client and underlying APIs
package region

import (
	"strings"
)

// Name is the name of a New Relic region
type Name string

// Region represents the members of the Region enumeration.
type Region struct {
	name                  string
	restBaseURL           string
	infrastructureBaseURL string
	syntheticsBaseURL     string
	nerdGraphBaseURL      string
}

// String returns a human readable value for the specified Region
func (r *Region) String() string {
	if r != nil && r.name != "" {
		return r.name
	}

	return "(Unknown)"
}

//
// NerdGraph - the future
//

// SetNerdGraphBaseURL Allows overriding the NerdGraph Base URL
func (r *Region) SetNerdGraphBaseURL(url string) {
	if r != nil && url != "" {
		r.nerdGraphBaseURL = url
	}
}

// NerdGraphURL returns the Full URL for Infrastructure REST API Calls, with any additional path elements appended
func (r *Region) NerdGraphURL(path ...string) string {
	if r == nil {
		return ""
	}

	elements := append([]string{r.nerdGraphBaseURL}, path...)

	return strings.Join(elements, "/")
}

//
// REST
//

// SetRestBaseURL Allows overriding the REST Base URL
func (r *Region) SetRestBaseURL(url string) {
	if r != nil && url != "" {
		r.restBaseURL = url
	}
}

// RestURL returns the Full URL for REST API Calls, with any additional path elements appended
func (r *Region) RestURL(path ...string) string {
	if r == nil {
		return ""
	}

	elements := append([]string{r.restBaseURL}, path...)

	return strings.Join(elements, "/")
}

//
// Infrastructure
//

// SetInfrastructureBaseURL Allows overriding the Infrastructure Base URL
func (r *Region) SetInfrastructureBaseURL(url string) {
	if r != nil && url != "" {
		r.infrastructureBaseURL = url
	}
}

// InfrastructureURL returns the Full URL for Infrastructure REST API Calls, with any additional path elements appended
func (r *Region) InfrastructureURL(path ...string) string {
	if r == nil {
		return ""
	}

	elements := append([]string{r.infrastructureBaseURL}, path...)

	return strings.Join(elements, "/")
}

//
// Synthetics
//

// SetSyntheticsBaseURL Allows overriding the Synthetics Base URL
func (r *Region) SetSyntheticsBaseURL(url string) {
	if r != nil && url != "" {
		r.syntheticsBaseURL = url
	}
}

// SyntheticsURL returns the Full URL for Infrastructure REST API Calls, with any additional path elements appended
func (r *Region) SyntheticsURL(path ...string) string {
	if r == nil {
		return ""
	}

	elements := append([]string{r.syntheticsBaseURL}, path...)

	return strings.Join(elements, "/")
}
