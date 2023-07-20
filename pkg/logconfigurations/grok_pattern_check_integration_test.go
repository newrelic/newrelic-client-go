//go:build integration
// +build integration

package logconfigurations

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationTestGrok_WithMatch(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		grok     = "%{INT:bytes_received}"
		logLines = []string{
			"{   \"host_ip\": \"43.3.120.2\",   \"bytes_received\": 2048,   \"bytes_sent\": 0 }",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Read
	res, err := client.GetTestGrok(testAccountID, grok, logLines)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res)
	require.Equal(t, len(*res), 1)
	for _, v := range *res {
		require.True(t, v.Matched)
	}

}

func TestIntegrationTestGrok_WithNotMatch(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		grok     = "%{IP:host_ip}"
		logLines = []string{
			"{  \"bytes_received\": 2048,   \"bytes_sent\": 0 }",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Read
	res, err := client.GetTestGrok(testAccountID, grok, logLines)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res)
	require.Equal(t, len(*res), 1)
	for _, v := range *res {
		require.False(t, v.Matched)
	}
}

func TestIntegrationTestGrok_WithInvalidPattern(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		grok     = "%{abcd}"
		logLines = []string{
			"{  \"bytes_received\": 2048,   \"bytes_sent\": 0 }",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Read
	res, err := client.GetTestGrok(testAccountID, grok, logLines)
	require.Error(t, err)
	require.Nil(t, res)
	require.Empty(t, res)
}
