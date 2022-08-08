//go:build integration
// +build integration

package workflows

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestMutationWorkflow(t *testing.T) {
	t.Parallel()

	n := newIntegrationTestClient(t)
	newrelicClient := newrelicIntegrationTestClient(t)

	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	// Create a destination to work with in this test
	testIntegrationDestinationNameRandStr := mock.RandSeq(5)
	destination := notifications.AiNotificationsDestinationInput{}
	destination.Type = notifications.AiNotificationsDestinationTypeTypes.WEBHOOK
	destination.Properties = []notifications.AiNotificationsPropertyInput{
		{
			Key:          "url",
			Value:        "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b",
			Label:        "",
			DisplayValue: "",
		},
	}
	destination.Auth = &notifications.AiNotificationsCredentialsInput{
		Type: notifications.AiNotificationsAuthTypeTypes.TOKEN,
		Token: notifications.AiNotificationsTokenAuthInput{
			Token:  "Token",
			Prefix: "Bearer",
		},
	}
	destination.Name = fmt.Sprintf("test-notifications-destination-%s", testIntegrationDestinationNameRandStr)

	// Test: Create Destination
	createDestinationResult, err := newrelicClient.Notifications.AiNotificationsCreateDestination(accountID, destination)
	require.NoError(t, err)
	require.NotNil(t, createDestinationResult)

	destinationID := createDestinationResult.Destination.ID

	// Create a channel to work with in this test
	testIntegrationChannelNameRandStr := mock.RandSeq(5)
	channel := notifications.AiNotificationsChannelInput{}
	channel.Type = notifications.AiNotificationsChannelTypeTypes.WEBHOOK
	channel.Product = notifications.AiNotificationsProductTypes.IINT
	channel.Properties = []notifications.AiNotificationsPropertyInput{
		{
			Key:          "payload",
			Value:        "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
			Label:        "Payload Template",
			DisplayValue: "",
		},
	}
	channel.DestinationId = createDestinationResult.Destination.ID
	channel.Name = fmt.Sprintf("test-notifications-channel-%s", testIntegrationChannelNameRandStr)

	// Test: Create Channel
	createChannelResult, err := newrelicClient.Notifications.AiNotificationsCreateChannel(accountID, channel)
	require.NoError(t, err)
	require.NotNil(t, createChannelResult)

	channelId := createChannelResult.Channel.ID

	// Create a workflow to work with in this test
	testIntegrationWorkflowNameRandStr := mock.RandSeq(5)
	workflow := AiWorkflowsCreateWorkflowInput{}
	workflow.WorkflowEnabled = false
	workflow.DestinationsEnabled = true
	workflow.EnrichmentsEnabled = true
	workflow.MutingRulesHandling = AiWorkflowsMutingRulesHandlingTypes.DONT_NOTIFY_FULLY_OR_PARTIALLY_MUTED_ISSUES
	workflow.Name = fmt.Sprintf("test-workflows-workflow-%s", testIntegrationWorkflowNameRandStr)
	workflow.Enrichments = &AiWorkflowsEnrichmentsInput{
		NRQL: []AiWorkflowsNRQLEnrichmentInput{{
			Name: "enrichment-test",
			Configuration: []AiWorkflowsNRQLConfigurationInput{{
				Query: "SELECT * FROM Logs",
			}},
		}},
	}
	workflow.IssuesFilter = AiWorkflowsFilterInput{
		Name: "filter-test",
		Type: AiWorkflowsFilterTypeTypes.FILTER,
		Predicates: []AiWorkflowsPredicateInput{{
			Attribute: "source",
			Operator:  AiWorkflowsOperatorTypes.CONTAINS,
			Values:    []string{"newrelic"},
		}},
	}
	workflow.DestinationConfigurations = []AiWorkflowsDestinationConfigurationInput{{
		ChannelId: channelId,
	}}

	// Test: Create Workflow
	createResult, err := n.AiWorkflowsCreateWorkflow(accountID, workflow)
	require.NoError(t, err)
	require.NotNil(t, createResult)

	id := createResult.Workflow.ID

	// Test: Get Workflow
	filters := ai.AiWorkflowsFilters{
		ID: id,
	}
	getWorkflowResult, err := n.GetWorkflows(accountID, "", filters)
	require.NoError(t, err)
	require.NotNil(t, getWorkflowResult)
	assert.Equal(t, 1, getWorkflowResult.TotalCount)

	// Test: Update Workflow
	updateWorkflow := AiWorkflowsUpdateWorkflowInput{}
	updateWorkflow.WorkflowEnabled = true
	updateWorkflow.DestinationsEnabled = false
	updateWorkflow.EnrichmentsEnabled = false
	updateWorkflow.MutingRulesHandling = AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES
	updateWorkflow.Enrichments = &AiWorkflowsUpdateEnrichmentsInput{
		NRQL: []AiWorkflowsNRQLUpdateEnrichmentInput{{
			Name: "enrichment-test-update",
			Configuration: []AiWorkflowsNRQLConfigurationInput{{
				Query: "SELECT * FROM Metric",
			}},
			ID: createResult.Workflow.Enrichments[0].ID,
		}},
	}
	updateWorkflow.IssuesFilter = AiWorkflowsUpdatedFilterInput{
		FilterInput: AiWorkflowsFilterInput{
			Name: "filter-test-update",
			Type: AiWorkflowsFilterTypeTypes.FILTER,
			Predicates: []AiWorkflowsPredicateInput{{
				Attribute: "source",
				Operator:  AiWorkflowsOperatorTypes.CONTAINS,
				Values:    []string{"servicenow"},
			}},
		},
		ID: createResult.Workflow.IssuesFilter.ID,
	}
	updateWorkflow.DestinationConfigurations = []AiWorkflowsDestinationConfigurationInput{{
		ChannelId: channelId,
	}}
	updateWorkflow.Name = fmt.Sprintf("test-workflows-update-workflow-%s", testIntegrationWorkflowNameRandStr)
	updateWorkflow.ID = id

	updateWorkflowResult, err := n.AiWorkflowsUpdateWorkflow(accountID, updateWorkflow)
	require.NoError(t, err)
	require.NotNil(t, updateWorkflowResult)

	// Test: Delete Workflow (with channel)
	deleteResult, err := n.AiWorkflowsDeleteWorkflow(accountID, id)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)

	// Delete Destination
	deleteDestinationResult, err := newrelicClient.notifications.AiNotificationsDeleteDestination(accountID, destinationID)
	require.NoError(t, err)
	require.NotNil(t, deleteDestinationResult)
}
