//go:build integration
// +build integration

package workflows

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Workflows {
	cfg := mock.NewIntegrationTestConfig(t)
	client := New(cfg)

	return client
}

func newNotificationsIntegrationTestClient(t *testing.T) notifications.Notifications {
	cfg := mock.NewIntegrationTestConfig(t)
	client := notifications.New(cfg)

	return client
}
