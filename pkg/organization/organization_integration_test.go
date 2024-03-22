//go:build integration
// +build integration

package organization

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetOrganizations(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.OrganizationGetOrganizations(
		context.Background(),
	)
	require.NoError(t, err)
	require.NotNil(t, response.Items)
}
