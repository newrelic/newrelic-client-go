package apiaccesskeys

import (
	"errors"
	"fmt"
	"github.com/newrelic/newrelic-client-go/internal/http"
)

// APIAccessKey represents a New Relic API access ingest or user key.
type APIAccessKey struct {
	ID         string `json:"id,omitempty"`
	Key        string `json:"key,omitempty"`
	Name       string `json:"name,omitempty"`
	Notes      string `json:"notes,omitempty"`
	Type       string `json:"type,omitempty"`
	AccountID  int    `json:"accountId,omitempty"`
	IngestType string `json:"ingestType,omitempty"`
	UserID     string `json:"userId,omitempty"`
}

// apiAccessKeyCreateResponse represents the JSON response returned from creating key(s).
type apiAccessKeyCreateResponse struct {
	APIAccessCreateKeys struct {
		CreatedKeys []APIAccessKey                    `json:"createdKeys,omitempty"`
		Errors      []apiAccessKeyMutationErrResponse `json:"errors,omitempty"`
	} `json:"apiAccessCreateKeys"`
}

// apiAccessKeyUpdateResponse represents the JSON response returned from updating key(s).
type apiAccessKeyUpdateResponse struct {
	APIAccessUpdateKeys struct {
		UpdatedKeys []APIAccessKey                    `json:"updatedKeys,omitempty"`
		Errors      []apiAccessKeyMutationErrResponse `json:"errors,omitempty"`
	} `json:"apiAccessUpdateKeys"`
}

// apiAccessKeyGetResponse represents the JSON response returned from getting an access key.
type apiAccessKeyGetResponse struct {
	Actor struct {
		APIAccess struct {
			Key APIAccessKey `json:"key,omitempty"`
		} `json:"apiAccess"`
	} `json:"actor"`
	http.GraphQLErrorResponse
}

// apiAccessKeyMutationErrResponse represents the generic error response returned from modifying API access keys.
type apiAccessKeyMutationErrResponse struct {
	AccountID       int    `json:"accountId,omitempty"`
	ID              string `json:"id,omitempty"`
	IngestErrorType string `json:"ingestErrorType,omitempty"`
	UserErrorType   string `json:"userErrorType,omitempty"`
	IngestType      string `json:"ingestType,omitempty"`
	Message         string `json:"message,omitempty"`
	Type            string `json:"type,omitempty"`
}

// apiAccessKeyDeleteResponse represents the JSON response returned from creating key(s).
type apiAccessKeyDeleteResponse struct {
	APIAccessDeleteKeys struct {
		DeletedKey []struct {
			ID string `json:"id"`
		} `json:"deletedKeys,omitempty"`
		Errors []apiAccessKeyMutationErrResponse `json:"errors,omitempty"`
	} `json:"apiAccessDeleteKeys"`
}

const (
	graphqlAPIAccessKeyBaseFields = `
		id
		key
		name
		notes
		type
		... on ApiAccessIngestKey {
			id
			name
			accountId
			ingestType
			key
			notes
			type
		}
		... on ApiAccessUserKey {
			id
			name
			accountId
			key
			notes
			type
			userId
		}
		... on ApiAccessKey {
			id
			name
			key
			notes
			type
		}`

	graphqlAPIAccessCreateKeyFields = `createdKeys {` + graphqlAPIAccessKeyBaseFields + `}`

	graphqlAPIAccessUpdatedKeyFields = `updatedKeys {` + graphqlAPIAccessKeyBaseFields + `}`

	graphqlAPIAccessKeyErrorFields = `errors {
		  message
		  type
		  ... on ApiAccessIngestKeyError {
			id
			ingestErrorType: errorType
			accountId
			ingestType
			message
			type
		  }
		  ... on ApiAccessKeyError {
			message
			type
		  }
		  ... on ApiAccessUserKeyError {
			id
			accountId
			userErrorType: errorType
			message
			type
			userId
		  }
		}
	`

	apiAccessKeyCreateKeys = `mutation($keys: ApiAccessCreateInput!) {
			apiAccessCreateKeys(keys: $keys) {` + graphqlAPIAccessCreateKeyFields + graphqlAPIAccessKeyErrorFields + `
		}}`

	apiAccessKeyGetKey = `query($id: ID!, $keyType: ApiAccessKeyType!) {
		actor {
			apiAccess {
				key(id: $id, keyType: $keyType) {` + graphqlAPIAccessKeyBaseFields + `}}}}`

	apiAccessKeyUpdateKeys = `mutation($keys: ApiAccessUpdateInput!) {
			apiAccessUpdateKeys(keys: $keys) {` + graphqlAPIAccessUpdatedKeyFields + graphqlAPIAccessKeyErrorFields + `
		}}`

	apiAccessKeyDeleteKeys = `mutation($keys: ApiAccessDeleteInput!) {
			apiAccessDeleteKeys(keys: $keys) {
				deletedKeys {
					id
				}` + graphqlAPIAccessKeyErrorFields + `}}`
)

// CreateAPIAccessKeysMutation create keys. You can create keys for multiple accounts at once.
func (a *APIAccessKeys) CreateAPIAccessKeysMutation(keys APIAccessCreateKeysInput) ([]APIAccessKey, error) {
	vars := map[string]interface{}{
		"keys": keys,
	}

	resp := apiAccessKeyCreateResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyCreateKeys, vars, &resp); err != nil {
		return nil, err
	}

	if resp.APIAccessCreateKeys.Errors != nil {
		return nil, errors.New(formatAPIAccessKeyMutationErrors(resp.APIAccessCreateKeys.Errors))
	}

	return resp.APIAccessCreateKeys.CreatedKeys, nil
}

// GetAPIAccessKeyMutation returns a single API access key.
func (a *APIAccessKeys) GetAPIAccessKeyMutation(key APIAccessGetInput) (*APIAccessKey, error) {
	vars := map[string]interface{}{
		"id":      key.ID,
		"keyType": key.KeyType,
	}

	resp := apiAccessKeyGetResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyGetKey, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Errors != nil {
		return nil, errors.New(resp.Error())
	}

	return &resp.Actor.APIAccess.Key, nil
}

// UpdateAPIAccessKeyMutation updates keys. You can update keys for multiple accounts at once.
func (a *APIAccessKeys) UpdateAPIAccessKeyMutation(keys APIAccessUpdateInput) ([]APIAccessKey, error) {
	vars := map[string]interface{}{
		"keys": keys,
	}

	resp := apiAccessKeyUpdateResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyUpdateKeys, vars, &resp); err != nil {
		return nil, err
	}

	if resp.APIAccessUpdateKeys.Errors != nil {
		return nil, errors.New(formatAPIAccessKeyMutationErrors(resp.APIAccessUpdateKeys.Errors))
	}

	return resp.APIAccessUpdateKeys.UpdatedKeys, nil
}

// DeleteAPIAccessKeyMutation deletes keys.
func (a *APIAccessKeys) DeleteAPIAccessKeyMutation(keys APIAccessDeleteInput) error {
	vars := map[string]interface{}{
		"keys": keys,
	}

	resp := apiAccessKeyDeleteResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyDeleteKeys, vars, &resp); err != nil {
		return err
	}

	if resp.APIAccessDeleteKeys.Errors != nil {
		return errors.New(formatAPIAccessKeyMutationErrors(resp.APIAccessDeleteKeys.Errors))
	}

	return nil
}

func formatAPIAccessKeyMutationErrors(errors []apiAccessKeyMutationErrResponse) string {
	errorString := ""
	for _, e := range errors {
		errorString += fmt.Sprintf("%v\n", e)
	}
	return errorString
}
