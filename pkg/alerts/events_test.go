//go:build unit
// +build unit

package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListAlertEventsOneObjectResponseJSON = `{
		"recent_events": [
			{
				"id": 123,
				"event_type": "event",
				"product": "product",
				"entity_type": "entity",
				"entity_group_id": 12345,
				"entity_id": 12345,
				"priority": "priority",
				"description": "description",
				"timestamp": 1575438237690,
				"incident_id": 12345
			}
		]
	}`
	testListAlertEventsMultipleObjectsResponseJSON = `{
		"recent_events": [
			{
				"id": 123,
				"event_type": "event",
				"product": "product",
				"entity_type": "entity",
				"entity_group_id": 12345,
				"entity_id": 12345,
				"priority": "priority",
				"description": "description",
				"timestamp": 1575438237690,
				"incident_id": 12345
			},{
				"id":28012408580,
				"event_type":"INCIDENT_OPEN",
				"description":"Some Alert Policy (Your Application Name violated Goroutines min count)",
				"timestamp":1575438237690,
				"incident_id":2217914781
      		}
		]
	}`
)

var testListAlertEventsCases = []struct {
	name       string
	body       string
	params     *ListAlertEventsParams
	expected   []*AlertEvent
	httpStatus int
	errStr     string
}{
	{
		name: "success one object",
		body: testListAlertEventsOneObjectResponseJSON,
		params: &ListAlertEventsParams{
			Product:       "test",
			EntityType:    "test",
			EntityGroupID: 12345,
			EntityID:      12345,
			EventType:     "test",
			IncidentID:    12345,
			Page:          1,
		},
		expected: []*AlertEvent{
			{
				ID:            123,
				EventType:     "event",
				Product:       "product",
				EntityType:    "entity",
				EntityGroupID: 12345,
				EntityID:      12345,
				Priority:      "priority",
				Description:   "description",
				Timestamp:     &testTimestamp,
				IncidentID:    12345,
			},
		},
		httpStatus: http.StatusOK,
		errStr:     "",
	},
	{
		name: "success multiple objects",
		body: testListAlertEventsMultipleObjectsResponseJSON,
		params: &ListAlertEventsParams{
			Product:       "test",
			EntityType:    "test",
			EntityGroupID: 12345,
			EntityID:      12345,
			EventType:     "test",
			IncidentID:    12345,
			Page:          1,
		},
		expected: []*AlertEvent{
			{
				ID:            123,
				EventType:     "event",
				Product:       "product",
				EntityType:    "entity",
				EntityGroupID: 12345,
				EntityID:      12345,
				Priority:      "priority",
				Description:   "description",
				Timestamp:     &testTimestamp,
				IncidentID:    12345,
			},
			{
				ID:          28012408580,
				EventType:   "INCIDENT_OPEN",
				Description: "Some Alert Policy (Your Application Name violated Goroutines min count)",
				Timestamp:   &testTimestamp,
				IncidentID:  2217914781,
			},
		},
		httpStatus: http.StatusOK,
		errStr:     "",
	},
	{
		name:       "success no records retrieved",
		body:       "{}",
		params:     nil,
		expected:   []*AlertEvent{},
		httpStatus: http.StatusOK,
		errStr:     "",
	},
	{
		name:       "failed unmarshal error",
		body:       "",
		params:     nil,
		expected:   nil,
		httpStatus: http.StatusOK,
		errStr:     "unexpected end of JSON input",
	},
}

func TestListAlertEvents(t *testing.T) {
	t.Parallel()
	for _, tc := range testListAlertEventsCases {
		t.Run(tc.name, func(t *testing.T) {
			alerts := newMockResponse(t, tc.body, tc.httpStatus)

			actual, err := alerts.ListAlertEvents(tc.params)

			if tc.errStr == "" {
				assert.NoError(t, err)
				if len(tc.expected) <= 1 {
					assert.Equal(t, tc.expected, actual)
				} else {
					for i := range tc.expected {
						assert.Equal(t, tc.expected[i], actual[i])
					}
				}
			} else {
				assert.Equal(t, tc.errStr, err.Error())
				assert.Empty(t, actual)
			}
		})
	}
}
