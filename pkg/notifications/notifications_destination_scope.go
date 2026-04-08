package notifications

import (
	"context"
	"fmt"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
)

// destinationMutationResponseFields is the shared GraphQL response fragment reused across
// create and update destination mutations to avoid duplication.
const destinationMutationResponseFields = `
	destination {
		accountId
		active
		auth {
			... on AiNotificationsBasicAuth {
			  authType
			  user
			}
			... on AiNotificationsOAuth2Auth {
			  accessTokenUrl
			  scope
			  refreshable
			  refreshInterval
			  prefix
			  clientId
			  authorizationUrl
			  authType
			}
			... on AiNotificationsTokenAuth {
			  authType
			  prefix
			}
			... on AiNotificationsCustomHeadersAuth {
			  authType
			  customHeaders {
				key
			  }
			}
		}
		createdAt
		id
		guid
		isUserAuthenticated
		lastSent
		name
		properties {
			displayValue
			key
			label
			value
		}
		scope {
			id
			type
		}
		secureUrl {
			prefix
		}
		status
		type
		updatedAt
		updatedBy
	}
	errors {
	  ... on AiNotificationsConstraintsError {
		constraints {
		  dependencies
		  name
		}
	  }
	  ... on AiNotificationsDataValidationError {
		details
		fields {
		  field
		  message
		}
	  }
	  ... on AiNotificationsResponseError {
		description
		details
		type
	  }
	  ... on AiNotificationsSuggestionError {
		description
		type
		details
	  }
	}
	error {
	  ... on AiNotificationsSuggestionError {
		description
		type
		details
	  }
	  ... on AiNotificationsResponseError {
		description
		type
		details
	  }
	  ... on AiNotificationsDataValidationError {
		details
		fields {
		  message
		  field
		}
	  }
	  ... on AiNotificationsConstraintsError {
		constraints {
		  name
		  dependencies
		}
	  }
	}
`

// ===== CREATE =====

// CreateDestinationWithScope creates a notification destination under either an account or
// an organization scope.
//
// # Migration from accountId
//
// The previous API accepted a standalone accountId integer. That parameter has been replaced
// by the scope object, which unifies account-level and organization-level targeting under a
// single argument.
//
// If you previously called:
//
//	client.AiNotificationsCreateDestination(accountID, destinationInput)
//
// Replace it with:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	client.CreateDestinationWithScope(destinationInput, scope)
//
// To create a destination at the organization level instead:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ORGANIZATION,
//	    ID:   "<organizationID>",
//	}
//	client.CreateDestinationWithScope(destinationInput, scope)
func (a *Notifications) CreateDestinationWithScope(
	destination AiNotificationsDestinationInput,
	scope *EntityScopeInput,
) (*AiNotificationsDestinationWithScopeResponse, error) {
	return a.CreateDestinationWithScopeWithContext(context.Background(), destination, scope)
}

// CreateDestinationWithScopeWithContext is the context-aware variant of CreateDestinationWithScope.
// Use this when you need request cancellation, deadlines, or tracing propagation.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	resp, err := client.CreateDestinationWithScopeWithContext(ctx, destinationInput, scope)
func (a *Notifications) CreateDestinationWithScopeWithContext(
	ctx context.Context,
	destination AiNotificationsDestinationInput,
	scope *EntityScopeInput,
) (*AiNotificationsDestinationWithScopeResponse, error) {

	if scope == nil {
		return nil, fmt.Errorf("scope is required")
	}

	resp := createDestinationWithScopeResponse{}
	vars := map[string]interface{}{
		"destination": destination,
		"scopeId":     scope.ID,
	}

	mutation := aiNotificationsCreateDestinationWithAccountScopeMutation
	if scope.Type == EntityScopeTypeInputTypes.ORGANIZATION {
		mutation = aiNotificationsCreateDestinationWithOrgScopeMutation
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, mutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AiNotificationsDestinationWithScopeResponse, nil
}

type createDestinationWithScopeResponse struct {
	AiNotificationsDestinationWithScopeResponse AiNotificationsDestinationWithScopeResponse `json:"AiNotificationsCreateDestination"`
}

const aiNotificationsCreateDestinationWithOrgScopeMutation = `mutation(
	$destination: AiNotificationsDestinationInput!,
	$scopeId: String!,
) { aiNotificationsCreateDestination(
	destination: $destination,
	scope: {type: ORGANIZATION, id: $scopeId},
) {` + destinationMutationResponseFields + `} }`

const aiNotificationsCreateDestinationWithAccountScopeMutation = `mutation(
	$destination: AiNotificationsDestinationInput!,
	$scopeId: String!,
) { aiNotificationsCreateDestination(
	destination: $destination,
	scope: {type: ACCOUNT, id: $scopeId},
) {` + destinationMutationResponseFields + `} }`

// ===== GET =====

// GetDestinationsWithScope fetches notification destinations under either an account or an
// organization scope, with support for filtering, sorting, and cursor-based pagination.
//
// # Migration from accountId
//
// The previous API accepted a standalone accountId integer. That parameter has been replaced
// by the scope object, which unifies account-level and organization-level targeting under a
// single argument.
//
// If you previously called:
//
//	client.GetDestinations(accountID, cursor, filters, sorter)
//
// Replace it with:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	client.GetDestinationsWithScope(cursor, filters, sorter, scope)
//
// To query destinations at the organization level instead:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ORGANIZATION,
//	    ID:   "<organizationID>",
//	}
//	client.GetDestinationsWithScope(cursor, filters, sorter, scope)
func (a *Notifications) GetDestinationsWithScope(
	cursor string,
	filters ai.AiNotificationsDestinationFilter,
	sorter AiNotificationsDestinationSorter,
	scope *EntityScopeInput,
) (*AiNotificationsDestinationsWithScopeResponse, error) {
	return a.GetDestinationsWithScopeWithContext(context.Background(), cursor, filters, sorter, scope)
}

// GetDestinationsWithScopeWithContext is the context-aware variant of GetDestinationsWithScope.
// Use this when you need request cancellation, deadlines, or tracing propagation.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	resp, err := client.GetDestinationsWithScopeWithContext(ctx, cursor, filters, sorter, scope)
func (a *Notifications) GetDestinationsWithScopeWithContext(
	ctx context.Context,
	cursor string,
	filters ai.AiNotificationsDestinationFilter,
	sorter AiNotificationsDestinationSorter,
	scope *EntityScopeInput,
) (*AiNotificationsDestinationsWithScopeResponse, error) {

	if scope == nil {
		return nil, fmt.Errorf("scope is required")
	}

	if scope.Type == EntityScopeTypeInputTypes.ORGANIZATION {
		return a.getDestinationsWithOrganizationScope(ctx, cursor, filters, sorter)
	}

	accountID, err := strconv.Atoi(scope.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid account scope ID %q: %w", scope.ID, err)
	}
	return a.getDestinationsWithAccountScope(ctx, accountID, cursor, filters, sorter)
}

func (a *Notifications) getDestinationsWithAccountScope(
	ctx context.Context,
	accountID int,
	cursor string,
	filters ai.AiNotificationsDestinationFilter,
	sorter AiNotificationsDestinationSorter,
) (*AiNotificationsDestinationsWithScopeResponse, error) {

	resp := destinationsWithScopeResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"cursor":    cursor,
		"filters":   filters,
		"sorter":    sorter,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getDestinationsWithAccountScopeQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Account.AiNotifications.Destinations, nil
}

func (a *Notifications) getDestinationsWithOrganizationScope(
	ctx context.Context,
	cursor string,
	filters ai.AiNotificationsDestinationFilter,
	sorter AiNotificationsDestinationSorter,
) (*AiNotificationsDestinationsWithScopeResponse, error) {

	resp := destinationsWithScopeResponse{}
	vars := map[string]interface{}{
		"cursor":  cursor,
		"filters": filters,
		"sorter":  sorter,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getDestinationsWithOrganizationScopeQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Organization.AiNotifications.Destinations, nil
}

// AiNotificationsDestinationWithScope extends AiNotificationsDestination with scope information.
type AiNotificationsDestinationWithScope struct {
	AiNotificationsDestination
	Scope *EntityScope `json:"scope,omitempty"`
}

// AiNotificationsDestinationWithScopeResponse is the mutation response type for scoped destination
// operations. It mirrors AiNotificationsDestinationResponse but surfaces the scope field.
type AiNotificationsDestinationWithScopeResponse struct {
	Destination AiNotificationsDestinationWithScope `json:"destination,omitempty"`
	Error       ai.AiNotificationsError             `json:"error,omitempty"`
	Errors      []ai.AiNotificationsError           `json:"errors"`
}

// AiNotificationsDestinationsWithScopeResponse is the response type for scoped destination queries.
type AiNotificationsDestinationsWithScopeResponse struct {
	Entities   []AiNotificationsDestinationWithScope `json:"entities"`
	Error      AiNotificationsResponseError          `json:"error,omitempty"`
	Errors     []AiNotificationsResponseError        `json:"errors"`
	NextCursor string                                `json:"nextCursor,omitempty"`
	TotalCount int                                   `json:"totalCount"`
}

// destinationsWithScopeResponse handles both account and organization scope paths.
type destinationsWithScopeResponse struct {
	Actor struct {
		Account struct {
			AiNotifications struct {
				Destinations AiNotificationsDestinationsWithScopeResponse `json:"destinations,omitempty"`
			} `json:"aiNotifications,omitempty"`
		} `json:"account,omitempty"`
		Organization struct {
			AiNotifications struct {
				Destinations AiNotificationsDestinationsWithScopeResponse `json:"destinations,omitempty"`
			} `json:"aiNotifications,omitempty"`
		} `json:"organization,omitempty"`
	} `json:"actor,omitempty"`
}

const getDestinationsWithAccountScopeQuery = `query($accountId: Int!, $filters: AiNotificationsDestinationFilter, $sorter: AiNotificationsDestinationSorter, $cursor: String) {
	actor {
		account(id: $accountId) {
			aiNotifications {
				destinations(filters: $filters, sorter: $sorter, cursor: $cursor) {
					error {
						description
						type
						details
					}
					totalCount
					entities {
						accountId
						active
						createdAt
						id
						guid
						lastSent
						name
						properties {
							value
							key
						}
						type
						updatedAt
						updatedBy
						auth {
							... on AiNotificationsBasicAuth {
								authType
								user
							}
							... on AiNotificationsOAuth2Auth {
								accessTokenUrl
								authType
								authorizationUrl
								clientId
								prefix
								refreshable
								refreshInterval
								scope
							}
							... on AiNotificationsTokenAuth {
								authType
								prefix
							}
							... on AiNotificationsCustomHeadersAuth {
								authType
								customHeaders {
									key
								}
							}
						}
						secureUrl {
							prefix
						}
						status
						scope {
							id
							type
						}
						isUserAuthenticated
					}
					nextCursor
				}
			}
		}
	}
}`

const getDestinationsWithOrganizationScopeQuery = `query($filters: AiNotificationsDestinationFilter, $sorter: AiNotificationsDestinationSorter, $cursor: String) {
	actor {
		organization {
			aiNotifications {
				destinations(filters: $filters, sorter: $sorter, cursor: $cursor) {
					error {
						description
						type
						details
					}
					totalCount
					entities {
						accountId
						active
						createdAt
						id
						guid
						lastSent
						name
						properties {
							value
							key
						}
						type
						updatedAt
						updatedBy
						auth {
							... on AiNotificationsBasicAuth {
								authType
								user
							}
							... on AiNotificationsOAuth2Auth {
								accessTokenUrl
								authType
								authorizationUrl
								clientId
								prefix
								refreshable
								refreshInterval
								scope
							}
							... on AiNotificationsTokenAuth {
								authType
								prefix
							}
							... on AiNotificationsCustomHeadersAuth {
								authType
								customHeaders {
									key
								}
							}
						}
						secureUrl {
							prefix
						}
						status
						scope {
							id
							type
						}
						isUserAuthenticated
					}
					nextCursor
				}
			}
		}
	}
}`

// ===== UPDATE =====

// UpdateDestinationWithScope updates an existing notification destination under either an
// account or an organization scope.
//
// # Migration from accountId
//
// The previous API accepted a standalone accountId integer. That parameter has been replaced
// by the scope object, which unifies account-level and organization-level targeting under a
// single argument.
//
// If you previously called:
//
//	client.AiNotificationsUpdateDestination(accountID, destinationUpdate, destinationID)
//
// Replace it with:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	client.UpdateDestinationWithScope(destinationID, destinationUpdate, scope)
//
// To update a destination at the organization level instead:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ORGANIZATION,
//	    ID:   "<organizationID>",
//	}
//	client.UpdateDestinationWithScope(destinationID, destinationUpdate, scope)
func (a *Notifications) UpdateDestinationWithScope(
	destinationId string,
	destination AiNotificationsDestinationUpdate,
	scope *EntityScopeInput,
) (*AiNotificationsDestinationWithScopeResponse, error) {
	return a.UpdateDestinationWithScopeWithContext(context.Background(), destinationId, destination, scope)
}

// UpdateDestinationWithScopeWithContext is the context-aware variant of UpdateDestinationWithScope.
// Use this when you need request cancellation, deadlines, or tracing propagation.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	resp, err := client.UpdateDestinationWithScopeWithContext(ctx, destinationID, destinationUpdate, scope)
func (a *Notifications) UpdateDestinationWithScopeWithContext(
	ctx context.Context,
	destinationId string,
	destination AiNotificationsDestinationUpdate,
	scope *EntityScopeInput,
) (*AiNotificationsDestinationWithScopeResponse, error) {

	if scope == nil {
		return nil, fmt.Errorf("scope is required")
	}

	resp := updateDestinationWithScopeResponse{}
	vars := map[string]interface{}{
		"destination":   destination,
		"destinationId": destinationId,
		"scopeId":       scope.ID,
	}

	mutation := aiNotificationsUpdateDestinationWithAccountScopeMutation
	if scope.Type == EntityScopeTypeInputTypes.ORGANIZATION {
		mutation = aiNotificationsUpdateDestinationWithOrgScopeMutation
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, mutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AiNotificationsDestinationWithScopeResponse, nil
}

type updateDestinationWithScopeResponse struct {
	AiNotificationsDestinationWithScopeResponse AiNotificationsDestinationWithScopeResponse `json:"AiNotificationsUpdateDestination"`
}

const aiNotificationsUpdateDestinationWithOrgScopeMutation = `mutation(
	$destination: AiNotificationsDestinationUpdate!,
	$destinationId: ID!,
	$scopeId: String!,
) { aiNotificationsUpdateDestination(
	destination: $destination,
	destinationId: $destinationId,
	scope: {type: ORGANIZATION, id: $scopeId},
) {` + destinationMutationResponseFields + `} }`

const aiNotificationsUpdateDestinationWithAccountScopeMutation = `mutation(
	$destination: AiNotificationsDestinationUpdate!,
	$destinationId: ID!,
	$scopeId: String!,
) { aiNotificationsUpdateDestination(
	destination: $destination,
	destinationId: $destinationId,
	scope: {type: ACCOUNT, id: $scopeId},
) {` + destinationMutationResponseFields + `} }`

// ===== DELETE =====

// DeleteDestinationWithScope deletes a notification destination under either an account or
// an organization scope.
//
// # Migration from accountId
//
// The previous API accepted a standalone accountId integer. That parameter has been replaced
// by the scope object, which unifies account-level and organization-level targeting under a
// single argument.
//
// If you previously called:
//
//	client.AiNotificationsDeleteDestination(accountID, destinationID)
//
// Replace it with:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	client.DeleteDestinationWithScope(destinationID, scope)
//
// To delete a destination at the organization level instead:
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ORGANIZATION,
//	    ID:   "<organizationID>",
//	}
//	client.DeleteDestinationWithScope(destinationID, scope)
func (a *Notifications) DeleteDestinationWithScope(
	destinationId string,
	scope *EntityScopeInput,
) (*AiNotificationsDeleteResponse, error) {
	return a.DeleteDestinationWithScopeWithContext(context.Background(), destinationId, scope)
}

// DeleteDestinationWithScopeWithContext is the context-aware variant of DeleteDestinationWithScope.
// Use this when you need request cancellation, deadlines, or tracing propagation.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	scope := &notifications.EntityScopeInput{
//	    Type: notifications.EntityScopeTypeInputTypes.ACCOUNT,
//	    ID:   strconv.Itoa(accountID),
//	}
//	resp, err := client.DeleteDestinationWithScopeWithContext(ctx, destinationID, scope)
func (a *Notifications) DeleteDestinationWithScopeWithContext(
	ctx context.Context,
	destinationId string,
	scope *EntityScopeInput,
) (*AiNotificationsDeleteResponse, error) {

	if scope == nil {
		return nil, fmt.Errorf("scope is required")
	}

	resp := deleteDestinationWithScopeResponse{}
	vars := map[string]interface{}{
		"destinationId": destinationId,
		"scopeId":       scope.ID,
	}

	mutation := aiNotificationsDeleteDestinationWithAccountScopeMutation
	if scope.Type == EntityScopeTypeInputTypes.ORGANIZATION {
		mutation = aiNotificationsDeleteDestinationWithOrgScopeMutation
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, mutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AiNotificationsDeleteResponse, nil
}

type deleteDestinationWithScopeResponse struct {
	AiNotificationsDeleteResponse AiNotificationsDeleteResponse `json:"AiNotificationsDeleteDestination"`
}

const aiNotificationsDeleteDestinationWithOrgScopeMutation = `mutation(
	$destinationId: ID!,
	$scopeId: String!,
) { aiNotificationsDeleteDestination(
	destinationId: $destinationId,
	scope: {type: ORGANIZATION, id: $scopeId},
) {
	error {
		description
		details
		type
	}
	errors {
		description
		details
		type
	}
	ids
} }`

const aiNotificationsDeleteDestinationWithAccountScopeMutation = `mutation(
	$destinationId: ID!,
	$scopeId: String!,
) { aiNotificationsDeleteDestination(
	destinationId: $destinationId,
	scope: {type: ACCOUNT, id: $scopeId},
) {
	error {
		description
		details
		type
	}
	errors {
		description
		details
		type
	}
	ids
} }`
