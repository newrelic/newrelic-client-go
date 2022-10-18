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

	actual, err := notifications.AiNotificationsCreateDestination(accountId, destinationInput)

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
		Errors:     []AiNotificationsResponseError{},
		Error:      AiNotificationsResponseError{},
		NextCursor: "",
		TotalCount: 1,
	}

	filters := ai.AiNotificationsDestinationFilter{
		ID: id,
	}
	sorter := AiNotificationsDestinationSorter{}

	actual, err := notifications.GetDestinations(accountId, "", filters, sorter)

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
		Errors: []AiNotificationsResponseError{},
	}

	actual, err := notifications.AiNotificationsDeleteDestination(accountId, id)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
