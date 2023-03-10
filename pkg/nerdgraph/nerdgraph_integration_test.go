//go:build integration
// +build integration

package nerdgraph

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationQuery(t *testing.T) {
	t.Parallel()

	client := newNerdGraphIntegrationTestClient(t)

	query := `{
		actor {
			user {
				email
				id
				name
			}
		}
	}`

	variables := map[string]interface{}{}

	actual, err := client.Query(query, variables)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

func TestIntegrationQueryWithVariables(t *testing.T) {
	t.Parallel()

	gqlClient := newNerdGraphIntegrationTestClient(t)

	query := `
		query($guid: EntityGuid!) {
			actor {
				entity(guid: $guid) {
					guid
					name
					domain
					entityType
				}
			}
		}
	`

	variables := map[string]interface{}{
		"guid": testhelpers.IntegrationTestApplicationEntityGUID,
	}

	actual, err := gqlClient.Query(query, variables)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

// nolint
func newNerdGraphIntegrationTestClient(t *testing.T) NerdGraph {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
