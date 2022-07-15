//go:build unit
// +build unit

package notifications

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/nrtime"
)

var (
	accountID           = 10867072
	testTimestampString = "2022-07-11T08:43:38.943355Z"
	testTimestamp       = nrtime.DateTime(testTimestampString)
	ID                  = "28932668-a894-41f9-bdd2-f7c094d164df"
	destinationId       = "b1e90a32-23b7-4028-b2c7-ffbdfe103852"

	testCreateChannelResponseJSON = `{
		"aiNotificationsCreateChannel": {
		  "channel": {
			"accountId": 10867072,
			"active": true,
			"createdAt": "2022-07-11T08:43:38.943355Z",
			"destinationId": "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
			"id": "28932668-a894-41f9-bdd2-f7c094d164df",
			"name": "test-notification-channel-1",
			"product": "IINT",
			"properties": [
			  {
				"displayValue": "",
				"key": "payload",
				"label": "Payload Template",
				"value": "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}"
			  }
			],
			"status": "DEFAULT",
			"type": "WEBHOOK",
			"updatedAt": "2022-07-11T08:43:38.943355Z",
			"updatedBy": 1547846
		  },
		  "error": null,
		  "errors": []
		}
	}`

	testUpdateChannelResponseJSON = `{
		"aiNotificationsUpdateChannel": {
		  "channel": {
			"accountId": 10867072,
			"active": false,
			"createdAt": "2022-07-11T08:43:38.943355Z",
			"destinationId": "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
			"id": "28932668-a894-41f9-bdd2-f7c094d164df",
			"name": "test-notification-channel-1-update",
			"product": "IINT",
			"properties": [
			  {
				"displayValue": "",
				"key": "payload",
				"label": "Payload Template",
				"value": "{\\n\\t\\\"id\\\": \\\"test-update\\\"\\n}"
			  }
			],
			"status": "DEFAULT",
			"type": "WEBHOOK",
			"updatedAt": "2022-07-11T08:43:38.943355Z",
			"updatedBy": 1547846
		  },
		  "error": null,
		  "errors": []
		}
	}`

	testDeleteChannelResponseJSON = `{
		"aiNotificationsDeleteChannel": {
		  "error": null,
		  "errors": [],
		  "ids": [
			"28932668-a894-41f9-bdd2-f7c094d164df"
		  ]
		}
	}`

	testGetChannelResponseJSON = `{
		"actor": {
		  "account": {
			"aiNotifications": {
			  "channels": {
				"entities": [
				  {
					"accountId": 10867072,
					"active": true,
					"createdAt": "2022-07-11T08:43:38.943355Z",
					"destinationId": "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
					"id": "28932668-a894-41f9-bdd2-f7c094d164df",
					"name": "test-notification-channel-1",
					"product": "IINT",
					"properties": [
					  {
						"displayValue": "",
						"key": "payload",
						"label": "Payload Template",
						"value": "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}"
					  }
					],
					"status": "DEFAULT",
					"type": "WEBHOOK",
					"updatedAt": "2022-07-11T08:43:38.943355Z",
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

func TestCreateChannel(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateChannelResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	channelInput := AiNotificationsChannelInput{
		DestinationId: destinationId,
		Name:          "test-notification-channel-1",
		Product:       AiNotificationsProductTypes.IINT,
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "payload",
				Label: "Payload Template",
				Value: "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
			},
		},
		Type: AiNotificationsChannelTypeTypes.WEBHOOK,
	}

	expected := &AiNotificationsChannelResponse{
		Channel: AiNotificationsChannel{
			AccountID:     accountID,
			Active:        true,
			CreatedAt:     testTimestamp,
			DestinationId: destinationId,
			ID:            ID,
			Name:          "test-notification-channel-1",
			Product:       AiNotificationsProductTypes.IINT,
			Properties: []AiNotificationsProperty{
				{
					DisplayValue: "",
					Key:          "payload",
					Label:        "Payload Template",
					Value:        "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
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

	actual, err := notifications.AiNotificationsCreateChannel(accountID, channelInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetChannel(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetChannelResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	expected := &AiNotificationsChannelsResponse{
		Entities: []AiNotificationsChannel{
			{
				AccountID:     accountID,
				Active:        true,
				CreatedAt:     testTimestamp,
				DestinationId: destinationId,
				ID:            ID,
				Name:          "test-notification-channel-1",
				Product:       AiNotificationsProductTypes.IINT,
				Properties: []AiNotificationsProperty{
					{
						DisplayValue: "",
						Key:          "payload",
						Label:        "Payload Template",
						Value:        "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
					},
				},
				Status:    AiNotificationsChannelStatusTypes.DEFAULT,
				Type:      AiNotificationsChannelTypeTypes.WEBHOOK,
				UpdatedAt: testTimestamp,
				UpdatedBy: 1547846,
			},
		},
		Error:      AiNotificationsResponseError{},
		Errors:     []AiNotificationsResponseError{},
		NextCursor: "",
		TotalCount: 1,
	}

	filters := ai.AiNotificationsChannelFilter{
		ID: ID,
	}
	sorter := AiNotificationsChannelSorter{}

	actual, err := notifications.GetChannels(accountID, "", filters, sorter)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateChannel(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateChannelResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusCreated)

	updateChannelInput := AiNotificationsChannelUpdate{
		Name: "test-notification-channel-1-update",
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "payload",
				Label: "Payload Template",
				Value: "{\\n\\t\\\"id\\\": \\\"test-update\\\"\\n}",
			},
		},
		Active: false,
	}

	expected := &AiNotificationsChannelResponse{
		Channel: AiNotificationsChannel{
			AccountID:     accountID,
			Active:        false,
			CreatedAt:     testTimestamp,
			DestinationId: destinationId,
			ID:            ID,
			Name:          "test-notification-channel-1-update",
			Product:       AiNotificationsProductTypes.IINT,
			Properties: []AiNotificationsProperty{
				{
					DisplayValue: "",
					Key:          "payload",
					Label:        "Payload Template",
					Value:        "{\\n\\t\\\"id\\\": \\\"test-update\\\"\\n}",
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

	actual, err := notifications.AiNotificationsUpdateChannel(accountID, updateChannelInput, ID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteChannel(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteChannelResponseJSON)
	notifications := newMockResponse(t, respJSON, http.StatusOK)

	expected := &AiNotificationsDeleteResponse{
		IDs:    []string{ID},
		Errors: []AiNotificationsResponseError{},
		Error:  AiNotificationsResponseError{},
	}

	actual, err := notifications.AiNotificationsDeleteChannel(accountID, ID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
