package usermanagement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationGetAllAuthDomains(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)
	authDomains, err := client.GetAuthenticationDomains("", nil)
	require.NoError(t, err)
	require.NotNil(t, authDomains)
}

func TestIntegrationGetAuthDomainWithID(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	authDomains, err := client.GetAuthenticationDomains("", []string{"0cc21d98-8dc2-484a-bb26-258e17ede584"})
	require.NoError(t, err)
	require.NotNil(t, authDomains)
	require.Equal(t, authDomains.TotalCount, 1)
}
func TestIntegrationGetAuthDomainWithInvalidID(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	authDomains, err := client.GetAuthenticationDomains("", []string{"**-8dc2-484a-bb26-258e17ede584"})
	require.NoError(t, err)
	require.NotNil(t, authDomains)
	require.Equal(t, authDomains.TotalCount, 0)
}
