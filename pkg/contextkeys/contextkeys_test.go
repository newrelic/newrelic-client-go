//go:build unit
// +build unit

package contextkeys

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetXAccountID_ReturnsAccountIDIfFound(t *testing.T) {
	testCtx := context.WithValue(context.Background(), keys.xAccountID, "some-account-id")
	xAccountID, ok := GetXAccountID(testCtx)

	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", xAccountID)
}

func TestGetXAccountID_DefaultReturnValueWhenNoAccountIDFound(t *testing.T) {
	defaultValue, ok := GetXAccountID(context.Background())

	assert.Equal(t, false, ok)
	assert.Equal(t, "", defaultValue)
}

func TestSetXAccountID(t *testing.T) {
	ctx := SetXAccountID(context.Background(), "some-account-id")
	result, ok := GetXAccountID(ctx)

	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", result)
}

func TestSetXAccountID_CreatesContextIfNil(t *testing.T) {
	ctx := SetXAccountID(nil, "some-account-id")
	assert.NotNil(t, ctx)

	result, ok := GetXAccountID(ctx)
	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", result)
}
