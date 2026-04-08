//go:build unit
// +build unit

package notifications

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
)

// --- JSON fixtures for scope-based tests ---

var (
	testCreateDestinationWithScopeResponseJSON = `{
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

	testGetDestinationsWithScopeResponseJSON = `{
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
										"key": "email",
										"value": "test@newrelic.com"
									}
								],
								"status": "DEFAULT",
								"type": "EMAIL",
								"updatedAt": "2022-07-10T11:10:43.123715Z",
								"updatedBy": 1547846,
								"scope": {
									"id": "1",
									"type": "ACCOUNT"
								}
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

	testGetDestinationsWithOrgScopeResponseJSON = `{
		"actor": {
			"organization": {
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
										"key": "email",
										"value": "test@newrelic.com"
									}
								],
								"status": "DEFAULT",
								"type": "EMAIL",
								"updatedAt": "2022-07-10T11:10:43.123715Z",
								"updatedBy": 1547846,
								"scope": {
									"id": "org-123",
									"type": "ORGANIZATION"
								}
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

	testUpdateDestinationWithScopeResponseJSON = `{
		"aiNotificationsUpdateDestination": {
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
				"name": "test-notification-destination-updated",
				"properties": [
					{
						"displayValue": null,
						"key": "email",
						"label": null,
						"value": "updated@newrelic.com"
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

	testDeleteDestinationWithScopeResponseJSON = `{
		"aiNotificationsDeleteDestination": {
			"error": null,
			"errors": [],
			"ids": [
				"7463c367-6d61-416b-9aac-47f4a285fe5a"
			]
		}
	}`
)

// --- Create Destination With Scope Tests ---

func TestCreateDestinationWithScope_OrgScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	destinationInput := AiNotificationsDestinationInput{
		Type: AiNotificationsDestinationTypeTypes.EMAIL,
		Name: "test-notification-destination-1",
		Properties: []AiNotificationsPropertyInput{
			{Key: "email", Value: "test@newrelic.com"},
		},
		Auth: &AiNotificationsCredentialsInput{
			Basic: AiNotificationsBasicAuthInput{User: user, Password: "Pass"},
			Type:  AiNotificationsAuthTypeTypes.BASIC,
		},
	}

	scope := &EntityScopeInput{
		ID:   "org-123",
		Type: EntityScopeTypeInputTypes.ORGANIZATION,
	}

	actual, err := notifications.AiNotificationsCreateDestinationWithScope(accountId, destinationInput, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, id, actual.Destination.ID)
}

func TestCreateDestinationWithScope_AccountScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	destinationInput := AiNotificationsDestinationInput{
		Type: AiNotificationsDestinationTypeTypes.EMAIL,
		Name: "test-notification-destination-1",
		Properties: []AiNotificationsPropertyInput{
			{Key: "email", Value: "test@newrelic.com"},
		},
		Auth: &AiNotificationsCredentialsInput{
			Basic: AiNotificationsBasicAuthInput{User: user, Password: "Pass"},
			Type:  AiNotificationsAuthTypeTypes.BASIC,
		},
	}

	scope := &EntityScopeInput{
		ID:   "12345",
		Type: EntityScopeTypeInputTypes.ACCOUNT,
	}

	actual, err := notifications.AiNotificationsCreateDestinationWithScope(accountId, destinationInput, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, id, actual.Destination.ID)
}

func TestCreateDestinationWithScope_NilScopePanics(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	destinationInput := AiNotificationsDestinationInput{
		Type: AiNotificationsDestinationTypeTypes.EMAIL,
		Name: "test-notification-destination-1",
		Properties: []AiNotificationsPropertyInput{
			{Key: "email", Value: "test@newrelic.com"},
		},
	}

	// nil scope causes panic because scope.ID is dereferenced before nil check
	assert.Panics(t, func() {
		notifications.AiNotificationsCreateDestinationWithScopeWithContext(
			context.Background(), accountId, destinationInput, nil,
		)
	})
}

// --- GetDestinationsWithScope Tests ---

func TestGetDestinationsWithScope_AccountScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{ID: id}
	sorter := AiNotificationsDestinationSorter{}
	scope := &EntityScopeInput{
		ID:   "1",
		Type: EntityScopeTypeInputTypes.ACCOUNT,
	}

	actual, err := notifications.GetDestinationsWithScope(context.Background(), "", filters, sorter, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 1, actual.TotalCount)
	assert.Equal(t, 1, len(actual.Entities))
	assert.Equal(t, id, actual.Entities[0].ID)
}

func TestGetDestinationsWithScope_OrgScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithOrgScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{ID: id}
	sorter := AiNotificationsDestinationSorter{}
	scope := &EntityScopeInput{
		ID:   "org-123",
		Type: EntityScopeTypeInputTypes.ORGANIZATION,
	}

	actual, err := notifications.GetDestinationsWithScope(context.Background(), "", filters, sorter, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 1, actual.TotalCount)
	assert.Equal(t, 1, len(actual.Entities))
	assert.Equal(t, id, actual.Entities[0].ID)
}

func TestGetDestinationsWithScope_InvalidAccountID(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{ID: id}
	sorter := AiNotificationsDestinationSorter{}
	scope := &EntityScopeInput{
		ID:   "not-a-number",
		Type: EntityScopeTypeInputTypes.ACCOUNT,
	}

	actual, err := notifications.GetDestinationsWithScope(context.Background(), "", filters, sorter, scope)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

// --- GetDestinationsWithAccountScopeWithContext Tests ---

func TestGetDestinationsWithAccountScopeWithContext(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{ID: id}
	sorter := AiNotificationsDestinationSorter{}

	actual, err := notifications.GetDestinationsWithAccountScopeWithContext(
		context.Background(), accountId, "", filters, sorter,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 1, actual.TotalCount)
	assert.Equal(t, id, actual.Entities[0].ID)
}

func TestGetDestinationsWithAccountScopeWithContext_WithCursor(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{Name: name}
	sorter := AiNotificationsDestinationSorter{}

	actual, err := notifications.GetDestinationsWithAccountScopeWithContext(
		context.Background(), accountId, "some-cursor", filters, sorter,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

// --- GetDestinationsWithOrganizationScopeWithContext Tests ---

func TestGetDestinationsWithOrganizationScopeWithContext(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithOrgScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{ID: id}

	actual, err := notifications.GetDestinationsWithOrganizationScopeWithContext(
		context.Background(), "", filters,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 1, actual.TotalCount)
	assert.Equal(t, id, actual.Entities[0].ID)
}

func TestGetDestinationsWithOrganizationScopeWithContext_WithCursor(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithOrgScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{Name: name}

	actual, err := notifications.GetDestinationsWithOrganizationScopeWithContext(
		context.Background(), "org-cursor", filters,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

// --- Update Destination With Scope Tests ---

func TestUpdateDestinationWithScope_OrgScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	destinationUpdate := AiNotificationsDestinationUpdate{
		Name:   "test-notification-destination-updated",
		Active: true,
	}
	scope := &EntityScopeInput{
		ID:   "org-123",
		Type: EntityScopeTypeInputTypes.ORGANIZATION,
	}

	actual, err := notifications.AiNotificationsUpdateDestinationWithScope(accountId, destinationUpdate, id, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, "test-notification-destination-updated", actual.Destination.Name)
}

func TestUpdateDestinationWithScope_AccountScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	destinationUpdate := AiNotificationsDestinationUpdate{
		Name:   "test-notification-destination-updated",
		Active: true,
	}
	scope := &EntityScopeInput{
		ID:   "12345",
		Type: EntityScopeTypeInputTypes.ACCOUNT,
	}

	actual, err := notifications.AiNotificationsUpdateDestinationWithScope(accountId, destinationUpdate, id, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, "test-notification-destination-updated", actual.Destination.Name)
}

func TestUpdateDestinationWithScope_NilScopePanics(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	destinationUpdate := AiNotificationsDestinationUpdate{
		Name:   "test-notification-destination-updated",
		Active: true,
	}

	// nil scope causes panic because scope.ID is dereferenced before nil check
	assert.Panics(t, func() {
		notifications.AiNotificationsUpdateDestinationWithScopeWithContext(
			context.Background(), accountId, destinationUpdate, id, nil,
		)
	})
}

// --- Delete Destination With Scope Tests ---

func TestDeleteDestinationWithScope_OrgScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	scope := &EntityScopeInput{
		ID:   "org-123",
		Type: EntityScopeTypeInputTypes.ORGANIZATION,
	}

	actual, err := notifications.AiNotificationsDeleteDestinationWithScope(accountId, id, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, []string{id}, actual.IDs)
}

func TestDeleteDestinationWithScope_AccountScope(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	scope := &EntityScopeInput{
		ID:   "12345",
		Type: EntityScopeTypeInputTypes.ACCOUNT,
	}

	actual, err := notifications.AiNotificationsDeleteDestinationWithScope(accountId, id, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, []string{id}, actual.IDs)
}

func TestDeleteDestinationWithScope_NilScopePanics(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	// nil scope causes panic because scope.ID is dereferenced before nil check
	assert.Panics(t, func() {
		notifications.AiNotificationsDeleteDestinationWithScopeWithContext(
			context.Background(), accountId, id, nil,
		)
	})
}

// --- Response struct deserialization tests ---

func TestDestinationsWithOrgScopeResponseDeserialization(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithOrgScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{}

	actual, err := notifications.GetDestinationsWithOrganizationScopeWithContext(
		context.Background(), "", filters,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 1, actual.TotalCount)
	entity := actual.Entities[0]
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, "org-123", entity.Scope.ID)
	assert.Equal(t, EntityScopeTypeInput("ORGANIZATION"), entity.Scope.Type)
}

func TestDestinationsWithAccountScopeResponseDeserialization(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationsWithScopeResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	filters := ai.AiNotificationsDestinationFilter{}
	sorter := AiNotificationsDestinationSorter{}

	actual, err := notifications.GetDestinationsWithAccountScopeWithContext(
		context.Background(), accountId, "", filters, sorter,
	)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 1, actual.TotalCount)
	entity := actual.Entities[0]
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, "1", entity.Scope.ID)
	assert.Equal(t, EntityScopeTypeInput("ACCOUNT"), entity.Scope.Type)
}
