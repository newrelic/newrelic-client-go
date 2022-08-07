//go:build unit
// +build unit

package workflows

import (
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/nrtime"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	timestampString = "2022-07-25T12:08:07.179638Z"
	timestamp       = nrtime.DateTime(timestampString)
	user            = "test-user"
	accountId       = 10867072
	channelId       = "0d11fd42-5919-4767-8cf5-e07cb71c1b04"
	id              = "03bd4929-3d86-4447-a077-a901b5d511ff"

	testCreateWorkflowResponseJSON = `{
		"aiWorkflowsCreateWorkflow": {
      	"errors": [],
		  "workflow": {
			"accountId": 10867072,
			"createdAt": "2022-07-25T12:08:07.179638Z",
			"destinationConfigurations": [
			  {
				"channelId": "0d11fd42-5919-4767-8cf5-e07cb71c1b04",
				"name": "EMPTY",
				"type": "EMAIL"
			  }
			],
			"destinationsEnabled": false,
			"enrichments": [
			  {
				"accountId": 10867072,
				"configurations": [
				  {
					"query": "SELECT * FROM Logs"
				  }
				],
				"createdAt": "2022-07-25T12:08:07.179638Z",
				"id": "79ce0157-9c20-4d95-839f-daeb070bebb1",
				"name": "Logs",
				"type": "NRQL",
				"updatedAt": "2022-07-25T12:08:07.179638Z"
			  }
			],
			"enrichmentsEnabled": false,
			"id": "03bd4929-3d86-4447-a077-a901b5d511ff",
			"issuesFilter": {
			  "accountId": 10867072,
			  "id": "1a7f6caf-5afa-4bd5-841c-56108c3b244d",
			  "name": "source",
			  "predicates": [
				{
				  "attribute": "source",
				  "operator": "CONTAINS",
				  "values": [
					"newrelic"
				  ]
				}
			  ],
			  "type": "FILTER"
			},
			"lastRun": null,
			"mutingRulesHandling": "DONT_NOTIFY_FULLY_MUTED_ISSUES",
			"name": "workflow-test",
			"updatedAt": "2022-07-25T12:08:07.179638Z",
			"workflowEnabled": false
		  }
		}
	}`

	testDeleteWorkflowResponseJSON = `{
		"aiWorkflowsDeleteWorkflow": {
		  "errors": [],
		  "id": "03bd4929-3d86-4447-a077-a901b5d511ff"
		}
	}`

	testGetWorkflowResponseJSON = `{
    "actor": {
      "account": {
        "aiWorkflows": {
          "workflows": {
            "entities": [
              {
                "accountId": 10867072,
                "createdAt": "2022-07-25T12:08:07.179638Z",
                "destinationConfigurations": [
                  {
                    "channelId": "0d11fd42-5919-4767-8cf5-e07cb71c1b04",
                    "name": "EMPTY",
                    "type": "EMAIL"
                  }
                ],
                "destinationsEnabled": false,
                "enrichments": [
                  {
                    "accountId": 10867072,
                    "configurations": [
                      {
                        "query": "SELECT * FROM Logs"
                      }
                    ],
                    "createdAt": "2022-07-25T12:08:07.179638Z",
                    "id": "79ce0157-9c20-4d95-839f-daeb070bebb1",
                    "name": "Logs",
                    "type": "NRQL",
                    "updatedAt": "2022-07-25T12:08:07.179638Z"
                  }
                ],
                "enrichmentsEnabled": false,
                "id": "03bd4929-3d86-4447-a077-a901b5d511ff",
                "issuesFilter": {
                  "accountId": 10867072,
                  "id": "1a7f6caf-5afa-4bd5-841c-56108c3b244d",
                  "name": "source",
                  "predicates": [
                    {
                      "attribute": "source",
                      "operator": "CONTAINS",
                      "values": [
                        "newrelic"
                      ]
                    }
                  ],
                  "type": "FILTER"
                },
                "lastRun": null,
                "mutingRulesHandling": "DONT_NOTIFY_FULLY_MUTED_ISSUES",
                "name": "workflow-test",
                "updatedAt": "2022-07-25T12:08:07.179638Z",
                "workflowEnabled": false
              }
            ],
            "nextCursor": null,
            "totalCount": 1
          }
        }
      }
	}}`
)

func TestCreateWorkflow(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateWorkflowResponseJSON)
	workflows := newMockResponse(t, respJSON, http.StatusCreated)

	enrichmentsInput := &AiWorkflowsEnrichmentsInput{
		NRQL: []AiWorkflowsNRQLEnrichmentInput{{
			Name:          "Logs",
			Configuration: []AiWorkflowsNRQLConfigurationInput{{Query: "SELECT * FROM Logs"}},
		}},
	}
	issuesFilterInput := AiWorkflowsFilterInput{
		Name: "source",
		Type: AiWorkflowsFilterTypeTypes.FILTER,
		Predicates: []AiWorkflowsPredicateInput{{
			Attribute: "source",
			Operator:  "CONTAINS",
			Values:    []string{"newrelic"},
		}},
	}
	workflowInput := AiWorkflowsCreateWorkflowInput{
		Name:                "workflow-test",
		MutingRulesHandling: AiWorkflowsMutingRulesHandlingTypes.DONT_NOTIFY_FULLY_MUTED_ISSUES,
		WorkflowEnabled:     false,
		DestinationsEnabled: false,
		EnrichmentsEnabled:  false,
		DestinationConfigurations: []AiWorkflowsDestinationConfigurationInput{{
			ChannelId: channelId,
		}},
		Enrichments:  enrichmentsInput,
		IssuesFilter: issuesFilterInput,
	}

	expectedDestinationConfiguration := []AiWorkflowsDestinationConfiguration{{
		ChannelId: channelId,
		Name:      "EMPTY",
		Type:      "EMAIL",
	}}
	expectedIssuedFilter := AiWorkflowsFilter{
		AccountID: accountId,
		ID:        "1a7f6caf-5afa-4bd5-841c-56108c3b244d",
		Name:      "source",
		Type:      AiWorkflowsFilterTypeTypes.FILTER,
		Predicates: []AiWorkflowsPredicate{{
			Attribute: "source",
			Operator:  "CONTAINS",
			Values:    []string{"newrelic"},
		}},
	}
	expectedEnrichmentConfigurations := []ai.AiWorkflowsConfiguration{
		{Query: "SELECT * FROM Logs"},
	}
	for _, config := range expectedEnrichmentConfigurations {
		config.ImplementsAiWorkflowsConfiguration()
	}

	expectedEnrichments := []AiWorkflowsEnrichment{{
		AccountID:      accountId,
		CreatedAt:      timestamp,
		UpdatedAt:      timestamp,
		ID:             "79ce0157-9c20-4d95-839f-daeb070bebb1",
		Name:           "Logs",
		Type:           AiWorkflowsEnrichmentTypeTypes.NRQL,
		Configurations: expectedEnrichmentConfigurations,
	}}
	expected := &AiWorkflowsCreateWorkflowResponse{
		Workflow: AiWorkflowsWorkflow{
			AccountID:                 accountId,
			CreatedAt:                 timestamp,
			UpdatedAt:                 timestamp,
			DestinationConfigurations: expectedDestinationConfiguration,
			Enrichments:               expectedEnrichments,
			IssuesFilter:              expectedIssuedFilter,
			MutingRulesHandling:       AiWorkflowsMutingRulesHandlingTypes.DONT_NOTIFY_FULLY_MUTED_ISSUES,
			ID:                        id,
			Name:                      "workflow-test",
			LastRun:                   "",
			WorkflowEnabled:           false,
			EnrichmentsEnabled:        false,
			DestinationsEnabled:       false,
		},
		Errors: []AiWorkflowsCreateResponseError{},
	}

	actual, err := workflows.AiWorkflowsCreateWorkflow(accountId, workflowInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetWorkflow(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testGetWorkflowResponseJSON)
	workflows := newMockResponse(t, respJSON, http.StatusOK)

	expectedDestinationConfiguration := []AiWorkflowsDestinationConfiguration{{
		ChannelId: channelId,
		Name:      "EMPTY",
		Type:      "EMAIL",
	}}
	expectedEnrichments := []AiWorkflowsEnrichment{{
		AccountID: accountId,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		ID:        "79ce0157-9c20-4d95-839f-daeb070bebb1",
		Name:      "Logs",
		Type:      AiWorkflowsEnrichmentTypeTypes.NRQL,
		Configurations: []ai.AiWorkflowsConfiguration{
			{Query: "SELECT * FROM Logs"},
		},
	}}
	expectedIssuedFilter := AiWorkflowsFilter{
		AccountID: accountId,
		ID:        "1a7f6caf-5afa-4bd5-841c-56108c3b244d",
		Name:      "source",
		Type:      AiWorkflowsFilterTypeTypes.FILTER,
		Predicates: []AiWorkflowsPredicate{{
			Attribute: "source",
			Operator:  "CONTAINS",
			Values:    []string{"newrelic"},
		}},
	}
	expected := &AiWorkflowsWorkflows{
		Entities: []AiWorkflowsWorkflow{
			{
				AccountID:                 accountId,
				CreatedAt:                 timestamp,
				UpdatedAt:                 timestamp,
				DestinationConfigurations: expectedDestinationConfiguration,
				Enrichments:               expectedEnrichments,
				IssuesFilter:              expectedIssuedFilter,
				MutingRulesHandling:       AiWorkflowsMutingRulesHandlingTypes.DONT_NOTIFY_FULLY_MUTED_ISSUES,
				ID:                        id,
				Name:                      "workflow-test",
				LastRun:                   "",
				WorkflowEnabled:           false,
				EnrichmentsEnabled:        false,
				DestinationsEnabled:       false,
			},
		},
		NextCursor: "",
		TotalCount: 1,
	}

	filters := ai.AiWorkflowsFilters{
		ID: id,
	}

	actual, err := workflows.GetWorkflows(accountId, "", filters)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestDeleteWorkflow(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testDeleteWorkflowResponseJSON)
	workflows := newMockResponse(t, respJSON, http.StatusOK)

	expected := &AiWorkflowsDeleteWorkflowResponse{
		ID:     id,
		Errors: []AiWorkflowsDeleteResponseError{},
	}

	actual, err := workflows.AiWorkflowsDeleteWorkflow(accountId, id)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
