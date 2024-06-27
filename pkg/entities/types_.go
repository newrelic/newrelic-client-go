package entities

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
)

// AiNotificationsAuth - Authentication interface
type AiNotificationsAuthInterface interface {
	ImplementsAiNotificationsAuth()
}

// UnmarshalAiNotificationsAuthInterface unmarshals the interface into the correct type
// based on __typename provided by GraphQL
func UnmarshalAiNotificationsAuthInterface(b []byte) (*AiNotificationsAuthInterface, error) {
	var err error

	var rawMessageAiNotificationsAuth map[string]*json.RawMessage
	err = json.Unmarshal(b, &rawMessageAiNotificationsAuth)
	if err != nil {
		return nil, err
	}

	// Nothing to unmarshal
	if len(rawMessageAiNotificationsAuth) < 1 {
		return nil, nil
	}

	var typeName string

	if rawTypeName, ok := rawMessageAiNotificationsAuth["__typename"]; ok {
		err = json.Unmarshal(*rawTypeName, &typeName)
		if err != nil {
			return nil, err
		}

		switch typeName {
		case "AiNotificationsBasicAuth":
			var interfaceType ai.AiNotificationsBasicAuth
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx AiNotificationsAuthInterface = &interfaceType

			return &xxx, nil
		case "AiNotificationsTokenAuth":
			var interfaceType ai.AiNotificationsTokenAuth
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx AiNotificationsAuthInterface = &interfaceType

			return &xxx, nil
		}
	} else {
		keys := []string{}
		for k := range rawMessageAiNotificationsAuth {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("interface AiNotificationsAuth did not include a __typename field for inspection: %s", keys)
	}

	return nil, fmt.Errorf("interface AiNotificationsAuth was not matched against all PossibleTypes: %s", typeName)
}

// AiNotificationsAuth - Authentication interface
type AiNotificationsAuth struct {
}

func (x *AiNotificationsAuth) ImplementsAiNotificationsAuth() {}

// AiWorkflowsConfiguration - Enrichment configuration object
type AiWorkflowsConfiguration struct {
}

func (x *AiWorkflowsConfiguration) ImplementsAiWorkflowsConfiguration() {}

// AiWorkflowsConfiguration - Enrichment configuration object
type AiWorkflowsConfigurationInterface interface {
	ImplementsAiWorkflowsConfiguration()
}

// UnmarshalAiWorkflowsConfigurationInterface unmarshals the interface into the correct type
// based on __typename provided by GraphQL
func UnmarshalAiWorkflowsConfigurationInterface(b []byte) (*AiWorkflowsConfigurationInterface, error) {
	var err error

	var rawMessageAiWorkflowsConfiguration map[string]*json.RawMessage
	err = json.Unmarshal(b, &rawMessageAiWorkflowsConfiguration)
	if err != nil {
		return nil, err
	}

	// Nothing to unmarshal
	if len(rawMessageAiWorkflowsConfiguration) < 1 {
		return nil, nil
	}

	var typeName string

	if rawTypeName, ok := rawMessageAiWorkflowsConfiguration["__typename"]; ok {
		err = json.Unmarshal(*rawTypeName, &typeName)
		if err != nil {
			return nil, err
		}

		switch typeName {
		case "AiWorkflowsNrqlConfiguration":
			var interfaceType ai.AiWorkflowsNRQLConfiguration
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx AiWorkflowsConfigurationInterface = &interfaceType

			return &xxx, nil
		}
	} else {
		keys := []string{}
		for k := range rawMessageAiWorkflowsConfiguration {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("interface AiWorkflowsConfiguration did not include a __typename field for inspection: %s", keys)
	}

	return nil, fmt.Errorf("interface AiWorkflowsConfiguration was not matched against all PossibleTypes: %s", typeName)
}
