//nolint:revive
package ai

import (
	"encoding/json"
	"fmt"
)

// AiNotificationsBasicAuth - Basic user and password authentication.
type AiNotificationsBasicAuth struct {
	// Username
	User string `json:"user"`
}

func (x *AiNotificationsBasicAuth) ImplementsAiNotificationsAuth() {}

// GetUser returns a pointer to the value of User from AiNotificationsBasicAuth
func (x AiNotificationsBasicAuth) GetUser() string {
	return x.User
}

// GetUser returns a pointer to the value of User from AiNotificationsBasicAuth
func (x AiNotificationsBasicAuth) GetPrefix() string {
	return ""
}

// AiNotificationsTokenAuth - Token based authentication
type AiNotificationsTokenAuth struct {
	// Token Prefix
	Prefix string `json:"prefix"`
}

func (x *AiNotificationsTokenAuth) ImplementsAiNotificationsAuth() {}

// GetUser returns a pointer to the value of User from AiNotificationsBasicAuth
func (x AiNotificationsTokenAuth) GetUser() string {
	return ""
}

// GetPrefix returns a pointer to the value of Prefix from AiNotificationsTokenAuth
func (x AiNotificationsTokenAuth) GetPrefix() string {
	return x.Prefix
}

// AiWorkflowsConfigurationDto - Enrichment configuration object
type AiWorkflowsConfigurationDtoInterface interface {
	ImplementsAiWorkflowsConfigurationDto()
}

// UnmarshalAiWorkflowsConfigurationDtoInterface unmarshals the interface into the correct type
// based on __typename provided by GraphQL
func UnmarshalAiWorkflowsConfigurationDtoInterface(b []byte) (*AiWorkflowsConfigurationDtoInterface, error) {
	var err error

	var rawMessageAiWorkflowsConfigurationDto map[string]*json.RawMessage
	err = json.Unmarshal(b, &rawMessageAiWorkflowsConfigurationDto)
	if err != nil {
		return nil, err
	}

	// Nothing to unmarshal
	if len(rawMessageAiWorkflowsConfigurationDto) < 1 {
		return nil, nil
	}

	var typeName string

	if rawTypeName, ok := rawMessageAiWorkflowsConfigurationDto["__typename"]; ok {
		err = json.Unmarshal(*rawTypeName, &typeName)
		if err != nil {
			return nil, err
		}

		switch typeName {
		case "AiWorkflowsNrqlConfigurationDto":
			var interfaceType AiWorkflowsNRQLConfigurationDto
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx AiWorkflowsConfigurationDtoInterface = &interfaceType

			return &xxx, nil
		}
	} else {
		keys := []string{}
		for k := range rawMessageAiWorkflowsConfigurationDto {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("interface AiWorkflowsConfigurationDto did not include a __typename field for inspection: %s", keys)
	}

	return nil, fmt.Errorf("interface AiWorkflowsConfigurationDto was not matched against all PossibleTypes: %s", typeName)
}

// AiNotificationsSuggestionError - Object for suggestion errors
type AiNotificationsSuggestionError struct {
	// SuggestionError description
	Description string `json:"description"`
	// SuggestionError details
	Details string `json:"details"`
	// SuggestionError type
	Type AiNotificationsErrorType `json:"type"`
}

// AiNotificationsErrorType - Error types
type AiNotificationsErrorType string

func (x *AiNotificationsSuggestionError) ImplementsAiNotificationsError() {}

// AiNotificationsDataValidationError - Object for validation errors
type AiNotificationsDataValidationError struct {
	// Top level error details
	Details string `json:"details"`
	// List of invalid fields
	Fields []AiNotificationsFieldError `json:"fields"`
}

// AiNotificationsFieldError - Invalid field object
type AiNotificationsFieldError struct {
	// Field name
	Field string `json:"field"`
	// Validation error
	Message string `json:"message"`
}

func (x *AiNotificationsDataValidationError) ImplementsAiNotificationsError() {}

// AiNotificationsConstraintError - Missing constraint error. Constraints can be retrieved using suggestion api
type AiNotificationsConstraintError struct {
	// Names of other constraints this constraint is dependent on
	Dependencies []string `json:"dependencies"`
	// Name of the missing constraint
	Name string `json:"name"`
}

func (x *AiNotificationsConstraintError) ImplementsAiNotificationsError() {}

// AiNotificationsResponseError - Response error object
type AiNotificationsResponseError struct {
	// Error description
	Description string `json:"description"`
	// Error details
	Details string `json:"details"`
	// Error type
	Type AiNotificationsErrorType `json:"type"`
}

func (x *AiNotificationsResponseError) ImplementsAiNotificationsError() {}

// AiNotificationsDestinationFilter - Filter destination object
type AiNotificationsDestinationFilter struct {
	// id
	ID string `json:"id,omitempty"`
}

// AiNotificationsChannelFilter - Filter channel object
type AiNotificationsChannelFilter struct {
	// id
	ID string `json:"id,omitempty"`
}
