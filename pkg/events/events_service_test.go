//go:build unit
// +build unit

package events

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v3/pkg/config"
)

var testAPIkey = "efgh56789"

func TestShouldCreateEventsService(t *testing.T) {
	t.Parallel()

	service, err := NewEventsService(config.ConfigPersonalAPIKey(testAPIkey))
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestShouldErrorCreateEventsService(t *testing.T) {
	t.Parallel()

	service, err := NewEventsService(config.ConfigPersonalAPIKey(""))
	assert.Error(t, err)
	assert.Nil(t, service)
}
