// +build integration

package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationSearchEntities(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	params := EntitySearchQueryBuilder{
		Name: "Dummy App",
	}

	actual, err := client.GetEntitySearch(
		EntitySearchOptions{},
		"",
		params,
		[]EntitySearchSortCriteria{},
	)

	require.NoError(t, err)
	require.Greater(t, len(actual.Results.Entities), 0)
}

func TestIntegrationSearchEntitiesByTags(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	params := EntitySearchQueryBuilder{
		Tags: []EntitySearchQueryBuilderTag{
			{
				Key:   "language",
				Value: "nodejs",
			},
		},
	}

	actual, err := client.GetEntitySearch(
		EntitySearchOptions{},
		"",
		params,
		[]EntitySearchSortCriteria{},
	)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

func TestIntegrationGetEntities(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	guids := []EntityGUID{"MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"}
	actual, err := client.GetEntities(guids)

	require.NoError(t, err)
	require.Greater(t, len((*actual)), 0)
}

func TestIntegrationGetEntity(t *testing.T) {
	t.Parallel()

	entityGUID := EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")
	client := newIntegrationTestClient(t)

	result, err := client.GetEntity(entityGUID)

	require.NoError(t, err)
	require.NotNil(t, result)

	actual := (*result).(*ApmApplicationEntity)

	// These are a bit fragile, if the above GUID ever changes...
	assert.Equal(t, 2520528, actual.AccountID)
	assert.Equal(t, "APM", actual.Domain)
	assert.Equal(t, EntityType("APM_APPLICATION_ENTITY"), actual.EntityType)
	assert.Equal(t, entityGUID, actual.GUID)
	assert.Equal(t, "Dummy App", actual.Name)
	assert.Equal(t, "https://one.newrelic.com/redirect/entity/"+string(entityGUID), actual.Permalink)
	assert.Equal(t, true, actual.Reporting)
}

// Looking at an APM Application, and the result set here.
func TestIntegrationGetEntity_ApmEntity(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	result, err := client.GetEntity("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")

	require.NoError(t, err)
	require.NotNil(t, result)

	actual := (*result).(*ApmApplicationEntity)

	// These are a bit fragile, if the above GUID ever changes...
	// from ApmApplicationEntity / ApmApplicationEntityOutline
	assert.Equal(t, 215037795, actual.ApplicationID)
	assert.Equal(t, EntityAlertSeverityTypes.NOT_ALERTING, actual.AlertSeverity)
	assert.Equal(t, "nodejs", actual.Language)
	assert.NotNil(t, actual.RunningAgentVersions)
	assert.NotNil(t, actual.RunningAgentVersions.MinVersion)
	assert.NotNil(t, actual.RunningAgentVersions.MaxVersion)
	assert.NotNil(t, actual.Settings)
	assert.NotNil(t, actual.Settings.ApdexTarget)
	assert.NotNil(t, actual.Settings.ServerSideConfig)

}

// Looking at a Browser Application, and the result set here.
func TestIntegrationGetEntity_BrowserEntity(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	result, err := client.GetEntity("MjUwODI1OXxCUk9XU0VSfEFQUExJQ0FUSU9OfDIwNDI2MTYyOA")

	require.NoError(t, err)
	require.NotNil(t, result)

	actual := (*result).(*BrowserApplicationEntity)

	// These are a bit fragile, if the above GUID ever changes...
	// from BrowserApplicationEntity / BrowserApplicationEntityOutline
	assert.Equal(t, 204261628, actual.ApplicationID)
	assert.Equal(t, 204261368, actual.ServingApmApplicationID)
	assert.Equal(t, EntityAlertSeverityTypes.NOT_CONFIGURED, actual.AlertSeverity)

}

// Looking at a Mobile Application, and the result set here.
func TestIntegrationGetEntity_MobileEntity(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	result, err := client.GetEntity("NDQ0NTN8TU9CSUxFfEFQUExJQ0FUSU9OfDE3ODg1NDI")

	require.NoError(t, err)
	require.NotNil(t, result)

	actual := (*result).(*MobileApplicationEntity)

	// These are a bit fragile, if the above GUID ever changes...
	// from MobileApplicationEntity / MobileApplicationEntityOutline
	assert.Equal(t, 1788542, actual.ApplicationID)
	assert.Equal(t, EntityAlertSeverityTypes.NOT_CONFIGURED, actual.AlertSeverity)

}

func newIntegrationTestClient(t *testing.T) Entities {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
