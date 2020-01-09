// +build unit

package alerts

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testPluginsConditionJSON = `{
		"id": 333444,
		"name": "Connected Clients (High)",
		"enabled": true,
		"entities": [
			"111222"
		],
		"metric_description": "Connected Clients",
		"metric": "Component/Connection/Clients[connections]",
		"value_function": "average",
		"runbook_url": "https://example.com/runbook",
		"terms": [
			{
				"duration": "5",
				"operator": "above",
				"priority": "critical",
				"threshold": "10",
				"time_function": "all"
			}
		],
		"plugin": {
			"id": "222555",
			"guid": "net.example.newrelic_redis_plugin"
		}
	}`

	testPluginsCondition = PluginsCondition{
		PolicyID:          123,
		ID:                333444,
		Name:              "Connected Clients (High)",
		Enabled:           true,
		Entities:          []string{"111222"},
		Metric:            "Component/Connection/Clients[connections]",
		MetricDescription: "Connected Clients",
		RunbookURL:        "https://example.com/runbook",
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    10,
				TimeFunction: "all",
			},
		},
		ValueFunction: "average",
		Plugin: AlertPlugin{
			ID:   "222555",
			GUID: "net.example.newrelic_redis_plugin",
		},
	}
)

func TestListPluginsConditions(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"plugins_conditions": [%s]}`, testPluginsConditionJSON)
	client := newMockResponse(t, responseJSON, http.StatusOK)

	expected := []*PluginsCondition{
		&testPluginsCondition,
	}

	actual, err := client.ListPluginsConditions(123)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetPluginsCondition(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"plugins_conditions": [%s]}`, testPluginsConditionJSON)
	client := newMockResponse(t, responseJSON, http.StatusOK)

	actual, err := client.GetPluginsCondition(123, 333444)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testPluginsCondition, actual)
}

func TestCreatePluginsCondition(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"plugins_condition": %s}`, testPluginsConditionJSON)
	client := newMockResponse(t, responseJSON, http.StatusCreated)

	actual, err := client.CreatePluginsCondition(testPluginsCondition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testPluginsCondition, actual)
}

func TestUpdatePluginsCondition(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"plugins_condition": %s}`, testPluginsConditionJSON)
	client := newMockResponse(t, responseJSON, http.StatusOK)

	actual, err := client.UpdatePluginsCondition(testPluginsCondition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testPluginsCondition, actual)
}

func TestDeletePluginsCondition(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"plugins_condition": %s}`, testPluginsConditionJSON)
	client := newMockResponse(t, responseJSON, http.StatusOK)

	expected := PluginsCondition{
		ID:                333444,
		Name:              "Connected Clients (High)",
		Enabled:           true,
		Entities:          []string{"111222"},
		Metric:            "Component/Connection/Clients[connections]",
		MetricDescription: "Connected Clients",
		RunbookURL:        "https://example.com/runbook",
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    10,
				TimeFunction: "all",
			},
		},
		ValueFunction: "average",
		Plugin: AlertPlugin{
			ID:   "222555",
			GUID: "net.example.newrelic_redis_plugin",
		},
	}

	actual, err := client.DeletePluginsCondition(333444)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &expected, actual)
}
