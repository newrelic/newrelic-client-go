// +build integration

package apm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationApplicationInstances(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	a, err := client.ListApplications(nil)
	require.NoError(t, err)

	_, err = client.GetApplication(a[0].ID)
	require.NoError(t, err)

	instanceParams := ListApplicationInstancesParams{
		IDs: []int{a[0].Links.InstanceIDs[0]},
	}

	instances, err := client.ListApplicationInstances(a[0].ID, &instanceParams)
	require.NoError(t, err)
	require.Equal(t, 1, len(instances))

	_, err = client.GetApplicationInstance(a[0].ID, instances[0].ID)
	require.NoError(t, err)
}
