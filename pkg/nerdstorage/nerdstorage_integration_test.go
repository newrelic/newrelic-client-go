//go:build integration
// +build integration

package nerdstorage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	testPackageID           = "8e57e72a-e3ac-4272-9bb8-aea1d71dde3d"
	testCollection          = "myCol"
	testDocumentID          = "myDoc"
	testAlternateDocumentID = "myOtherDoc"
	testDocument            = struct {
		MyField string
	}{
		MyField: "myVal",
	}
	testGetDocumentInput = GetDocumentInput{
		PackageID:  testPackageID,
		Collection: testCollection,
		DocumentID: testDocumentID,
	}
	testGetCollectionInput = GetCollectionInput{
		PackageID:  testPackageID,
		Collection: testCollection,
	}
	testWriteInput = WriteDocumentInput{
		PackageID:  testPackageID,
		Collection: testCollection,
		DocumentID: testDocumentID,
		Document:   testDocument,
	}
	testDeleteDocumentInput = DeleteDocumentInput{
		PackageID:  testPackageID,
		Collection: testCollection,
		DocumentID: testDocumentID,
	}
	testDeleteCollectionInput = DeleteCollectionInput{
		PackageID:  testPackageID,
		Collection: testCollection,
	}
)

func TestIntegrationNerdStorageWithAccountScope(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	document, err := client.WriteDocumentWithAccountScope(testAccountID, testWriteInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	testAlternateWriteInput := testWriteInput
	testAlternateWriteInput.DocumentID = testAlternateDocumentID

	document, err = client.WriteDocumentWithAccountScope(testAccountID, testAlternateWriteInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	collection, err := client.GetCollectionWithAccountScope(testAccountID, testGetCollectionInput)
	require.NoError(t, err)
	require.NotNil(t, collection)

	document, err = client.GetDocumentWithAccountScope(testAccountID, testGetDocumentInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	ok, err := client.DeleteDocumentWithAccountScope(testAccountID, testDeleteDocumentInput)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = client.DeleteCollectionWithAccountScope(testAccountID, testDeleteCollectionInput)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestIntegrationNerdStorageWithUserScope(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	document, err := client.WriteDocumentWithUserScope(testWriteInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	testAlternateWriteInput := testWriteInput
	testAlternateWriteInput.DocumentID = testAlternateDocumentID

	document, err = client.WriteDocumentWithUserScope(testAlternateWriteInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	collection, err := client.GetCollectionWithUserScope(testGetCollectionInput)
	require.NoError(t, err)
	require.NotNil(t, collection)

	document, err = client.GetDocumentWithUserScope(testGetDocumentInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	ok, err := client.DeleteDocumentWithUserScope(testDeleteDocumentInput)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = client.DeleteCollectionWithUserScope(testDeleteCollectionInput)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestIntegrationNerdStorageWithEntityScope(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	document, err := client.WriteDocumentWithEntityScope(testhelpers.IntegrationTestApplicationEntityGUIDNew, testWriteInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	testAlternateWriteInput := testWriteInput
	testAlternateWriteInput.DocumentID = testAlternateDocumentID

	document, err = client.WriteDocumentWithEntityScope(testhelpers.IntegrationTestApplicationEntityGUIDNew, testAlternateWriteInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	collection, err := client.GetCollectionWithEntityScope(testhelpers.IntegrationTestApplicationEntityGUIDNew, testGetCollectionInput)
	require.NoError(t, err)
	require.NotNil(t, collection)

	document, err = client.GetDocumentWithEntityScope(testhelpers.IntegrationTestApplicationEntityGUIDNew, testGetDocumentInput)
	require.NoError(t, err)
	require.NotNil(t, document)

	ok, err := client.DeleteDocumentWithEntityScope(testhelpers.IntegrationTestApplicationEntityGUIDNew, testDeleteDocumentInput)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = client.DeleteCollectionWithEntityScope(testhelpers.IntegrationTestApplicationEntityGUIDNew, testDeleteCollectionInput)
	require.NoError(t, err)
	require.True(t, ok)
}

func newIntegrationTestClient(t *testing.T) NerdStorage {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
