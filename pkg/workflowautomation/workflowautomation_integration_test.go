//go:build integration
// +build integration

package workflowautomation

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/organization"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// generateWorkflowYAML creates a workflow YAML with a dynamic name
func generateWorkflowYAML(name, description string) string {
	return fmt.Sprintf(`name: %s
description: %s
workflowInputs:
  testInput:
    type: String
    defaultValue: "test-value"
    validations: []
steps:
  - name: testStep
    type: action
    action: newrelic.nerdgraph.execute
    version: 1
    inputs:
      graphql: |
        query {
          actor {
            user {
              email
              id
              name
            }
          }
        }
      variables: {}`, name, description)
}

// generateUpdatedWorkflowYAML creates an updated workflow YAML with a dynamic name
func generateUpdatedWorkflowYAML(name, description string) string {
	return fmt.Sprintf(`name: %s
description: %s
workflowInputs:
  testInput:
    type: String
    defaultValue: "updated-test-value"
    validations: []
  additionalInput:
    type: String
    defaultValue: "new-input"
    validations: []
steps:
  - name: testStepUpdated
    type: action
    action: newrelic.nerdgraph.execute
    version: 1
    inputs:
      graphql: |
        query {
          actor {
            user {
              email
              id
              name
            }
          }
        }
      variables: {}`, name, description)
}

var (
	organizationID = os.Getenv("INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID")
)

func newIntegrationTestClient(t *testing.T) Workflowautomation {
	cfg := mock.NewIntegrationTestConfig(t)
	return New(cfg)
}

func newOrganizationIntegrationTestClient(t *testing.T) organization.OrganizationManagement {
	cfg := mock.NewIntegrationTestConfig(t)
	return organization.New(cfg)
}

// TestIntegrationAccountScopedWorkflow tests the full lifecycle of an account-scoped workflow
func TestIntegrationAccountScopedWorkflow(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	workflowName := fmt.Sprintf("test-account-workflow-%s", mock.RandSeq(5))

	// Step 1: Create account-scoped workflow
	createInput := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue(generateWorkflowYAML(workflowName, "Test workflow for integration testing")),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", accountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{
		{
			Key:    "environment",
			Values: []string{"test"},
		},
	}

	createResult, err := client.WorkflowAutomationCreateWorkflowDefinition(createInput, scope, tags)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.NotNil(t, createResult.Definition)
	require.NotEmpty(t, createResult.Definition.DefinitionId)
	require.Equal(t, workflowName, createResult.Definition.Name)
	require.Equal(t, 1, createResult.Definition.Version)

	// Ensure cleanup happens even if test fails
	defer func() {
		deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
			Name: createResult.Definition.Name,
		}
		_, _ = client.WorkflowAutomationDeleteWorkflowDefinition(deleteInput, scope)
	}()

	// Step 2: Read the workflow
	readResult, err := client.GetWorkflow(accountID, createResult.Definition.Name, 0)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.NotNil(t, readResult.Definition)
	require.Equal(t, createResult.Definition.DefinitionId, readResult.Definition.DefinitionId)
	require.Equal(t, createResult.Definition.Name, readResult.Definition.Name)
	require.Equal(t, createResult.Definition.Version, readResult.Definition.Version)

	// Step 3: Update the workflow
	updateInput := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: SecureValue(generateUpdatedWorkflowYAML(workflowName, "Updated test workflow for integration testing")),
	}
	updateResult, err := client.WorkflowAutomationUpdateWorkflowDefinition(updateInput, scope, tags)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.NotNil(t, updateResult.Definition)
	require.Equal(t, createResult.Definition.DefinitionId, updateResult.Definition.DefinitionId)
	require.Equal(t, 2, updateResult.Definition.Version) // Version should increment

	// Step 4: Read again to verify update
	readAfterUpdate, err := client.GetWorkflow(accountID, createResult.Definition.Name, 0)
	require.NoError(t, err)
	require.NotNil(t, readAfterUpdate)
	require.Equal(t, 2, readAfterUpdate.Definition.Version)

	// Step 5: Delete the workflow
	deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: createResult.Definition.Name,
	}
	deleteResult, err := client.WorkflowAutomationDeleteWorkflowDefinition(deleteInput, scope)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.NotNil(t, deleteResult.Definition)
	require.Equal(t, createResult.Definition.Name, deleteResult.Definition.Name)

	// Step 6: Verify deletion - workflow should not exist
	_, err = client.GetWorkflow(accountID, createResult.Definition.Name, 0)
	require.Error(t, err) // Should get an error when trying to read deleted workflow
}

// TestIntegrationOrganizationScopedWorkflow tests the full lifecycle of an organization-scoped workflow
func TestIntegrationOrganizationScopedWorkflow(t *testing.T) {
	t.Parallel()

	if organizationID == "" {
		t.Skip("Skipping test: INTEGRATION_TESTING_NEW_RELIC_ORGANIZATION_ID environment variable not set")
	}

	client := newIntegrationTestClient(t)
	orgClient := newOrganizationIntegrationTestClient(t)
	workflowName := fmt.Sprintf("test-org-workflow-%s", mock.RandSeq(5))

	// Step 1: Create organization-scoped workflow
	createInput := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue(generateWorkflowYAML(workflowName, "Test org workflow for integration testing")),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   organizationID,
		Type: WorkflowAutomationScopeTypeTypes.ORGANIZATION,
	}
	tags := []WorkflowAutomationTag{
		{
			Key:    "environment",
			Values: []string{"test", "integration"},
		},
		{
			Key:    "team",
			Values: []string{"platform"},
		},
	}

	createResult, err := client.WorkflowAutomationCreateWorkflowDefinition(createInput, scope, tags)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.NotNil(t, createResult.Definition)
	require.NotEmpty(t, createResult.Definition.DefinitionId)
	require.Equal(t, workflowName, createResult.Definition.Name)

	// Ensure cleanup happens even if test fails
	defer func() {
		deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
			Name: createResult.Definition.Name,
		}
		_, _ = client.WorkflowAutomationDeleteWorkflowDefinition(deleteInput, scope)
	}()

	// Step 2: Read the workflow from organization package
	readResult, err := orgClient.GetWorkflow(createResult.Definition.Name, 0)
	require.NoError(t, err)
	require.NotNil(t, readResult)
	require.NotNil(t, readResult.Definition)
	require.Equal(t, createResult.Definition.DefinitionId, readResult.Definition.DefinitionId)
	require.Equal(t, createResult.Definition.Name, readResult.Definition.Name)

	// Step 3: Update the workflow with new tags
	updateInput := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: SecureValue(generateUpdatedWorkflowYAML(workflowName, "Updated test org workflow for integration testing")),
	}
	updatedTags := []WorkflowAutomationTag{
		{
			Key:    "environment",
			Values: []string{"test", "integration", "updated"},
		},
	}
	updateResult, err := client.WorkflowAutomationUpdateWorkflowDefinition(updateInput, scope, updatedTags)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.NotNil(t, updateResult.Definition)
	require.Equal(t, createResult.Definition.DefinitionId, updateResult.Definition.DefinitionId)

	// Step 4: Read again from organization package to verify update
	readAfterUpdate, err := orgClient.GetWorkflow(createResult.Definition.Name, 0)
	require.NoError(t, err)
	require.NotNil(t, readAfterUpdate)
	require.Equal(t, updateResult.Definition.Version, readAfterUpdate.Definition.Version)

	// Step 5: Delete the workflow
	deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: createResult.Definition.Name,
	}
	deleteResult, err := client.WorkflowAutomationDeleteWorkflowDefinition(deleteInput, scope)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.NotNil(t, deleteResult.Definition)
	require.Equal(t, createResult.Definition.Name, deleteResult.Definition.Name)
	require.Equal(t, organizationID, deleteResult.Definition.Scope.ID)
	require.Equal(t, WorkflowAutomationScopeTypeTypes.ORGANIZATION, deleteResult.Definition.Scope.Type)

	// Step 6: Verify deletion - workflow should not exist
	_, err = orgClient.GetWorkflow(createResult.Definition.Name, 0)
	require.Error(t, err) // Should get an error when trying to read deleted workflow
}

// TestIntegrationWorkflowAutomation_GetNonExistentWorkflow tests error handling for non-existent workflows
func TestIntegrationWorkflowAutomation_GetNonExistentWorkflow(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Try to read a workflow that doesn't exist
	nonExistentWorkflowName := fmt.Sprintf("non-existent-workflow-%s", mock.RandSeq(10))
	_, err = client.GetWorkflow(accountID, nonExistentWorkflowName, 0)
	require.Error(t, err)
}

// TestIntegrationWorkflowAutomation_DeleteNonExistentWorkflow tests error handling when deleting non-existent workflows
func TestIntegrationWorkflowAutomation_DeleteNonExistentWorkflow(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Try to delete a workflow that doesn't exist
	nonExistentWorkflowName := fmt.Sprintf("non-existent-workflow-%s", mock.RandSeq(10))
	deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: nonExistentWorkflowName,
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", accountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}

	_, err = client.WorkflowAutomationDeleteWorkflowDefinition(deleteInput, scope)
	require.Error(t, err)
}

// TestIntegrationWorkflowAutomation_CreateWithEmptyYAML tests error handling for empty YAML
func TestIntegrationWorkflowAutomation_CreateWithEmptyYAML(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Try to create a workflow with empty YAML
	createInput := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue(""),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", accountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	_, err = client.WorkflowAutomationCreateWorkflowDefinition(createInput, scope, tags)
	require.Error(t, err)
}

// TestIntegrationWorkflowAutomation_UpdateNonExistentWorkflow tests error handling when updating non-existent workflows
func TestIntegrationWorkflowAutomation_UpdateNonExistentWorkflow(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Try to update a workflow that doesn't exist
	nonExistentName := fmt.Sprintf("non-existent-workflow-%s", mock.RandSeq(10))
	updateInput := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: SecureValue(generateUpdatedWorkflowYAML(nonExistentName, "Non-existent workflow")),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", accountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	// Note: The workflow name would be extracted from the YAML
	_, err = client.WorkflowAutomationUpdateWorkflowDefinition(updateInput, scope, tags)
	require.Error(t, err)
}

// TestIntegrationWorkflowAutomation_CreateAndUpdateMultipleTimes tests multiple updates to the same workflow
func TestIntegrationWorkflowAutomation_CreateAndUpdateMultipleTimes(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	workflowName := fmt.Sprintf("test-multi-update-workflow-%s", mock.RandSeq(5))

	// Create workflow
	createInput := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue(generateWorkflowYAML(workflowName, "Test workflow for multiple updates")),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", accountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	createResult, err := client.WorkflowAutomationCreateWorkflowDefinition(createInput, scope, tags)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Ensure cleanup
	defer func() {
		deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
			Name: createResult.Definition.Name,
		}
		_, _ = client.WorkflowAutomationDeleteWorkflowDefinition(deleteInput, scope)
	}()

	// Perform multiple updates
	for i := 1; i <= 3; i++ {
		description := fmt.Sprintf("Updated test workflow iteration %d", i)
		updateInput := WorkflowAutomationUpdateWorkflowDefinitionInput{
			Yaml: SecureValue(generateUpdatedWorkflowYAML(workflowName, description)),
		}
		updateResult, err := client.WorkflowAutomationUpdateWorkflowDefinition(updateInput, scope, tags)
		require.NoError(t, err)
		require.NotNil(t, updateResult)
		require.Equal(t, i+1, updateResult.Definition.Version) // Version should increment with each update
	}

	// Read final version
	finalRead, err := client.GetWorkflow(accountID, workflowName, 0)
	require.NoError(t, err)
	require.Equal(t, 4, finalRead.Definition.Version) // Original version 1 + 3 updates
}

// TestIntegrationWorkflowAutomation_WithContextMethods tests the WithContext variants
func TestIntegrationWorkflowAutomation_WithContextMethods(t *testing.T) {
	t.Parallel()

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	workflowName := fmt.Sprintf("test-context-workflow-%s", mock.RandSeq(5))

	// Test CreateWithContext
	createInput := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue(generateWorkflowYAML(workflowName, "Test workflow with context methods")),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", accountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	ctx := context.Background()
	createResult, err := client.WorkflowAutomationCreateWorkflowDefinitionWithContext(ctx, createInput, scope, tags)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Ensure cleanup
	defer func() {
		deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
			Name: createResult.Definition.Name,
		}
		_, _ = client.WorkflowAutomationDeleteWorkflowDefinitionWithContext(ctx, deleteInput, scope)
	}()

	// Test GetWorkflowWithContext
	readResult, err := client.GetWorkflowWithContext(ctx, accountID, createResult.Definition.Name, 0)
	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test UpdateWithContext
	updateInput := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: SecureValue(generateUpdatedWorkflowYAML(workflowName, "Updated workflow with context methods")),
	}
	updateResult, err := client.WorkflowAutomationUpdateWorkflowDefinitionWithContext(ctx, updateInput, scope, tags)
	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test DeleteWithContext
	deleteInput := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: createResult.Definition.Name,
	}
	deleteResult, err := client.WorkflowAutomationDeleteWorkflowDefinitionWithContext(ctx, deleteInput, scope)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}
