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
