//go:build integration
// +build integration

package usermanagement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAuthenticationDomains(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetAuthenticationDomains(
		"",
		[]string{authenticationDomainId},
	)
	require.NoError(t, err)
	require.NotNil(t, response.TotalCount)
}

func TestGetAuthenticationDomainsError(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetAuthenticationDomains(
		"",
		[]string{mockAuthenticationDomainId},
	)
	require.NoError(t, err)
	require.Zero(t, response.TotalCount)
	require.Zero(t, len(response.AuthenticationDomains))
}
