// +build integration

package entities

import (
	"fmt"
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSearchEntities(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	params := SearchEntitiesParams{
		Name: "Dummy App",
	}

	actual, err := client.SearchEntities(params)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

func TestIntegrationSearchEntitiesRaw(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	params := SearchEntitiesParams{
		Name: "Dummy App",
	}

	actual, err := client.SearchEntitiesRaw(params, []string{"guid", "name"})

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)

	for _, entity := range actual {
		// Ensure these fields are populated since we
		// specified these fields to be returned in the query
		require.NotEmpty(t, entity.GUID)
		require.NotEmpty(t, entity.Name)

		// Ensure this field is not returned populated
		// since we didn't specify this field to be returned
		require.Empty(t, entity.Domain)
	}
}

func TestIntegrationSearchEntitiesRawQuery(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	params := SearchEntitiesParams{
		Name:   "Dummy App",
		Domain: "BROWSER",
	}

	graphqlQuery := `
		query($queryBuilder: EntitySearchQueryBuilder, $cursor: String) {
				actor {
						entitySearch(queryBuilder: $queryBuilder)  {
								results(cursor: $cursor) {
										nextCursor
										entities {
											guid
											name
										}
								}
						}
				}
		}
	`
	var nextCursor *string
	variables := map[string]interface{}{
		"queryBuilder": params,
		"cursor":       nextCursor,
	}

	actual, err := client.SearchEntitiesRawQuery(graphqlQuery, variables)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)

	for _, entity := range actual {
		// Ensure these fields are populated since we
		// specified these fields to be returned in the query
		require.NotEmpty(t, entity.GUID)
		require.NotEmpty(t, entity.Name)

		// Ensure this field is not returned populated
		// since we didn't specify this field to be returned
		require.Empty(t, entity.Domain)
	}
}

func TestIntegrationSearchEntitiesByTags(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	params := SearchEntitiesParams{
		Tags: &TagValue{
			Key:   "language",
			Value: "nodejs",
		},
	}

	actual, err := client.SearchEntities(params)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

func TestIntegrationGetEntities(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	guids := []string{"MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"}
	actual, err := client.GetEntities(guids)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

func TestIntegrationGetEntity(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	actual, err := client.GetEntity("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")

	require.NoError(t, err)
	require.NotNil(t, actual)
}

// nolint
func newIntegrationTestClient(t *testing.T) Entities {
	personalAPIKey := os.Getenv("NEWRELIC_PERSONAL_API_KEY")

	if personalAPIKey == "" {
		t.Skipf("acceptance testing for graphql requires your personal API key")
	}

	return New(config.Config{
		PersonalAPIKey: personalAPIKey,
		UserAgent:      "newrelic/newrelic-client-go",
		LogLevel:       "debug",
	})
}
