package workflowautomation

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCreateWorkflowDefinitionResponseJSON = `{
		"workflowAutomationCreateWorkflowDefinition": {
			"definition": {
				"definitionId": "12345",
				"description": "Test workflow description",
				"name": "test-workflow",
				"version": 1,
				"yaml": "workflow:\n  name: test-workflow\n  steps:\n    - name: test-step"
			}
		}
	}`

	testUpdateWorkflowDefinitionResponseJSON = `{
		"workflowAutomationUpdateWorkflowDefinition": {
			"definition": {
				"definitionId": "12345",
				"description": "Updated workflow description",
				"name": "test-workflow",
				"version": 2,
				"yaml": "workflow:\n  name: test-workflow\n  steps:\n    - name: updated-step"
			}
		}
	}`

	testDeleteWorkflowDefinitionResponseJSON = `{
		"workflowAutomationDeleteWorkflowDefinition": {
			"definition": {
				"description": "Test workflow description",
				"name": "test-workflow",
				"scope": {
					"id": "12345",
					"type": "ACCOUNT"
				},
				"version": 1
			}
		}
	}`

	testGetWorkflowResponseJSON = `{
		"actor": {
			"account": {
				"workflowAutomation": {
					"workflow": {
						"definition": {
							"definitionId": "12345",
							"description": "Test workflow description",
							"name": "test-workflow",
							"version": 1,
							"yaml": "workflow:\n  name: test-workflow\n  steps:\n    - name: test-step"
						}
					}
				}
			}
		}
	}`
)

func TestWorkflowAutomationCreateWorkflowDefinition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{"data":%s}`, testCreateWorkflowDefinitionResponseJSON)
	workflowautomation := newMockResponse(t, respJSON, http.StatusCreated)

	definition := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: "workflow:\n  name: test-workflow\n  steps:\n    - name: test-step",
	}
	scope := WorkflowAutomationScopeInput{
		ID:   strconv.Itoa(12345),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{
		{
			Key:    "env",
			Values: []string{"test"},
		},
	}

	expected := &WorkflowAutomationCreateWorkflowDefinitionResponse{
		Definition: WorkflowAutomationWorkflowDefinition{
			DefinitionId: "12345",
			Description:  "Test workflow description",
			Name:         "test-workflow",
			Version:      1,
			Yaml:         "workflow:\n  name: test-workflow\n  steps:\n    - name: test-step",
		},
	}

	actual, err := workflowautomation.WorkflowAutomationCreateWorkflowDefinition(definition, scope, tags)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Definition.DefinitionId, actual.Definition.DefinitionId)
	assert.Equal(t, expected.Definition.Name, actual.Definition.Name)
	assert.Equal(t, expected.Definition.Version, actual.Definition.Version)
	assert.Equal(t, expected.Definition.Description, actual.Definition.Description)
	assert.Equal(t, expected.Definition.Yaml, actual.Definition.Yaml)
}

func TestWorkflowAutomationUpdateWorkflowDefinition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{"data":%s}`, testUpdateWorkflowDefinitionResponseJSON)
	workflowautomation := newMockResponse(t, respJSON, http.StatusOK)

	definition := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: "workflow:\n  name: test-workflow\n  steps:\n    - name: updated-step",
	}
	scope := WorkflowAutomationScopeInput{
		ID:   strconv.Itoa(12345),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{
		{
			Key:    "env",
			Values: []string{"test"},
		},
	}

	expected := &WorkflowAutomationUpdateWorkflowDefinitionResponse{
		Definition: WorkflowAutomationWorkflowDefinition{
			DefinitionId: "12345",
			Description:  "Updated workflow description",
			Name:         "test-workflow",
			Version:      2,
			Yaml:         "workflow:\n  name: test-workflow\n  steps:\n    - name: updated-step",
		},
	}

	actual, err := workflowautomation.WorkflowAutomationUpdateWorkflowDefinition(definition, scope, tags)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Definition.DefinitionId, actual.Definition.DefinitionId)
	assert.Equal(t, expected.Definition.Name, actual.Definition.Name)
	assert.Equal(t, expected.Definition.Version, actual.Definition.Version)
	assert.Equal(t, expected.Definition.Description, actual.Definition.Description)
	assert.Equal(t, expected.Definition.Yaml, actual.Definition.Yaml)
}

func TestWorkflowAutomationDeleteWorkflowDefinition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{"data":%s}`, testDeleteWorkflowDefinitionResponseJSON)
	workflowautomation := newMockResponse(t, respJSON, http.StatusOK)

	definition := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: "test-workflow",
	}
	scope := WorkflowAutomationScopeInput{
		ID:   strconv.Itoa(12345),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}

	expected := &WorkflowAutomationDeleteWorkflowDefinitionResponse{
		Definition: WorkflowAutomationWorkflowDefinitionOutline{
			Name:        "test-workflow",
			Description: "Test workflow description",
			Version:     1,
			Scope: WorkflowAutomationScope{
				ID:   "12345",
				Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
			},
		},
	}

	actual, err := workflowautomation.WorkflowAutomationDeleteWorkflowDefinition(definition, scope)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Definition.Name, actual.Definition.Name)
	assert.Equal(t, expected.Definition.Description, actual.Definition.Description)
	assert.Equal(t, expected.Definition.Version, actual.Definition.Version)
	assert.Equal(t, expected.Definition.Scope.ID, actual.Definition.Scope.ID)
	assert.Equal(t, expected.Definition.Scope.Type, actual.Definition.Scope.Type)
}

func TestGetWorkflow(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{"data":%s}`, testGetWorkflowResponseJSON)
	workflowautomation := newMockResponse(t, respJSON, http.StatusOK)

	accountID := 12345
	name := "test-workflow"
	version := 1

	expected := &WorkflowAutomationWorkflowResponse{
		Definition: WorkflowAutomationWorkflowDefinition{
			DefinitionId: "12345",
			Description:  "Test workflow description",
			Name:         "test-workflow",
			Version:      1,
			Yaml:         "workflow:\n  name: test-workflow\n  steps:\n    - name: test-step",
		},
	}

	actual, err := workflowautomation.GetWorkflow(accountID, name, version)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Definition.DefinitionId, actual.Definition.DefinitionId)
	assert.Equal(t, expected.Definition.Name, actual.Definition.Name)
	assert.Equal(t, expected.Definition.Version, actual.Definition.Version)
	assert.Equal(t, expected.Definition.Description, actual.Definition.Description)
	assert.Equal(t, expected.Definition.Yaml, actual.Definition.Yaml)
}

// Test error scenarios
func TestWorkflowAutomationCreateWorkflowDefinitionError(t *testing.T) {
	t.Parallel()
	respJSON := `{"errors": [{"message": "Invalid workflow YAML"}]}`
	workflowautomation := newMockResponse(t, respJSON, http.StatusBadRequest)

	definition := WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: "invalid yaml",
	}
	scope := WorkflowAutomationScopeInput{
		ID:   strconv.Itoa(12345),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	actual, err := workflowautomation.WorkflowAutomationCreateWorkflowDefinition(definition, scope, tags)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestWorkflowAutomationUpdateWorkflowDefinitionError(t *testing.T) {
	t.Parallel()
	respJSON := `{"errors": [{"message": "Workflow not found"}]}`
	workflowautomation := newMockResponse(t, respJSON, http.StatusNotFound)

	definition := WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: "workflow:\n  name: test-workflow",
	}
	scope := WorkflowAutomationScopeInput{
		ID:   strconv.Itoa(12345),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}
	tags := []WorkflowAutomationTag{}

	actual, err := workflowautomation.WorkflowAutomationUpdateWorkflowDefinition(definition, scope, tags)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestWorkflowAutomationDeleteWorkflowDefinitionError(t *testing.T) {
	t.Parallel()
	respJSON := `{"errors": [{"message": "Workflow not found"}]}`
	workflowautomation := newMockResponse(t, respJSON, http.StatusNotFound)

	definition := WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: "non-existent-workflow",
	}
	scope := WorkflowAutomationScopeInput{
		ID:   strconv.Itoa(12345),
		Type: WorkflowAutomationScopeTypeTypes.ACCOUNT,
	}

	actual, err := workflowautomation.WorkflowAutomationDeleteWorkflowDefinition(definition, scope)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestGetWorkflowError(t *testing.T) {
	t.Parallel()
	respJSON := `{"errors": [{"message": "Workflow not found"}]}`
	workflowautomation := newMockResponse(t, respJSON, http.StatusNotFound)

	accountID := 12345
	name := "non-existent-workflow"
	version := 1

	actual, err := workflowautomation.GetWorkflow(accountID, name, version)

	assert.Error(t, err)
	assert.Nil(t, actual)
}
