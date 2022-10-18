//go:build integration
// +build integration

package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationUsers(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: Get
	user, err := client.GetUser()

	require.NoError(t, err)
	require.NotNil(t, user)

	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Email)
	assert.Greater(t, user.ID, 0)
}

func newIntegrationTestClient(t *testing.T) Users {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
