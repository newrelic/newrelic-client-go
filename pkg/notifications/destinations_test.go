//go:build unit
// +build unit

package notifications

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	accountId           = 1
	testTimestampString = "2021-07-08T12:30:00-07:00"
	testTimestamp, err1 = time.Parse(time.RFC3339, "2021-07-08T12:30:00-07:00")
	user                = "test-user"
	id                  = "6f820700-6a3b-4c0f-9bc7-2c322a1455e6"

	testCreateDestinationResponseJSON = `{
		"AiNotificationsCreateDestination": {
			"destination": {
				"id": "6f820700-6a3b-4c0f-9bc7-2c322a1455e6",
				"name": "test-notification-destination-1",
				"createdAt": "2021-07-08T12:30:00-07:00",
				"updatedAt": "2021-07-08T12:30:00-07:00",
				"accountId": 1,
				"active": true,
				"auth": {
				  "authType": "BASIC",
				  "user": "test-user"
				},
				"lastSent": "2021-07-08T12:30:00-07:00",
				"properties": [{
						"displayValue": "",
						"key": "email",
						"label": "",
						"value": "dshemesh@newrelic.com"
				}],
				"status": "default",
				"type": "EMAIL",
				"updatedBy": 1547846
			}
		}
	}`

	testDeleteDestinationResponseJSON = `{
		"aiNotificationsDeleteDestination": {
			"ids": [
				"6f820700-6a3b-4c0f-9bc7-2c322a1455e6"
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
								"id": "6f820700-6a3b-4c0f-9bc7-2c322a1455e6",
								"name": "test-notification-destination-1",
								"createdAt": "2021-07-08T12:30:00-07:00",
								"updatedAt": "2021-07-08T12:30:00-07:00",
								"accountId": 1,
								"active": true,
								"auth": {
								  "authType": "BASIC",
								  "user": "test-user"
								},
								"lastSent": "2021-07-08T12:30:00-07:00",
								"properties": [{
										"displayValue": "",
										"key": "email",
										"label": "",
										"value": "dshemesh@newrelic.com"
								}],
								"status": "default",
								"type": "EMAIL",
								"updatedBy": 1547846
							}
						]
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

	destinationInput := DestinationInput{
		Type: DestinationTypes.Email,
		Name: "test-notification-destination-1",
		Properties: []PropertyInput{
			{
				Key:   "email",
				Value: "dshemesh@newrelic.com",
			},
		},
		Auth: AiNotificationsCredentialsInput{
			Basic: BasicAuth{
				User:     user,
				Password: "Pass",
			},
			Type: AuthTypes.Basic,
		},
	}

	expected := &Destination{
		ID:        UUID(id),
		Name:      "test-notification-destination-1",
		Type:      DestinationTypes.Email,
		CreatedAt: testTimestamp,
		UpdatedAt: testTimestamp,
		UpdatedBy: 1547846,
		AccountId: 1,
		Status:    DestinationStatuses.Default,
		Active:    true,
		LastSent:  testTimestamp,
		Auth: Auth{
			AuthType: &AuthTypes.Basic,
			User:     &user,
		},
		Properties: []Property{
			{
				DisplayValue: "",
				Key:          "email",
				Label:        "",
				Value:        "dshemesh@newrelic.com",
			},
		},
	}

	actual, err := notifications.CreateDestinationMutation(accountId, destinationInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetChannel(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	expected := &Destination{
		ID:        UUID(id),
		Name:      "test-notification-destination-1",
		Type:      DestinationTypes.Email,
		CreatedAt: testTimestamp,
		UpdatedAt: testTimestamp,
		UpdatedBy: 1547846,
		AccountId: 1,
		Status:    DestinationStatuses.Default,
		Active:    true,
		LastSent:  testTimestamp,
		Auth: Auth{
			AuthType: &AuthTypes.Basic,
			User:     &user,
		},
		Properties: []Property{
			{
				DisplayValue: "",
				Key:          "email",
				Label:        "",
				Value:        "dshemesh@newrelic.com",
			},
		},
	}

	actual, err := notifications.GetDestination(accountId, UUID(id))

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteDestination(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	expected := []string{id}

	actual, err := notifications.DeleteDestinationMutation(accountId, UUID(id))

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
