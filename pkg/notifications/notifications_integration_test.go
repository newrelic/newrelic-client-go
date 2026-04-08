//go:build integration
// +build integration

package notifications

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
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
	destination.Auth = &AiNotificationsCredentialsInput{
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

	// Test: Get Destination by id
	filters := ai.AiNotificationsDestinationFilter{
		ID: createResult.Destination.ID,
	}
	sorter := AiNotificationsDestinationSorter{}
	getDestinationResult, err := n.GetDestinations(accountID, "", filters, sorter)
	require.NoError(t, err)
	require.NotNil(t, getDestinationResult)
	assert.Equal(t, 1, getDestinationResult.TotalCount)
	require.NotEmpty(t, getDestinationResult.Entities[0].GUID)

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
	updateDestination.Auth = &AiNotificationsCredentialsInput{
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

func TestNotificationMutationDestination_FilterByName(t *testing.T) {
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
	destination.Auth = &AiNotificationsCredentialsInput{
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

	// Test: Get Destination by name
	filtersByName := ai.AiNotificationsDestinationFilter{
		Name: createResult.Destination.Name,
	}
	sorter := AiNotificationsDestinationSorter{}
	getDestinationByNameResult, err := n.GetDestinations(accountID, "", filtersByName, sorter)
	require.NoError(t, err)
	require.NotNil(t, getDestinationByNameResult)
	assert.Equal(t, 1, getDestinationByNameResult.TotalCount)

	// Test: Delete
	deleteResult, err := n.AiNotificationsDeleteDestination(accountID, createResult.Destination.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

func TestNotificationMutationDestination_FilterByExactName(t *testing.T) {
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
	destination.Auth = &AiNotificationsCredentialsInput{
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

	// Test: Get Destination by exact name
	filtersByExactName := ai.AiNotificationsDestinationFilter{
		ExactName: createResult.Destination.Name,
	}
	sorter := AiNotificationsDestinationSorter{}
	getDestinationByExactNameResult, err := n.GetDestinations(accountID, "", filtersByExactName, sorter)
	require.NoError(t, err)
	require.NotNil(t, getDestinationByExactNameResult)
	assert.Equal(t, 1, getDestinationByExactNameResult.TotalCount)

	// Test: Delete
	deleteResult, err := n.AiNotificationsDeleteDestination(accountID, createResult.Destination.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

func TestNotificationMutationDestination_CustomHeaderAuth(t *testing.T) {
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
	destination.Auth = &AiNotificationsCredentialsInput{
		Type: AiNotificationsAuthTypeTypes.CUSTOM_HEADERS,
		CustomHeaders: &AiNotificationsCustomHeadersAuthInput{
			[]AiNotificationsCustomHeaderInput{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
		},
	}
	destination.Name = fmt.Sprintf("test-notifications-destination-%s", testIntegrationDestinationNameRandStr)

	// Test: Create
	createResult, err := n.AiNotificationsCreateDestination(accountID, destination)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.NotEmpty(t, createResult.Destination.Auth)
	require.Equal(t, ai.AiNotificationsAuthType("CUSTOM_HEADERS"), createResult.Destination.Auth.AuthType)
	require.Equal(t, 3, len(createResult.Destination.Auth.CustomHeaders))
	require.Equal(t, "key1", createResult.Destination.Auth.CustomHeaders[0].Key)
	require.Equal(t, "key2", createResult.Destination.Auth.CustomHeaders[1].Key)
	require.Equal(t, "key3", createResult.Destination.Auth.CustomHeaders[2].Key)

	// Test: Get Destination by id
	filters := ai.AiNotificationsDestinationFilter{
		ID: createResult.Destination.ID,
	}
	sorter := AiNotificationsDestinationSorter{}
	getDestinationResult, err := n.GetDestinations(accountID, "", filters, sorter)
	require.NoError(t, err)
	require.NotNil(t, getDestinationResult)
	assert.Equal(t, 1, getDestinationResult.TotalCount)
	require.NotEmpty(t, getDestinationResult.Entities[0].GUID)
	require.Equal(t, ai.AiNotificationsAuthType("CUSTOM_HEADERS"), getDestinationResult.Entities[0].Auth.AuthType)
	require.Equal(t, 3, len(getDestinationResult.Entities[0].Auth.CustomHeaders))
	require.Equal(t, "key1", getDestinationResult.Entities[0].Auth.CustomHeaders[0].Key)
	require.Equal(t, "key2", getDestinationResult.Entities[0].Auth.CustomHeaders[1].Key)
	require.Equal(t, "key3", getDestinationResult.Entities[0].Auth.CustomHeaders[2].Key)

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
	updateDestination.Auth = &AiNotificationsCredentialsInput{
		Type: AiNotificationsAuthTypeTypes.CUSTOM_HEADERS,
		CustomHeaders: &AiNotificationsCustomHeadersAuthInput{
			[]AiNotificationsCustomHeaderInput{
				{Key: "key1", Value: "value1"},
				{Key: "key4", Value: "value4"},
			},
		},
	}
	updateDestination.Name = fmt.Sprintf("test-notifications-update-destination-%s", testIntegrationDestinationNameRandStr)

	updateDestinationResult, err := n.AiNotificationsUpdateDestination(accountID, updateDestination, createResult.Destination.ID)
	require.NoError(t, err)
	require.NotNil(t, updateDestinationResult)
	require.Equal(t, ai.AiNotificationsAuthType("CUSTOM_HEADERS"), updateDestinationResult.Destination.Auth.AuthType)
	require.Equal(t, 2, len(updateDestinationResult.Destination.Auth.CustomHeaders))
	require.Equal(t, "key1", updateDestinationResult.Destination.Auth.CustomHeaders[0].Key)
	require.Equal(t, "key4", updateDestinationResult.Destination.Auth.CustomHeaders[1].Key)

	// Test: Delete
	deleteResult, err := n.AiNotificationsDeleteDestination(accountID, createResult.Destination.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

func TestNotificationMutationDestination_secureUrl(t *testing.T) {
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
	destination.Properties = []AiNotificationsPropertyInput{}
	destination.SecureURL = &AiNotificationsSecureURLInput{
		Prefix:       "https://webhook.site",
		SecureSuffix: "/94193c01-4a81-4782-8f1b-554d5230395b",
	}
	destination.Auth = &AiNotificationsCredentialsInput{
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
	require.NotNil(t, createResult.Destination.SecureURL)
	require.Equal(t, createResult.Destination.SecureURL.Prefix, "https://webhook.site")
	require.NotEmpty(t, createResult.Destination.Auth)
	require.Equal(t, ai.AiNotificationsAuthType("TOKEN"), createResult.Destination.Auth.AuthType)

	// Test: Get Destination by id
	filters := ai.AiNotificationsDestinationFilter{
		ID: createResult.Destination.ID,
	}
	sorter := AiNotificationsDestinationSorter{}
	getDestinationResult, err := n.GetDestinations(accountID, "", filters, sorter)
	require.NoError(t, err)
	require.NotNil(t, getDestinationResult)
	assert.Equal(t, 1, getDestinationResult.TotalCount)
	require.NotEmpty(t, getDestinationResult.Entities[0].GUID)
	require.NotNil(t, getDestinationResult.Entities[0].SecureURL)
	require.Equal(t, getDestinationResult.Entities[0].SecureURL.Prefix, "https://webhook.site")

	// Test: Update Destination
	updateDestination := AiNotificationsDestinationUpdate{}
	updateDestination.Active = false
	updateDestination.Properties = []AiNotificationsPropertyInput{}
	updateDestination.SecureURL = &AiNotificationsSecureURLUpdate{
		Prefix:       "https://webhook2.site",
		SecureSuffix: "/59bb0d7a-1708-481a-a178-9161416f8ba6",
	}
	updateDestination.Auth = &AiNotificationsCredentialsInput{
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
	require.NotNil(t, updateDestinationResult.Destination.SecureURL)
	require.Equal(t, updateDestinationResult.Destination.SecureURL.Prefix, "https://webhook2.site")

	// Test: Delete
	deleteResult, err := n.AiNotificationsDeleteDestination(accountID, createResult.Destination.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

func TestNotificationMutationChannel(t *testing.T) {
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
	destination.Auth = &AiNotificationsCredentialsInput{
		Type: AiNotificationsAuthTypeTypes.TOKEN,
		Token: AiNotificationsTokenAuthInput{
			Token:  "Token",
			Prefix: "Bearer",
		},
	}
	destination.Name = fmt.Sprintf("test-notifications-destination-%s", testIntegrationDestinationNameRandStr)

	// Test: Create Destination
	createDestinationResult, err := n.AiNotificationsCreateDestination(accountID, destination)
	require.NoError(t, err)
	require.NotNil(t, createDestinationResult)

	destinationId := createDestinationResult.Destination.ID

	// Create a channel to work with in this test
	testIntegrationChannelNameRandStr := mock.RandSeq(5)
	channel := AiNotificationsChannelInput{}
	channel.Type = AiNotificationsChannelTypeTypes.WEBHOOK
	channel.Product = AiNotificationsProductTypes.IINT
	channel.Properties = []AiNotificationsPropertyInput{
		{
			Key:          "headers",
			Value:        "{}",
			Label:        "Custom headers",
			DisplayValue: "",
		},
		{
			Key:          "payload",
			Value:        "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
			Label:        "Payload Template",
			DisplayValue: "",
		},
	}
	channel.DestinationId = destinationId
	channel.Name = fmt.Sprintf("test-notifications-channel-%s", testIntegrationChannelNameRandStr)

	// Test: Create Channel
	createResult, err := n.AiNotificationsCreateChannel(accountID, channel)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get Channel
	filters := ai.AiNotificationsChannelFilter{
		ID: createResult.Channel.ID,
	}
	sorter := AiNotificationsChannelSorter{}

	getChannelResult, err := n.GetChannels(accountID, "", filters, sorter)
	require.NoError(t, err)
	require.NotNil(t, getChannelResult)
	assert.Equal(t, 1, getChannelResult.TotalCount)

	// Test: Update Channel
	updateChannel := AiNotificationsChannelUpdate{}
	updateChannel.Active = false
	updateChannel.Properties = []AiNotificationsPropertyInput{
		{
			Key:          "headers",
			Value:        "{}",
			Label:        "Custom headers",
			DisplayValue: "",
		},
		{
			Key:          "payload",
			Value:        "{\\n\\t\\\"id\\\": \\\"test-update\\\"\\n}",
			Label:        "Payload Template",
			DisplayValue: "",
		},
	}
	updateChannel.Name = fmt.Sprintf("test-notifications-update-channel-%s", testIntegrationChannelNameRandStr)

	updateChannelResult, err := n.AiNotificationsUpdateChannel(accountID, updateChannel, createResult.Channel.ID)
	require.NoError(t, err)
	require.NotNil(t, updateChannelResult)

	// Test: Delete Channel
	deleteResult, err := n.AiNotificationsDeleteChannel(accountID, createResult.Channel.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)

	// Test: Delete Destination
	deleteDestinationResult, err := n.AiNotificationsDeleteDestination(accountID, destinationId)
	require.NoError(t, err)
	require.NotNil(t, deleteDestinationResult)
}

// ---------------------------------------------------------------------------
// Scope-based destination tests
// ---------------------------------------------------------------------------

// newAccountScope is a helper that builds an ACCOUNT-type EntityScopeInput from an int account ID.
// It mirrors the migration pattern documented on CreateDestinationWithScope / GetDestinationsWithScope.
func newAccountScope(accountID int) *EntityScopeInput {
	return &EntityScopeInput{
		Type: EntityScopeTypeInputTypes.ACCOUNT,
		ID:   strconv.Itoa(accountID),
	}
}

// newWebhookDestinationInput is a helper that returns a minimal webhook destination input
// with TOKEN auth, used across scope-based tests.
func newWebhookDestinationInput(nameSuffix string) AiNotificationsDestinationInput {
	return AiNotificationsDestinationInput{
		Type: AiNotificationsDestinationTypeTypes.WEBHOOK,
		Name: fmt.Sprintf("test-notifications-destination-%s", nameSuffix),
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "url",
				Value: "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b",
			},
		},
		Auth: &AiNotificationsCredentialsInput{
			Type: AiNotificationsAuthTypeTypes.TOKEN,
			Token: AiNotificationsTokenAuthInput{
				Token:  "Token",
				Prefix: "Bearer",
			},
		},
	}
}

// TestNotificationScopedDestination_AccountScope exercises the full Create → Get → Update → Delete
// lifecycle using the scope-based API with an ACCOUNT scope.
//
// This is the scope-API equivalent of TestNotificationMutationDestination. Customers who
// previously passed accountId directly should use EntityScopeTypeInputTypes.ACCOUNT with their
// account ID as the scope ID string.
func TestNotificationScopedDestination_AccountScope(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	scope := newAccountScope(accountID)
	randStr := mock.RandSeq(5)
	destination := newWebhookDestinationInput(randStr)

	// Test: Create
	createResult, err := n.CreateDestinationWithScope(destination, scope)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.NotEmpty(t, createResult.Destination.Auth)
	require.Equal(t, ai.AiNotificationsAuthType("TOKEN"), createResult.Destination.Auth.AuthType)
	// Scope must be returned on the created destination
	require.NotNil(t, createResult.Destination.Scope)
	assert.Equal(t, EntityScopeTypeInputTypes.ACCOUNT, createResult.Destination.Scope.Type)
	assert.Equal(t, strconv.Itoa(accountID), createResult.Destination.Scope.ID)

	// Test: Get by ID using scope
	filters := ai.AiNotificationsDestinationFilter{ID: createResult.Destination.ID}
	sorter := AiNotificationsDestinationSorter{}
	getResult, err := n.GetDestinationsWithScope("", filters, sorter, scope)
	require.NoError(t, err)
	require.NotNil(t, getResult)
	assert.Equal(t, 1, getResult.TotalCount)
	require.NotEmpty(t, getResult.Entities[0].GUID)
	// Scope must be populated on each returned entity
	require.NotNil(t, getResult.Entities[0].Scope)
	assert.Equal(t, EntityScopeTypeInputTypes.ACCOUNT, getResult.Entities[0].Scope.Type)

	// Test: Update using scope
	updateDestination := AiNotificationsDestinationUpdate{
		Active: false,
		Name:   fmt.Sprintf("test-notifications-update-destination-%s", randStr),
		Properties: []AiNotificationsPropertyInput{
			{
				Key:   "url",
				Value: "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b",
			},
		},
		Auth: &AiNotificationsCredentialsInput{
			Type: AiNotificationsAuthTypeTypes.TOKEN,
			Token: AiNotificationsTokenAuthInput{
				Token:  "TokenUpdate",
				Prefix: "BearerUpdate",
			},
		},
	}
	updateResult, err := n.UpdateDestinationWithScope(createResult.Destination.ID, updateDestination, scope)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.NotNil(t, updateResult.Destination.Scope)
	assert.Equal(t, EntityScopeTypeInputTypes.ACCOUNT, updateResult.Destination.Scope.Type)
	assert.Equal(t, fmt.Sprintf("test-notifications-update-destination-%s", randStr), updateResult.Destination.Name)

	// Test: Delete using scope
	deleteResult, err := n.DeleteDestinationWithScope(createResult.Destination.ID, scope)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	assert.Contains(t, deleteResult.IDs, createResult.Destination.ID)
}

// TestNotificationScopedDestination_AccountScope_FilterByName verifies that
// GetDestinationsWithScope correctly filters by destination name under an ACCOUNT scope.
func TestNotificationScopedDestination_AccountScope_FilterByName(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	scope := newAccountScope(accountID)
	randStr := mock.RandSeq(5)
	destination := newWebhookDestinationInput(randStr)

	// Test: Create
	createResult, err := n.CreateDestinationWithScope(destination, scope)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get by name using scope
	filtersByName := ai.AiNotificationsDestinationFilter{Name: createResult.Destination.Name}
	sorter := AiNotificationsDestinationSorter{}
	getResult, err := n.GetDestinationsWithScope("", filtersByName, sorter, scope)
	require.NoError(t, err)
	require.NotNil(t, getResult)
	assert.Equal(t, 1, getResult.TotalCount)
	assert.Equal(t, createResult.Destination.Name, getResult.Entities[0].Name)

	// Test: Delete using scope
	deleteResult, err := n.DeleteDestinationWithScope(createResult.Destination.ID, scope)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

// TestNotificationScopedDestination_AccountScope_FilterByExactName verifies that
// GetDestinationsWithScope correctly filters by exact destination name under an ACCOUNT scope.
func TestNotificationScopedDestination_AccountScope_FilterByExactName(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	scope := newAccountScope(accountID)
	randStr := mock.RandSeq(5)
	destination := newWebhookDestinationInput(randStr)

	// Test: Create
	createResult, err := n.CreateDestinationWithScope(destination, scope)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Get by exact name using scope
	filtersByExactName := ai.AiNotificationsDestinationFilter{ExactName: createResult.Destination.Name}
	sorter := AiNotificationsDestinationSorter{}
	getResult, err := n.GetDestinationsWithScope("", filtersByExactName, sorter, scope)
	require.NoError(t, err)
	require.NotNil(t, getResult)
	assert.Equal(t, 1, getResult.TotalCount)
	assert.Equal(t, createResult.Destination.Name, getResult.Entities[0].Name)

	// Test: Delete using scope
	deleteResult, err := n.DeleteDestinationWithScope(createResult.Destination.ID, scope)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}

// TestNotificationScopedDestination_NilScopeReturnsError verifies that all four scope-based
// operations return a clear error immediately when scope is nil, without making a network call.
func TestNotificationScopedDestination_NilScopeReturnsError(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)

	_, err := n.CreateDestinationWithScope(AiNotificationsDestinationInput{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scope is required")

	_, err = n.GetDestinationsWithScope("", ai.AiNotificationsDestinationFilter{}, AiNotificationsDestinationSorter{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scope is required")

	_, err = n.UpdateDestinationWithScope("someId", AiNotificationsDestinationUpdate{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scope is required")

	_, err = n.DeleteDestinationWithScope("someId", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scope is required")
}

// TestNotificationScopedDestination_InvalidAccountScopeIDReturnsError verifies that
// GetDestinationsWithScope returns a descriptive error when the scope ID cannot be parsed
// as an integer for ACCOUNT scope (the NerdGraph account query requires a numeric ID).
func TestNotificationScopedDestination_InvalidAccountScopeIDReturnsError(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)

	invalidScope := &EntityScopeInput{
		Type: EntityScopeTypeInputTypes.ACCOUNT,
		ID:   "not-a-number",
	}

	_, err := n.GetDestinationsWithScope("", ai.AiNotificationsDestinationFilter{}, AiNotificationsDestinationSorter{}, invalidScope)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid account scope ID")
	assert.Contains(t, err.Error(), "not-a-number")
}
