//go:build integration
// +build integration

package pathpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationPathPoint_CRUD(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	flowName := "newrelic-client-go test-pathpoint-" + mock.RandSeq(5)

	// Build a minimal but realistic flow: one stage with one level and one step.
	flowInput := PathPointFlowInput{
		Name:            flowName,
		Category:        "Testing",
		Description:     "Integration test flow",
		HealthRollup:    PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP,
		RefreshInterval: PathPointRefreshIntervalTypes.FIVE_MINUTES,
		Stages: []PathPointStageInput{
			{
				Name:         "Test Stage",
				HealthRollup: PathPointStageHealthRollupTypes.AUTOMATIC_ROLL_UP,
				Levels: []PathPointLevelInput{
					{
						Steps: []PathPointStepInput{
							{
								Name:           "Test Step",
								ScopedAccounts: []int{accountID},
								EntitySearchQuery: PathPointSignalQueryInput{
									Query: "domain = 'APM' AND type = 'APPLICATION'",
								},
								Config: PathPointStepStatusThresholdInput{
									HealthRollup:   PathPointStepHealthRollupTypes.WORST_STATUS_WINS,
									ThresholdType:  PathPointThresholdTypeTypes.FIXED,
									ThresholdValue: 1,
								},
							},
						},
					},
				},
			},
		},
	}

	scope := PathPointScopeInput{
		ID:   accountID,
		Type: PathPointScopeTypeTypes.ACCOUNT,
	}

	// Test: PathPointCreate
	created, err := client.PathPointCreate(flowInput, scope)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.GUID)

	cleanedUp := false
	defer func() {
		if !cleanedUp {
			_, err := client.PathPointDelete(created.GUID)
			if err != nil {
				t.Logf("error cleaning up pathpoint flow %s: %s", created.GUID, err)
			}
		}
	}()

	assert.Equal(t, flowName, created.Name)
	assert.Equal(t, "Testing", created.Category)
	assert.Equal(t, "Integration test flow", created.Description)
	assert.Equal(t, PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP, created.HealthRollup)
	assert.Equal(t, PathPointRefreshIntervalTypes.FIVE_MINUTES, created.RefreshInterval)
	require.Equal(t, 1, created.Stages.TotalCount)

	// Test: GetFlow
	fetched, err := client.GetFlow(accountID, created.GUID)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	assert.Equal(t, created.GUID, fetched.GUID)
	assert.Equal(t, flowName, fetched.Name)
	assert.Equal(t, "Testing", fetched.Category)
	assert.Equal(t, "Integration test flow", fetched.Description)
	assert.Equal(t, PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP, fetched.HealthRollup)
	require.Equal(t, 1, fetched.Stages.TotalCount)
	require.Len(t, fetched.Stages.Items, 1)
	assert.Equal(t, "Test Stage", fetched.Stages.Items[0].Name)

	// Test: PathPointUpdate — rename the flow, change category/description, and change refresh interval
	updateInput := PathPointFlowUpdateInput{
		Name:            flowName + "-updated",
		Category:        "Testing-Updated",
		Description:     "Integration test flow updated",
		HealthRollup:    PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP,
		RefreshInterval: PathPointRefreshIntervalTypes.TEN_MINUTES,
		Version:         fetched.Version,
		Stages: []PathPointStageUpdateInput{
			{
				ID:           fetched.Stages.Items[0].ID,
				Name:         "Updated Stage",
				HealthRollup: PathPointStageHealthRollupTypes.AUTOMATIC_ROLL_UP,
				Levels: []PathPointLevelUpdateInput{
					{
						ID: fetched.Stages.Items[0].Levels.Items[0].ID,
						Steps: []PathPointStepUpdateInput{
							{
								ID:             fetched.Stages.Items[0].Levels.Items[0].Steps.Items[0].ID,
								Name:           "Test Step",
								ScopedAccounts: []int{accountID},
								EntitySearchQuery: PathPointSignalQueryInput{
									Query: "domain = 'APM' AND type = 'APPLICATION'",
								},
								Config: PathPointStepStatusThresholdInput{
									HealthRollup:   PathPointStepHealthRollupTypes.WORST_STATUS_WINS,
									ThresholdType:  PathPointThresholdTypeTypes.FIXED,
									ThresholdValue: 1,
								},
							},
						},
					},
				},
			},
		},
	}

	updated, err := client.PathPointUpdate(created.GUID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, created.GUID, updated.GUID)
	assert.Equal(t, flowName+"-updated", updated.Name)
	assert.Equal(t, "Testing-Updated", updated.Category)
	assert.Equal(t, "Integration test flow updated", updated.Description)
	assert.Equal(t, PathPointRefreshIntervalTypes.TEN_MINUTES, updated.RefreshInterval)
	require.Equal(t, 1, updated.Stages.TotalCount)
	assert.Equal(t, "Updated Stage", updated.Stages.Items[0].Name)

	// Test: PathPointDelete
	deleteResult, err := client.PathPointDelete(created.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	assert.Equal(t, created.GUID, deleteResult.GUID)
	cleanedUp = true
}

func TestIntegrationPathPoint_GetFlow_InvalidGUID(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	result, err := client.GetFlow(accountID, "invalid-guid")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegrationPathPoint_Update_InvalidGUID(t *testing.T) {
	t.Parallel()

	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	update := PathPointFlowUpdateInput{
		Name:         "Should Not Update",
		HealthRollup: PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP,
	}

	result, err := client.PathPointUpdate("invalid-guid", update)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegrationPathPoint_Delete_AlreadyDeleted(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Create a flow to delete.
	flowInput := PathPointFlowInput{
		Name:            "newrelic-client-go test-double-delete-" + mock.RandSeq(5),
		HealthRollup:    PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP,
		RefreshInterval: PathPointRefreshIntervalTypes.FIVE_MINUTES,
	}
	scope := PathPointScopeInput{ID: accountID, Type: PathPointScopeTypeTypes.ACCOUNT}

	created, err := client.PathPointCreate(flowInput, scope)
	require.NoError(t, err)
	require.NotNil(t, created)

	// First delete should succeed.
	_, err = client.PathPointDelete(created.GUID)
	require.NoError(t, err)

	// Second delete on the same GUID should return an error.
	result, err := client.PathPointDelete(created.GUID)
	require.Error(t, err)
	assert.Nil(t, result)
}
