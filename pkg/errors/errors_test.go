//go:build unit
// +build unit

package errors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorNotFound(t *testing.T) {
	t.Parallel()

	var e NotFound

	assert.Equal(t, "resource not found", e.Error())
}

func TestErrorNotFoundf(t *testing.T) {
	t.Parallel()

	e := NewNotFoundf("test %s", "format")

	assert.Equal(t, "test format", e.Error())
}

func TestErrorInvalidInput(t *testing.T) {
	t.Parallel()

	var e InvalidInput

	assert.Equal(t, "invalid input error", e.Error())
}

func TestErrorUnexpectedStatusCode(t *testing.T) {
	t.Parallel()

	e := NewUnexpectedStatusCode(99, "wat")

	assert.Equal(t, "99 response returned: wat", e.Error())
}

func TestErrorUnexpectedStatusCodef(t *testing.T) {
	t.Parallel()

	e := NewUnexpectedStatusCodef(99, "really bad %s", "request")

	assert.Equal(t, "99 response returned: really bad request", e.Error())
}

func TestErrorUnauthorized(t *testing.T) {
	t.Parallel()

	e := NewUnauthorizedError()

	assert.Equal(t, 401, e.statusCode)
	assert.True(t, strings.Contains(e.Error(), "Invalid credentials provided"))
}

func TestErrorMaxRetriesReached(t *testing.T) {
	t.Parallel()

	e := NewMaxRetriesReached("2")

	assert.Equal(t, e.Error(), "maximum retries reached: 2")
}

func TestErrorMaxRetriesReachedf(t *testing.T) {
	t.Parallel()

	e := NewMaxRetriesReachedf("2 + %d", 2)

	assert.Equal(t, "maximum retries reached: 2 + 2", e.Error())
}

func TestInvalidInput(t *testing.T) {
	t.Parallel()

	e := NewInvalidInput("")

	assert.Equal(t, e.Error(), "invalid input error")
}

func TestInvalidInputf(t *testing.T) {
	t.Parallel()

	e := NewInvalidInputf("test %s", "format")

	assert.Equal(t, "test format", e.Error())
}

func TestInvalidInputWithOptionalMsg(t *testing.T) {
	t.Parallel()

	e := NewInvalidInput("oopsies")

	assert.Equal(t, e.Error(), "oopsies")
}

func TestPaymentRequiredError(t *testing.T) {
	t.Parallel()

	e := NewPaymentRequiredError()

	assert.Equal(t, "Payment Required", e.Error())
}
