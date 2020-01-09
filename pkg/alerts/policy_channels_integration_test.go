// +build integration

package alerts

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPolicyChannels(t *testing.T) {
	t.Parallel()

	var (
		testPolicyNameRandStr = nr.RandSeq(5)
		testIntegrationPolicy = Policy{
			IncidentPreference: "PER_POLICY",
			Name:               fmt.Sprintf("test-alert-policy-%s", testPolicyNameRandStr),
		}
		testIntegrationChannel = Channel{
			Name: fmt.Sprintf("test-alert-channel-%s", testPolicyNameRandStr),
			Type: "slack",
			Configuration: ChannelConfiguration{
				URL:     "https://example-org.slack.com",
				Channel: testPolicyNameRandStr,
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	policyResp, err := client.CreatePolicy(testIntegrationPolicy)
	policy := *policyResp

	require.NoError(t, err)

	channelResp, err := client.CreateChannel(testIntegrationChannel)
	channel := *channelResp

	require.NoError(t, err)

	// Teardown
	defer func() {
		_, err = client.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("Error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}

		_, err = client.DeleteChannel(channel.ID)
		if err != nil {
			t.Logf("Error cleaning up alert channel %d (%s): %s", channel.ID, channel.Name, err)
		}
	}()

	// Test: Update
	updateResult, err := client.UpdatePolicyChannels(policy.ID, []int{channel.ID})

	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	deleteResult, err := client.DeletePolicyChannel(policy.ID, updateResult.ChannelIDs[0])

	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}
