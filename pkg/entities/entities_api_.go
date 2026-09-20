package entities

import (
	"context"
	"fmt"
	"strings"
)

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

// GetEntitySearchByQueryWithCursor is like GetEntitySearchByQuery but paginates
// via the results cursor returned in prior responses. Pass an empty cursor on
// the first call, then feed each subsequent call the value from the previous
// response's Results.NextCursor. When NextCursor is empty, pagination is done.
//
// See GetEntitySearchByQuery for details on the query, options, and sortBy
// arguments.
func (a *Entities) GetEntitySearchByQueryWithCursor(
	options EntitySearchOptions,
	query string,
	sortBy []EntitySearchSortCriteria,
	cursor string,
) (*EntitySearch, error) {
	return a.GetEntitySearchByQueryWithCursorWithContext(context.Background(),
		options,
		query,
		sortBy,
		cursor,
	)
}

// GetEntitySearchByQueryWithCursorWithContext is like
// GetEntitySearchByQueryWithCursor but accepts a context.Context.
func (a *Entities) GetEntitySearchByQueryWithCursorWithContext(
	ctx context.Context,
	options EntitySearchOptions,
	query string,
	sortBy []EntitySearchSortCriteria,
	cursor string,
) (*EntitySearch, error) {

	resp := entitySearchResponse{}
	vars := map[string]interface{}{
		"options": options,
		"query":   query,
		"sortBy":  sortBy,
		"cursor":  nilIfEmpty(cursor),
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getEntitySearchByQueryWithCursor, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.EntitySearch, nil
}

// GetEntitySearchWithCursor is like GetEntitySearch but paginates via the
// results cursor returned in prior responses. Pass an empty cursor on the
// first call, then feed each subsequent call the value from the previous
// response's Results.NextCursor. When NextCursor is empty, pagination is done.
//
// See GetEntitySearch for details on the other arguments.
func (a *Entities) GetEntitySearchWithCursor(
	options EntitySearchOptions,
	query string,
	queryBuilder EntitySearchQueryBuilder,
	sortBy []EntitySearchSortCriteria,
	sortByWithDirection []SortCriterionWithDirection,
	cursor string,
) (*EntitySearch, error) {
	return a.GetEntitySearchWithCursorWithContext(context.Background(),
		options,
		query,
		queryBuilder,
		sortBy,
		sortByWithDirection,
		cursor,
	)
}

// GetEntitySearchWithCursorWithContext is like GetEntitySearchWithCursor but
// accepts a context.Context.
func (a *Entities) GetEntitySearchWithCursorWithContext(
	ctx context.Context,
	options EntitySearchOptions,
	query string,
	queryBuilder EntitySearchQueryBuilder,
	sortBy []EntitySearchSortCriteria,
	sortByWithDirection []SortCriterionWithDirection,
	cursor string,
) (*EntitySearch, error) {

	resp := entitySearchResponse{}
	vars := map[string]interface{}{
		"options":             options,
		"query":               query,
		"queryBuilder":        queryBuilder,
		"sortBy":              sortBy,
		"sortByWithDirection": sortByWithDirection,
		"cursor":              nilIfEmpty(cursor),
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getEntitySearchWithCursor, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.EntitySearch, nil
}

// The two *WithCursor query strings are derived from the existing non-cursor
// queries by injecting a `$cursor: String` variable and passing it into the
// `results(cursor: $cursor)` sub-field. Keeping them derived avoids duplicating
// the hundreds of lines of entity selection sets already maintained in the
// generated file. If a regeneration ever changes the anchor patterns below,
// addCursorArg panics at package init so the drift is caught in CI rather than
// silently returning un-paginated data at runtime.
var (
	getEntitySearchByQueryWithCursor = addCursorArg(getEntitySearchByQuery)
	getEntitySearchWithCursor        = addCursorArg(getEntitySearchQuery)
)

// addCursorArg rewrites an entitySearch query to accept a `$cursor: String`
// variable and forwards it to the `results` sub-field. See the block comment
// above for context.
func addCursorArg(query string) string {
	const (
		// The two anchors appear exactly once in each entitySearch query const
		// (getEntitySearchByQuery in this file, and getEntitySearchQuery in
		// entities_api.go). If tutone regeneration reshapes either query, the
		// checks below fail fast at package init.
		varAnchor    = ") { actor { entitySearch("
		resultAnchor = "\n\tresults {\n"
		varInsert    = "\t$cursor: String,\n"
		resultInsert = "\n\tresults(cursor: $cursor) {\n"
	)

	if strings.Count(query, varAnchor) != 1 {
		panic(fmt.Sprintf("entities: cannot add cursor argument, %q anchor not found exactly once", varAnchor))
	}
	if strings.Count(query, resultAnchor) != 1 {
		panic(fmt.Sprintf("entities: cannot add cursor argument, %q anchor not found exactly once", resultAnchor))
	}

	query = strings.Replace(query, varAnchor, varInsert+varAnchor, 1)
	query = strings.Replace(query, resultAnchor, resultInsert, 1)
	return query
}

// nilIfEmpty returns a *string that is nil when s is empty, otherwise pointing
// at s. NerdGraph treats a nil (null) cursor as "start from the beginning";
// passing an empty string would send `cursor: ""` which some resolvers reject.
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
