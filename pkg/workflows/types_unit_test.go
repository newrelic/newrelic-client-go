//go:build unit
// +build unit

package workflows

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// verify that an empty update input is serialized into an empty json
func TestAiWorkflowsUpdateWorkflowResponse_EmptyInput_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiWorkflowsUpdateWorkflowInput{
		ID: "10",
	}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(t, "{\"id\":\"10\"}", string(serialized))
}
