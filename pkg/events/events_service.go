package events

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

type EventsAPI interface {
	CreateEvent(accountID int, event interface{}) error
	CreateEventWithContext(ctx context.Context, accountID int, event interface{}) error
}

func NewEventsService(opts ...config.ConfigOption) (*Events, error) {
	cfg := config.New()

	err := cfg.Init(opts)
	if err != nil {
		return nil, err
	}

	events := New(cfg)
	return &events, nil
}
