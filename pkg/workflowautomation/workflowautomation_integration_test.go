package workflowautomation

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationWorkflowAutomationCreateWorkflowDefinition(t *testing.T) {
	t.Skipf("skipping create workflow test case - requires manual setup and cleanup")
	t.Parallel()

	client := newWorkflowautomationTestClient(t)

	definition := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue(`name: test-workflow-tf-` + mock.RandSeq(5) + `
description: This is a test workflow created by terraform
steps:
  - name: stepchanged
    type: wait
    seconds: 10`),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{
		{
			Key:    "env",
			Values: []string{"test"},
		},
	}

	actual, err := client.WorkflowAutomationCreateWorkflowDefinition(definition, scope, tags)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.Definition.DefinitionId)
	require.NotEmpty(t, actual.Definition.Name)
	require.Equal(t, 1, actual.Definition.Version)
}

func TestIntegrationWorkflowAutomationUpdateWorkflowDefinition(t *testing.T) {
	t.Skipf("skipping update workflow test case - requires manual setup and cleanup")
	t.Parallel()

	client := newWorkflowautomationTestClient(t)

	definition := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: SecureValue(`name: existing-test-workflow
description: This is a test workflow created by terraform (updated)
steps:
  - name: stepchanged
    type: wait
    seconds: 15`),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{
		{
			Key:    "env",
			Values: []string{"test", "integration"},
		},
	}

	actual, err := client.WorkflowAutomationUpdateWorkflowDefinition(definition, scope, tags)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.Definition.DefinitionId)
	require.NotEmpty(t, actual.Definition.Name)
	require.Greater(t, actual.Definition.Version, 1)
}

func TestIntegrationWorkflowAutomationDeleteWorkflowDefinition(t *testing.T) {
	t.Skipf("skipping delete workflow test case - requires manual setup")
	t.Parallel()

	client := newWorkflowautomationTestClient(t)

	definition := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: "test-workflow-to-delete",
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}

	actual, err := client.WorkflowAutomationDeleteWorkflowDefinition(definition, scope)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, definition.Name, actual.Definition.Name)
	require.NotEmpty(t, actual.Definition.Scope.ID)
	require.Equal(t, WorkflowAutomationScopeTypeTypes.ACCOUNT, actual.Definition.Scope.Type)
}

func TestIntegrationGetWorkflow(t *testing.T) {
	t.Skipf("skipping get workflow test case - requires manual setup")
	t.Parallel()

	client := newWorkflowautomationTestClient(t)
	accountID := 12345 // Replace with actual account ID
	workflowName := "existing-workflow"
	version := 1

	actual, err := client.GetWorkflow(accountID, workflowName, version)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.Definition.DefinitionId)
	require.Equal(t, workflowName, actual.Definition.Name)
	require.Equal(t, version, actual.Definition.Version)
}

func TestIntegrationWorkflowAutomation_CreateUpdateDeleteWorkflow(t *testing.T) {
	t.Skipf("skipping full workflow lifecycle test - requires manual setup and cleanup")
	t.Parallel()

	client := newWorkflowautomationTestClient(t)
	workflowName := "test-workflow-" + mock.RandSeq(5)

	// Create Workflow
	createDefinition := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue(`name: ` + workflowName + `
description: This is a test workflow created by terraform
steps:
  - name: stepchanged
    type: wait
    seconds: 10`),
	}
	createScope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	createTags := []WorkflowAutomationTag{
		{
			Key:    "env",
			Values: []string{"test"},
		},
	}

	createResponse, err := client.WorkflowAutomationCreateWorkflowDefinition(createDefinition, createScope, createTags)

	require.NoError(t, err)
	require.NotNil(t, createResponse)
	require.NotNil(t, createResponse.Definition.DefinitionId)
	require.Equal(t, workflowName, createResponse.Definition.Name)
	require.Equal(t, 1, createResponse.Definition.Version)

	time.Sleep(time.Second * 2) // Wait for creation to propagate

	// Get Workflow
	getResponse, err := client.GetWorkflow(
		mock.IntegrationTestAccountID,
		workflowName,
		createResponse.Definition.Version,
	)

	require.NoError(t, err)
	require.NotNil(t, getResponse)
	require.Equal(t, createResponse.Definition.DefinitionId, getResponse.Definition.DefinitionId)
	require.Equal(t, createResponse.Definition.Name, getResponse.Definition.Name)

	time.Sleep(time.Second * 2)

	// Update Workflow
	updateDefinition := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: SecureValue(`name: ` + workflowName + `
description: This is a test workflow created by terraform (updated)
steps:
  - name: stepchanged
    type: wait
    seconds: 20`),
	}
	updateScope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	updateTags := []WorkflowAutomationTag{
		{
			Key:    "env",
			Values: []string{"test", "updated"},
		},
	}

	updateResponse, err := client.WorkflowAutomationUpdateWorkflowDefinition(updateDefinition, updateScope, updateTags)

	require.NoError(t, err)
	require.NotNil(t, updateResponse)
	require.Equal(t, createResponse.Definition.DefinitionId, updateResponse.Definition.DefinitionId)
	require.Equal(t, workflowName, updateResponse.Definition.Name)
	require.Equal(t, 2, updateResponse.Definition.Version) // Version should increment

	time.Sleep(time.Second * 2)

	// Delete Workflow
	deleteDefinition := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: workflowName,
	}
	deleteScope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}

	deleteResponse, err := client.WorkflowAutomationDeleteWorkflowDefinition(deleteDefinition, deleteScope)

	require.NoError(t, err)
	require.NotNil(t, deleteResponse)
	require.Equal(t, workflowName, deleteResponse.Definition.Name)
}

// Test error scenarios
func TestIntegrationWorkflowAutomationCreateWorkflowDefinitionError(t *testing.T) {
	t.Parallel()

	client := newWorkflowautomationTestClient(t)

	// Test with invalid YAML
	definition := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: SecureValue("invalid: yaml: structure: [[["),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	actual, err := client.WorkflowAutomationCreateWorkflowDefinition(definition, scope, tags)

	require.Error(t, err)
	require.Nil(t, actual)
}

func TestIntegrationWorkflowAutomationUpdateWorkflowDefinitionError(t *testing.T) {
	t.Parallel()

	client := newWorkflowautomationTestClient(t)

	// Test with non-existent workflow
	definition := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: SecureValue(`name: non-existent-workflow-` + mock.RandSeq(10) + `
description: This workflow does not exist
steps:
  - name: stepchanged
    type: wait
    seconds: 10`),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	actual, err := client.WorkflowAutomationUpdateWorkflowDefinition(definition, scope, tags)

	require.Error(t, err)
	require.Nil(t, actual)
}

func TestIntegrationWorkflowAutomationDeleteWorkflowDefinitionError(t *testing.T) {
	t.Parallel()

	client := newWorkflowautomationTestClient(t)

	// Test with non-existent workflow
	definition := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: "non-existent-workflow-" + mock.RandSeq(10),
	}
	scope := WorkflowAutomationScopeInput{
		ID:   fmt.Sprintf("%d", mock.IntegrationTestAccountID),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}

	actual, err := client.WorkflowAutomationDeleteWorkflowDefinition(definition, scope)

	require.Error(t, err)
	require.Nil(t, actual)
}

func TestIntegrationGetWorkflowError(t *testing.T) {
	t.Parallel()

	client := newWorkflowautomationTestClient(t)
	accountID := 12345 // Invalid account ID
	workflowName := "non-existent-workflow-" + mock.RandSeq(10)
	version := 1

	actual, err := client.GetWorkflow(accountID, workflowName, version)

	require.Error(t, err)
	require.Nil(t, actual)
}
