//go:build integration
// +build integration

package nrqldroprules

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationDropRules(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	rand := mock.RandSeq(5)

	inputs := []NRQLDropRulesCreateDropRuleInput{
		{
			Description: "Primary test drop rule " + rand,
			NRQL:        "SELECT * FROM Log WHERE container_name = 'primary_noise_" + rand + "'",
			Action:      NRQLDropRulesActionTypes.DROP_DATA,
		},
		{
			Description: "Secondary test drop rule " + rand,
			NRQL:        "SELECT * FROM Log WHERE container_name = 'secondary_noise_" + rand + "'",
			Action:      NRQLDropRulesActionTypes.DROP_DATA,
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

	// Map description -> submitted input for parity checks
	inputByDesc := map[string]NRQLDropRulesCreateDropRuleInput{}
	for _, in := range inputs {
		inputByDesc[in.Description] = in
	}

	for _, s := range createRes.Successes {
		createdIDs = append(createdIDs, s.ID)

		submitted, ok := inputByDesc[s.Description]
		require.Truef(t, ok, "unexpected created rule description %q", s.Description)

		// Parity checks
		require.Equal(t, submitted.Action, s.Action, "action mismatch")
		require.Equal(t, submitted.Description, s.Description, "description mismatch")
		require.Equal(t, submitted.NRQL, s.NRQL, "nrql mismatch")

		// Response-only critical attribute
		// this doesn't work yet with create, to be revisited later
		// require.NotEmpty(t, s.PipelineCloudRuleEntityId, "PipelineCloudRuleEntityId must be populated")
		require.NotEmpty(t, s.ID, "id must be set")
	}

	// LIST
	time.Sleep(10 * time.Second)
	listRes, err := client.GetList(accountID)
	require.NoError(t, err, "list should succeed")
	require.NotNil(t, listRes)
	require.GreaterOrEqual(t, len(listRes.Rules), len(inputs), "list should contain at least created rules")

	// Build lookup from list
	listByID := map[string]NRQLDropRulesDropRule{}
	for _, r := range listRes.Rules {
		listByID[r.ID] = r
	}

	for _, id := range createdIDs {
		r, ok := listByID[id]
		require.Truef(t, ok, "expected created rule id %s in list", id)
		require.NotEmpty(t, r.PipelineCloudRuleEntityId, "listed rule %s missing PipelineCloudRuleEntityId", id)
	}

	// DELETE (explicit test coverage)
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

func TestIntegrationDropRules_Fail(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	rand := mock.RandSeq(5)

	// Intentionally invalid NRQL / input variants
	inputs := []NRQLDropRulesCreateDropRuleInput{
		{
			Description: "Invalid empty NRQL " + rand,
			NRQL:        "", // empty
			Action:      NRQLDropRulesActionTypes.DROP_DATA,
		},
		{
			Description: "Invalid malformed NRQL " + rand,
			NRQL:        "SELECT FROM", // malformed
			Action:      NRQLDropRulesActionTypes.DROP_DATA,
		},
		{
			Description: "Invalid wrong verb NRQL " + rand,
			NRQL:        "DELETE FROM Log WHERE foo = 'bar'", // not SELECT *
			Action:      NRQLDropRulesActionTypes.DROP_DATA,
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
		require.Contains(t, descSet, f.Submitted.Description, "unexpected failure description")
		require.NotEmpty(t, f.Error.Description, "error description should be populated")
		require.NotEmpty(t, f.Error.Reason, "error reason should be populated")
		switch f.Error.Reason {
		case NRQLDropRulesErrorReasonTypes.INVALID_QUERY,
			NRQLDropRulesErrorReasonTypes.INVALID_INPUT:
		default:
			t.Logf("warning: unexpected error reason %s for %s", f.Error.Reason, f.Submitted.Description)
		}
	}

	for range res.Failures {
		require.Zero(t, len(res.Successes))
	}
}

func newIntegrationTestClient(t *testing.T) Nrqldroprules {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
