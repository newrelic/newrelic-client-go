//go:build unit
// +build unit

package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testSearchCompoundConditionsResponseJSON = `{
		"data": {
			"actor": {
				"account": {
					"alerts": {
						"compoundConditions": {
							"totalCount": 1,
							"nextCursor": null,
							"items": [
								{
									"id": "456",
									"name": "Compound Condition Test",
									"enabled": true,
									"policyId": "333",
									"facetMatchingBehavior": "FACETS_MATCH",
									"runbookUrl": "https://example.com/runbook",
									"thresholdDuration": 60,
									"triggerExpression": "a and b",
									"componentConditions": [
										{
											"id": "1",
											"alias": "a"
										}
									]
								}
							]
						}
					}
				}
			}
		}
	}`

	testCompoundConditionCreateResponseJSON = `{
		"data": {
			"alertsCompoundConditionCreate": {
				"id": "456",
				"name": "Compound Condition Test",
				"enabled": true,
				"policyId": "333",
				"facetMatchingBehavior": "FACETS_MATCH",
				"runbookUrl": "https://example.com/runbook",
				"thresholdDuration": 60,
				"triggerExpression": "a and b",
				"componentConditions": [
					{
						"id": "1",
						"alias": "a"
					}
				]
			}
		}
	}`

	testCompoundConditionUpdateResponseJSON = `{
		"data": {
			"alertsCompoundConditionUpdate": {
				"id": "456",
				"name": "Updated Compound Condition",
				"enabled": false,
				"policyId": "444",
				"facetMatchingBehavior": "FACETS_IGNORED",
				"runbookUrl": "https://example.com/updated-runbook",
				"thresholdDuration": 60,
				"triggerExpression": "a or b",
				"componentConditions": [
					{
						"id": "1",
						"alias": "a"
					},
					{
						"id": "2",
						"alias": "b"
					}
				]
			}
		}
	}`

	testCompoundConditionDeleteResponseJSON = `{
		"data": {
			"alertsCompoundConditionDelete": {
				"id": "456"
			}
		}
	}`
)

func TestSearchCompoundConditions(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testSearchCompoundConditionsResponseJSON, http.StatusOK)

	eq := "456"
	filter := &AlertsCompoundConditionFilterInput{
		Id: &AlertsCompoundConditionIDFilter{
			Eq: &eq,
		},
	}
	sort := []AlertsCompoundConditionSortInput{
		{
			Key:       string(AlertsCompoundConditionSortKeyTypes.NAME),
			Direction: string(AlertsCompoundConditionSortDirectionTypes.ASCENDING),
		},
	}

	expected := []*CompoundCondition{
		{
			ID:                    "456",
			Name:                  "Compound Condition Test",
			Enabled:               true,
			PolicyID:              "333",
			FacetMatchingBehavior: string(AlertsFacetMatchingBehaviorTypes.FACETS_MATCH),
			RunbookURL:            "https://example.com/runbook",
			ThresholdDuration:     60,
			TriggerExpression:     "a and b",
			ComponentConditions: []ComponentCondition{
				{
					ID:    "1",
					Alias: "a",
				},
			},
		},
	}

	actual, err := alerts.SearchCompoundConditions(123456, filter, sort, nil)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateCompoundCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testCompoundConditionCreateResponseJSON, http.StatusCreated)

	condition := CompoundConditionCreateInput{
		Name:                  "Compound Condition Test",
		Enabled:               true,
		FacetMatchingBehavior: stringPtr(string(AlertsFacetMatchingBehaviorTypes.FACETS_MATCH)),
		RunbookURL:            stringPtr("https://example.com/runbook"),
		ThresholdDuration:     intPtr(60),
		TriggerExpression:     "a and b",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    "1",
				Alias: "a",
			},
		},
	}

	expected := &CompoundCondition{
		ID:                    "456",
		Name:                  "Compound Condition Test",
		Enabled:               true,
		PolicyID:              "333",
		FacetMatchingBehavior: string(AlertsFacetMatchingBehaviorTypes.FACETS_MATCH),
		RunbookURL:            "https://example.com/runbook",
		ThresholdDuration:     60,
		TriggerExpression:     "a and b",
		ComponentConditions: []ComponentCondition{
			{
				ID:    "1",
				Alias: "a",
			},
		},
	}

	actual, err := alerts.CreateCompoundCondition(123456, "333", condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestUpdateCompoundCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testCompoundConditionUpdateResponseJSON, http.StatusOK)

	condition := CompoundConditionUpdateInput{
		Name:                  "Updated Compound Condition",
		Enabled:               boolPtr(false),
		PolicyID:              stringPtr("444"),
		FacetMatchingBehavior: stringPtr(string(AlertsFacetMatchingBehaviorTypes.FACETS_IGNORED)),
		RunbookURL:            stringPtr("https://example.com/updated-runbook"),
		ThresholdDuration:     intPtr(60),
		TriggerExpression:     "a or b",
		ComponentConditions: []ComponentConditionInput{
			{
				ID:    "1",
				Alias: "a",
			},
			{
				ID:    "2",
				Alias: "b",
			},
		},
	}

	expected := &CompoundCondition{
		ID:                    "456",
		Name:                  "Updated Compound Condition",
		Enabled:               false,
		PolicyID:              "444",
		FacetMatchingBehavior: string(AlertsFacetMatchingBehaviorTypes.FACETS_IGNORED),
		RunbookURL:            "https://example.com/updated-runbook",
		ThresholdDuration:     60,
		TriggerExpression:     "a or b",
		ComponentConditions: []ComponentCondition{
			{
				ID:    "1",
				Alias: "a",
			},
			{
				ID:    "2",
				Alias: "b",
			},
		},
	}

	actual, err := alerts.UpdateCompoundCondition(123456, "456", condition)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteCompoundCondition(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testCompoundConditionDeleteResponseJSON, http.StatusOK)

	expected := "456"

	actual, err := alerts.DeleteCompoundCondition(123456, "456")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
