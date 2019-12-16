package infrastructure

import (
	"strconv"
)

type listAlertConditionsResponse struct {
	AlertConditions []AlertCondition `json:"data,omitempty"`
}

// ListAlertConditions is used to retrieve New Relic Infrastructure alert conditions.
func (i *Infrastructure) ListAlertConditions(policyID int) ([]AlertCondition, error) {
	res := listAlertConditionsResponse{}
	paramsMap := map[string]string{
		"policy_id": strconv.Itoa(policyID),
		"limit":     "1",
	}

	responses, err := i.client.GetMultiple("/alerts/conditions", &paramsMap, &res)

	alertConditions := []AlertCondition{}
	for _, r := range responses {
		if response, ok := r.(*listAlertConditionsResponse); ok {
			alertConditions = append(alertConditions, response.AlertConditions...)
		}
	}

	if err != nil {
		return nil, err
	}

	return alertConditions, nil
}
