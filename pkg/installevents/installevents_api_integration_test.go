//go:build integration
// +build integration

package installevents

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestInstallationCreateRecipeEvent(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	status := InstallationRecipeStatus{
		CliVersion: "0.0.1",
		Status:     InstallationRecipeStatusTypeTypes.AVAILABLE,
	}

	response, err := client.InstallationCreateRecipeEvent(testAccountID, status)
	require.NoError(t, err)
	require.NotNil(t, response)
}

func TestInstallationCreateRecipeEvent_ShouldSendMetadata(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	status := InstallationRecipeStatus{
		CliVersion: "0.0.1",
		Name:       "test",
		Status:     InstallationRecipeStatusTypeTypes.AVAILABLE,
		Metadata: map[string]interface{}{
			"someKey": "some value",
		},
	}

	response, err := client.InstallationCreateRecipeEvent(testAccountID, status)

	require.NoError(t, err)
	require.NotNil(t, response.Metadata)

	if metaValue, ok := response.Metadata["someKey"].(string); ok {
		require.Equal(t, "some value", metaValue)
	}
}

func newIntegrationTestClient(t *testing.T) Installevents {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
