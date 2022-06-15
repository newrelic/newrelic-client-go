package notifications

import (
	"context"
	_ "github.com/newrelic/newrelic-client-go/pkg/errors"
)

// CreateDestinationMutation creates a destination within a given account.
func (n *Notifications) CreateDestinationMutation(accountID int, destination DestinationInput) (*Destination, error) {
	return n.CreateDestinationMutationWithContext(context.Background(), accountID, destination)
}

// CreateDestinationMutationWithContext creates a destination within a given account.
func (n *Notifications) CreateDestinationMutationWithContext(ctx context.Context, accountID int, destinationInput DestinationInput) (*Destination, error) {
	vars := map[string]interface{}{
		"accountID":        accountID,
		"destinationInput": destinationInput,
	}

	resp := createDestinationResponse{}

	if err := n.client.NerdGraphQueryWithContext(ctx, notificationsCreateDestination, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AiNotificationsCreateDestination.Destination, nil
}

// ListDestinations returns all notifications destinations.
func (n *Notifications) ListDestinations(accountID int) ([]*Destination, error) {
	return n.ListDestinationsWithContext(context.Background(), accountID)
}

// ListDestinationsWithContext returns all notifications destinations.
func (n *Notifications) ListDestinationsWithContext(ctx context.Context, accountID int) ([]*Destination, error) {
	resp := listDestinationsResponse{}

	vars := map[string]interface{}{
		"accountID": accountID,
	}

	if err := n.client.NerdGraphQueryWithContext(ctx, notificationsGetDestinations, vars, &resp); err != nil {
		return nil, err
	}

	return resp.Actor.Account.AiNotifications.Destinations.Entities, nil
}

// GetDestination returns a specific notification destination by ID for a given account.
func (n *Notifications) GetDestination(accountID int, id UUID) (*Destination, error) {
	return n.GetDestinationWithContext(context.Background(), accountID, id)
}

// GetDestinationWithContext returns a specific alert channel by ID for a given account.
func (n *Notifications) GetDestinationWithContext(ctx context.Context, accountID int, id UUID) (*Destination, error) {
	resp := listDestinationsResponse{}

	filters := map[string]interface{}{
		"id": id,
	}

	vars := map[string]interface{}{
		"accountID": accountID,
		"filters":   filters,
	}

	if err := n.client.NerdGraphQueryWithContext(ctx, notificationsGetDestination, vars, &resp); err != nil {
		return nil, err
	}

	return resp.Actor.Account.AiNotifications.Destinations.Entities[0], nil
}

// DeleteDestinationMutation deletes the destination with the specified ID.
func (n *Notifications) DeleteDestinationMutation(accountID int, id UUID) ([]string, error) {
	return n.DeleteDestinationMutationWithContext(context.Background(), accountID, id)
}

// DeleteDestinationMutationWithContext deletes the destination with the specified ID.
func (n *Notifications) DeleteDestinationMutationWithContext(ctx context.Context, accountID int, id UUID) ([]string, error) {
	resp := deleteDestinationResponse{}

	vars := map[string]interface{}{
		"accountID":     accountID,
		"destinationID": id,
	}

	if err := n.client.NerdGraphQueryWithContext(ctx, notificationsDeleteDestination, vars, &resp); err != nil {
		return nil, err
	}

	return resp.AiNotificationsDeleteDestination.Ids, nil
}
