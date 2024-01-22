//go:build integration
// +build integration

package usermanagement

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchAuthenticationDomain(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetAuthenticationDomains(
		"",
		[]string{"84cb286a-8eb0-4478-b469-cdf2ccfef553"},
	)

	fmt.Println(response)
	require.NoError(t, err)
}

func TestFetchGroupByID(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetGroups(
		[]string{"84cb286a-8eb0-4478-b469-cdf2ccfef553"},
		[]string{"f347272d-baa4-4052-b823-618b3bf5748d"},
		"something",
	)

	fmt.Println(response)
	require.NoError(t, err)
}

func TestFetchGroupByName(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetGroups(
		[]string{"84cb286a-8eb0-4478-b469-cdf2ccfef553"},
		[]string{},
		"New Group",
	)

	fmt.Println(response)
	require.NoError(t, err)
}

func TestFetchGroupsWithUsers(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	response, err := client.GetGroupsWithUsers(
		[]string{"84cb286a-8eb0-4478-b469-cdf2ccfef553"},
		[]string{"f347272d-baa4-4052-b823-618b3bf5748d"},
	)

	fmt.Println(response)
	require.NoError(t, err)
}
