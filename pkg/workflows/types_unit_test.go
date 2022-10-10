//go:build unit
// +build unit

package workflows

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verify that boolean values with "empty" values (== false) are not ignored by json marshaller
// If `false` is not included into input, it would be impossible to change any boolean flag to `false` in an update
func TestAiWorkflowsUpdateWorkflowResponse_OptionalBooleans_JsonFormat(t *testing.T) {
	t.Parallel()
	var falseValue = false
	var input = AiWorkflowsUpdateWorkflowInput{
		ID:                  "10",
		WorkflowEnabled:     &falseValue,
		EnrichmentsEnabled:  &falseValue,
		DestinationsEnabled: &falseValue,
	}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"destinationsEnabled\":false,\"enrichmentsEnabled\":false,\"id\":\"10\",\"workflowEnabled\":false}",
		string(serialized),
	)
}

// Verify that if the user wants to erase all destinations, we explicitly pass an empty destination list
// While it might not be impossible to remove all destinations from a workflow, we should at least forward the user intent
// to out backend API. This way the user would see a meaningful error.
// If we just silently omit empty arrays, the changes just won't be applied while reporting a success, which is really confusing
func TestAiWorkflowsUpdateWorkflowResponse_EmptyDestinations_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiWorkflowsUpdateWorkflowInput{
		ID:                        "10",
		DestinationConfigurations: &[]AiWorkflowsDestinationConfigurationInput{},
	}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"destinationConfigurations\":[],\"id\":\"10\"}",
		string(serialized),
	)
}

// Verify that if the user wants to erase all workflow name, we explicitly pass an empty value
// While it might not be impossible to remove a workflow name workflow, we should at least forward the user intent
// to out backend API. This way the user would see a meaningful error.
// If we just silently omit empty name update, the changes just won't be applied while reporting a success, which is really confusing
func TestAiWorkflowsUpdateWorkflowResponse_EmptyName_JsonFormat(t *testing.T) {
	t.Parallel()
	var emptyStr = ""
	var input = AiWorkflowsUpdateWorkflowInput{
		ID:   "10",
		Name: &emptyStr,
	}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"id\":\"10\",\"name\":\"\"}",
		string(serialized),
	)
}

// Verify that it is possible to pass an empty value for muting rules
func TestAiWorkflowsUpdateWorkflowResponse_EmptyMutingRules_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiWorkflowsUpdateWorkflowInput{
		ID:                  "10",
		MutingRulesHandling: "",
	}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"id\":\"10\"}",
		string(serialized),
	)
}

// Verify that an empty update input is serialized into an empty json
func TestAiWorkflowsUpdateWorkflowResponse_EmptyInput_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiWorkflowsUpdateWorkflowInput{
		ID: "10",
	}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(t, "{\"id\":\"10\"}", string(serialized))
}
