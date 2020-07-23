package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Events is used to send custom events to NRDB.
type Events struct {
	client http.Client
	config config.Config
	logger logging.Logger
}

// New is used to create a new Events client instance.
func New(cfg config.Config) Events {
	cfg.Compression = config.Compression.Gzip

	client := http.NewClient(cfg)
	client.SetAuthStrategy(&http.InsightsInsertKeyAuthorizer{})

	pkg := Events{
		client: client,
		config: cfg,
		logger: cfg.GetLogger(),
	}

	return pkg
}

// CreateEvent reports a custom event to New Relic.
func (e *Events) CreateEvent(accountID int, event interface{}) error {
	var jsonData []byte

	switch event := event.(type) {
	case []byte:
		jsonData = event
	case string:
		jsonData = []byte(event)
	default:
		var err error
		jsonData, err = json.Marshal(event)
		if err != nil {
			return fmt.Errorf("error marshaling event data: %s", err.Error())
		}
	}

	if !strings.Contains(string(jsonData), "eventType") {
		return fmt.Errorf("event data must contain eventType field. %s", jsonData)
	}

	resp := &createEventResponse{}
	_, err := e.client.Post(e.config.Region().InsightsURL(accountID), nil, jsonData, resp)

	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New("failed creating custom event")
	}

	return nil
}

type createEventResponse struct {
	Success bool   `json:"success"`
	UUID    string `json:"uuid"`
}
