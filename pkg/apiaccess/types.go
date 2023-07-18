// Code generated by tutone: DO NOT EDIT
package apiaccess

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/newrelic-client-go/v2/pkg/accounts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/users"
)

// APIAccessIngestKeyErrorType - The type of error.
type APIAccessIngestKeyErrorType string

var APIAccessIngestKeyErrorTypeTypes = struct {
	// Occurs when the user issuing the mutation does not have sufficient permissions to perform the action for a key.
	FORBIDDEN APIAccessIngestKeyErrorType
	// Occurs when the action taken on a key did not successfully pass validation.
	INVALID APIAccessIngestKeyErrorType
	// Occurs when the requested key `id` was not found.
	NOT_FOUND APIAccessIngestKeyErrorType
}{
	// Occurs when the user issuing the mutation does not have sufficient permissions to perform the action for a key.
	FORBIDDEN: "FORBIDDEN",
	// Occurs when the action taken on a key did not successfully pass validation.
	INVALID: "INVALID",
	// Occurs when the requested key `id` was not found.
	NOT_FOUND: "NOT_FOUND",
}

// APIAccessIngestKeyType - The type of ingest key, which dictates what types of agents can use it to report.
type APIAccessIngestKeyType string

var APIAccessIngestKeyTypeTypes = struct {
	// Ingest keys of type `BROWSER` mean browser agents will use them to report data to New Relic.
	BROWSER APIAccessIngestKeyType
	// For ingest keys of type `LICENSE`: APM and Infrastructure agents use the key to report data to New Relic.
	LICENSE APIAccessIngestKeyType
}{
	// Ingest keys of type `BROWSER` mean browser agents will use them to report data to New Relic.
	BROWSER: "BROWSER",
	// For ingest keys of type `LICENSE`: APM and Infrastructure agents use the key to report data to New Relic.
	LICENSE: "LICENSE",
}

// APIAccessKeyType - The type of key.
type APIAccessKeyType string

var APIAccessKeyTypeTypes = struct {
	// An ingest key is used by New Relic agents to authenticate with New Relic and send data to the assigned account.
	INGEST APIAccessKeyType
	// A user key is used by New Relic users to authenticate with New Relic and to interact with the New Relic GraphQL API .
	USER APIAccessKeyType
}{
	// An ingest key is used by New Relic agents to authenticate with New Relic and send data to the assigned account.
	INGEST: "INGEST",
	// A user key is used by New Relic users to authenticate with New Relic and to interact with the New Relic GraphQL API .
	USER: "USER",
}

// APIAccessUserKeyErrorType - The type of error.
type APIAccessUserKeyErrorType string

var APIAccessUserKeyErrorTypeTypes = struct {
	// Occurs when the user issuing the mutation does not have sufficient permissions to perform the action for a key.
	FORBIDDEN APIAccessUserKeyErrorType
	// Occurs when the action taken on a key did not successfully pass validation.
	INVALID APIAccessUserKeyErrorType
	// Occurs when the requested key `id` was not found.
	NOT_FOUND APIAccessUserKeyErrorType
}{
	// Occurs when the user issuing the mutation does not have sufficient permissions to perform the action for a key.
	FORBIDDEN: "FORBIDDEN",
	// Occurs when the action taken on a key did not successfully pass validation.
	INVALID: "INVALID",
	// Occurs when the requested key `id` was not found.
	NOT_FOUND: "NOT_FOUND",
}

// APIAccessActorStitchedFields -
type APIAccessActorStitchedFields struct {
	// Fetch a single key by ID and type.
	Key APIAccessKeyInterface `json:"key,omitempty"`
	// A list of keys scoped to the current actor and filter arguments. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).
	KeySearch APIAccessKeySearchResult `json:"keySearch,omitempty"`
}

// special
func (x *APIAccessActorStitchedFields) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	for k, v := range objMap {
		if v == nil {
			continue
		}

		switch k {
		case "key":
			if v == nil {
				continue
			}
			xxx, err := UnmarshalAPIAccessKeyInterface(*v)
			if err != nil {
				return err
			}

			if xxx != nil {
				x.Key = *xxx
			}
		case "keySearch":
			err = json.Unmarshal(*v, &x.KeySearch)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// APIAccessCreateIngestKeyInput - The input for any ingest keys you want to create. Each ingest key must have a type that communicates what kind of data it is for. You can optionally add a name or notes to your key, which can be updated later.
type APIAccessCreateIngestKeyInput struct {
	// The account ID indicating which account you want to make the key for. This cannot be updated once created.
	AccountID int `json:"accountId"`
	// The type of ingest key you want to create. This cannot be updated once created.
	IngestType APIAccessIngestKeyType `json:"ingestType"`
	// The name of the key. This can be updated later.
	Name string `json:"name,omitempty"`
	// Any notes about this ingest key. This can be updated later.
	Notes string `json:"notes,omitempty"`
}

// APIAccessCreateInput - The input object to create one or more keys.
type APIAccessCreateInput struct {
	// Ingest keys are used by agents to report data about your applications to New Relic. Each ingest key input entered here must have a type that communicates what kind of data it is for. You can optionally add a name or notes to your key, which can be updated later.
	Ingest []APIAccessCreateIngestKeyInput `json:"ingest,omitempty"`
	// Create user keys. You can optionally add a name or notes to your key, which can be updated later.
	User []APIAccessCreateUserKeyInput `json:"user,omitempty"`
}

// APIAccessCreateKeyResponse - The response of the create keys mutation.
type APIAccessCreateKeyResponse struct {
	// Lists all successfully created keys.
	CreatedKeys []APIAccessKeyInterface `json:"createdKeys,omitempty"`
	// Lists all errors for keys that could not be created. Each error maps to a single key input.
	Errors []APIAccessKeyErrorInterface `json:"errors,omitempty"`
}

// special
func (x *APIAccessCreateKeyResponse) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	for k, v := range objMap {
		if v == nil {
			continue
		}

		switch k {
		case "createdKeys":
			if v == nil {
				continue
			}
			var rawMessageCreatedKeys []*json.RawMessage
			err = json.Unmarshal(*v, &rawMessageCreatedKeys)
			if err != nil {
				return err
			}

			for _, m := range rawMessageCreatedKeys {
				xxx, err := UnmarshalAPIAccessKeyInterface(*m)
				if err != nil {
					return err
				}

				if xxx != nil {
					x.CreatedKeys = append(x.CreatedKeys, *xxx)
				}
			}
		case "errors":
			if v == nil {
				continue
			}
			var rawMessageErrors []*json.RawMessage
			err = json.Unmarshal(*v, &rawMessageErrors)
			if err != nil {
				return err
			}

			for _, m := range rawMessageErrors {
				xxx, err := UnmarshalAPIAccessKeyErrorInterface(*m)
				if err != nil {
					return err
				}

				if xxx != nil {
					x.Errors = append(x.Errors, *xxx)
				}
			}
		}
	}

	return nil
}

// APIAccessCreateUserKeyInput - The input for any ingest keys you want to create. Each ingest key must have a type that communicates what kind of data it is for. You can optionally add a name or notes to your key, which can be updated later.
type APIAccessCreateUserKeyInput struct {
	// The account ID indicating which account you want to make the key for. This cannot be updated once created.
	AccountID int `json:"accountId"`
	// The name of the key. This can be updated later.
	Name string `json:"name,omitempty"`
	// Any notes about this ingest key. This can be updated later.
	Notes string `json:"notes,omitempty"`
	// The user ID indicating which user you want to make the key for. This cannot be updated once created.
	UserID int `json:"userId"`
}

// APIAccessDeleteInput - The input to delete keys.
type APIAccessDeleteInput struct {
	// A list of the ingest key `id`s that you want to delete.
	IngestKeyIDs []string `json:"ingestKeyIds,omitempty"`
	// A list of the user key `id`s that you want to delete.
	UserKeyIDs []string `json:"userKeyIds,omitempty"`
}

// APIAccessDeleteKeyResponse - The response of the key delete mutation.
type APIAccessDeleteKeyResponse struct {
	// The `id`s of the successfully deleted ingest keys and any errors that occurred when deleting keys.
	DeletedKeys []APIAccessDeletedKey `json:"deletedKeys,omitempty"`
	// Lists all errors for keys that could not be deleted. Each error maps to a single key input.
	Errors []APIAccessKeyErrorInterface `json:"errors,omitempty"`
}

// special
func (x *APIAccessDeleteKeyResponse) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	for k, v := range objMap {
		if v == nil {
			continue
		}

		switch k {
		case "deletedKeys":
			err = json.Unmarshal(*v, &x.DeletedKeys)
			if err != nil {
				return err
			}
		case "errors":
			if v == nil {
				continue
			}
			var rawMessageErrors []*json.RawMessage
			err = json.Unmarshal(*v, &rawMessageErrors)
			if err != nil {
				return err
			}

			for _, m := range rawMessageErrors {
				xxx, err := UnmarshalAPIAccessKeyErrorInterface(*m)
				if err != nil {
					return err
				}

				if xxx != nil {
					x.Errors = append(x.Errors, *xxx)
				}
			}
		}
	}

	return nil
}

// APIAccessDeletedKey - The deleted key response of the key delete mutation.
type APIAccessDeletedKey struct {
	// The `id` of the deleted key.
	ID string `json:"id,omitempty"`
}

// APIAccessIngestKey - An ingest key.
type APIAccessIngestKey struct {
	// The account this key is in.
	Account accounts.AccountReference `json:"account,omitempty"`
	// The account attached to the ingest key. Agents using this key will report to the account the key belongs to.
	AccountID int `json:"accountId,omitempty"`
	// The UNIX epoch when the key was created, in seconds.
	CreatedAt EpochSeconds `json:"createdAt,omitempty"`
	// The ID of the ingest key. This can be used to identify a key without revealing the key itself (used to update and delete).
	ID string `json:"id,omitempty"`
	// The type of ingest key, which dictates what types of agents can use it to report.
	IngestType APIAccessIngestKeyType `json:"ingestType,omitempty"`
	// The keystring of the key.
	Key string `json:"key,omitempty"`
	// The name of the key.
	Name string `json:"name,omitempty"`
	// Any notes can be attached to an key.
	Notes string `json:"notes,omitempty"`
	// The type of key, indicating what New Relic APIs it can be used to access.
	Type APIAccessKeyType `json:"type,omitempty"`
}

func (x *APIAccessIngestKey) ImplementsAPIAccessKey() {}

// APIAccessIngestKeyError - An ingest key error. Each error maps to a single key input.
type APIAccessIngestKeyError struct {
	// The account ID of the key.
	AccountID int `json:"accountId,omitempty"`
	// The error type of the error.
	ErrorType APIAccessIngestKeyErrorType `json:"errorType,omitempty"`
	// The `id` of the key being updated.
	ID string `json:"id,omitempty"`
	// The ingest type of the key.
	IngestType APIAccessIngestKeyType `json:"ingestType,omitempty"`
	// A message about why the key creation failed.
	Message string `json:"message,omitempty"`
	// The type of the key.
	Type APIAccessKeyType `json:"type,omitempty"`
}

func (x *APIAccessIngestKeyError) ImplementsAPIAccessKeyError() {}

// APIAccessKey - A key for accessing New Relic APIs.
type APIAccessKey struct {
	// The UNIX epoch when the key was created, in seconds.
	CreatedAt EpochSeconds `json:"createdAt,omitempty"`
	// The ID of the key. This can be used to identify a key without revealing the key itself (used to update and delete).
	ID string `json:"id,omitempty"`
	// The keystring of the key.
	Key string `json:"key,omitempty"`
	// The name of the key. This can be used a short identifier for easy reference.
	Name string `json:"name,omitempty"`
	// Any notes can be attached to a key. This is intended for more a more detailed description of the key use if desired.
	Notes string `json:"notes,omitempty"`
	// The type of key, indicating what New Relic APIs it can be used to access.
	Type APIAccessKeyType `json:"type,omitempty"`
}

func (x *APIAccessKey) ImplementsAPIAccessKey() {}

// APIAccessKeyError - A key error. Each error maps to a single key input.
type APIAccessKeyError struct {
	// A message about why the key creation failed.
	Message string `json:"message,omitempty"`
	// The type of the key.
	Type APIAccessKeyType `json:"type,omitempty"`
}

func (x *APIAccessKeyError) ImplementsAPIAccessKeyError() {}

// APIAccessKeySearchQuery - Parameters by which to filter the search.
type APIAccessKeySearchQuery struct {
	// Criteria by which to narrow the scope of keys to be returned.
	Scope APIAccessKeySearchScope `json:"scope,omitempty"`
	// A list of key types to be included in the search. If no types are provided, all types will be returned by default.
	Types []APIAccessKeyType `json:"types"`
}

// APIAccessKeySearchResult - A list of all keys scoped to the current actor.
type APIAccessKeySearchResult struct {
	// The total number of keys found in scope, irrespective of pagination.
	Count int `json:"count,omitempty"`
	// A list of all keys scoped to the current actor.
	Keys []APIAccessKeyInterface `json:"keys,omitempty"`
	// The next cursor, used for pagination. If a cursor is present, it means more keys can be fetched.
	NextCursor string `json:"nextCursor,omitempty"`
}

// special
func (x *APIAccessKeySearchResult) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	for k, v := range objMap {
		if v == nil {
			continue
		}

		switch k {
		case "count":
			err = json.Unmarshal(*v, &x.Count)
			if err != nil {
				return err
			}
		case "keys":
			if v == nil {
				continue
			}
			var rawMessageKeys []*json.RawMessage
			err = json.Unmarshal(*v, &rawMessageKeys)
			if err != nil {
				return err
			}

			for _, m := range rawMessageKeys {
				xxx, err := UnmarshalAPIAccessKeyInterface(*m)
				if err != nil {
					return err
				}

				if xxx != nil {
					x.Keys = append(x.Keys, *xxx)
				}
			}
		case "nextCursor":
			err = json.Unmarshal(*v, &x.NextCursor)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// APIAccessKeySearchScope - The scope of keys to be returned. Note that some filters only apply to certain key types.
type APIAccessKeySearchScope struct {
	// A list of key account IDs.
	AccountIDs []int `json:"accountIds,omitempty"`
	// The ingest type of the key. Only applies to ingest keys, and does not affect user key filtering.
	IngestTypes []APIAccessIngestKeyType `json:"ingestTypes,omitempty"`
	// A list of key user ids. Only applies to user keys, and does not affect ingest key filtering.
	UserIDs []int `json:"userIds,omitempty"`
}

// APIAccessUpdateIngestKeyInput - The `id` and data to update one or more keys.
type APIAccessUpdateIngestKeyInput struct {
	// The `id` of the key you want to update.
	KeyID string `json:"keyId"`
	// The name you want to assign to the key.
	Name string `json:"name,omitempty"`
	// The notes you want to assign to the key.
	Notes string `json:"notes,omitempty"`
}

// APIAccessUpdateInput - The `id` and data to update one or more keys.
type APIAccessUpdateInput struct {
	// A list of the configurations of each ingest key you want to update.
	Ingest []APIAccessUpdateIngestKeyInput `json:"ingest,omitempty"`
	// A list of the configurations of each user key you want to update.
	User []APIAccessUpdateUserKeyInput `json:"user,omitempty"`
}

// APIAccessUpdateKeyResponse - The response of the update keys mutation.
type APIAccessUpdateKeyResponse struct {
	// Lists all errors for keys that could not be updated. Each error maps to a single key input.
	Errors []APIAccessKeyErrorInterface `json:"errors,omitempty"`
	// Lists all successfully updated keys.
	UpdatedKeys []APIAccessKeyInterface `json:"updatedKeys,omitempty"`
}

// special
func (x *APIAccessUpdateKeyResponse) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	for k, v := range objMap {
		if v == nil {
			continue
		}

		switch k {
		case "errors":
			if v == nil {
				continue
			}
			var rawMessageErrors []*json.RawMessage
			err = json.Unmarshal(*v, &rawMessageErrors)
			if err != nil {
				return err
			}

			for _, m := range rawMessageErrors {
				xxx, err := UnmarshalAPIAccessKeyErrorInterface(*m)
				if err != nil {
					return err
				}

				if xxx != nil {
					x.Errors = append(x.Errors, *xxx)
				}
			}
		case "updatedKeys":
			if v == nil {
				continue
			}
			var rawMessageUpdatedKeys []*json.RawMessage
			err = json.Unmarshal(*v, &rawMessageUpdatedKeys)
			if err != nil {
				return err
			}

			for _, m := range rawMessageUpdatedKeys {
				xxx, err := UnmarshalAPIAccessKeyInterface(*m)
				if err != nil {
					return err
				}

				if xxx != nil {
					x.UpdatedKeys = append(x.UpdatedKeys, *xxx)
				}
			}
		}
	}

	return nil
}

// APIAccessUpdateUserKeyInput - The `id` and data to update one or more keys.
type APIAccessUpdateUserKeyInput struct {
	// The `id` of the key you want to update.
	KeyID string `json:"keyId"`
	// The name you want to assign to the key.
	Name string `json:"name,omitempty"`
	// The notes you want to assign to the key.
	Notes string `json:"notes,omitempty"`
}

// APIAccessUserKey - A user key.
type APIAccessUserKey struct {
	// The account this key is in.
	Account accounts.AccountReference `json:"account,omitempty"`
	// The account ID of the key.
	AccountID int `json:"accountId,omitempty"`
	// The UNIX epoch when the key was created, in seconds.
	CreatedAt EpochSeconds `json:"createdAt,omitempty"`
	// The ID of the user key. This can be used to identify a key without revealing the key itself (used to update and delete).
	ID string `json:"id,omitempty"`
	// The keystring of the key.
	Key string `json:"key,omitempty"`
	// The name of the key.
	Name string `json:"name,omitempty"`
	// Any notes can be attached to a key.
	Notes string `json:"notes,omitempty"`
	// The type of key, indicating what New Relic APIs it can be used to access.
	Type APIAccessKeyType `json:"type,omitempty"`
	// The user this key belongs to.
	User users.UserReference `json:"user,omitempty"`
	// The user ID of the key.
	UserID int `json:"userId,omitempty"`
}

func (x *APIAccessUserKey) ImplementsAPIAccessKey() {}

// APIAccessUserKeyError - A user key error. Each error maps to a single key input.
type APIAccessUserKeyError struct {
	// The account ID of the key.
	AccountID int `json:"accountId,omitempty"`
	// The error type of the error.
	ErrorType APIAccessUserKeyErrorType `json:"errorType,omitempty"`
	// The `id` of the key being updated.
	ID string `json:"id,omitempty"`
	// A message about why the key creation failed.
	Message string `json:"message,omitempty"`
	// The type of the key.
	Type APIAccessKeyType `json:"type,omitempty"`
	// The user ID of the key.
	UserID int `json:"userId,omitempty"`
}

func (x *APIAccessUserKeyError) ImplementsAPIAccessKeyError() {}

// EpochSeconds - The `EpochSeconds` scalar represents the number of seconds since the Unix epoch
type EpochSeconds int

// APIAccessKey - A key for accessing New Relic APIs.
type APIAccessKeyInterface interface {
	ImplementsAPIAccessKey()
}

// UnmarshalAPIAccessKeyInterface unmarshals the interface into the correct type
// based on __typename provided by GraphQL
func UnmarshalAPIAccessKeyInterface(b []byte) (*APIAccessKeyInterface, error) {
	var err error

	var rawMessageAPIAccessKey map[string]*json.RawMessage
	err = json.Unmarshal(b, &rawMessageAPIAccessKey)
	if err != nil {
		return nil, err
	}

	// Nothing to unmarshal
	if len(rawMessageAPIAccessKey) < 1 {
		return nil, nil
	}

	var typeName string

	if rawTypeName, ok := rawMessageAPIAccessKey["__typename"]; ok {
		err = json.Unmarshal(*rawTypeName, &typeName)
		if err != nil {
			return nil, err
		}

		switch typeName {
		case "ApiAccessIngestKey":
			var interfaceType APIAccessIngestKey
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx APIAccessKeyInterface = &interfaceType

			return &xxx, nil
		case "ApiAccessUserKey":
			var interfaceType APIAccessUserKey
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx APIAccessKeyInterface = &interfaceType

			return &xxx, nil
		}
	} else {
		keys := []string{}
		for k := range rawMessageAPIAccessKey {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("interface APIAccessKey did not include a __typename field for inspection: %s", keys)
	}

	return nil, fmt.Errorf("interface APIAccessKey was not matched against all PossibleTypes: %s", typeName)
}

// APIAccessKeyError - A key error. Each error maps to a single key input.
type APIAccessKeyErrorInterface interface {
	ImplementsAPIAccessKeyError()
	GetError() error
}

type APIAccessKeyErrorResponse struct {
	// The message with the error cause.
	Message string `json:"message,omitempty"`
	// Type of error.
	Type               string                      `json:"type,omitempty"`
	UserKeyErrorType   APIAccessUserKeyErrorType   `json:"userErrorType,omitempty"`
	IngestKeyErrorType APIAccessIngestKeyErrorType `json:"ingestErrorType,omitempty"`
	Id                 string                      `json:"id,omitempty"`
}

// UnmarshalAPIAccessKeyErrorInterface unmarshals the interface into the correct type
// based on __typename provided by GraphQL
func UnmarshalAPIAccessKeyErrorInterface(b []byte) (*APIAccessKeyErrorInterface, error) {
	var err error

	var rawMessageAPIAccessKeyError map[string]*json.RawMessage
	err = json.Unmarshal(b, &rawMessageAPIAccessKeyError)
	if err != nil {
		return nil, err
	}

	// Nothing to unmarshal
	if len(rawMessageAPIAccessKeyError) < 1 {
		return nil, nil
	}

	var typeName string

	if rawTypeName, ok := rawMessageAPIAccessKeyError["__typename"]; ok {
		err = json.Unmarshal(*rawTypeName, &typeName)
		if err != nil {
			return nil, err
		}

		switch typeName {
		case "ApiAccessIngestKeyError":
			var interfaceType APIAccessIngestKeyError
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx APIAccessKeyErrorInterface = &interfaceType

			return &xxx, nil
		case "ApiAccessUserKeyError":
			var interfaceType APIAccessUserKeyError
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx APIAccessKeyErrorInterface = &interfaceType

			return &xxx, nil
		}
	} else {
		keys := []string{}
		for k := range rawMessageAPIAccessKeyError {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("interface APIAccessKeyError did not include a __typename field for inspection: %s", keys)
	}

	return nil, fmt.Errorf("interface APIAccessKeyError was not matched against all PossibleTypes: %s", typeName)
}
