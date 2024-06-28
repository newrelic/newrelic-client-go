// Code generated by tutone: DO NOT EDIT
package workflows

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
)

// Create a new Workflow with issues filter, enrichments and destinations
func (a *Workflows) AiWorkflowsCreateWorkflow(
	accountID int,
	createWorkflowData AiWorkflowsCreateWorkflowInput,
) (*AiWorkflowsCreateWorkflowResponse, error) {
	return a.AiWorkflowsCreateWorkflowWithContext(context.Background(),
		accountID,
		createWorkflowData,
	)
}

// Create a new Workflow with issues filter, enrichments and destinations
func (a *Workflows) AiWorkflowsCreateWorkflowWithContext(
	ctx context.Context,
	accountID int,
	createWorkflowData AiWorkflowsCreateWorkflowInput,
) (*AiWorkflowsCreateWorkflowResponse, error) {

	resp := AiWorkflowsCreateWorkflowQueryResponse{}
	vars := map[string]interface{}{
		"accountId":          accountID,
		"createWorkflowData": createWorkflowData,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AiWorkflowsCreateWorkflowMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AiWorkflowsCreateWorkflowResponse, nil
}

type AiWorkflowsCreateWorkflowQueryResponse struct {
	AiWorkflowsCreateWorkflowResponse AiWorkflowsCreateWorkflowResponse `json:"AiWorkflowsCreateWorkflow"`
}

const AiWorkflowsCreateWorkflowMutation = `mutation(
	$accountId: Int!,
	$createWorkflowData: AiWorkflowsCreateWorkflowInput!,
) { aiWorkflowsCreateWorkflow(
	accountId: $accountId,
	createWorkflowData: $createWorkflowData,
) {
	errors {
		description
		type
	}
	workflow {
		accountId
		createdAt
		destinationConfigurations {
			channelId
			name
			notificationTriggers
			type
			updateOriginalMessage
		}
		destinationsEnabled
		enrichments {
			accountId
			configurations {
                ... on AiWorkflowsNrqlConfiguration {
                  query
                }
		  	}
			createdAt
			id
			name
			type
			updatedAt
		}
		enrichmentsEnabled
		guid
		id
		issuesFilter {
			accountId
			id
			name
			predicates {
				attribute
				operator
				values
			}
			type
		}
		lastRun
		mutingRulesHandling
		name
		updatedAt
		workflowEnabled
	}
} }`

// Delete a workflow and all it's sub entities: filter, enrichments and destinations
func (a *Workflows) AiWorkflowsDeleteWorkflow(
	accountID int,
	deleteChannels bool,
	iD string,
) (*AiWorkflowsDeleteWorkflowResponse, error) {
	return a.AiWorkflowsDeleteWorkflowWithContext(context.Background(),
		accountID,
		deleteChannels,
		iD,
	)
}

// Delete a workflow and all it's sub entities: filter, enrichments and destinations
func (a *Workflows) AiWorkflowsDeleteWorkflowWithContext(
	ctx context.Context,
	accountID int,
	deleteChannels bool,
	iD string,
) (*AiWorkflowsDeleteWorkflowResponse, error) {

	resp := AiWorkflowsDeleteWorkflowQueryResponse{}
	vars := map[string]interface{}{
		"accountId":      accountID,
		"deleteChannels": deleteChannels,
		"id":             iD,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AiWorkflowsDeleteWorkflowMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AiWorkflowsDeleteWorkflowResponse, nil
}

type AiWorkflowsDeleteWorkflowQueryResponse struct {
	AiWorkflowsDeleteWorkflowResponse AiWorkflowsDeleteWorkflowResponse `json:"AiWorkflowsDeleteWorkflow"`
}

const AiWorkflowsDeleteWorkflowMutation = `mutation(
	$accountId: Int!,
	$deleteChannels: Boolean!,
	$id: ID!,
) { aiWorkflowsDeleteWorkflow(
	accountId: $accountId,
	deleteChannels: $deleteChannels,
	id: $id,
) {
	errors {
		description
		type
	}
	id
} }`

// Update Workflow with issues filter, enrichments and destinations
func (a *Workflows) AiWorkflowsUpdateWorkflow(
	accountID int,
	deleteUnusedChannels bool,
	updateWorkflowData AiWorkflowsUpdateWorkflowInput,
) (*AiWorkflowsUpdateWorkflowResponse, error) {
	return a.AiWorkflowsUpdateWorkflowWithContext(context.Background(),
		accountID,
		deleteUnusedChannels,
		updateWorkflowData,
	)
}

// Update Workflow with issues filter, enrichments and destinations
func (a *Workflows) AiWorkflowsUpdateWorkflowWithContext(
	ctx context.Context,
	accountID int,
	deleteUnusedChannels bool,
	updateWorkflowData AiWorkflowsUpdateWorkflowInput,
) (*AiWorkflowsUpdateWorkflowResponse, error) {

	resp := AiWorkflowsUpdateWorkflowQueryResponse{}
	vars := map[string]interface{}{
		"accountId":            accountID,
		"deleteUnusedChannels": deleteUnusedChannels,
		"updateWorkflowData":   updateWorkflowData,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AiWorkflowsUpdateWorkflowMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AiWorkflowsUpdateWorkflowResponse, nil
}

type AiWorkflowsUpdateWorkflowQueryResponse struct {
	AiWorkflowsUpdateWorkflowResponse AiWorkflowsUpdateWorkflowResponse `json:"AiWorkflowsUpdateWorkflow"`
}

const AiWorkflowsUpdateWorkflowMutation = `mutation(
	$accountId: Int!,
	$deleteUnusedChannels: Boolean!,
	$updateWorkflowData: AiWorkflowsUpdateWorkflowInput!,
) { aiWorkflowsUpdateWorkflow(
	accountId: $accountId,
	deleteUnusedChannels: $deleteUnusedChannels,
	updateWorkflowData: $updateWorkflowData,
) {
	errors {
		description
		type
	}
	workflow {
		accountId
		createdAt
		destinationConfigurations {
			channelId
			name
			notificationTriggers
			type
			updateOriginalMessage
		}
		destinationsEnabled
		enrichments {
			accountId
			configurations {
                ... on AiWorkflowsNrqlConfiguration {
                  query
                }
		  	}
			createdAt
			id
			name
			type
			updatedAt
		}
		enrichmentsEnabled
		guid
		id
		issuesFilter {
			accountId
			id
			name
			predicates {
				attribute
				operator
				values
			}
			type
		}
		lastRun
		mutingRulesHandling
		name
		updatedAt
		workflowEnabled
	}
} }`

// Returns a list of workflows with pagination cursor according to account id and filters
func (a *Workflows) GetWorkflows(
	accountID int,
	cursor string,
	filters ai.AiWorkflowsFilters,
) (*AiWorkflowsWorkflows, error) {
	return a.GetWorkflowsWithContext(context.Background(),
		accountID,
		cursor,
		filters,
	)
}

// Returns a list of workflows with pagination cursor according to account id and filters
func (a *Workflows) GetWorkflowsWithContext(
	ctx context.Context,
	accountID int,
	cursor string,
	filters ai.AiWorkflowsFilters,
) (*AiWorkflowsWorkflows, error) {

	resp := workflowsResponse{}
	vars := map[string]interface{}{
		"accountID": accountID,
		"cursor":    cursor,
		"filters":   filters,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getWorkflowsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Account.AiWorkflows.Workflows, nil
}

const getWorkflowsQuery = `query(
	$accountID: Int!, $filters: AiWorkflowsFilters,
) { actor { account(id: $accountID) { aiWorkflows { workflows(filters: $filters) {
	entities {
		accountId
		createdAt
		destinationConfigurations {
			channelId
			name
			notificationTriggers
			type
			updateOriginalMessage
		}
		destinationsEnabled
		enrichments {
			accountId
			configurations {
                ... on AiWorkflowsNrqlConfiguration {
                  query
                }
		  	}
			createdAt
			id
			name
			type
			updatedAt
		}
		enrichmentsEnabled
		guid
		id
		issuesFilter {
			accountId
			id
			name
			predicates {
				attribute
				operator
				values
			}
			type
		}
		lastRun
		mutingRulesHandling
		name
		updatedAt
		workflowEnabled
	}
	nextCursor
	totalCount
} } } } }`
