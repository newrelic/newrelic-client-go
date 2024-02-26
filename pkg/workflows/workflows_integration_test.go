//go:build integration
// +build integration

package workflows

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/v3/pkg/ai"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v3/pkg/notifications"

	mock "github.com/newrelic/newrelic-client-go/v3/pkg/testhelpers"
)

func TestIntegrationCreateWorkflow(t *testing.T) {
	t.Parallel()
	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	// Create a destination to work with in this test
	destination, channel := createTestChannel(t, accountID)
	defer cleanupDestination(t, destination)

	notificationTriggers := []AiWorkflowsNotificationTrigger{"ACTIVATED"}

	// Create a workflow to work with in this test
	workflowInput := generateCreateWorkflowInput(channel, notificationTriggers)

	n := newIntegrationTestClient(t)
	createResult, err := n.AiWorkflowsCreateWorkflow(accountID, workflowInput)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	defer cleanupWorkflow(t, &createResult.Workflow)
	var createdWorkflow = createResult.Workflow

	// compare plain fields
	require.Equal(t, workflowInput.WorkflowEnabled, createdWorkflow.WorkflowEnabled)
	require.Equal(t, workflowInput.DestinationsEnabled, createdWorkflow.DestinationsEnabled)
	require.Equal(t, workflowInput.EnrichmentsEnabled, createdWorkflow.EnrichmentsEnabled)
	require.Equal(t, workflowInput.Name, createdWorkflow.Name)
	require.Equal(t, workflowInput.MutingRulesHandling, createdWorkflow.MutingRulesHandling)

	// compare filter input to actual filter
	require.Equal(t, workflowInput.IssuesFilter.Name, createdWorkflow.IssuesFilter.Name)
	require.Equal(t, workflowInput.IssuesFilter.Type, createdWorkflow.IssuesFilter.Type)
	require.Equal(t, 1, len(createdWorkflow.IssuesFilter.Predicates))
	require.Equal(t, workflowInput.IssuesFilter.Predicates[0].Attribute, createdWorkflow.IssuesFilter.Predicates[0].Attribute)
	require.Equal(t, workflowInput.IssuesFilter.Predicates[0].Values, createdWorkflow.IssuesFilter.Predicates[0].Values)
	require.Equal(t, workflowInput.IssuesFilter.Predicates[0].Operator, createdWorkflow.IssuesFilter.Predicates[0].Operator)

	// compare enrichments
	require.Equal(t, len(workflowInput.Enrichments.NRQL), len(createdWorkflow.Enrichments))
	require.Equal(t, workflowInput.Enrichments.NRQL[0].Name, createdWorkflow.Enrichments[0].Name)
	require.Equal(t, workflowInput.Enrichments.NRQL[0].Configuration[0].Query, createdWorkflow.Enrichments[0].Configurations[0].Query)

	// compare destinations
	require.Equal(t, len(workflowInput.DestinationConfigurations), len(createdWorkflow.DestinationConfigurations))
	require.Equal(t, workflowInput.DestinationConfigurations[0].ChannelId, createdWorkflow.DestinationConfigurations[0].ChannelId)
	require.Equal(t, workflowInput.DestinationConfigurations[0].NotificationTriggers, createdWorkflow.DestinationConfigurations[0].NotificationTriggers)
}

func TestIntegrationCreateWorkflowWithoutNotificationTriggers(t *testing.T) {
	t.Parallel()
	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	// Create a destination to work with in this test
	destination, channel := createTestChannel(t, accountID)
	defer cleanupDestination(t, destination)

	// Create a workflow to work with in this test
	workflowInput := generateCreateWorkflowInput(channel, nil)

	n := newIntegrationTestClient(t)
	createResult, err := n.AiWorkflowsCreateWorkflow(accountID, workflowInput)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	defer cleanupWorkflow(t, &createResult.Workflow)
	var createdWorkflow = createResult.Workflow

	require.Equal(t, []AiWorkflowsNotificationTrigger{"ACKNOWLEDGED", "ACTIVATED", "CLOSED", "PRIORITY_CHANGED"}, createdWorkflow.DestinationConfigurations[0].NotificationTriggers)
}

func TestIntegrationDeleteWorkflow_CanDeleteChannels(t *testing.T) {
	t.Parallel()

	// Create stuff to delete
	workflow, destination, channel := createTestWorkflow(t)
	defer cleanupDestination(t, destination)
	defer cleanupChannel(t, channel)
	workflowsClient := newIntegrationTestClient(t)

	// Test: Delete Workflow (with channel)
	deleteResult, err := workflowsClient.AiWorkflowsDeleteWorkflow(workflow.AccountID, true, workflow.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	requireDoesNotExist(t, workflow)
	requireChannelDoesNotExist(t, channel)
}

func TestIntegrationDeleteWorkflow_CanLeaveChannelsAlive(t *testing.T) {
	t.Parallel()

	// Create stuff to delete
	workflow, destination, channel := createTestWorkflow(t)
	defer cleanupDestination(t, destination)
	defer cleanupChannel(t, channel)
	workflowsClient := newIntegrationTestClient(t)

	// Test: Delete Workflow (with channel)
	deleteResult, err := workflowsClient.AiWorkflowsDeleteWorkflow(workflow.AccountID, false, workflow.ID)
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	requireChannelExists(t, channel)
}

func TestIntegrationUpdateWorkflow_EmptyUpdate(t *testing.T) {
	t.Parallel()

	// Create stuff to update
	workflow, destination, channel := createTestWorkflow(t)
	defer cleanupDestination(t, destination)
	defer cleanupWorkflow(t, workflow)
	defer cleanupChannel(t, channel)
	workflowsClient := newIntegrationTestClient(t)

	updatedWorkflow, err := workflowsClient.AiWorkflowsUpdateWorkflow(workflow.AccountID, true, AiWorkflowsUpdateWorkflowInput{
		ID: workflow.ID,
	})

	require.NoError(t, err)
	require.NotNil(t, updatedWorkflow)
	require.Equal(t, workflow.ID, updatedWorkflow.Workflow.ID)
}

func TestIntegrationUpdateWorkflow_UpdateDestinations_CanLeaveOldChannelAlive(t *testing.T) {
	t.Parallel()

	// Create stuff to update
	workflow, destination, channel := createTestWorkflow(t)
	defer cleanupWorkflow(t, workflow)
	defer cleanupChannel(t, channel)
	defer cleanupDestination(t, destination)
	workflowsClient := newIntegrationTestClient(t)

	newDestination, newChannel := createTestChannel(t, workflow.AccountID)
	defer cleanupDestination(t, newDestination)
	defer cleanupChannel(t, newChannel)

	updatedWorkflow, err := workflowsClient.AiWorkflowsUpdateWorkflow(workflow.AccountID, false, AiWorkflowsUpdateWorkflowInput{
		ID: workflow.ID,
		DestinationConfigurations: &[]AiWorkflowsDestinationConfigurationInput{{
			ChannelId:            newChannel.ID,
			NotificationTriggers: []AiWorkflowsNotificationTrigger{"ACTIVATED", "CLOSED"},
		}},
	})

	require.NoError(t, err)
	require.NotNil(t, updatedWorkflow)
	requireChannelExists(t, newChannel)
	requireChannelExists(t, channel)
}

func TestIntegrationUpdateWorkflow_UpdateDestinations_CanDeleteOldChannel(t *testing.T) {
	t.Parallel()

	// Create stuff to update
	workflow, destination, channel := createTestWorkflow(t)
	defer cleanupWorkflow(t, workflow)
	defer cleanupChannel(t, channel)
	defer cleanupDestination(t, destination)
	workflowsClient := newIntegrationTestClient(t)

	newDestination, newChannel := createTestChannel(t, workflow.AccountID)
	defer cleanupDestination(t, newDestination)
	defer cleanupChannel(t, newChannel)

	updatedWorkflow, err := workflowsClient.AiWorkflowsUpdateWorkflow(workflow.AccountID, true, AiWorkflowsUpdateWorkflowInput{
		ID: workflow.ID,
		DestinationConfigurations: &[]AiWorkflowsDestinationConfigurationInput{{
			ChannelId:            newChannel.ID,
			NotificationTriggers: []AiWorkflowsNotificationTrigger{"ACTIVATED", "CLOSED"},
		}},
	})

	require.NoError(t, err)
	require.NotNil(t, updatedWorkflow)
	requireChannelExists(t, newChannel)
	requireChannelDoesNotExist(t, channel)
}

func TestIntegrationUpdateWorkflow_UpdateEverything(t *testing.T) {
	t.Parallel()

	// Create stuff to update
	workflow, destination, channel := createTestWorkflow(t)
	newDestination, newChannel := createTestChannel(t, workflow.AccountID)
	defer cleanupDestination(t, destination)
	defer cleanupDestination(t, newDestination)
	defer cleanupWorkflow(t, workflow)
	defer cleanupChannel(t, channel)
	defer cleanupChannel(t, newChannel)

	workflowsClient := newIntegrationTestClient(t)

	// Update multiple fields
	falseValue := false
	newName := fmt.Sprintf("test-workflows-update-workflow-%s", mock.RandSeq(5))
	workflowInput := AiWorkflowsUpdateWorkflowInput{
		ID:                  workflow.ID,
		WorkflowEnabled:     &falseValue,
		DestinationsEnabled: &falseValue,
		EnrichmentsEnabled:  &falseValue,
		MutingRulesHandling: AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES,
		Enrichments: &AiWorkflowsUpdateEnrichmentsInput{
			NRQL: []AiWorkflowsNRQLUpdateEnrichmentInput{{
				Name: "enrichment-test-update",
				Configuration: []AiWorkflowsNRQLConfigurationInput{{
					Query: "SELECT * FROM AnotherMetric",
				}},
				// TODO: we absolutely should not require the user to know the old ID
				ID: workflow.Enrichments[0].ID,
			}},
		},
		IssuesFilter: &AiWorkflowsUpdatedFilterInput{
			FilterInput: AiWorkflowsFilterInput{
				Name: "filter-test-update",
				Type: AiWorkflowsFilterTypeTypes.FILTER,
				Predicates: []AiWorkflowsPredicateInput{{
					Attribute: "accumulations.sources",
					Operator:  AiWorkflowsOperatorTypes.CONTAINS,
					Values:    []string{"servicenow"},
				}},
			},
			// TODO: we absolutely should not require the user to know the old ID
			ID: workflow.IssuesFilter.ID,
		},
		DestinationConfigurations: &[]AiWorkflowsDestinationConfigurationInput{{
			ChannelId:            newChannel.ID,
			NotificationTriggers: []AiWorkflowsNotificationTrigger{"ACTIVATED", "CLOSED"},
		}},
		Name: &newName,
	}

	workflowUpdateResult, err := workflowsClient.AiWorkflowsUpdateWorkflow(workflow.AccountID, true, workflowInput)

	require.NoError(t, err)
	require.NotNil(t, workflowUpdateResult)
	updatedWorkflow := workflowUpdateResult.Workflow

	// compare plain fields
	require.Equal(t, *workflowInput.WorkflowEnabled, updatedWorkflow.WorkflowEnabled)
	require.Equal(t, *workflowInput.DestinationsEnabled, updatedWorkflow.DestinationsEnabled)
	require.Equal(t, *workflowInput.EnrichmentsEnabled, updatedWorkflow.EnrichmentsEnabled)
	require.Equal(t, *workflowInput.Name, updatedWorkflow.Name)
	require.Equal(t, workflowInput.MutingRulesHandling, updatedWorkflow.MutingRulesHandling)

	// compare filter input to actual filter
	require.Equal(t, workflowInput.IssuesFilter.FilterInput.Name, updatedWorkflow.IssuesFilter.Name)
	require.Equal(t, workflowInput.IssuesFilter.FilterInput.Type, updatedWorkflow.IssuesFilter.Type)
	require.Equal(t, 1, len(updatedWorkflow.IssuesFilter.Predicates))
	require.Equal(t, workflowInput.IssuesFilter.FilterInput.Predicates[0].Attribute, updatedWorkflow.IssuesFilter.Predicates[0].Attribute)
	require.Equal(t, workflowInput.IssuesFilter.FilterInput.Predicates[0].Values, updatedWorkflow.IssuesFilter.Predicates[0].Values)
	require.Equal(t, workflowInput.IssuesFilter.FilterInput.Predicates[0].Operator, updatedWorkflow.IssuesFilter.Predicates[0].Operator)

	// compare enrichments
	require.Equal(t, len(workflowInput.Enrichments.NRQL), len(updatedWorkflow.Enrichments))
	require.Equal(t, workflowInput.Enrichments.NRQL[0].Name, updatedWorkflow.Enrichments[0].Name)
	require.Equal(t, workflowInput.Enrichments.NRQL[0].Configuration[0].Query, updatedWorkflow.Enrichments[0].Configurations[0].Query)

	// compare destinations
	require.Equal(t, len(*workflowInput.DestinationConfigurations), len(updatedWorkflow.DestinationConfigurations))
	require.Equal(t, (*workflowInput.DestinationConfigurations)[0].ChannelId, updatedWorkflow.DestinationConfigurations[0].ChannelId)
	require.Equal(t, (*workflowInput.DestinationConfigurations)[0].NotificationTriggers, updatedWorkflow.DestinationConfigurations[0].NotificationTriggers)
}

func TestIntegrationUpdateWorkflow_RemoveEnrichments(t *testing.T) {
	t.Parallel()

	// Create stuff to update
	workflow, destination, _ := createTestWorkflow(t)
	defer cleanupDestination(t, destination)
	defer cleanupWorkflow(t, workflow)

	// just assert that the created workflow is enabled
	require.Equal(t, true, workflow.WorkflowEnabled)

	workflowsClient := newIntegrationTestClient(t)
	emptyEnrichments := AiWorkflowsUpdateEnrichmentsInput{NRQL: []AiWorkflowsNRQLUpdateEnrichmentInput{}}
	updatedWorkflow, err := workflowsClient.AiWorkflowsUpdateWorkflow(workflow.AccountID, true, AiWorkflowsUpdateWorkflowInput{
		ID:          workflow.ID,
		Enrichments: &emptyEnrichments,
	})

	require.NoError(t, err)
	require.NotNil(t, updatedWorkflow)
	require.Empty(t, updatedWorkflow.Workflow.Enrichments)
}

func TestIntegrationUpdateWorkflow_DisableWorkflow(t *testing.T) {
	t.Parallel()

	// Create stuff to update
	workflow, destination, _ := createTestWorkflow(t)
	defer cleanupDestination(t, destination)
	defer cleanupWorkflow(t, workflow)

	// just assert that the created workflow is enabled
	require.Equal(t, true, workflow.WorkflowEnabled)

	workflowsClient := newIntegrationTestClient(t)
	falseValue := false
	updatedWorkflow, err := workflowsClient.AiWorkflowsUpdateWorkflow(workflow.AccountID, true, AiWorkflowsUpdateWorkflowInput{
		ID:              workflow.ID,
		WorkflowEnabled: &falseValue,
	})

	require.NoError(t, err)
	require.NotNil(t, updatedWorkflow)
	require.Equal(t, false, updatedWorkflow.Workflow.WorkflowEnabled)
}

func TestIntegrationUpdateWorkflow_AddNotificationTriggers(t *testing.T) {
	t.Parallel()

	// Create stuff to update
	workflow, destination, _ := createTestWorkflow(t)
	defer cleanupDestination(t, destination)
	defer cleanupWorkflow(t, workflow)

	// just assert that the created workflow is enabled
	require.Equal(t, true, workflow.WorkflowEnabled)

	destinationWithNotificationTriggers := workflow.DestinationConfigurations
	destinationWithNotificationTriggers[0].NotificationTriggers = []AiWorkflowsNotificationTrigger{"ACTIVATED"}

	workflowsClient := newIntegrationTestClient(t)
	updatedWorkflow, err := workflowsClient.AiWorkflowsUpdateWorkflow(workflow.AccountID, true, AiWorkflowsUpdateWorkflowInput{
		ID: workflow.ID,
		DestinationConfigurations: &[]AiWorkflowsDestinationConfigurationInput{{
			ChannelId:            destinationWithNotificationTriggers[0].ChannelId,
			NotificationTriggers: destinationWithNotificationTriggers[0].NotificationTriggers,
		}},
	})

	require.NoError(t, err)
	require.NotNil(t, updatedWorkflow)
	require.Equal(t, []AiWorkflowsNotificationTrigger{"ACTIVATED"}, updatedWorkflow.Workflow.DestinationConfigurations[0].NotificationTriggers)
}

func TestIntegrationGetWorkflow(t *testing.T) {
	workflow, destination, _ := createTestWorkflow(t)
	defer cleanupDestination(t, destination)
	defer cleanupWorkflow(t, workflow)

	workflowsClient := newIntegrationTestClient(t)
	workflows, err := workflowsClient.GetWorkflows(workflow.AccountID, "", ai.AiWorkflowsFilters{
		ID: workflow.ID,
	})

	require.NoError(t, err)
	require.NotNil(t, workflows)
	require.Len(t, workflows.Entities, 1)
	require.Equal(t, workflow, &workflows.Entities[0])
	require.NotEmpty(t, workflows.Entities[0].GUID)
}

func TestIntegrationGetWorkflow_WorkflowDoesNotExist(t *testing.T) {
	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}
	workflowsClient := newIntegrationTestClient(t)
	workflows, err := workflowsClient.GetWorkflows(accountID, "", ai.AiWorkflowsFilters{
		ID: "214ecde0-135d-4d00-83af-195d6ad07985", // random UUID that does not exist
	})
	require.NoError(t, err)
	require.Equal(t, 0, workflows.TotalCount)
	require.Equal(t, 0, len(workflows.Entities))
}

func TestIntegrationCreateDisabledWorkflow(t *testing.T) {
	t.Parallel()
	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	// Create a destination to work with in this test
	destination, channel := createTestChannel(t, accountID)
	defer cleanupDestination(t, destination)

	// Create a workflow to work with in this test
	workflowInput := generateCreateWorkflowInput(channel, []AiWorkflowsNotificationTrigger{"ACTIVATED"})
	workflowInput.WorkflowEnabled = false
	workflowInput.EnrichmentsEnabled = false
	workflowInput.DestinationsEnabled = false

	n := newIntegrationTestClient(t)
	createResult, err := n.AiWorkflowsCreateWorkflow(accountID, workflowInput)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	defer cleanupWorkflow(t, &createResult.Workflow)
	var createdWorkflow = createResult.Workflow

	// compare plain fields
	require.Equal(t, false, createdWorkflow.WorkflowEnabled)
	require.Equal(t, false, createdWorkflow.DestinationsEnabled)
	require.Equal(t, false, createdWorkflow.EnrichmentsEnabled)
}

func createTestWorkflow(t *testing.T) (*AiWorkflowsWorkflow, *notifications.AiNotificationsDestination, *notifications.AiNotificationsChannel) {
	accountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	destination, channel := createTestChannel(t, accountID)

	notificationTriggers := []AiWorkflowsNotificationTrigger{"ACTIVATED"}

	workflowsClient := newIntegrationTestClient(t)
	workflowInput := generateCreateWorkflowInput(channel, notificationTriggers)
	createResult, err := workflowsClient.AiWorkflowsCreateWorkflow(accountID, workflowInput)
	if err != nil {
		cleanupChannel(t, channel)
		cleanupDestination(t, destination)
	}
	require.NoError(t, err)

	return &createResult.Workflow, destination, channel

}

func generateCreateWorkflowInput(channel *notifications.AiNotificationsChannel, notificationTriggers []AiWorkflowsNotificationTrigger) AiWorkflowsCreateWorkflowInput {
	enrichmentsInput := AiWorkflowsEnrichmentsInput{
		NRQL: []AiWorkflowsNRQLEnrichmentInput{{
			Name: "enrichment-test",
			Configuration: []AiWorkflowsNRQLConfigurationInput{{
				Query: "SELECT * FROM Logs",
			}},
		}},
	}
	filterInput := AiWorkflowsFilterInput{
		Name: "filter-test",
		Type: AiWorkflowsFilterTypeTypes.FILTER,
		Predicates: []AiWorkflowsPredicateInput{{
			Attribute: "accumulations.sources",
			Operator:  AiWorkflowsOperatorTypes.CONTAINS,
			Values:    []string{"newrelic"},
		}},
	}
	destinationsInput := []AiWorkflowsDestinationConfigurationInput{{
		ChannelId:            channel.ID,
		NotificationTriggers: notificationTriggers,
	}}

	return AiWorkflowsCreateWorkflowInput{
		WorkflowEnabled:           true,
		DestinationsEnabled:       true,
		EnrichmentsEnabled:        true,
		MutingRulesHandling:       AiWorkflowsMutingRulesHandlingTypes.DONT_NOTIFY_FULLY_OR_PARTIALLY_MUTED_ISSUES,
		Name:                      fmt.Sprintf("test-workflows-workflow-%s", mock.RandSeq(5)),
		Enrichments:               &enrichmentsInput,
		IssuesFilter:              filterInput,
		DestinationConfigurations: destinationsInput,
	}
}

func createTestChannel(t *testing.T, accountID int) (*notifications.AiNotificationsDestination, *notifications.AiNotificationsChannel) {
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

	client := newNotificationsIntegrationTestClient(t)
	createDestinationResult, err := client.AiNotificationsCreateDestination(accountID, destination)
	require.NoError(t, err)

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
	createChannelResult, err := client.AiNotificationsCreateChannel(accountID, channel)
	if err != nil {
		cleanupDestination(t, &createDestinationResult.Destination)
	}
	require.NoError(t, err)

	return &createDestinationResult.Destination, &createChannelResult.Channel
}

func cleanupDestination(t *testing.T, destination *notifications.AiNotificationsDestination) {
	client := newNotificationsIntegrationTestClient(t)
	_, err := client.AiNotificationsDeleteDestination(destination.AccountID, destination.ID)
	require.NoError(t, err)
}

func cleanupChannel(t *testing.T, channel *notifications.AiNotificationsChannel) {
	client := newNotificationsIntegrationTestClient(t)
	_, err := client.AiNotificationsDeleteChannel(channel.AccountID, channel.ID)
	require.NoError(t, err)
}

func cleanupWorkflow(t *testing.T, workflow *AiWorkflowsWorkflow) {
	client := newIntegrationTestClient(t)
	_, err := client.AiWorkflowsDeleteWorkflow(workflow.AccountID, true, workflow.ID)
	require.NoError(t, err)
}

func requireDoesNotExist(t *testing.T, workflow *AiWorkflowsWorkflow) {
	workflowsClient := newIntegrationTestClient(t)
	workflows, err := workflowsClient.GetWorkflows(workflow.AccountID, "", ai.AiWorkflowsFilters{
		ID: workflow.ID,
	})
	require.NoError(t, err)
	require.Equal(t, 0, workflows.TotalCount)
	require.Equal(t, 0, len(workflows.Entities))
}

func requireChannelExists(t *testing.T, channel *notifications.AiNotificationsChannel) {
	workflowsClient := newNotificationsIntegrationTestClient(t)
	channels, err := workflowsClient.GetChannels(channel.AccountID, "", ai.AiNotificationsChannelFilter{
		ID: channel.ID,
	}, notifications.AiNotificationsChannelSorter{})
	require.NoError(t, err)
	require.Equal(t, 1, channels.TotalCount)
	require.Equal(t, 1, len(channels.Entities))
	require.Equal(t, channel.ID, channels.Entities[0].ID)
}

func requireChannelDoesNotExist(t *testing.T, channel *notifications.AiNotificationsChannel) {
	workflowsClient := newNotificationsIntegrationTestClient(t)
	channels, err := workflowsClient.GetChannels(channel.AccountID, "", ai.AiNotificationsChannelFilter{
		ID: channel.ID,
	}, notifications.AiNotificationsChannelSorter{})
	require.NoError(t, err)
	require.Equal(t, 0, channels.TotalCount)
	require.Equal(t, 0, len(channels.Entities))
}
