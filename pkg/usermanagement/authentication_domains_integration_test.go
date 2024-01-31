//go:build integration
// +build integration

package usermanagement

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestGetAuthenticationDomains(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetAuthenticationDomains(
		"",
		[]string{""},
	)
	require.NoError(t, err)
	require.NotNil(t, response.TotalCount)
}

func TestGetAuthenticationDomainsNone(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetAuthenticationDomains(
		"",
		[]string{"NOT-A-REAL-ID"},
	)
	require.NoError(t, err)
	require.Zero(t, response.TotalCount)
	require.Zero(t, len(response.AuthenticationDomains))
}

func newIntegrationTestClient(t *testing.T) UserManagement {
	tc := mock.NewIntegrationTestConfig(t)
	return New(tc)
}
