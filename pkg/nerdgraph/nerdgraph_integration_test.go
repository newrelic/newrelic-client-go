// +build integration

package nerdgraph

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
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
	personalAPIKey := os.Getenv("NEWRELIC_PERSONAL_API_KEY")

	if personalAPIKey == "" {
		t.Skipf("acceptance testing for NerdGraph requires your personal API key")
	}

	return New(config.Config{
		PersonalAPIKey: personalAPIKey,
		UserAgent:      "newrelic/newrelic-client-go",
	})
}
