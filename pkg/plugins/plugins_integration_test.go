//go:build integration
// +build integration

package plugins

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationPlugins(t *testing.T) {
	t.Parallel()

	tc := mock.NewIntegrationTestConfig(t)
	client := New(tc)

	// Test: List
	listResult, err := client.ListPlugins(nil)

	require.NoError(t, err)

	if len(listResult) == 0 {
		t.Skip("Skipping `GetPlugin()` integration test due to zero plugins found")
		return
	}

	// Test: Get
	qp := GetPluginParams{
		Detailed: true,
	}
	getResult, err := client.GetPlugin(listResult[0].ID, &qp)

	require.NoError(t, err)
	require.NotNil(t, getResult)
	require.NotNil(t, getResult.Details)
}
