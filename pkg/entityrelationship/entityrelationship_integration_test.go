//go:build integration
// +build integration

package entityrelationship

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

var (
	sourceEntityGuidVal = "MzgwNjUyNnxFWFR8U0VSVklDRV9MRVZFTHw1ODA4MDM"
	targetEntityGuidVal = "MzgwNjUyNnxFWFR8U0VSVklDRV9MRVZFTHw1NzE0Nzk"
)

func newIntegrationTestClient(t *testing.T) Entityrelationship {
	tc := testhelpers.NewIntegrationTestConfig(t)
	return New(tc)
}

func TestIntegrationEntityRelationship_CreateAndDelete(t *testing.T) {
	t.Parallel()
	client := newIntegrationTestClient(t)

	// Define test data
	sourceEntityGUID := common.EntityGUID(sourceEntityGuidVal)
	targetEntityGUID := common.EntityGUID(targetEntityGuidVal)
	relationshipType := EntityRelationshipEdgeType("TRIGGERS")

	// Create or replace entity relationship
	createResult, err := client.EntityRelationshipUserDefinedCreateOrReplace(
		sourceEntityGUID,
		targetEntityGUID,
		relationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.Empty(t, createResult.Errors)

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

	// Validate the relationship creation
	require.Empty(t, createResult.Errors)

	// Delete the entity relationship
	deleteResult, err := client.EntityRelationshipUserDefinedDelete(
		sourceEntityGUID,
		targetEntityGUID,
		relationshipType,
	)

	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.Empty(t, deleteResult.Errors)
}
