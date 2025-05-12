//go:build integration
// +build integration

package entityrelationship

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func newIntegrationTestClient(t *testing.T) Entityrelationship {
	tc := testhelpers.NewIntegrationTestConfig(t)
	return New(tc)
}

func TestIntegrationEntityRelationship_CRUD(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)
	entitiesClient := entities.New(testhelpers.NewIntegrationTestConfig(t))

	// Define test data
	sourceEntityGUID := common.EntityGUID(testhelpers.EntityRelationshipTestSourceEntityGUID)
	targetEntityGUID := common.EntityGUID(testhelpers.EntityRelationshipTestTargetEntityGUID)
	relationshipType := EntityRelationshipEdgeType("TRIGGERS")

	// Create or replace entity relationship (initial creation)
	createResult, err := client.EntityRelationshipUserDefinedCreateOrReplace(
		sourceEntityGUID,
		targetEntityGUID,
		relationshipType,
	)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.Empty(t, createResult.Errors)

	// Duplicate CreateOrReplace to test update functionality
	updateResult, err := client.EntityRelationshipUserDefinedCreateOrReplace(
		sourceEntityGUID,
		targetEntityGUID,
		relationshipType,
	)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Empty(t, updateResult.Errors)

	// Validate the relationship creation by reading it
	entity, err := entitiesClient.GetEntity(sourceEntityGUID)
	require.NoError(t, err)
	require.NotNil(t, entity)

	// Check if the relationship exists in the read result
	relatedEntities := (*entity).(*entities.ExternalEntity).RelatedEntities
	found := false

	for _, relationship := range relatedEntities.Results {
		if userDefinedEdge, ok := relationship.(*entities.EntityRelationshipUserDefinedEdge); ok {
			// Now 'userDefinedEdge' is a pointer to an EntityRelationshipUserDefinedEdge and you can access its fields.
			if userDefinedEdge.Target.GUID == targetEntityGUID && EntityRelationshipEdgeType(userDefinedEdge.Type) == relationshipType {
				found = true
				break
			}
		}
	}
	require.True(t, found, "Expected relationship not found in read operation")

	// Defer block to clean up by deleting the created relationship
	defer func() {
		deleteResult, err := client.EntityRelationshipUserDefinedDelete(
			sourceEntityGUID,
			targetEntityGUID,
			relationshipType,
		)
		require.NoError(t, err)
		require.NotNil(t, deleteResult)
		require.Empty(t, deleteResult.Errors)
	}()

	// Delete the entity relationship
	deleteResult, err := client.EntityRelationshipUserDefinedDelete(
		sourceEntityGUID,
		targetEntityGUID,
		relationshipType,
	)
	require.NotNil(t, deleteResult)
	require.Empty(t, deleteResult.Errors)
}
