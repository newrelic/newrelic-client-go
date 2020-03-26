// +build sometimes

package nerdgraph

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/config"
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
		"guid": "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1",
	}

	actual, err := gqlClient.Query(query, variables)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

// nolint
func newNerdGraphIntegrationTestClient(t *testing.T) NerdGraph {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing for NerdGraph requires NEW_RELIC_API_KEY to be set")
	}

	return New(config.Config{
		PersonalAPIKey: apiKey,
		UserAgent:      "newrelic/newrelic-client-go",
	})
}
