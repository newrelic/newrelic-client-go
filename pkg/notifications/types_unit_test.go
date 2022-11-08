//go:build unit
// +build unit

package notifications

import (
"encoding/json"
"testing"

"github.com/stretchr/testify/assert"
)


// Verify that boolean values with "empty" values (== false) are not ignored by json marshaller
// If `false` is not included into input, it would be impossible to change any boolean flag to `false` in an update
func TestAiNotificationsChannelFilter_OptionalBooleans_JsonFormat(t *testing.T) {
	t.Parallel()
	var falseValue = false
	var input = AiNotificationsChannelFilter{
		Name: "test-notification-channel-1-update",
		Active: &falseValue,
	}
	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"active\":false,\"name\":\"test-notification-channel-1-update\",\"property\":{\"key\":\"\",\"value\":\"\"}}",
		string(serialized),
	)
}

// Verify that an empty update input is serialized into an empty json
func TestAiNotificationsChannelFilter_EmptyInput_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiNotificationsChannelFilter{
	}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(t, "{\"property\":{\"key\":\"\",\"value\":\"\"}}", string(serialized))
}

// Verify that boolean values with "empty" values (== false) are not ignored by json marshaller
// If `false` is not included into input, it would be impossible to change any boolean flag to `false` in an update
func TestAiNotificationsChannelUpdate_OptionalBooleans_JsonFormat(t *testing.T) {
	t.Parallel()
	var falseValue = false
	var input = AiNotificationsChannelUpdate{
		Name: "test-notification-channel-1-update",
		Properties: []AiNotificationsPropertyInput{},
		Active: &falseValue,
	}
	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"active\":false,\"name\":\"test-notification-channel-1-update\"}",
		string(serialized),
	)
}

// Verify that an empty update input is serialized into an empty json
func TestAiNotificationsChannelUpdate_EmptyInput_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiNotificationsChannelUpdate{}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(t, "{}", string(serialized))
}

// Verify that boolean values with "empty" values (== false) are not ignored by json marshaller
// If `false` is not included into input, it would be impossible to change any boolean flag to `false` in an update
func TestAiNotificationsDestinationFilter_OptionalBooleans_JsonFormat(t *testing.T) {
	t.Parallel()
	var falseValue = false
	var input = AiNotificationsDestinationFilter{
		Name: "test-notification-channel-1-update",
		Active: &falseValue,
	}
	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"active\":false,\"name\":\"test-notification-channel-1-update\",\"property\":{\"key\":\"\",\"value\":\"\"}}",
		string(serialized),
	)
}

// Verify that an empty update input is serialized into an empty json
func TestAiNotificationsDestinationFilter_EmptyInput_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiNotificationsDestinationFilter{}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(t, "{\"property\":{\"key\":\"\",\"value\":\"\"}}", string(serialized))
}

// Verify that boolean values with "empty" values (== false) are not ignored by json marshaller
// If `false` is not included into input, it would be impossible to change any boolean flag to `false` in an update
func TestAiNotificationsDestinationUpdate_OptionalBooleans_JsonFormat(t *testing.T) {
	t.Parallel()
	var falseValue = false
	var input = AiNotificationsDestinationUpdate{
		Name: "test-notification-channel-1-update",
		Properties: []AiNotificationsPropertyInput{},
		Active: &falseValue,
		DisableAuth: &falseValue,
	}
	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"active\":false,\"disableAuth\":false,\"name\":\"test-notification-channel-1-update\"}",
		string(serialized),
	)
}

// Verify that an empty update input is serialized into an empty json
func TestAiNotificationsDestinationUpdate_EmptyInput_JsonFormat(t *testing.T) {
	t.Parallel()
	var input = AiNotificationsDestinationUpdate{}

	var serialized, err = json.Marshal(input)

	assert.NoError(t, err)
	assert.Equal(t, "{}", string(serialized))
}
