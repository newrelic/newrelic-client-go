//go:build unit
// +build unit

package contextkeys

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSetAccountID_SetsTheAccountID(t *testing.T) {
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

func TestContextKeys_SetAccountID_SetsValues(t *testing.T) {
	x := contextKeys{}
	ctx := x.SetAccountID(context.Background(), "some-account-id")
	result, ok := x.GetAccountID(ctx)

	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", result)
}

func TestContextKeys_SetAccountID_CreatesContextIfNil(t *testing.T) {
	x := contextKeys{}
	ctx := x.SetAccountID(nil, "some-account-id")
	assert.NotNil(t, ctx)

	result, ok := x.GetAccountID(ctx)
	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", result)
}

func TestContextKeys_GetAccountID_ReturnsAccountIDIfFound(t *testing.T) {
	x := contextKeys{}
	testCtx := context.WithValue(context.Background(), keys.accountID, "some-account-id")
	accountID, ok := x.GetAccountID(testCtx)

	assert.Equal(t, true, ok)
	assert.Equal(t, "some-account-id", accountID)
}

func TestContextKeys_GetAccountID_DefaultReturnValueWhenNoAccountIDFound(t *testing.T) {
	x := contextKeys{}
	defaultValue, ok := x.GetAccountID(context.Background())

	assert.Equal(t, false, ok)
	assert.Equal(t, "", defaultValue)
}
