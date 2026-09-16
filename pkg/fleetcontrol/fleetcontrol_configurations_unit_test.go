//go:build unit
// +build unit

package fleetcontrol

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testConfigurationVersionsMockResponse mirrors the actual shape returned by
// the blob-service API for AgentConfigurationVersions: entityGuid/blobId are
// camelCase, not entity_guid/blob_id.
const testConfigurationVersionsMockResponse = `{
	"versions": [
		{
			"entityGuid": "test-version-entity-guid",
			"blobId": "test-blob-id",
			"version": "1",
			"timestamp": "2026-09-16T16:05:52.098Z"
		}
	],
	"cursor": null
}`

func TestUnitFleetControlGetConfigurationVersions_ParsesCamelCaseFields(t *testing.T) {
	t.Parallel()

	client := newMockResponse(t, testConfigurationVersionsMockResponse, http.StatusOK)

	result, err := client.FleetControlGetConfigurationVersions("test-config-entity-guid", testOrganizationID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Versions, 1)

	assert.Equal(t, "test-version-entity-guid", result.Versions[0].EntityGUID)
	assert.Equal(t, "test-blob-id", result.Versions[0].BlobID)
	assert.Equal(t, "1", result.Versions[0].Version)
}
