//go:build integration
// +build integration

package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
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

	params = EntitySearchQueryBuilder{
		Name: "WebPortal",
	}

	actual, err = client.GetEntitySearch(
		EntitySearchOptions{},
		"",
		params,
		[]EntitySearchSortCriteria{},
	)

	require.NoError(t, err)
	require.Greater(t, len(actual.Results.Entities), 0)
}

func TestIntegrationSearchEntitiesByQuery(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	query := "domain = 'APM' AND type = 'APPLICATION' and name = 'Dummy App'"

	actual, err := client.GetEntitySearchByQuery(
		EntitySearchOptions{},
		query,
		[]EntitySearchSortCriteria{},
	)

	require.NoError(t, err)
	require.Greater(t, len(actual.Results.Entities), 0)
}

func TestIntegrationSearchEntities_domain(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	domains := []EntitySearchQueryBuilderDomain{
		EntitySearchQueryBuilderDomainTypes.APM,
		EntitySearchQueryBuilderDomainTypes.BROWSER,
		EntitySearchQueryBuilderDomainTypes.INFRA,
		EntitySearchQueryBuilderDomainTypes.MOBILE,
		EntitySearchQueryBuilderDomainTypes.SYNTH,
	}

	for _, d := range domains {
		params := EntitySearchQueryBuilder{
			Domain: d,
		}

		result, err := client.GetEntitySearch(
			EntitySearchOptions{},
			"",
			params,
			[]EntitySearchSortCriteria{},
		)

		require.NoError(t, err)
		require.Greater(t, len(result.Results.Entities), 0)
	}
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

	// GUID of Dummy App
	guids := []common.EntityGUID{"MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"}
	actual, err := client.GetEntities(guids)

	if e, ok := err.(*http.GraphQLErrorResponse); ok {
		if !e.IsDeprecated() {
			require.NoError(t, e)
		}
	}
	require.Greater(t, len((*actual)), 0)
}

func TestIntegrationGetEntity(t *testing.T) {
	t.Parallel()

	// GUID of Dummy App
	entityGUID := common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")
	client := newIntegrationTestClient(t)

	result, err := client.GetEntity(entityGUID)

	if e, ok := err.(*http.GraphQLErrorResponse); ok {
		if !e.IsDeprecated() {
			require.NoError(t, e)
		}
	}
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

	// GUID of Dummy App
	result, err := client.GetEntity("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")

	if e, ok := err.(*http.GraphQLErrorResponse); ok {
		if !e.IsDeprecated() {
			require.NoError(t, e)
		}
	}
	require.NotNil(t, result)

	actual := (*result).(*ApmApplicationEntity)

	// NOT_ALERTING or CRITICAL alert status can be expected
	acceptableAlertStatuses := []EntityAlertSeverity{
		EntityAlertSeverityTypes.NOT_ALERTING,
		EntityAlertSeverityTypes.CRITICAL,
	}

	// These are a bit fragile, if the above GUID ever changes...
	// from ApmApplicationEntity / ApmApplicationEntityOutline
	assert.Equal(t, 215037795, actual.ApplicationID)
	assert.Contains(t, acceptableAlertStatuses, actual.AlertSeverity)
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

	if e, ok := err.(*http.GraphQLErrorResponse); ok {
		if !e.IsDeprecated() {
			require.NoError(t, e)
		}
	}
	require.NotNil(t, result)

	ref := *result
	actual, ok := ref.(*BrowserApplicationEntity)

	if actual == nil || !ok {
		t.Skip("Skipping `TestIntegrationGetEntity_BrowserEntity` integration test due to missing test entity. This entity might have been deleted from this test account.")
		return
	}

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

	if e, ok := err.(*http.GraphQLErrorResponse); ok {
		if !e.IsDeprecated() {
			require.NoError(t, e)
		}
	}
	require.NotNil(t, (*result))

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
