//go:build integration
// +build integration

package notifications

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestNotificationMutationDestination_GraphQL(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)

	// Notifications test account
	accountID := 10867072

	// Create a destination to work with in this test
	testIntegrationDestinationNameRandStr := nr.RandSeq(5)
	destination := DestinationInput{}
	destination.Type = DestinationTypes.Webhook
	destination.Properties = []PropertyInput{
		{
			Key:   "url",
			Value: "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b",
		},
	}
	destination.Auth = AiNotificationsCredentialsInput{
		Type: AuthTypes.Token,
		Token: TokenAuth{
			Token:  "Token",
			Prefix: "Bearer",
		},
	}
	destination.Name = fmt.Sprintf("test-notifications-destination-%s", testIntegrationDestinationNameRandStr)

	// Test: Get List
	listResult, err := n.ListDestinations(accountID)
	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Create
	createResult, err := n.CreateDestinationMutation(accountID, destination)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get Destination
	getDestinationResult, err := n.GetDestination(accountID, createResult.ID)
	require.NoError(t, err)
	require.NotNil(t, getDestinationResult)

	// Test: Delete
	deleteResult, err := n.DeleteDestinationMutation(accountID, createResult.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}
