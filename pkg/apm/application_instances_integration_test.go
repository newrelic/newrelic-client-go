//go:build integration
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

	var app *Application
	for _, app = range a {
		if len(app.Links.InstanceIDs) > 0 {
			break
		}
	}

	if len(app.Links.InstanceIDs) == 0 {
		t.Skip("no applications found with instances")
	}

	instanceParams := ListApplicationInstancesParams{
		IDs: []int{app.Links.InstanceIDs[0]},
	}

	instances, err := client.ListApplicationInstances(app.ID, &instanceParams)
	require.NoError(t, err)
	require.Equal(t, 1, len(instances))

	_, err = client.GetApplicationInstance(app.ID, instances[0].ID)
	require.NoError(t, err)
}
