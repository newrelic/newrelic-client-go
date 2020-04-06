// +build unit

package region

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidError(t *testing.T) {
	t.Parallel()

	err := InvalidError{}
	assert.EqualError(t, err, "invalid region")

	err = InvalidError{Message: "asdf"}
	assert.EqualError(t, err, "invalid region: asdf")

	// Custom func for nils
	err = ErrorNil()
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid region: value is nil")
}

func TestUnknownError(t *testing.T) {
	t.Parallel()

	err := UnknownError{}
	assert.EqualError(t, err, "unknown region")

	err = UnknownError{Message: "test"}
	assert.EqualError(t, err, "unknown region: test")
}

func TestUnknownUsingDefaultError(t *testing.T) {
	t.Parallel()

	err := UnknownUsingDefaultError{}
	assert.EqualError(t, err, "unknown region, using default: "+Default.String())

	err = UnknownUsingDefaultError{Message: "test"}
	assert.EqualError(t, err, "unknown region: test, using default: "+Default.String())
}
