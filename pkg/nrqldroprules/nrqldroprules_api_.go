package nrqldroprules

import (
	"context"
	"strconv"
)

// Get only particular rule using Rule id from all drop rules list for a given account
func (a *Nrqldroprules) GetDropRuleByID(
	accountID int,
	dropRuleID int,
) (*NRQLDropRulesDropRule, error) {
	dropRuleResult, err := a.GetListWithContext(context.Background(),
		accountID,
	)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(dropRuleResult.Rules); i++ {
		if dropRuleResult.Rules[i].ID == strconv.Itoa(dropRuleID) {
			return &dropRuleResult.Rules[i], nil
		}
	}
	return nil, err
}
