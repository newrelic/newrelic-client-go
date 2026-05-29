//go:build integration
// +build integration

package keytransaction

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationKeyTransaction_All(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)

	// Pre-test cleanup: delete any leftover key transactions from previous runs.
	entitiesClient := newIntegrationTestClient_Entities(t)
	testAccountID, _ := testhelpers.GetTestAccountID()
	query := fmt.Sprintf("type = 'KEY_TRANSACTION' and accountId = '%d'", testAccountID)

	existingEntities, err := entitiesClient.GetEntitySearchByQuery(
		entities.EntitySearchOptions{},
		query,
		[]entities.EntitySearchSortCriteria{},
	)
	if err == nil {
		for _, entity := range existingEntities.Results.Entities {
			name := entity.GetName()
			if name != "" && strings.Contains(name, "nr-test") {
				fmt.Println("Deleting leftover key transaction:", name)
				if _, deleteErr := client.KeyTransactionDelete(EntityGUID(string(entity.GetGUID()))); deleteErr != nil {
					fmt.Printf("Error deleting leftover key transaction %s: %v\n", name, deleteErr)
				}
			}
		}
	}

	// Register best-effort cleanup before creating, so a mid-test failure
	// never leaves a dangling key transaction.
	var createdGUID EntityGUID
	deleted := false
	defer func() {
		if deleted || createdGUID == "" {
			return
		}
		_, _ = client.KeyTransactionDelete(createdGUID) // best-effort cleanup
	}()

	// Create
	createResult, err := client.KeyTransactionCreate(
		10,
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		10,
		testhelpers.IntegrationTestApplicationEntityNameNew,
		testKeyTransactionName,
	)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.Equal(t, testKeyTransactionName, createResult.Name)

	createdGUID = createResult.GUID

	// The KeyTransactionUpdate mutation resolves a deeply nested entity including
	// serviceLevel.indicators. On a brand-new entity the NerdGraph backend needs
	// time to index those fields — without this sleep the update reliably hits
	// "An error occurred resolving this field" under CI load and exhausts all
	// retries. Five seconds matches the pattern used in the servicelevel tests.
	fmt.Println("waiting 5 seconds for key transaction entity to be indexed...")
	time.Sleep(5 * time.Second)

	// Duplicate create — expected to fail (same metricName)
	_, err = client.KeyTransactionCreate(
		10,
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		10,
		testhelpers.IntegrationTestApplicationEntityNameNew,
		testKeyTransactionName,
	)
	require.Error(t, err)

	// Update
	updateResult, err := client.KeyTransactionUpdate(
		11,
		11,
		createdGUID,
		testKeyTransactionName+"-updated",
	)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, float64(11), updateResult.ApdexTarget)
	require.Equal(t, float64(11), updateResult.BrowserApdexTarget)
	require.Equal(t, testKeyTransactionName+"-updated", updateResult.Name)

	// Delete
	deleteResult, err := client.KeyTransactionDelete(createdGUID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.Equal(t, true, deleteResult.Success)
	deleted = true
}
