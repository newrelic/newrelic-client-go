//go:build unit
// +build unit

package notifications

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
)

var (
	accountId           = 1
	testTimestampString = "2022-07-07T16:51:29.469355Z"
	testTimestamp, err1 = time.Parse(time.RFC3339, testTimestampString)
	id                  = "742e6f13-2de1-4627-a6dc-0e62b2caa1cf"
	destinationId       = "4a43c21f-427d-4a44-9061-29c455d8df3b"

	testCreateChannelResponseJSON = `{
		"aiNotificationsCreateChannel": {
		  "channel": {
			"accountId": 1,
			"active": true,
			"createdAt": "2022-07-07T16:51:29.469355Z",
			"destinationId": "4a43c21f-427d-4a44-9061-29c455d8df3b",
			"id": "742e6f13-2de1-4627-a6dc-0e62b2caa1cf",
			"name": "test-notification-channel-1",
			"product": "IINT",
			"properties": [
			  {
				"displayValue": "",
				"key": "payload",
				"label": "Payload Template",
				"value": "{\\n\\t\\\"id\\\": \\\"blabla\\\"\\n}"
			  },
			  {
				"displayValue": "",
				"key": "headers",
				"label": "Custom headers",
				"value": "{}"
			  }
			],
			"status": "DEFAULT",
			"type": "WEBHOOK",
			"updatedAt": "2022-07-07T16:51:29.469355Z",
			"updatedBy": 1547846
		  },
		  "error": null,
		  "errors": []
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

func TestCreateChannel(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateChannelResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	channelInput := AiNotificationsChannelInput{
		DestinationId: "4a43c21f-427d-4a44-9061-29c455d8df3b",
		Name:          "test-notification-channel-1",
		Product:       AiNotificationsProductTypes.IINT,
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "payload",
				Label: "Payload Template",
				Value: "{\\n\\t\\\"id\\\": \\\"blabla\\\"\\n}",
			},
			{
				Key:   "headers",
				Label: "Custom headers",
				Value: "{}",
			},
		},
		Type: AiNotificationsChannelTypeTypes.WEBHOOK,
	}

	expected := &AiNotificationsChannelResponse{
		Channel: AiNotificationsChannel{
			AccountId:     accountId,
			Active:        true,
			CreatedAt:     testTimestamp,
			DestinationId: destinationId,
			ID:            id,
			Name:          "test-notification-channel-1",
			Product:       AiNotificationsProductTypes.IINT,
			Properties: []AiNotificationsProperty{
				{
					DisplayValue: "",
					Key:          "headers",
					Label:        "Custom headers",
					Value:        "{}",
				},
				{
					DisplayValue: "",
					Key:          "payload",
					Label:        "Payload Template",
					Value:        "{\\n\\t\\\"id\\\": \\\"blabla\\\"\\n}",
				},
			},
			Status:    AiNotificationsChannelStatusTypes.DEFAULT,
			Type:      AiNotificationsChannelTypeTypes.WEBHOOK,
			UpdatedAt: testTimestamp,
			UpdatedBy: 1547846,
		},
		Error:  ai.AiNotificationsError{},
		Errors: []ai.AiNotificationsError{},
	}

	actual, err := notifications.AiNotificationsCreateChannel(accountId, channelInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

//func TestGetChannel(t *testing.T) {
//	t.Parallel()
//	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetDestinationResponseJSON)
//	notifications := newMockResponse(t, respJSON, http.StatusOK)
//
//	expected := &Destination{
//		ID:        UUID(id),
//		Name:      "test-notification-destination-1",
//		Type:      DestinationTypes.Email,
//		CreatedAt: testTimestamp,
//		UpdatedAt: testTimestamp,
//		UpdatedBy: 1547846,
//		AccountId: 1,
//		Status:    DestinationStatuses.Default,
//		Active:    true,
//		LastSent:  testTimestamp,
//		Auth: Auth{
//			AuthType: &AuthTypes.Basic,
//			User:     &user,
//		},
//		Properties: []Property{
//			{
//				DisplayValue: "",
//				Key:          "email",
//				Label:        "",
//				Value:        "dshemesh@newrelic.com",
//			},
//		},
//	}
//
//	actual, err := notifications.GetDestination(accountId, UUID(id))
//
//	assert.NoError(t, err)
//	assert.NotNil(t, actual)
//	assert.Equal(t, expected, actual)
//}
//
//func TestUpdateChannel(t *testing.T) {
//	t.Parallel()
//	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationResponseJSON)
//	notifications := newMockResponse(t, respJSON, http.StatusOK)
//
//	expected := []string{id}
//
//	actual, err := notifications.DeleteDestinationMutation(accountId, UUID(id))
//
//	assert.NoError(t, err)
//	assert.NotNil(t, actual)
//	assert.Equal(t, expected, actual)
//}
//
//func TestDeleteChannel(t *testing.T) {
//	t.Parallel()
//	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteDestinationResponseJSON)
//	notifications := newMockResponse(t, respJSON, http.StatusOK)
//
//	expected := []string{id}
//
//	actual, err := notifications.DeleteDestinationMutation(accountId, UUID(id))
//
//	assert.NoError(t, err)
//	assert.NotNil(t, actual)
//	assert.Equal(t, expected, actual)
//}
