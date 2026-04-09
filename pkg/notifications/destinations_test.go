//go:build unit
// +build unit

package notifications

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

var (
	timestampString = "2022-07-10T11:10:43.123715Z"
	timestamp       = nrtime.DateTime(timestampString)
	user            = "test-user"
	accountId       = 1
	id              = "7463c367-6d61-416b-9aac-47f4a285fe5a"
	guid            = "MXxBSU9QU3xERVNUSU5BVElPTnw3NDYzYzM2Ny02ZDYxLTQxNmItOWFhYy00N2Y0YTI4NWZlNWE="
	name            = "test-notification-destination-1"

	testCreateDestinationResponseJSON = `{
	 "aiNotificationsCreateDestination": {
		  "destination": {
			"accountId": 1,
			"active": true,
			"auth": {
			  "authType": "BASIC",
			  "user": "test-user"
			},
			"createdAt": "2022-07-10T11:10:43.123715Z",
			"id": "7463c367-6d61-416b-9aac-47f4a285fe5a",
			"guid": "MXxBSU9QU3xERVNUSU5BVElPTnw3NDYzYzM2Ny02ZDYxLTQxNmItOWFhYy00N2Y0YTI4NWZlNWE=",
			"isUserAuthenticated": false,
			"lastSent": "2022-07-10T11:10:43.123715Z",
			"name": "test-notification-destination-1",
			"properties": [
			  {
				"displayValue": null,
				"key": "email",
				"label": null,
				"value": "test@newrelic.com"
			  }
			],
			"status": "DEFAULT",
			"type": "EMAIL",
			"updatedAt": "2022-07-10T11:10:43.123715Z",
			"updatedBy": 1547846
		  },
		  "error": null,
		  "errors": []
		}
	}`

	testDeleteDestinationResponseJSON = `{
		"aiNotificationsDeleteDestination": {
		  "error": null,
		  "errors": [],
		  "ids": [
			"7463c367-6d61-416b-9aac-47f4a285fe5a"
		  ]
		}
	}`

	orgId = "org-123"

	testUpdateDestinationResponseJSON = `{
	 "aiNotificationsUpdateDestination": {
		  "destination": {
			"accountId": 1,
			"active": false,
			"auth": {
			  "authType": "BASIC",
			  "user": "test-user"
			},
			"createdAt": "2022-07-10T11:10:43.123715Z",
			"id": "7463c367-6d61-416b-9aac-47f4a285fe5a",
			"guid": "MXxBSU9QU3xERVNUSU5BVElPTnw3NDYzYzM2Ny02ZDYxLTQxNmItOWFhYy00N2Y0YTI4NWZlNWE=",
			"isUserAuthenticated": false,
			"lastSent": "2022-07-10T11:10:43.123715Z",
			"name": "test-notification-destination-1-updated",
			"properties": [
			  {
				"displayValue": null,
				"key": "email",
				"label": null,
				"value": "updated@newrelic.com"
			  }
			],
			"scope": {"id": "1", "type": "ACCOUNT"},
			"status": "DEFAULT",
			"type": "EMAIL",
			"updatedAt": "2022-07-10T11:10:43.123715Z",
			"updatedBy": 1547846
		  },
		  "error": null,
		  "errors": []
		}
	}`

	testUpdateDestinationOrgScopeResponseJSON = `{
	 "aiNotificationsUpdateDestination": {
		  "destination": {
			"accountId": 0,
			"active": false,
			"auth": {
			  "authType": "BASIC",
			  "user": "test-user"
			},
			"createdAt": "2022-07-10T11:10:43.123715Z",
			"id": "7463c367-6d61-416b-9aac-47f4a285fe5a",
			"guid": "MXxBSU9QU3xERVNUSU5BVElPTnw3NDYzYzM2Ny02ZDYxLTQxNmItOWFhYy00N2Y0YTI4NWZlNWE=",
			"isUserAuthenticated": false,
			"lastSent": "2022-07-10T11:10:43.123715Z",
			"name": "test-notification-destination-1-updated",
			"properties": [
			  {
				"displayValue": null,
				"key": "email",
				"label": null,
				"value": "updated@newrelic.com"
			  }
			],
			"scope": {"id": "org-123", "type": "ORGANIZATION"},
			"status": "DEFAULT",
			"type": "EMAIL",
			"updatedAt": "2022-07-10T11:10:43.123715Z",
			"updatedBy": 1547846
		  },
		  "error": null,
		  "errors": []
		}
	}`

	testCreateDestinationOrgScopeResponseJSON = `{
	 "aiNotificationsCreateDestination": {
		  "destination": {
			"accountId": 0,
			"active": true,
			"auth": {
			  "authType": "BASIC",
			  "user": "test-user"
			},
			"createdAt": "2022-07-10T11:10:43.123715Z",
			"id": "7463c367-6d61-416b-9aac-47f4a285fe5a",
			"guid": "MXxBSU9QU3xERVNUSU5BVElPTnw3NDYzYzM2Ny02ZDYxLTQxNmItOWFhYy00N2Y0YTI4NWZlNWE=",
			"isUserAuthenticated": false,
			"lastSent": "2022-07-10T11:10:43.123715Z",
			"name": "test-notification-destination-1",
			"properties": [
			  {
				"displayValue": null,
				"key": "email",
				"label": null,
				"value": "test@newrelic.com"
			  }
			],
			"scope": {"id": "org-123", "type": "ORGANIZATION"},
			"status": "DEFAULT",
			"type": "EMAIL",
			"updatedAt": "2022-07-10T11:10:43.123715Z",
			"updatedBy": 1547846
		  },
		  "error": null,
		  "errors": []
		}
	}`

	testGetDestinationResponseJSON = `{
    "actor": {
      "account": {
        "aiNotifications": {
          "destinations": {
            "entities": [
              {
                "accountId": 1,
                "active": true,
                "auth": {
                  "authType": "BASIC",
                  "user": "test-user"
                },
                "createdAt": "2022-07-10T11:10:43.123715Z",
                "id": "7463c367-6d61-416b-9aac-47f4a285fe5a",
                "guid": "MXxBSU9QU3xERVNUSU5BVElPTnw3NDYzYzM2Ny02ZDYxLTQxNmItOWFhYy00N2Y0YTI4NWZlNWE=",
                "isUserAuthenticated": false,
                "lastSent": "2022-07-10T11:10:43.123715Z",
                "name": "test-notification-destination-1",
                "properties": [
                  {
                    "displayValue": null,
                    "key": "email",
                    "label": null,
                    "value": "test@newrelic.com"
                  }
                ],
                "status": "DEFAULT",
                "type": "EMAIL",
                "updatedAt": "2022-07-10T11:10:43.123715Z",
                "updatedBy": 1547846
              }
            ],
            "error": null,
            "errors": [],
            "nextCursor": null,
            "totalCount": 1
			}
          }
        }
      }
	}`

	testGetDestinationOrgScopeResponseJSON = `{
    "actor": {
      "organization": {
        "aiNotifications": {
          "destinations": {
            "entities": [
              {
                "accountId": 0,
                "active": true,
                "auth": {
                  "authType": "BASIC",
                  "user": "test-user"
                },
                "createdAt": "2022-07-10T11:10:43.123715Z",
                "id": "7463c367-6d61-416b-9aac-47f4a285fe5a",
                "guid": "MXxBSU9QU3xERVNUSU5BVElPTnw3NDYzYzM2Ny02ZDYxLTQxNmItOWFhYy00N2Y0YTI4NWZlNWE=",
                "isUserAuthenticated": false,
                "lastSent": "2022-07-10T11:10:43.123715Z",
                "name": "test-notification-destination-1",
                "properties": [
                  {
                    "displayValue": null,
                    "key": "email",
                    "label": null,
                    "value": "test@newrelic.com"
                  }
                ],
                "scope": {"id": "org-123", "type": "ORGANIZATION"},
                "status": "DEFAULT",
                "type": "EMAIL",
                "updatedAt": "2022-07-10T11:10:43.123715Z",
                "updatedBy": 1547846
              }
            ],
            "error": null,
            "errors": [],
            "nextCursor": null,
            "totalCount": 1
          }
        }
      }
    }
	}`
)

func TestCreateDestination(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	destinationInput := AiNotificationsDestinationInput{
		Type: AiNotificationsDestinationTypeTypes.EMAIL,
		Name: "test-notification-destination-1",
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "email",
				Value: "test@newrelic.com",
			},
		},
		Auth: &AiNotificationsCredentialsInput{
			Basic: AiNotificationsBasicAuthInput{
				User:     user,
				Password: "Pass",
			},
			Type: AiNotificationsAuthTypeTypes.BASIC,
		},
	}

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	expected := &AiNotificationsDestinationResponse{
		Destination: AiNotificationsDestination{
			AccountID:           accountId,
			Active:              true,
			Auth:                auth,
			CreatedAt:           timestamp,
			ID:                  id,
			GUID:                EntityGUID(guid),
			IsUserAuthenticated: false,
			LastSent:            timestamp,
			Name:                "test-notification-destination-1",
			Properties: []AiNotificationsProperty{
				{
					DisplayValue: "",
					Key:          "email",
					Label:        "",
					Value:        "test@newrelic.com",
				},
			},
			Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
			Type:      AiNotificationsDestinationTypeTypes.EMAIL,
			UpdatedAt: timestamp,
			UpdatedBy: 1547846,
		},
		Errors: []ai.AiNotificationsError{},
	}

	actual, err := notifications.AiNotificationsCreateDestination(accountId, destinationInput, &AiNotificationsEntityScopeInput{})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetDestinations(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	expected := &AiNotificationsDestinationsResponse{
		Entities: []AiNotificationsDestination{
			{
				AccountID:           accountId,
				Active:              true,
				Auth:                auth,
				CreatedAt:           timestamp,
				ID:                  id,
				GUID:                EntityGUID(guid),
				IsUserAuthenticated: false,
				LastSent:            timestamp,
				Name:                "test-notification-destination-1",
				Properties: []AiNotificationsProperty{
					{
						DisplayValue: "",
						Key:          "email",
						Label:        "",
						Value:        "test@newrelic.com",
					},
				},
				Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
				Type:      AiNotificationsDestinationTypeTypes.EMAIL,
				UpdatedAt: timestamp,
				UpdatedBy: 1547846,
			},
		},
		Errors:     []ai.AiNotificationsResponseError{},
		Error:      ai.AiNotificationsResponseError{},
		NextCursor: "",
		TotalCount: 1,
	}

	filters := ai.AiNotificationsDestinationFilter{
		ID: id,
	}
	sorter := AiNotificationsDestinationSorter{}

	actual, err := notifications.GetDestinationsAccount(accountId, "", filters, sorter)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetDestinationsByName(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	expected := &AiNotificationsDestinationsResponse{
		Entities: []AiNotificationsDestination{
			{
				AccountID:           accountId,
				Active:              true,
				Auth:                auth,
				CreatedAt:           timestamp,
				ID:                  id,
				GUID:                EntityGUID(guid),
				IsUserAuthenticated: false,
				LastSent:            timestamp,
				Name:                "test-notification-destination-1",
				Properties: []AiNotificationsProperty{
					{
						DisplayValue: "",
						Key:          "email",
						Label:        "",
						Value:        "test@newrelic.com",
					},
				},
				Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
				Type:      AiNotificationsDestinationTypeTypes.EMAIL,
				UpdatedAt: timestamp,
				UpdatedBy: 1547846,
			},
		},
		Errors:     []ai.AiNotificationsResponseError{},
		Error:      ai.AiNotificationsResponseError{},
		NextCursor: "",
		TotalCount: 1,
	}

	filters := ai.AiNotificationsDestinationFilter{
		Name: name,
	}
	sorter := AiNotificationsDestinationSorter{}

	actual, err := notifications.GetDestinationsAccount(accountId, "", filters, sorter)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetDestinationsByExactName(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	expected := &AiNotificationsDestinationsResponse{
		Entities: []AiNotificationsDestination{
			{
				AccountID:           accountId,
				Active:              true,
				Auth:                auth,
				CreatedAt:           timestamp,
				ID:                  id,
				GUID:                EntityGUID(guid),
				IsUserAuthenticated: false,
				LastSent:            timestamp,
				Name:                "test-notification-destination-1",
				Properties: []AiNotificationsProperty{
					{
						DisplayValue: "",
						Key:          "email",
						Label:        "",
						Value:        "test@newrelic.com",
					},
				},
				Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
				Type:      AiNotificationsDestinationTypeTypes.EMAIL,
				UpdatedAt: timestamp,
				UpdatedBy: 1547846,
			},
		},
		Errors:     []ai.AiNotificationsResponseError{},
		Error:      ai.AiNotificationsResponseError{},
		NextCursor: "",
		TotalCount: 1,
	}

	filters := ai.AiNotificationsDestinationFilter{
		ExactName: name,
	}
	sorter := AiNotificationsDestinationSorter{}

	actual, err := notifications.GetDestinationsAccount(accountId, "", filters, sorter)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateDestination(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	updateInput := AiNotificationsDestinationUpdate{
		Active: false,
		Name:   "test-notification-destination-1-updated",
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "email",
				Value: "updated@newrelic.com",
			},
		},
		Auth: &AiNotificationsCredentialsInput{
			Basic: AiNotificationsBasicAuthInput{
				User:     user,
				Password: "Pass",
			},
			Type: AiNotificationsAuthTypeTypes.BASIC,
		},
	}

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	expected := &AiNotificationsDestinationResponse{
		Destination: AiNotificationsDestination{
			AccountID:           accountId,
			Active:              false,
			Auth:                auth,
			CreatedAt:           timestamp,
			ID:                  id,
			GUID:                EntityGUID(guid),
			IsUserAuthenticated: false,
			LastSent:            timestamp,
			Name:                "test-notification-destination-1-updated",
			Properties: []AiNotificationsProperty{
				{
					DisplayValue: "",
					Key:          "email",
					Label:        "",
					Value:        "updated@newrelic.com",
				},
			},
			Scope:     AiNotificationsEntityScope{ID: "1", Type: AiNotificationsEntityScopeTypeTypes.ACCOUNT},
			Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
			Type:      AiNotificationsDestinationTypeTypes.EMAIL,
			UpdatedAt: timestamp,
			UpdatedBy: 1547846,
		},
		Errors: []ai.AiNotificationsError{},
	}

	// Account-scoped update: pass accountId, empty scope
	actual, err := notifications.AiNotificationsUpdateDestination(accountId, updateInput, id, &AiNotificationsEntityScopeInput{})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateDestinationWithOrgScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateDestinationOrgScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	updateInput := AiNotificationsDestinationUpdate{
		Active: false,
		Name:   "test-notification-destination-1-updated",
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "email",
				Value: "updated@newrelic.com",
			},
		},
		Auth: &AiNotificationsCredentialsInput{
			Basic: AiNotificationsBasicAuthInput{
				User:     user,
				Password: "Pass",
			},
			Type: AiNotificationsAuthTypeTypes.BASIC,
		},
	}

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	orgScope := AiNotificationsEntityScopeInput{
		ID:   orgId,
		Type: AiNotificationsEntityScopeTypeInputTypes.ORGANIZATION,
	}

	expected := &AiNotificationsDestinationResponse{
		Destination: AiNotificationsDestination{
			AccountID:           0,
			Active:              false,
			Auth:                auth,
			CreatedAt:           timestamp,
			ID:                  id,
			GUID:                EntityGUID(guid),
			IsUserAuthenticated: false,
			LastSent:            timestamp,
			Name:                "test-notification-destination-1-updated",
			Properties: []AiNotificationsProperty{
				{
					DisplayValue: "",
					Key:          "email",
					Label:        "",
					Value:        "updated@newrelic.com",
				},
			},
			Scope:     AiNotificationsEntityScope{ID: orgId, Type: AiNotificationsEntityScopeTypeTypes.ORGANIZATION},
			Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
			Type:      AiNotificationsDestinationTypeTypes.EMAIL,
			UpdatedAt: timestamp,
			UpdatedBy: 1547846,
		},
		Errors: []ai.AiNotificationsError{},
	}

	// Org-scoped update: accountId is 0, scope carries the org ID
	actual, err := notifications.AiNotificationsUpdateDestination(0, updateInput, id, orgScope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateDestinationWithOrgScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateDestinationOrgScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	destinationInput := AiNotificationsDestinationInput{
		Type: AiNotificationsDestinationTypeTypes.EMAIL,
		Name: "test-notification-destination-1",
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "email",
				Value: "test@newrelic.com",
			},
		},
		Auth: &AiNotificationsCredentialsInput{
			Basic: AiNotificationsBasicAuthInput{
				User:     user,
				Password: "Pass",
			},
			Type: AiNotificationsAuthTypeTypes.BASIC,
		},
	}

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	// Org-scoped: accountId is 0, scope carries the org ID
	orgScope := AiNotificationsEntityScopeInput{
		ID:   orgId,
		Type: AiNotificationsEntityScopeTypeInputTypes.ORGANIZATION,
	}

	expected := &AiNotificationsDestinationResponse{
		Destination: AiNotificationsDestination{
			AccountID:           0,
			Active:              true,
			Auth:                auth,
			CreatedAt:           timestamp,
			ID:                  id,
			GUID:                EntityGUID(guid),
			IsUserAuthenticated: false,
			LastSent:            timestamp,
			Name:                "test-notification-destination-1",
			Properties: []AiNotificationsProperty{
				{
					DisplayValue: "",
					Key:          "email",
					Label:        "",
					Value:        "test@newrelic.com",
				},
			},
			Scope:     AiNotificationsEntityScope{ID: orgId, Type: AiNotificationsEntityScopeTypeTypes.ORGANIZATION},
			Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
			Type:      AiNotificationsDestinationTypeTypes.EMAIL,
			UpdatedAt: timestamp,
			UpdatedBy: 1547846,
		},
		Errors: []ai.AiNotificationsError{},
	}

	actual, err := notifications.AiNotificationsCreateDestination(0, destinationInput, orgScope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetDestinationsOrganization(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationOrgScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	auth.ImplementsAiNotificationsAuth()

	expected := &AiNotificationsDestinationsResponse{
		Entities: []AiNotificationsDestination{
			{
				AccountID:           0,
				Active:              true,
				Auth:                auth,
				CreatedAt:           timestamp,
				ID:                  id,
				GUID:                EntityGUID(guid),
				IsUserAuthenticated: false,
				LastSent:            timestamp,
				Name:                "test-notification-destination-1",
				Properties: []AiNotificationsProperty{
					{
						DisplayValue: "",
						Key:          "email",
						Label:        "",
						Value:        "test@newrelic.com",
					},
				},
				Scope:     AiNotificationsEntityScope{ID: orgId, Type: AiNotificationsEntityScopeTypeTypes.ORGANIZATION},
				Status:    AiNotificationsDestinationStatusTypes.DEFAULT,
				Type:      AiNotificationsDestinationTypeTypes.EMAIL,
				UpdatedAt: timestamp,
				UpdatedBy: 1547846,
			},
		},
		Errors:     []ai.AiNotificationsResponseError{},
		Error:      ai.AiNotificationsResponseError{},
		NextCursor: "",
		TotalCount: 1,
	}

	filters := ai.AiNotificationsDestinationFilter{
		ID: id,
	}
	sorter := AiNotificationsDestinationSorter{}

	// Org-scoped query: no accountId, uses organization path
	actual, err := notifications.GetDestinationsOrganization("", filters, sorter)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteDestination(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	expected := &AiNotificationsDeleteResponse{
		IDs:    []string{id},
		Errors: []ai.AiNotificationsResponseError{},
	}

	// Account-scoped delete: pass accountId, empty scope
	actual, err := notifications.AiNotificationsDeleteDestination(accountId, id, &AiNotificationsEntityScopeInput{})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteDestinationWithOrgScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	expected := &AiNotificationsDeleteResponse{
		IDs:    []string{id},
		Errors: []ai.AiNotificationsResponseError{},
	}

	// Org-scoped delete: accountId is 0, scope carries the org ID
	orgScope := AiNotificationsEntityScopeInput{
		ID:   orgId,
		Type: AiNotificationsEntityScopeTypeInputTypes.ORGANIZATION,
	}

	actual, err := notifications.AiNotificationsDeleteDestination(0, id, orgScope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
