// +build unit

package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListAlertEventsResponseJSON = `{
		"alert_events": [
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
)

func TestListAlertEvents(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListAlertEventsResponseJSON, http.StatusOK)

	expected := []*AlertEvent{
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
	}

	params := &ListAlertEventsParams{
		Product:       "test",
		EntityType:    "test",
		EntityGroupID: 12345,
		EntityID:      12345,
		EventType:     "test",
		IncidentID:    12345,
		Page:          1,
	}

	actual, err := alerts.ListAlertEvents(params)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
