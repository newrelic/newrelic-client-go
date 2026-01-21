//go:build integration
// +build integration

package fleetcontrol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCreateBlob(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createBlobResponse, err := client.FleetControlCreateBlob(
		"random body",
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": "{\"name\": \"Random Build v2\", \"agentType\": \"NRInfra\", \"managedEntityType\": \"KUBERNETESCLUSTER\"}",
			},
		},
		"fb33fea3-4d7e-4736-9701-acb59a634fdf",
	)

	require.NoError(t, err)
	fmt.Println(createBlobResponse)
	// require.NotNil(t, createUserResponse.CreatedUser.ID)
}
