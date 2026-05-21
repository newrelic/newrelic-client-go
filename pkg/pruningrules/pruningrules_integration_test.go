//go:build integration
// +build integration

package pruningrules

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationPruningRules(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	rand := mock.RandSeq(5)

	inputs := []NRQLDropRulesCreateDropRuleInput{
		{
			Description: "Primary pruning rule " + rand,
			NRQL:        "SELECT collector.name FROM Metric WHERE metricName = 'test.pruning.primary." + rand + "'",
			Action:      NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES,
		},
		{
			Description: "Secondary pruning rule " + rand,
			NRQL:        "SELECT collector.name FROM Metric WHERE metricName = 'test.pruning.secondary." + rand + "'",
			Action:      NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES,
		},
	}

	client := newIntegrationTestClient(t)

	// CREATE
	createRes, err := client.NRQLDropRulesCreate(accountID, inputs)
	require.NoError(t, err, "create should succeed")
	require.NotNil(t, createRes)
	require.Len(t, createRes.Failures, 0, "expected zero failures")
	require.Len(t, createRes.Successes, len(inputs), "all rules should be created")

	createdIDs := make([]string, 0, len(createRes.Successes))
	deleted := false
	defer func() {
		if deleted || len(createdIDs) == 0 {
			return
		}
		_, _ = client.NRQLDropRulesDelete(accountID, createdIDs) // best-effort cleanup
	}()

	inputByDesc := map[string]NRQLDropRulesCreateDropRuleInput{}
	for _, in := range inputs {
		inputByDesc[in.Description] = in
	}

	for _, s := range createRes.Successes {
		createdIDs = append(createdIDs, s.ID)

		submitted, ok := inputByDesc[s.Description]
		require.Truef(t, ok, "unexpected created rule description %q", s.Description)

		require.Equal(t, NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES, s.Action, "action must be DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES")
		require.Equal(t, submitted.Description, s.Description)
		require.Equal(t, submitted.NRQL, s.NRQL)
		require.NotEmpty(t, s.ID, "id must be set")
	}

	// LIST — allow propagation
	time.Sleep(10 * time.Second)
	listRes, err := client.GetList(accountID)
	require.NoError(t, err, "list should succeed")
	require.NotNil(t, listRes)
	require.GreaterOrEqual(t, len(listRes.Rules), len(inputs))

	listByID := map[string]NRQLDropRulesDropRule{}
	for _, r := range listRes.Rules {
		listByID[r.ID] = r
	}

	for _, id := range createdIDs {
		r, ok := listByID[id]
		require.Truef(t, ok, "expected created rule id %s in list", id)
		require.NotEmpty(t, r.PipelineCloudRuleEntityId, "listed rule %s missing PipelineCloudRuleEntityId", id)
	}

	// GET BY ID
	rule, err := client.GetPruningRuleByID(accountID, createdIDs[0])
	require.NoError(t, err, "GetPruningRuleByID should succeed")
	require.Equal(t, createdIDs[0], rule.ID)
	require.Equal(t, NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES, rule.Action)

	// DELETE
	deleteRes, err := client.NRQLDropRulesDelete(accountID, createdIDs)
	require.NoError(t, err, "delete should succeed")
	require.NotNil(t, deleteRes)
	require.Len(t, deleteRes.Failures, 0, "no delete failures expected")
	require.Len(t, deleteRes.Successes, len(createdIDs), "all created rules should be deleted")
	for _, s := range deleteRes.Successes {
		require.NotEmpty(t, s.PipelineCloudRuleEntityId, "deleted success missing PipelineCloudRuleEntityId")
	}
	deleted = true
}

func TestIntegrationPruningRules_Fail(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	rand := mock.RandSeq(5)

	inputs := []NRQLDropRulesCreateDropRuleInput{
		{
			Description: "Invalid empty NRQL " + rand,
			NRQL:        "",
			Action:      NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES,
		},
		{
			// FROM non-Metric type is invalid for DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES
			Description: "Invalid non-Metric type " + rand,
			NRQL:        "SELECT container_name FROM Log WHERE container_name = 'noise_" + rand + "'",
			Action:      NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES,
		},
		{
			// SELECT * is invalid for DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES (requires specific attributes)
			Description: "Invalid SELECT star " + rand,
			NRQL:        "SELECT * FROM Metric WHERE metricName = 'test.pruning." + rand + "'",
			Action:      NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES,
		},
	}

	client := newIntegrationTestClient(t)

	res, err := client.NRQLDropRulesCreate(accountID, inputs)
	require.NoError(t, err, "API call itself should succeed")
	require.NotNil(t, res)
	require.Len(t, res.Successes, 0, "no invalid rule should succeed")
	require.Len(t, res.Failures, len(inputs), "all invalid submissions should fail")

	descSet := map[string]struct{}{}
	for _, in := range inputs {
		descSet[in.Description] = struct{}{}
	}

	for _, f := range res.Failures {
		require.NotZero(t, f.Submitted.AccountID)
		require.NotEmpty(t, f.Submitted.Description)
		require.Contains(t, descSet, f.Submitted.Description)
		require.NotEmpty(t, f.Error.Description)
		require.NotEmpty(t, f.Error.Reason)
	}
}

func TestIntegrationPruningRules_GetByIDNotFound(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	_, err = client.GetPruningRuleByID(accountID, "nonexistent-id-000000")
	require.Error(t, err, "should return error for nonexistent ID")
	require.Contains(t, err.Error(), "not found")
}

func newIntegrationTestClient(t *testing.T) Pruningrules {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
