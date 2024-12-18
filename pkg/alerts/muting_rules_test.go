//go:build unit
// +build unit

package alerts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testMutingRuleListResponseJSON = `{
		"actor": {
			"account": {
				"alerts": {
					"mutingRules": [
						{
							"id": "123",
							"accountId": 400304,
							"condition": {
								"conditions": [
									{
										"attribute": "conditionName",
										"operator": "EQUALS",
										"values": [
											"please not me"
										]
									}
								],
								"operator": "AND"
							},
							"createdAt": "2021-01-12T00:50:39.533Z",
							"createdByUser": {
								"email": "testemail@newrelic.com",
								"gravatar": "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
								"id": 1,
								"name": "Test User"
							},
							"description": null,
							"enabled": true,
							"name": "Test Muting Rule",
							"schedule": {
								"endRepeat": null,
								"endTime": "2021-07-08T14:30:00-07:00",
								"nextEndTime": "2021-07-08T14:30:00-07:00",
								"nextStartTime": "2021-07-08T12:30:00-07:00",
								"repeat": "DAILY",
								"repeatCount": 10,
								"startTime": "2021-07-08T12:30:00-07:00",
								"timeZone": "America/Los_Angeles",
								"weeklyRepeatDays": null
							},
							"actionOnMutingRuleWindowEnded": "CLOSE_ISSUES_ON_INACTIVE",
							"status": "INACTIVE",
							"updatedAt": "2021-01-12T00:50:39.533Z",
							"updatedByUser": {
								"email": "testemail@newrelic.com",
								"gravatar": "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
								"id": 1,
								"name": "Test User"
							}
						}
					]
				}
			}
		}
	}`

	testMutingRuleGetResponseJSON = `{
		"actor": {
			"account": {
				"alerts": {
					"mutingRule": {
						"id": "123",
						"accountId": 400304,
						"condition": {
							"conditions": [
								{
									"attribute": "conditionName",
									"operator": "EQUALS",
									"values": [
										"please not me"
									]
								}
							],
							"operator": "AND"
						},
						"createdAt": "2021-01-12T00:50:39.533Z",
						"createdByUser": {
							"email": "testemail@newrelic.com",
							"gravatar": "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
							"id": 1,
							"name": "Test User"
						},
						"description": null,
						"enabled": true,
						"name": "Test Muting Rule",
						"schedule": {
							"endRepeat": null,
							"endTime": "2021-07-08T14:30:00-07:00",
							"nextEndTime": "2021-07-08T14:30:00-07:00",
							"nextStartTime": "2021-07-08T12:30:00-07:00",
							"repeat": "DAILY",
							"repeatCount": 10,
							"startTime": "2021-07-08T12:30:00-07:00",
							"timeZone": "America/Los_Angeles",
							"weeklyRepeatDays": null
						},
						"status": "INACTIVE",
						"updatedAt": "2021-01-12T00:50:39.533Z",
						"updatedByUser": {
							"email": "testemail@newrelic.com",
							"gravatar": "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
							"id": 1,
							"name": "Test User"
						},
						"actionOnMutingRuleWindowEnded": "CLOSE_ISSUES_ON_INACTIVE"
					}
				}
			}
		}
	}`

	testMutingRuleCreateResponseJSON = `{
		"alertsMutingRuleCreate": {
			"id": "123",
			"accountId": 400304,
			"condition": {
				"conditions": [
					{
						"attribute": "conditionName",
						"operator": "EQUALS",
						"values": [
							"please not me"
						]
					}
				],
				"operator": "AND"
			},
			"createdAt": "2021-01-12T00:50:39.533Z",
			"createdByUser": {
				"email": "testemail@newrelic.com",
				"gravatar": "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
				"id": 1,
				"name": "Test User"
			},
			"description": null,
			"enabled": true,
			"name": "Test Muting Rule",
			"schedule": {
				"endRepeat": null,
				"endTime": "2021-07-08T14:30:00-07:00",
				"nextEndTime": "2021-07-08T14:30:00-07:00",
				"nextStartTime": "2021-07-08T12:30:00-07:00",
				"repeat": "DAILY",
				"repeatCount": 10,
				"startTime": "2021-07-08T12:30:00-07:00",
				"timeZone": "America/Los_Angeles",
				"weeklyRepeatDays": null
			},
			"status": "INACTIVE",
			"updatedAt": "2021-01-12T00:50:39.533Z",
			"updatedByUser": {
				"email": "testemail@newrelic.com",
				"gravatar": "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
				"id": 1,
				"name": "Test User"
			},
			"actionOnMutingRuleWindowEnded": "CLOSE_ISSUES_ON_INACTIVE"
		}
	}`

	testMutingRuleUpdateResponseJSON = `{
		"alertsMutingRuleUpdate": {
			"id": "123",
			"accountId": 400304,
			"schedule": {
				"endRepeat": "2021-08-08T12:30:00-07:00",
				"startTime": "2021-07-08T12:30:00-07:00",
				"repeat": null,
				"repeatCount": null,
				"weeklyRepeatDays": null
			},
			"actionOnMutingRuleWindowEnded": "DO_NOTHING"
		}
	}`

	testMutingRuleDeleteResponseJSON = `{
		"alertsMutingRuleDelete": {
			"id": "123"
		}
	}`
)

func TestListMutingRules(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testMutingRuleListResponseJSON)
	alerts := newMockResponse(t, respJSON, http.StatusOK)
	accountID := 400304
	startTime, err1 := time.Parse(time.RFC3339, "2021-07-08T12:30:00-07:00")
	if err1 != nil {
		t.Fatal(err1)
	}
	endTime, err2 := time.Parse(time.RFC3339, "2021-07-08T14:30:00-07:00")
	if err2 != nil {
		t.Fatal(err2)
	}
	repeatCount := 10

	expected := []MutingRule{
		{
			ID:        123,
			AccountID: accountID,
			Condition: MutingRuleConditionGroup{
				Conditions: []MutingRuleCondition{
					{
						Attribute: "conditionName",
						Operator:  "EQUALS",
						Values:    []string{"please not me"},
					},
				},
				Operator: "AND",
			},
			CreatedAt: "2021-01-12T00:50:39.533Z",
			CreatedByUser: ByUser{
				Email:    "testemail@newrelic.com",
				Gravatar: "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
				ID:       1,
				Name:     "Test User",
			},
			Description: "",
			Enabled:     true,
			Name:        "Test Muting Rule",
			Schedule: &MutingRuleSchedule{
				EndRepeat:        nil,
				EndTime:          &endTime,
				Repeat:           &MutingRuleScheduleRepeatTypes.DAILY,
				RepeatCount:      &repeatCount,
				StartTime:        &startTime,
				TimeZone:         "America/Los_Angeles",
				WeeklyRepeatDays: nil,
			},
			UpdatedAt: "2021-01-12T00:50:39.533Z",
			UpdatedByUser: ByUser{
				Email:    "testemail@newrelic.com",
				Gravatar: "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
				ID:       1,
				Name:     "Test User",
			},
			ActionOnMutingRuleWindowEnded: "CLOSE_ISSUES_ON_INACTIVE",
		},
	}

	actual, err := alerts.ListMutingRules(accountID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetMutingRule(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testMutingRuleGetResponseJSON)
	alerts := newMockResponse(t, respJSON, http.StatusOK)
	accountID := 400304
	ruleID := 123
	startTime, err1 := time.Parse(time.RFC3339, "2021-07-08T12:30:00-07:00")
	if err1 != nil {
		t.Fatal(err1)
	}
	endTime, err2 := time.Parse(time.RFC3339, "2021-07-08T14:30:00-07:00")
	if err2 != nil {
		t.Fatal(err2)
	}
	repeatCount := 10

	expected := MutingRule{
		ID:        123,
		AccountID: accountID,
		Condition: MutingRuleConditionGroup{
			Conditions: []MutingRuleCondition{
				{
					Attribute: "conditionName",
					Operator:  "EQUALS",
					Values:    []string{"please not me"},
				},
			},
			Operator: "AND",
		},
		CreatedAt: "2021-01-12T00:50:39.533Z",
		CreatedByUser: ByUser{
			Email:    "testemail@newrelic.com",
			Gravatar: "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
			ID:       1,
			Name:     "Test User",
		},
		Description: "",
		Enabled:     true,
		Name:        "Test Muting Rule",
		Schedule: &MutingRuleSchedule{
			EndRepeat:        nil,
			EndTime:          &endTime,
			Repeat:           &MutingRuleScheduleRepeatTypes.DAILY,
			RepeatCount:      &repeatCount,
			StartTime:        &startTime,
			TimeZone:         "America/Los_Angeles",
			WeeklyRepeatDays: nil,
		},
		UpdatedAt: "2021-01-12T00:50:39.533Z",
		UpdatedByUser: ByUser{
			Email:    "testemail@newrelic.com",
			Gravatar: "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
			ID:       1,
			Name:     "Test User",
		},
		ActionOnMutingRuleWindowEnded: "CLOSE_ISSUES_ON_INACTIVE",
	}

	actual, err := alerts.GetMutingRule(accountID, ruleID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestCreateMutingRule(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testMutingRuleCreateResponseJSON)
	alerts := newMockResponse(t, respJSON, http.StatusCreated)
	accountID := 400304
	startTime, err1 := time.Parse(time.RFC3339, "2021-07-08T12:30:00-07:00")
	if err1 != nil {
		t.Fatal(err1)
	}
	endTime, err2 := time.Parse(time.RFC3339, "2021-07-08T14:30:00-07:00")
	if err2 != nil {
		t.Fatal(err2)
	}
	repeatCount := 10

	expected := MutingRule{
		ID:        123,
		AccountID: accountID,
		Condition: MutingRuleConditionGroup{
			Conditions: []MutingRuleCondition{
				{
					Attribute: "conditionName",
					Operator:  "EQUALS",
					Values:    []string{"please not me"},
				},
			},
			Operator: "AND",
		},
		CreatedAt: "2021-01-12T00:50:39.533Z",
		CreatedByUser: ByUser{
			Email:    "testemail@newrelic.com",
			Gravatar: "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
			ID:       1,
			Name:     "Test User",
		},
		Description: "",
		Enabled:     true,
		Name:        "Test Muting Rule",
		Schedule: &MutingRuleSchedule{
			EndRepeat:        nil,
			EndTime:          &endTime,
			Repeat:           &MutingRuleScheduleRepeatTypes.DAILY,
			RepeatCount:      &repeatCount,
			StartTime:        &startTime,
			TimeZone:         "America/Los_Angeles",
			WeeklyRepeatDays: nil,
		},
		UpdatedAt: "2021-01-12T00:50:39.533Z",
		UpdatedByUser: ByUser{
			Email:    "testemail@newrelic.com",
			Gravatar: "https://secure.gravatar.com/avatar/692dc9742bd717014494f5093faff304",
			ID:       1,
			Name:     "Test User",
		},
		ActionOnMutingRuleWindowEnded: "CLOSE_ISSUES_ON_INACTIVE",
	}

	actual, err := alerts.CreateMutingRule(accountID, MutingRuleCreateInput{})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestUpdateMutingRule(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testMutingRuleUpdateResponseJSON)
	alerts := newMockResponse(t, respJSON, http.StatusCreated)
	accountID := 400304
	ruleID := 123
	startTime, err1 := time.Parse(time.RFC3339, "2021-07-08T12:30:00-07:00")
	if err1 != nil {
		t.Fatal(err1)
	}
	endRepeat, err1 := time.Parse(time.RFC3339, "2021-08-08T12:30:00-07:00")
	if err1 != nil {
		t.Fatal(err1)
	}

	expected := MutingRule{
		ID:        123,
		AccountID: accountID,
		Schedule: &MutingRuleSchedule{
			EndRepeat: &endRepeat,
			StartTime: &startTime,
		},
		ActionOnMutingRuleWindowEnded: "DO_NOTHING",
	}

	actual, err := alerts.UpdateMutingRule(accountID, ruleID, MutingRuleUpdateInput{})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestDeleteMutingRule(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testMutingRuleDeleteResponseJSON)
	alerts := newMockResponse(t, respJSON, http.StatusOK)
	accountID := 400304
	ruleID := 123

	err := alerts.DeleteMutingRule(accountID, ruleID)

	assert.NoError(t, err)
}

var (
	location, _ = time.LoadLocation("America/Los_Angeles")

	naiveDateTimeTests = []struct {
		in           time.Time
		out          string
		errorMessage string
	}{
		{time.Date(
			2006, 01, 02, 15, 04, 05, 0, time.UTC), "\"2006-01-02T15:04:05\"", ""},
		{time.Date(
			2006, 01, 02, 15, 04, 05, 0, location), "", "json: error calling MarshalJSON for type alerts.NaiveDateTime: time offset -28800 not allowed. You can call .UTC() on the time provided to reset the offset"},
	}
)

func TestNaiveDateTimeMarshaling(t *testing.T) {
	for _, tt := range naiveDateTimeTests {
		tt := tt
		expected := tt.out

		t.Run(tt.in.String(), func(t *testing.T) {
			t.Parallel()
			naiveDateTime := NaiveDateTime{tt.in}
			actual, err := json.Marshal(naiveDateTime)

			if string(actual) != expected {
				t.Errorf("expected %q, but got %q", expected, string(actual))
			}

			actualError := ""

			if err != nil {
				actualError = err.Error()
			}

			if tt.errorMessage != actualError {
				t.Errorf("expected %q, but got %q", tt.errorMessage, actualError)
			}
		})
	}
}
