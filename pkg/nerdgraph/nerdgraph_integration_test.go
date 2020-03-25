// +build sometimes

package nerdgraph

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestSchemaQuery(t *testing.T) {
	n := newNerdGraphIntegrationTestClient(t)

	schema, err := n.QuerySchema()
	require.NoError(t, err)

	namesToResolve := []string{
		// "AlertsMutingRuleInput",
		"AlertsPoliciesSearchCriteriaInput",
		// "AlertsMutingRuleConditionInput",
		"AlertsPoliciesSearchResultSet",
	}

	types, err := ResolveSchemaTypes(*schema, namesToResolve)
	require.NoError(t, err)
	// fmt.Printf("Types:\n%s", types)

	f, err := os.Create("types.go")
	require.NoError(t, err)

	_, err = f.WriteString("package nerdgraph\n")
	require.NoError(t, err)

	defer f.Close()

	keys := make([]string, 0, len(types))
	for k := range types {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("\n%s", types[k])
		_, err := f.WriteString(types[k])
		require.NoError(t, err)
	}

}

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
