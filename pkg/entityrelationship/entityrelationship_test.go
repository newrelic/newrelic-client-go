//go:build unit
// +build unit

package entityrelationship

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	testSourceEntityGUID = common.EntityGUID("source-entity-guid")
	testTargetEntityGUID = common.EntityGUID("target-entity-guid")
	testRelationshipType = EntityRelationshipEdgeType("MONITORS")

	testCreateOrReplaceResponseJSON = `
	{
		"EntityRelationshipUserDefinedCreateOrReplace": {
			"errors": []
		}
	}`

	testDeleteResponseJSON = `
	{
		"EntityRelationshipUserDefinedDelete": {
			"errors": []
		}
	}`
)

func TestEntityRelationshipUserDefinedCreateOrReplace(t *testing.T) {
	t.Parallel()
	entityRelationship := newMockResponse(t, testCreateOrReplaceResponseJSON, http.StatusOK)

	result, err := entityRelationship.EntityRelationshipUserDefinedCreateOrReplace(
		testSourceEntityGUID,
		testTargetEntityGUID,
		testRelationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Errors)
}

func TestEntityRelationshipUserDefinedCreateOrReplaceWithContext(t *testing.T) {
	t.Parallel()
	entityRelationship := newMockResponse(t, testCreateOrReplaceResponseJSON, http.StatusOK)

	result, err := entityRelationship.EntityRelationshipUserDefinedCreateOrReplaceWithContext(
		context.Background(),
		testSourceEntityGUID,
		testTargetEntityGUID,
		testRelationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Errors)
}

func TestEntityRelationshipUserDefinedUpdate(t *testing.T) {
	t.Parallel()
	entityRelationship := newMockResponse(t, testCreateOrReplaceResponseJSON, http.StatusOK)

	result, err := entityRelationship.EntityRelationshipUserDefinedCreateOrReplace(
		testSourceEntityGUID,
		testTargetEntityGUID,
		testRelationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Errors)
}

func TestEntityRelationshipUserDefinedUpdateWithContext(t *testing.T) {
	t.Parallel()
	entityRelationship := newMockResponse(t, testCreateOrReplaceResponseJSON, http.StatusOK)

	result, err := entityRelationship.EntityRelationshipUserDefinedCreateOrReplaceWithContext(
		context.Background(),
		testSourceEntityGUID,
		testTargetEntityGUID,
		testRelationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Errors)
}

func TestEntityRelationshipUserDefinedDelete(t *testing.T) {
	t.Parallel()
	entityRelationship := newMockResponse(t, testDeleteResponseJSON, http.StatusOK)

	result, err := entityRelationship.EntityRelationshipUserDefinedDelete(
		testSourceEntityGUID,
		testTargetEntityGUID,
		testRelationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Errors)
}

func TestEntityRelationshipUserDefinedDeleteWithContext(t *testing.T) {
	t.Parallel()
	entityRelationship := newMockResponse(t, testDeleteResponseJSON, http.StatusOK)

	result, err := entityRelationship.EntityRelationshipUserDefinedDeleteWithContext(
		context.Background(),
		testSourceEntityGUID,
		testTargetEntityGUID,
		testRelationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Errors)
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Entityrelationship {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}
