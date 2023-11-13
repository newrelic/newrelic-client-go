//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateEventsService(t *testing.T) {
	t.Parallel()

	client, _ := New(ConfigPersonalAPIKey(testAPIkey))
	service := NewEventsService(client)
	assert.Equal(t, service, &client.Events)
}
