// Package entities provides a programmatic API for interacting with New Relic One entities.
package entities

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logging"
)

// Entities is used to communicate with the New Relic Entities product.
type Entities struct {
	client http.Client
	logger logging.Logger
}

// New returns a new client for interacting with New Relic One entities.
func New(config config.Config) Entities {
	return Entities{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}

// Search for entities using a custom query.
//
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
