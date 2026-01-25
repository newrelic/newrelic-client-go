//go:build integration
// +build integration

package fleetcontrol

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCreateConfigurationAndVersion(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	configurationVersionOneBody := "Body for Version 1"
	configurationVersionTwoBody := "Body for Version 2"

	createConfigurationResponse, err := client.FleetControlCreateConfiguration(
		configurationVersionOneBody,
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": "{\"name\": \"Random Build v100 Test\", \"agentType\": \"NRInfra\", \"managedEntityType\": \"KUBERNETESCLUSTER\"}",
			},
		},
		testOrganizationID,
	)

	require.NoError(t, err)
	fmt.Println(createConfigurationResponse.ConfigurationEntityGUID)
	fmt.Println(createConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)

	require.NotNil(t, createConfigurationResponse.ConfigurationEntityGUID)
	require.NotNil(t, createConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.Equal(t, createConfigurationResponse.ConfigurationVersion.ConfigurationVersionNumber, 1)

	addVersionToConfigurationResponse, err := client.FleetControlCreateConfiguration(
		configurationVersionTwoBody,
		map[string]interface{}{
			"x-newrelic-client-go-custom-headers": map[string]string{
				"Newrelic-Entity": fmt.Sprintf("{\"agentConfiguration\": \"%s\"}", createConfigurationResponse.ConfigurationEntityGUID),
			},
		},
		testOrganizationID,
	)

	require.NoError(t, err)
	fmt.Println(addVersionToConfigurationResponse.ConfigurationEntityGUID)
	fmt.Println(addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.NotNil(t, addVersionToConfigurationResponse.ConfigurationEntityGUID)
	require.NotNil(t, addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.NotEqual(t, createConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID, addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID)
	require.Equal(t, addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionNumber, 2)

	time.Sleep(time.Second * 10)

	getConfigurationResponse, err := client.FleetControlGetConfiguration(
		createConfigurationResponse.ConfigurationEntityGUID,
		testOrganizationID,
		GetConfigurationModeTypes.ConfigEntity,
		-1,
	)

	require.NoError(t, err)
	require.Equal(t, *getConfigurationResponse, GetConfigurationResponse(configurationVersionTwoBody))

	getConfigurationVersionOneResponse, err := client.FleetControlGetConfiguration(
		createConfigurationResponse.ConfigurationEntityGUID,
		testOrganizationID,
		GetConfigurationModeTypes.ConfigEntity,
		1,
	)

	require.NoError(t, err)
	require.Equal(t, *getConfigurationVersionOneResponse, GetConfigurationResponse(configurationVersionOneBody))

	getConfigurationVersionTwoResponse, err := client.FleetControlGetConfiguration(
		addVersionToConfigurationResponse.ConfigurationVersion.ConfigurationVersionEntityGUID,
		testOrganizationID,
		GetConfigurationModeTypes.ConfigVersionEntity,
		-1,
	)

	require.NoError(t, err)
	require.Equal(t, *getConfigurationVersionTwoResponse, GetConfigurationResponse(configurationVersionTwoBody))

	// Test getting all configuration versions
	getConfigurationVersionsResponse, err := client.FleetControlGetConfigurationVersions(
		createConfigurationResponse.ConfigurationEntityGUID,
		testOrganizationID,
	)

	require.NoError(t, err)
	require.NotNil(t, getConfigurationVersionsResponse)
	require.NotNil(t, getConfigurationVersionsResponse.Versions)
	require.GreaterOrEqual(t, len(getConfigurationVersionsResponse.Versions), 2)

	// Verify that we have at least 2 versions
	var foundVersionOne, foundVersionTwo bool
	for _, version := range getConfigurationVersionsResponse.Versions {
		if version.Version == "1" {
			foundVersionOne = true
		}
		if version.Version == "2" {
			foundVersionTwo = true
		}
	}
	require.True(t, foundVersionOne, "Expected to find version 1")
	require.True(t, foundVersionTwo, "Expected to find version 2")

}

func TestIntegrationDeleteBlob(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	createBlobResponse, err := client.FleetControlDeleteConfiguration(
		"NDgyOTY3M3xOR0VQfEFHRU5UX0NPTkZJR1VSQVRJT058MDE5YjBkMWUtMzBiNS03NGYwLTk2M2EtMjk1NzZjNWUwNjEx",
		testOrganizationID,
	)

	fmt.Println(createBlobResponse)
	require.NoError(t, err)

	// require.NotNil(t, createUserResponse.CreatedUser.ID)
}

func TestIntegrationDeleteConfigurationVersion(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	err = client.FleetControlDeleteConfigurationVersion(
		"NDgyOTY3M3xOR0VQfEFHRU5UX0NPTkZJR1VSQVRJT05fVkVSU0lPTnwwMTliZTljZC1jNDljLTdjZTgtOWJjOS03Y2UyNTVjYWIzMjI",
		testOrganizationID,
	)

	require.NoError(t, err)

	// require.NotNil(t, createUserResponse.CreatedUser.ID)
}
