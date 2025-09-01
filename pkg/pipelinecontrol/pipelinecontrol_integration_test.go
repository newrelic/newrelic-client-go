package pipelinecontrol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationEntityManagement_PipelineCloudRule_CRUD(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)
	accountID, _ := testhelpers.GetTestAccountID()

	// Generate a unique name for the test rule
	ruleName := fmt.Sprintf("test-rule-%s", testhelpers.RandSeq(5))
	createInput := EntityManagementPipelineCloudRuleEntityCreateInput{
		Name:        ruleName,
		Description: "A test pipeline cloud rule from integration testing.",
		NRQL:        nrdb.NRQL("DELETE FROM Log where name = 'test-client-go-terraform'"),
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   fmt.Sprintf("%d", accountID),
		},
	}

	// 1. Create the entity
	createResult, err := client.EntityManagementCreatePipelineCloudRule(createInput)
	require.NotNil(t, createResult)
	require.NoError(t, err)

	require.NotEmpty(t, createResult.Entity.ID)
	require.Equal(t, ruleName, createResult.Entity.Name)

	// Defer the deletion to ensure cleanup even if assertions fail
	defer func() {
		_, deleteErr := client.EntityManagementDelete(createResult.Entity.ID, createResult.Entity.Metadata.Version)
		require.NoError(t, deleteErr, "Failed to clean up entity %s", createResult.Entity.ID)
	}()

	// 2. Read the entity to verify creation
	getResult, err := client.GetEntity(createResult.Entity.ID)
	require.NoError(t, err)
	require.NotNil(t, getResult)

	// Type assert the result to access specific fields
	ruleEntity, ok := (*getResult).(*EntityManagementPipelineCloudRuleEntity)
	require.True(t, ok, "Fetched entity was not of the expected type")
	require.Equal(t, createResult.Entity.ID, ruleEntity.ID)
	require.Equal(t, createInput.Name, ruleEntity.Name)
	require.Equal(t, createInput.Description, ruleEntity.Description)
	require.Equal(t, createInput.NRQL, ruleEntity.NRQL)

	// 3. Delete the entity (this is handled by the deferred function)
	// The test will complete, and the deferred function will execute for cleanup.
}
