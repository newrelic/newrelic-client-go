// +build integration

package synthetics

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationGetMonitorLocations(t *testing.T) {
	t.Parallel()

	tc := mock.NewIntegrationTestConfig(t)

	synthetics := New(tc)

	locations, err := synthetics.GetMonitorLocations()
	require.NoError(t, err)
	require.Greater(t, len(locations), 0)
}
