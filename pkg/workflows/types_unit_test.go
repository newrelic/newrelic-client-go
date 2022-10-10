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
