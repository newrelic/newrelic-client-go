//go:build unit
// +build unit

package contextkeys

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAccountID_ReturnsAccountIDIfFound(t *testing.T) {
	testCtx := context.WithValue(context.Background(), keys.accountID, "some-account-id")
	accountID, ok := GetAccountID(testCtx)

	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", accountID)
}

func TestGetAccountID_DefaultReturnValueWhenNoAccountIDFound(t *testing.T) {
	defaultValue, ok := GetAccountID(context.Background())

	assert.Equal(t, false, ok)
	assert.Equal(t, "", defaultValue)
}

func TestSetAccountID(t *testing.T) {
	ctx := SetAccountID(context.Background(), "some-account-id")
	result, ok := GetAccountID(ctx)

	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", result)
}

func TestSetAccountID_CreatesContextIfNil(t *testing.T) {
	ctx := SetAccountID(nil, "some-account-id")
	assert.NotNil(t, ctx)

	result, ok := GetAccountID(ctx)
	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", result)
}
