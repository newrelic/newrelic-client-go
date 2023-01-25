//go:build integration
// +build integration

package usermanagement

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationGroup(t *testing.T) {
	t.Parallel()

	var (
		rand                 = mock.RandSeq(5)
		testName             = "test_" + rand
		id                   = "0cc21d98-8dc2-484a-bb26-258e17ede584"
		testCreateGroupInput = UserManagementCreateGroup{
			AuthenticationDomainId: id,
			DisplayName:            testName,
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.UserManagementCreateGroup(testCreateGroupInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)
	require.Equal(t, created.Group.DisplayName, testName)

	// Test: Update
	updated, err := client.UserManagementUpdateGroup(UserManagementUpdateGroup{
		ID:          created.Group.ID,
		DisplayName: testName + "_update",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotEmpty(t, updated)
	require.Equal(t, updated.Group.DisplayName, testName+"_update")

	//Test delete
	deleted, err := client.UserManagementDeleteGroup(UserManagementDeleteGroup{
		ID: created.Group.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.NotEmpty(t, deleted)
}
