package newrelic

import (
	"context"
)

type EventsAPI interface {
	CreateEvent(accountID int, event interface{}) error
	CreateEventWithContext(ctx context.Context, accountID int, event interface{}) error
}

func NewEventsService(client *NewRelic) EventsAPI {
	return &client.Events
}
