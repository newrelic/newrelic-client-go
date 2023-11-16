//go:build integration
// +build integration

package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/internal/http"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
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
	guids := []common.EntityGUID{testhelpers.IntegrationTestApplicationEntityGUID}
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
	entityGUID := common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)
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
	assert.Equal(t, 3806526, actual.AccountID)
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
	result, err := client.GetEntity(testhelpers.IntegrationTestApplicationEntityGUID)

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
	assert.Equal(t, 573482638, actual.ApplicationID)
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

	result, err := client.GetEntity("MzgwNjUyNnxNT0JJTEV8QVBQTElDQVRJT058NjAxMzc1OTAx")

	if e, ok := err.(*http.GraphQLErrorResponse); ok {
		if !e.IsDeprecated() {
			require.NoError(t, e)
		}
	}
	require.NotNil(t, (*result))

	actual := (*result).(*MobileApplicationEntity)

	// These are a bit fragile, if the above GUID ever changes...
	// from MobileApplicationEntity / MobileApplicationEntityOutline
	assert.Equal(t, 601375901, actual.ApplicationID)
	assert.Equal(t, EntityAlertSeverityTypes.NOT_CONFIGURED, actual.AlertSeverity)
}

func TestIntegrationGetEntity_SyntheticsEntity(t *testing.T) {
	t.Parallel()
	syntheticsEntityMonitorGUID := "MzgwNjUyNnxTWU5USHxNT05JVE9SfGFmZmM0MTRiLTVhNmMtNGI5NS1iMzYwLThhNmQ2ZTkzOTM3Yw"
	client := newIntegrationTestClient(t)

	result, err := client.GetEntity(common.EntityGUID(syntheticsEntityMonitorGUID))
	if err != nil || result == nil {
		t.Skipf("Entity not found with GUID: %s. Skipping entity integration test for synthetics entity.", syntheticsEntityMonitorGUID)
	}

	if e, ok := err.(*http.GraphQLErrorResponse); ok {
		if !e.IsDeprecated() {
			require.NoError(t, e)
		}
	}
	require.NotNil(t, result)

	entity := (*result).(*SyntheticMonitorEntity)
	require.NotNil(t, entity)

	deviceOrientation := FindTagByKey(entity.Tags, "deviceOrientation")
	runtimeType := FindTagByKey(entity.Tags, "runtimeType")
	runtimeTypeVersion := FindTagByKey(entity.Tags, "runtimeTypeVersion")
	require.Greater(t, len(deviceOrientation), 0)
	require.Greater(t, len(runtimeType), 0)
	require.Greater(t, len(runtimeTypeVersion), 0)
}

func newIntegrationTestClient(t *testing.T) Entities {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
