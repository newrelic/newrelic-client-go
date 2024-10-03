//go:build integration

package agent

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func newAgentIntegrationTestClient(t *testing.T) Agent {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}

func TestIntegrationGetCurrentAgentRelease(t *testing.T) {
	t.Parallel()

	client := newAgentIntegrationTestClient(t)

	agentName := AgentReleasesFilterTypes.GO

	getResult, err := client.GetCurrentAgentRelease(agentName)
	require.NoError(t, err)
	require.NotNil(t, getResult)

	require.NotNil(t, getResult.Bugs)
	require.NotNil(t, getResult.Date)
	require.NotNil(t, getResult.DownloadLink)
	require.NotNil(t, getResult.EolDate)
	require.NotNil(t, getResult.Features)
	require.NotNil(t, getResult.Security)
	require.NotNil(t, getResult.Slug)
	require.NotNil(t, getResult.Version)
}

func TestIntegrationGetCurrentAgentRelease_Invalid(t *testing.T) {
	t.Parallel()

	client := newAgentIntegrationTestClient(t)

	var agentName AgentReleasesFilter = "ASDF"

	getResult, err := client.GetCurrentAgentRelease(agentName)
	require.Error(t, err)
	require.Nil(t, getResult)
}
