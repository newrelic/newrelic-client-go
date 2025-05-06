//go:build integration
// +build integration

package keytransaction

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationKeyTransaction_All(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)

	// creating a key transaction
	// this is expected to throw no error, and successfully create the key transaction
	createKeyTransactionTestResult, err := client.KeyTransactionCreate(
		10,
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		10,
		testhelpers.IntegrationTestApplicationEntityNameNew,
		testKeyTransactionName,
	)
	fmt.Println("createKeyTransactionTestResult", createKeyTransactionTestResult)
	fmt.Println("err", err)
	time.Sleep(3 * time.Second)

	require.NoError(t, err)
	require.NotNil(t, createKeyTransactionTestResult)
	require.Equal(t, testKeyTransactionName, createKeyTransactionTestResult.Name)

	// defer block to delete the created key transaction, at the end of execution of this test
	defer func() {
		deletedResult, err := client.KeyTransactionDelete(createKeyTransactionTestResult.GUID)
		require.NoError(t, err)
		require.NotNil(t, deletedResult)
		require.Equal(t, deletedResult.Success, true)
	}()

	// updating the key transaction created
	// this is expected to throw no error, and successfully update the key transaction
	updateKeyTransactionTestResult, err := client.KeyTransactionUpdate(
		11,
		11,
		createKeyTransactionTestResult.GUID,
		testKeyTransactionName+"-updated",
	)
	time.Sleep(3 * time.Second)
	fmt.Println("updateKeyTransactionTestResult", updateKeyTransactionTestResult)
	fmt.Println("err", err)

	require.NoError(t, err)
	require.NotNil(t, updateKeyTransactionTestResult)
	require.Equal(t, updateKeyTransactionTestResult.ApdexTarget, float64(11))
	require.Equal(t, updateKeyTransactionTestResult.BrowserApdexTarget, float64(11))
	require.Equal(t, testKeyTransactionName+"-updated", updateKeyTransactionTestResult.Name)

	// deleting the key transaction created
	// this is expected to throw no error and delete the created key transaction
	deletedResult, err := client.KeyTransactionDelete(createKeyTransactionTestResult.GUID)
	time.Sleep(3 * time.Second)
	fmt.Println("deletedResult", deletedResult)
	fmt.Println("err", err)
	require.NoError(t, err)
	require.NotNil(t, deletedResult)
	require.Equal(t, deletedResult.Success, true)
}

func TestIntegrationKeyTransaction_CreateDuplicateError(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)

	// creating a key transaction
	// this is expected to throw no error, and successfully create the key transaction
	createdResult, err := client.KeyTransactionCreate(
		10,
		testhelpers.IntegrationTestApplicationEntityGUIDNew,
		10,
		testhelpers.IntegrationTestApplicationEntityNameNew,
		testKeyTransactionName,
	)

	require.NoError(t, err)
	require.NotNil(t, createdResult)
	require.Equal(t, testKeyTransactionName, createdResult.Name)

	// defer block to delete the created key transaction, at the end of execution of this test
	defer func() {
		deletedResult, err := client.KeyTransactionDelete(createdResult.GUID)
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
}
