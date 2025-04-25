//go:build integration
// +build integration

package keytransaction

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationKeyTransaction_All(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)

	// cleanup task prior to the creation of the integration test
	entitiesClient := newIntegrationTestClient_Entities(t)

	testAccountID, _ := testhelpers.GetTestAccountID()
	query := fmt.Sprintf("type = 'KEY_TRANSACTION' and accountId = '%d'", testAccountID)

	entities, err := entitiesClient.GetEntitySearchByQuery(
		entities.EntitySearchOptions{},
		query,
		[]entities.EntitySearchSortCriteria{},
	)

	for _, entity := range entities.Results.Entities {
		keyTransactionName := entity.GetName()
		if keyTransactionName != "" && strings.Contains(keyTransactionName, "nr-test") {
			_, deleteErr := client.KeyTransactionDelete(EntityGUID(string(entity.GetGUID())))
			fmt.Println("Deleting key transaction ", entity.GetName())
			if deleteErr != nil {
				fmt.Printf("Error deleting key transaction %s: %v\n", entity.GetName(), err)
			}
		}
	}

	// creating a key transaction
	// this is expected to throw no error, and successfully create the key transaction
	createKeyTransactionTestResult, err := client.KeyTransactionCreate(
		10,
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		10,
		testhelpers.IntegrationTestApplicationEntityNameNew,
		testKeyTransactionName,
	)

	require.NoError(t, err)
	require.NotNil(t, createKeyTransactionTestResult)
	require.Equal(t, testKeyTransactionName, createKeyTransactionTestResult.Name)

	// defer block to delete the created key transaction, at the end of execution of this test
	// the cleanup task in the beginning of the test should do this, but we've noticed
	// way too much flaky behaviour and unnecessary failures, hence adding this as a double check
	defer func() {
		deletedResult, err := client.KeyTransactionDelete(createKeyTransactionTestResult.GUID)
		require.NoError(t, err)
		require.NotNil(t, deletedResult)
		require.Equal(t, deletedResult.Success, true)
	}()

	// attempt to create the same key transaction again
	// this is expected to throw an error, as multiple key transactions cannot be created with the same metricName
	_, err = client.KeyTransactionCreate(
		10,
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		10,
		testhelpers.IntegrationTestApplicationEntityNameNew,
		testKeyTransactionName,
	)

	require.Error(t, err)

	// updating the key transaction created
	// this is expected to throw no error, and successfully update the key transaction
	updateKeyTransactionTestResult, err := client.KeyTransactionUpdate(
		11,
		11,
		createKeyTransactionTestResult.GUID,
		testKeyTransactionName+"-updated",
	)

	require.NoError(t, err)
	require.NotNil(t, updateKeyTransactionTestResult)
	require.Equal(t, updateKeyTransactionTestResult.ApdexTarget, float64(11))
	require.Equal(t, updateKeyTransactionTestResult.BrowserApdexTarget, float64(11))
	require.Equal(t, testKeyTransactionName+"-updated", updateKeyTransactionTestResult.Name)

	// deleting the key transaction created
	// this is expected to throw no error and delete the created key transaction
	deletedResult, err := client.KeyTransactionDelete(createKeyTransactionTestResult.GUID)
	require.NoError(t, err)
	require.NotNil(t, deletedResult)
	require.Equal(t, deletedResult.Success, true)
}
