//go:build integration
// +build integration

package installevents

import (
	log "github.com/sirupsen/logrus"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
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

	log.Infof("Sending this InstallationRecipeStatus from the test: %v", status)
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

	log.Infof("Sending this InstallationRecipeStatus from the test: %v", status)
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
