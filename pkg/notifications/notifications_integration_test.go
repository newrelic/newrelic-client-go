//go:build integration
// +build integration

package notifications

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestNotificationMutationDestination(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	// Create a destination to work with in this test
	testIntegrationDestinationNameRandStr := mock.RandSeq(5)
	destination := AiNotificationsDestinationInput{}
	destination.Type = AiNotificationsDestinationTypeTypes.WEBHOOK
	destination.Properties = []AiNotificationsPropertyInput{
		{
			Key:          "url",
			Value:        "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b",
			Label:        "",
			DisplayValue: "",
		},
	}
	destination.Auth = AiNotificationsCredentialsInput{
		Type: AiNotificationsAuthTypeTypes.TOKEN,
		Token: AiNotificationsTokenAuthInput{
			Token:  "Token",
			Prefix: "Bearer",
		},
	}
	destination.Name = fmt.Sprintf("test-notifications-destination-%s", testIntegrationDestinationNameRandStr)

	// Test: Create
	createResult, err := n.AiNotificationsCreateDestination(accountID, destination)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.NotEmpty(t, createResult.Destination.Auth)
	require.Equal(t, ai.AiNotificationsAuthType("TOKEN"), createResult.Destination.Auth.AuthType)

	// Test: Get Destination
	filters := ai.AiNotificationsDestinationFilter{
		ID: createResult.Destination.ID,
	}
	sorter := AiNotificationsDestinationSorter{}
	getDestinationResult, err := n.GetDestinations(accountID, "", filters, sorter)
	require.NoError(t, err)
	require.NotNil(t, getDestinationResult)
	assert.Equal(t, 1, getDestinationResult.TotalCount)

	// Test: Update Destination
	updateDestination := AiNotificationsDestinationUpdate{}
	updateDestination.Active = false
	updateDestination.Properties = []AiNotificationsPropertyInput{
		{
			Key:          "url",
			Value:        "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b",
			Label:        "",
			DisplayValue: "",
		},
	}
	updateDestination.Auth = AiNotificationsCredentialsInput{
		Type: AiNotificationsAuthTypeTypes.TOKEN,
		Token: AiNotificationsTokenAuthInput{
			Token:  "TokenUpdate",
			Prefix: "BearerUpdate",
		},
	}
	updateDestination.Name = fmt.Sprintf("test-notifications-update-destination-%s", testIntegrationDestinationNameRandStr)

	updateDestinationResult, err := n.AiNotificationsUpdateDestination(accountID, updateDestination, createResult.Destination.ID)
	require.NoError(t, err)
	require.NotNil(t, updateDestinationResult)

	// Test: Delete
	deleteResult, err := n.AiNotificationsDeleteDestination(accountID, createResult.Destination.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}
