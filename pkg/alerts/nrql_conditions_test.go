//go:build unit
// +build unit

package alerts

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testListNrqlConditionsResponseJSON = `{
		"nrql_conditions": [
			{
				"type": "static",
				"id": 12345,
				"name": "NRQL Test Alert",
				"enabled": true,
				"violation_time_limit_seconds": 3600,
				"terms": [
					{
						"duration": "5",
						"operator": "above",
						"priority": "critical",
						"threshold": "1",
						"time_function": "all"
					}
				],
				"nrql": {
					"query": "SELECT count(*) FROM Transactions",
					"since_value": "3"
				},
				"entity_guid": "NDAwQzA0rEFST1BTfENPTkRJVElPTnw3MzUzNjL"
			}
		]
	}`

	testNrqlConditionJSON = `{
		"nrql_condition": {
			"type": "static",
			"id": 12345,
			"name": "NRQL Test Alert",
			"enabled": true,
			"violation_time_limit_seconds": 3600,
			"terms": [
				{
					"duration": "5",
					"operator": "above",
					"priority": "critical",
					"threshold": "1",
					"time_function": "all"
				}
			],
			"nrql": {
				"query": "SELECT count(*) FROM Transactions",
				"since_value": "3"
			},
			"entity_guid": "NDAwQzA0rEFST1BTfENPTkRJVElPTnw3MzUzNjL"
		}
	}`

	testNrqlConditionCreateJSON = `{
		"nrql_condition": {
			"type": "static",
			"id": 12345,
			"name": "NRQL Test Alert",
			"enabled": true,
			"violation_time_limit_seconds": 3600,
			"terms": [
				{
					"duration": "5",
					"operator": "above",
					"priority": "critical",
					"threshold": "1",
					"time_function": "all"
				}
			],
			"nrql": {
				"query": "SELECT count(*) FROM Transactions",
				"since_value": "3"
			}
		}
	}`

	testNrqlConditionUpdatedJSON = `{
		"nrql_condition": {
			"type": "static",
			"id": 12345,
			"name": "NRQL Test Alert Updated",
			"enabled": false,
			"violation_time_limit_seconds": 3600,
			"terms": [
				{
					"duration": "5",
					"operator": "below",
					"priority": "critical",
					"threshold": "1",
					"time_function": "all"
				}
			],
			"nrql": {
				"query": "SELECT count(*) FROM Transactions",
				"since_value": "3"
			},
			"runbook_url": "https://www.example.com/docs"
		}
	}`

	testNrqlConditionEntityGUID = common.EntityGUID("NDAwQzA0rEFST1BTfENPTkRJVElPTnw3MzUzNjL")
)

func TestListNrqlConditions(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListNrqlConditionsResponseJSON, http.StatusOK)

	expected := []*NrqlCondition{
		{
			Nrql: NrqlQuery{
				Query:      "SELECT count(*) FROM Transactions",
				SinceValue: "3",
			},
			Terms: []ConditionTerm{
				{
					Duration:     5,
					Operator:     "above",
					Priority:     "critical",
					Threshold:    1,
					TimeFunction: "all",
				},
			},
			Type:                "static",
			Name:                "NRQL Test Alert",
			RunbookURL:          "",
			ID:                  12345,
			ViolationCloseTimer: 3600,
			Enabled:             true,
			EntityGUID:          &testNrqlConditionEntityGUID,
		},
	}

	actual, err := alerts.ListNrqlConditions(123)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetNrqlCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListNrqlConditionsResponseJSON, http.StatusOK)

	expected := &NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
		EntityGUID:          &testNrqlConditionEntityGUID,
	}

	actual, err := alerts.GetNrqlCondition(123, 12345)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateNrqlCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testNrqlConditionCreateJSON, http.StatusCreated)
	policyID := 333

	condition := NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
	}

	expected := &condition

	actual, err := alerts.CreateNrqlCondition(policyID, condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateNrqlCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testNrqlConditionUpdatedJSON, http.StatusCreated)

	condition := NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
		EntityGUID:          &testNrqlConditionEntityGUID,
	}

	expected := &NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "below",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert Updated",
		RunbookURL:          "https://www.example.com/docs",
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             false,
	}

	actual, err := alerts.UpdateNrqlCondition(condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteNrqlCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testNrqlConditionJSON, http.StatusOK)
	expected := &NrqlCondition{
		Nrql: NrqlQuery{
			Query:      "SELECT count(*) FROM Transactions",
			SinceValue: "3",
		},
		Terms: []ConditionTerm{
			{
				Duration:     5,
				Operator:     "above",
				Priority:     "critical",
				Threshold:    1,
				TimeFunction: "all",
			},
		},
		Type:                "static",
		Name:                "NRQL Test Alert",
		RunbookURL:          "",
		ID:                  12345,
		ViolationCloseTimer: 3600,
		Enabled:             true,
		EntityGUID:          &testNrqlConditionEntityGUID,
	}

	actual, err := alerts.DeleteNrqlCondition(12345)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateNrqlConditionStaticMutation_TermsWireSemantics(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		terms       *[]NrqlConditionTerm
		expectKey   bool
		expectValue interface{}
	}{
		"nil terms omits the field from the request": {
			terms:     nil,
			expectKey: false,
		},
		"pointer to an empty slice sends an explicit empty list": {
			terms:       &[]NrqlConditionTerm{},
			expectKey:   true,
			expectValue: []interface{}{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var capturedVariables map[string]interface{}

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body struct {
					Variables map[string]interface{} `json:"variables"`
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				capturedVariables = body.Variables

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":{"alertsNrqlConditionStaticUpdate":{"id":"123"}}}`))
			})

			alertsClient := newTestClient(t, handler)

			input := NrqlConditionUpdateInput{
				NrqlConditionUpdateBase: NrqlConditionUpdateBase{
					Name:  "test",
					Terms: tc.terms,
				},
			}

			_, err := alertsClient.UpdateNrqlConditionStaticMutation(123456, "123", input)
			require.NoError(t, err)

			condition, ok := capturedVariables["condition"].(map[string]interface{})
			require.True(t, ok)

			value, hasKey := condition["terms"]
			assert.Equal(t, tc.expectKey, hasKey)
			if tc.expectKey {
				assert.Equal(t, tc.expectValue, value)
			}
		})
	}
}
