//go:build unit
// +build unit

package entityrelationship

import (
	"context"
	"net/http"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/stretchr/testify/require"
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
	alerts := newMockResponse(t, testCreateOrReplaceResponseJSON, http.StatusOK)

	result, err := alerts.EntityRelationshipUserDefinedCreateOrReplace(
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
	alerts := newMockResponse(t, testCreateOrReplaceResponseJSON, http.StatusOK)

	result, err := alerts.EntityRelationshipUserDefinedCreateOrReplaceWithContext(
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
	alerts := newMockResponse(t, testDeleteResponseJSON, http.StatusOK)

	result, err := alerts.EntityRelationshipUserDefinedDelete(
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
	alerts := newMockResponse(t, testDeleteResponseJSON, http.StatusOK)

	result, err := alerts.EntityRelationshipUserDefinedDeleteWithContext(
		context.Background(),
		testSourceEntityGUID,
		testTargetEntityGUID,
		testRelationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Errors)
}
