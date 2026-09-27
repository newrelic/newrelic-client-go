package entities

import (
	"context"
)

// maxEntitySearchPages caps the number of paginated entity search pages fetched
// by GetAllEntitySearchGUIDsByQueryWithContext. At 200 entities per page this
// allows up to 10,000 tag-matched entities per team — far beyond any realistic
// production scenario — while bounding the number of API calls to prevent
// runaway performance on plans for very broadly-tagged teams.
const maxEntitySearchPages = 50

// GetAllEntitySearchGUIDsByQueryWithContext pages through results of an entity
// search query, following nextCursor up to maxEntitySearchPages pages, and
// returns the accumulated entity GUIDs. The second return value is true when
// the result set was truncated (i.e. more pages exist beyond the cap).
//
// Use this instead of GetEntitySearchByQueryWithContext when the result set may
// exceed 200 (the per-page default). Callers should surface a note to users
// when truncated=true so they know the list is incomplete.
func (a *Entities) GetAllEntitySearchGUIDsByQueryWithContext(
	ctx context.Context,
	query string,
) (guids []string, truncated bool, err error) {
	// cursor must be nil (not "") on the first call — the API rejects empty
	// string as "Invalid cursor value". Subsequent pages use the string returned
	// by nextCursor.
	var cursor interface{} = nil

	for page := 0; page < maxEntitySearchPages; page++ {
		resp := entitySearchResponse{}
		vars := map[string]interface{}{
			"query":  query,
			"cursor": cursor,
		}
		if err = a.client.NerdGraphQueryWithContext(ctx, getEntitySearchByQueryWithCursor, vars, &resp); err != nil {
			return nil, false, err
		}
		for _, e := range resp.Actor.EntitySearch.Results.Entities {
			guids = append(guids, string(e.GetGUID()))
		}
		next := resp.Actor.EntitySearch.Results.NextCursor
		if next == "" {
			return guids, false, nil
		}
		cursor = next
	}

	// Reached the page cap — more results exist but were not fetched.
	return guids, true, nil
}

const getEntitySearchByQueryWithCursor = `query(
	$query: String,
	$cursor: String,
) { actor { entitySearch(
	query: $query,
) {
	results(cursor: $cursor) {
		entities {
			guid
		}
		nextCursor
	}
} } }`

// Search for entities using a custom query.
// For more details on how to create a custom query
// and what entity data you can request, visit our
// [entity docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/use-new-relic-graphql-api-query-entities).
//
// Note: you must supply either a `query` OR a `queryBuilder` argument, not both.
func (a *Entities) GetEntitySearchByQuery(
	options EntitySearchOptions,
	query string,
	sortBy []EntitySearchSortCriteria,
) (*EntitySearch, error) {
	return a.GetEntitySearchByQueryWithContext(context.Background(),
		options,
		query,
		sortBy,
	)
}

// Search for entities using a custom query.
//
// For more details on how to create a custom query
// and what entity data you can request, visit our
// [entity docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/use-new-relic-graphql-api-query-entities).
//
// Note: you must supply either a `query` OR a `queryBuilder` argument, not both.
func (a *Entities) GetEntitySearchByQueryWithContext(
	ctx context.Context,
	options EntitySearchOptions,
	query string,
	sortBy []EntitySearchSortCriteria,
) (*EntitySearch, error) {

	resp := entitySearchResponse{}
	vars := map[string]interface{}{
		"options": options,
		"query":   query,
		"sortBy":  sortBy,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getEntitySearchByQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.EntitySearch, nil
}

const getEntitySearchByQuery = `query(
	$query: String,
) { actor { entitySearch(
	query: $query,
) {
	count
	query
	results {
		entities {
			__typename
			accountId
			alertSeverity
			domain
			entityType
			guid
			indexedAt
			name
			permalink
			reporting
			tags {
				key
				values
			}
			type
			... on ApmApplicationEntityOutline {
				__typename
				applicationId
				language
			}
			... on ApmDatabaseInstanceEntityOutline {
				__typename
				host
				portOrPath
				vendor
			}
			... on ApmExternalServiceEntityOutline {
				__typename
				host
			}
			... on BrowserApplicationEntityOutline {
				__typename
				agentInstallType
				applicationId
				servingApmApplicationId
			}
			... on DashboardEntityOutline {
				__typename
				createdAt
				dashboardParentGuid
				permissions
				updatedAt
			}
			... on ExternalEntityOutline {
				__typename
			}
			... on GenericEntityOutline {
				__typename
				tags {
					key
					values
				}
			}
			... on GenericInfrastructureEntityOutline {
				__typename
				integrationTypeCode
			}
			... on InfrastructureAwsLambdaFunctionEntityOutline {
				__typename
				integrationTypeCode
				runtime
			}
			... on InfrastructureHostEntityOutline {
				__typename
			}
			... on MobileApplicationEntityOutline {
				__typename
				applicationId
			}
			... on SecureCredentialEntityOutline {
				__typename
				description
				secureCredentialId
				updatedAt
			}
			... on SyntheticMonitorEntityOutline {
				__typename
				monitorId
				monitorType
				monitoredUrl
				period
			}
			... on ThirdPartyServiceEntityOutline {
				__typename
			}
			... on UnavailableEntityOutline {
				__typename
			}
			... on WorkloadEntityOutline {
				__typename
				createdAt
				updatedAt
			}
		}
		nextCursor
	}
	types {
		count
		domain
		entityType
		type
	}
} } }`
