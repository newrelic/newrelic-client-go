//go:build integration
// +build integration

package changetracking

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestChangeTrackingCreateDeployment_Basic(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingDeploymentInput{
		Changelog:      "test",
		Commit:         "12345a",
		DeepLink:       "newrelic-client-go",
		DeploymentType: ChangeTrackingDeploymentTypeTypes.BASIC,
		Description:    "This is a test description",
		EntityGUID:     common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew),
		GroupId:        "deployment",
		Timestamp:      nrtime.EpochMilliseconds(time.Now()),
		User:           "newrelic-go-client",
		Version:        "0.0.1",
	}

	res, err := a.ChangeTrackingCreateDeployment(
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
		input,
	)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.Equal(t, res.EntityGUID, input.EntityGUID)
}

// func TestChangeTrackingCreateDeployment_CustomAttributes(t *testing.T) {
// 	skipMsg := fmt.Sprintf("Skipping %s until custom attributes are out of limited preview.", t.Name())
// 	t.Skip(skipMsg)
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)

// 	var customAttributes = `{"a":"1","b":"two","c":"1.5","d":"true"}`
// 	attrs := make(map[string]interface{})
// 	err := json.Unmarshal([]byte(customAttributes), &attrs)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	input := ChangeTrackingDeploymentInput{
// 		Changelog:        "test",
// 		Commit:           "12345a",
// 		CustomAttributes: ChangeTrackingRawCustomAttributesMap(attrs),
// 		DeepLink:         "newrelic-client-go",
// 		DeploymentType:   ChangeTrackingDeploymentTypeTypes.BASIC,
// 		Description:      "This is a test description",
// 		EntityGUID:       common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew),
// 		GroupId:          "deployment",
// 		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
// 		User:             "newrelic-go-client",
// 		Version:          "0.0.1",
// 	}

// 	res, err := a.ChangeTrackingCreateDeployment(
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 		input,
// 	)
// 	require.NoError(t, err)

// 	require.NotNil(t, res)
// 	require.Equal(t, res.EntityGUID, input.EntityGUID)
// }

func TestChangeTrackingCreateDeployment_TimestampError(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingDeploymentInput{
		Changelog:      "test",
		Commit:         "12345a",
		DeepLink:       "newrelic-client-go",
		DeploymentType: ChangeTrackingDeploymentTypeTypes.BASIC,
		Description:    "This is a test description",
		EntityGUID:     common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID),
		GroupId:        "deployment",
		Timestamp:      nrtime.EpochMilliseconds(time.UnixMilli(0)),
		User:           "newrelic-go-client",
		Version:        "0.0.1",
	}

	res, err := a.ChangeTrackingCreateDeployment(
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
		input,
	)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestChangeTrackingCreateDeployment_OlderThan24HoursTimestampError(t *testing.T) {
	t.Parallel()
	now := time.Now()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingDeploymentInput{
		Changelog:      "test",
		Commit:         "12345a",
		DeepLink:       "newrelic-client-go",
		DeploymentType: ChangeTrackingDeploymentTypeTypes.BASIC,
		Description:    "This is a test description",
		EntityGUID:     common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID),
		GroupId:        "deployment",
		Timestamp: nrtime.EpochMilliseconds(
			time.Date(
				now.Year(),
				now.Month(),
				now.Day()-2,
				now.Hour()-3,
				now.Minute()-30,
				0,
				0,
				time.Local,
			),
		),
		User:    "newrelic-go-client",
		Version: "0.0.1",
	}

	res, err := a.ChangeTrackingCreateDeployment(
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
		input,
	)
	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile("not be more than 24 hours"), err.Error())
	require.Nil(t, res)
}

func TestChangeTrackingCreateDeployment_TimestampZeroNanosecondsTest(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)
	now := time.Now()

	input := ChangeTrackingDeploymentInput{
		Changelog:      "test",
		Commit:         "12345a",
		DeepLink:       "newrelic-client-go",
		DeploymentType: ChangeTrackingDeploymentTypeTypes.BASIC,
		Description:    "This is a test description",
		EntityGUID:     common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew),
		GroupId:        "deployment",
		Timestamp: nrtime.EpochMilliseconds(
			time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour()-3,
				now.Minute()-30,
				0,
				0,
				time.Local,
			),
		),
		User:    "newrelic-go-client",
		Version: "0.0.1",
	}

	res, err := a.ChangeTrackingCreateDeployment(
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
		input,
	)
	require.NoError(t, err)
	require.NotNil(t, res.EntityGUID)
	require.Equal(t, res.EntityGUID, input.EntityGUID)
}

func TestChangeTrackingCreateDeployment_TimestampNonZeroNanosecondsTest(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)
	now := time.Now()

	input := ChangeTrackingDeploymentInput{
		Changelog:      "test",
		Commit:         "12345a",
		DeepLink:       "newrelic-client-go",
		DeploymentType: ChangeTrackingDeploymentTypeTypes.BASIC,
		Description:    "This is a test description",
		EntityGUID:     common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew),
		GroupId:        "deployment",
		Timestamp: nrtime.EpochMilliseconds(
			time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour()-6,
				now.Minute()-30,
				0,
				231567,
				time.Local,
			),
		),
		User:    "newrelic-go-client",
		Version: "0.0.1",
	}

	res, err := a.ChangeTrackingCreateDeployment(
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
		input,
	)
	require.NoError(t, err)
	require.NotNil(t, res.EntityGUID)
	require.Equal(t, res.EntityGUID, input.EntityGUID)
}
////////
////////
// func TestChangeTrackingCreateEvent_Basic(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test change tracking event",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "event-group",
// 		ShortDescription: "Test event",
// 		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
// 		User:             "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "DEPLOYMENT",
// 				Type:     "BASIC",
// 			},
// 			CategoryFields: ChangeTrackingCategoryFieldsInput{
// 				Deployment: ChangeTrackingDeploymentFieldsInput{
// 					Version:   "1.0.0",
// 					Changelog: "test deployment changelog",
// 					Commit:    "12345a",
// 					DeepLink:  "https://example.com/deployment",
// 				},
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.ALLOW_CUSTOM_CATEGORY_OR_TYPE}},
// 	)
// 	require.NoError(t, err)

// 	require.NotNil(t, res)
// 	require.NotNil(t, res.ChangeTrackingEvent)

// 	// Type assert to access the fields
// 	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
// 		require.NotEmpty(t, event.ChangeTrackingId)
// 	}
// }

// func TestChangeTrackingCreateEvent_FeatureFlag(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test feature flag change tracking event",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "feature-flag-group",
// 		ShortDescription: "Feature flag test event",
// 		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
// 		User:             "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "FEATURE_FLAG",
// 				Type:     "BASIC",
// 			},
// 			CategoryFields: ChangeTrackingCategoryFieldsInput{
// 				FeatureFlag: ChangeTrackingFeatureFlagFieldsInput{
// 					FeatureFlagId: "test-feature-flag-123",
// 				},
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.NoError(t, err)

// 	require.NotNil(t, res)
// 	require.NotNil(t, res.ChangeTrackingEvent)

// 	// Type assert to access the fields
// 	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
// 		require.NotEmpty(t, event.ChangeTrackingId)
// 	}
// }

// func TestChangeTrackingCreateEvent_Operational(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test operational change tracking event",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "operational-group",
// 		ShortDescription: "Server reboot event",
// 		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
// 		User:             "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "OPERATIONAL",
// 				Type:     "SERVER_REBOOT",
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.NoError(t, err)

// 	require.NotNil(t, res)
// 	require.NotNil(t, res.ChangeTrackingEvent)

// 	// Type assert to access the fields
// 	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
// 		require.NotEmpty(t, event.ChangeTrackingId)
// 	}
// }

// func TestChangeTrackingCreateEvent_TimestampError(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test change tracking event with invalid timestamp",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "event-group",
// 		ShortDescription: "Test event",
// 		Timestamp:        nrtime.EpochMilliseconds(time.UnixMilli(0)),
// 		User:             "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "DEPLOYMENT",
// 				Type:     "BASIC",
// 			},
// 			CategoryFields: ChangeTrackingCategoryFieldsInput{
// 				Deployment: ChangeTrackingDeploymentFieldsInput{
// 					Version:   "1.0.0",
// 					Changelog: "test deployment changelog",
// 					Commit:    "12345a",
// 					DeepLink:  "https://example.com/deployment",
// 				},
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.Error(t, err)
// 	require.Nil(t, res)
// }

// func TestChangeTrackingCreateEvent_OlderThan24HoursTimestampError(t *testing.T) {
// 	t.Parallel()
// 	now := time.Now()

// 	a := newIntegrationTestClient(t)

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test change tracking event with timestamp older than 24 hours",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "event-group",
// 		ShortDescription: "Test event",
// 		Timestamp: nrtime.EpochMilliseconds(
// 			time.Date(
// 				now.Year(),
// 				now.Month(),
// 				now.Day()-2,
// 				now.Hour()-3,
// 				now.Minute()-30,
// 				0,
// 				0,
// 				time.Local,
// 			),
// 		),
// 		User: "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "DEPLOYMENT",
// 				Type:     "BASIC",
// 			},
// 			CategoryFields: ChangeTrackingCategoryFieldsInput{
// 				Deployment: ChangeTrackingDeploymentFieldsInput{
// 					Version:   "1.0.0",
// 					Changelog: "test deployment changelog",
// 					Commit:    "12345a",
// 					DeepLink:  "https://example.com/deployment",
// 				},
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.Error(t, err)
// 	require.Regexp(t, regexp.MustCompile("not be more than 24 hours"), err.Error())
// 	require.Nil(t, res)
// }

// func TestChangeTrackingCreateEvent_TimestampZeroNanosecondsTest(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)
// 	now := time.Now()

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test change tracking event with zero nanoseconds timestamp",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "event-group",
// 		ShortDescription: "Test event",
// 		Timestamp: nrtime.EpochMilliseconds(
// 			time.Date(
// 				now.Year(),
// 				now.Month(),
// 				now.Day(),
// 				now.Hour()-3,
// 				now.Minute()-30,
// 				0,
// 				0,
// 				time.Local,
// 			),
// 		),
// 		User: "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "DEPLOYMENT",
// 				Type:     "BASIC",
// 			},
// 			CategoryFields: ChangeTrackingCategoryFieldsInput{
// 				Deployment: ChangeTrackingDeploymentFieldsInput{
// 					Version:   "1.0.0",
// 					Changelog: "test deployment changelog",
// 					Commit:    "12345a",
// 					DeepLink:  "https://example.com/deployment",
// 				},
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.NoError(t, err)
// 	require.NotNil(t, res.ChangeTrackingEvent)

// 	// Type assert to access the fields
// 	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
// 		require.NotEmpty(t, event.ChangeTrackingId)
// 	}
// }

// func TestChangeTrackingCreateEvent_TimestampNonZeroNanosecondsTest(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)
// 	now := time.Now()

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test change tracking event with non-zero nanoseconds timestamp",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "event-group",
// 		ShortDescription: "Test event",
// 		Timestamp: nrtime.EpochMilliseconds(
// 			time.Date(
// 				now.Year(),
// 				now.Month(),
// 				now.Day(),
// 				now.Hour()-6,
// 				now.Minute()-30,
// 				0,
// 				231567,
// 				time.Local,
// 			),
// 		),
// 		User: "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "DEPLOYMENT",
// 				Type:     "BASIC",
// 			},
// 			CategoryFields: ChangeTrackingCategoryFieldsInput{
// 				Deployment: ChangeTrackingDeploymentFieldsInput{
// 					Version:   "1.0.0",
// 					Changelog: "test deployment changelog",
// 					Commit:    "12345a",
// 					DeepLink:  "https://example.com/deployment",
// 				},
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.NoError(t, err)
// 	require.NotNil(t, res.ChangeTrackingEvent)

// 	// Type assert to access the fields
// 	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
// 		require.NotEmpty(t, event.ChangeTrackingId)
// 	}
// }

// func TestChangeTrackingCreateEvent_BusinessEvent(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test business event change tracking event",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "business-event-group",
// 		ShortDescription: "Marketing campaign event",
// 		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
// 		User:             "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "BUSINESS_EVENT",
// 				Type:     "MARKETING_CAMPAIGN",
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.NoError(t, err)

// 	require.NotNil(t, res)
// 	require.NotNil(t, res.ChangeTrackingEvent)

// 	// Type assert to access the fields
// 	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
// 		require.NotEmpty(t, event.ChangeTrackingId)
// 	}
// }

// func TestChangeTrackingCreateEvent_DeploymentLifecycle(t *testing.T) {
// 	t.Parallel()

// 	a := newIntegrationTestClient(t)

// 	input := ChangeTrackingCreateEventInput{
// 		Description: "This is a test deployment lifecycle change tracking event",
// 		EntitySearch: ChangeTrackingEntitySearchInput{
// 			Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
// 		},
// 		GroupId:          "deployment-lifecycle-group",
// 		ShortDescription: "Artifact deployment event",
// 		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
// 		User:             "newrelic-go-client",
// 		CategoryAndTypeData: ChangeTrackingCategoryRelatedInput{
// 			Kind: ChangeTrackingCategoryAndTypeInput{
// 				Category: "DEPLOYMENT_LIFECYCLE",
// 				Type:     "ARTIFACT_DEPLOYMENT",
// 			},
// 		},
// 	}

// 	res, err := a.ChangeTrackingCreateEvent(
// 		input,
// 		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
// 	)
// 	require.NoError(t, err)

// 	require.NotNil(t, res)
// 	require.NotNil(t, res.ChangeTrackingEvent)

// 	// Type assert to access the fields
// 	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
// 		require.NotEmpty(t, event.ChangeTrackingId)
// 	}
// }



////////

func TestChangeTrackingCreateEvent_Basic(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingCreateEventInput{
		Description: "This is a test change tracking event",
		EntitySearch: ChangeTrackingEntitySearchInput{
			// Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
			Query: "name = 'My PHP Application'",
		},
		ShortDescription: "Test event",
		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
		User:             "newrelic-go-client",
		CategoryAndTypeData: &ChangeTrackingCategoryRelatedInput{
			Kind: &ChangeTrackingCategoryAndTypeInput{
				Category: "DEPLOYMENT",
				Type:     "BASIC",
			},
			CategoryFields: &ChangeTrackingCategoryFieldsInput{
				Deployment: &ChangeTrackingDeploymentFieldsInput{
					Version:   "1.0.0",
					Changelog: "test deployment changelog",
					Commit:    "12345a",
					DeepLink:  "https://example.com/deployment",
				},
			},
		},
	}

	res, err := a.ChangeTrackingCreateEvent(
		input,
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.ALLOW_CUSTOM_CATEGORY_OR_TYPE}},
	)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.NotNil(t, res.ChangeTrackingEvent)

	// Type assert to access the fields
	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
		require.NotEmpty(t, event.ChangeTrackingId)
	}
}

func TestChangeTrackingCreateEvent_FeatureFlag(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingCreateEventInput{
		Description: "This is a test feature flag change tracking event",
		EntitySearch: ChangeTrackingEntitySearchInput{
			// Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
			Query: "name = 'My PHP Application'",

		},
		GroupId:          "feature-flag-group",
		ShortDescription: "Feature flag test event",
		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
		User:             "newrelic-go-client",
		CategoryAndTypeData: &ChangeTrackingCategoryRelatedInput{
			Kind: &ChangeTrackingCategoryAndTypeInput{
				Category: "FEATURE_FLAG",
				Type:     "BASIC",
			},
			CategoryFields: &ChangeTrackingCategoryFieldsInput{
				FeatureFlag: &ChangeTrackingFeatureFlagFieldsInput{
					FeatureFlagId: "test-feature-flag-123",
				},
			},
		},
	}

	res, err := a.ChangeTrackingCreateEvent(
		input,
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
	)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.NotNil(t, res.ChangeTrackingEvent)

	// Type assert to access the fields
	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
		require.NotEmpty(t, event.ChangeTrackingId)
	}
}

func TestChangeTrackingCreateEvent_Operational(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingCreateEventInput{
		Description: "This is a test operational change tracking event",
		EntitySearch: ChangeTrackingEntitySearchInput{
			// Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
			Query: "name = 'My PHP Application'",
		},
		GroupId:          "operational-group",
		ShortDescription: "Server reboot event",
		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
		User:             "newrelic-go-client",
		CategoryAndTypeData: &ChangeTrackingCategoryRelatedInput{
			Kind: &ChangeTrackingCategoryAndTypeInput{
				Category: "OPERATIONAL",
				Type:     "SERVER_REBOOT",
			},
		},
	}

	res, err := a.ChangeTrackingCreateEvent(
		input,
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
	)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.NotNil(t, res.ChangeTrackingEvent)

	// Type assert to access the fields
	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
		require.NotEmpty(t, event.ChangeTrackingId)
	}
}

func TestChangeTrackingCreateEvent_TimestampError(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingCreateEventInput{
		Description: "This is a test change tracking event with invalid timestamp",
		EntitySearch: ChangeTrackingEntitySearchInput{
			// Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
			Query: "name = 'My PHP Application'",
		},
		GroupId:          "event-group",
		ShortDescription: "Test event",
		Timestamp:        nrtime.EpochMilliseconds(time.UnixMilli(0)),
		User:             "newrelic-go-client",
		CategoryAndTypeData: &ChangeTrackingCategoryRelatedInput{
			Kind: &ChangeTrackingCategoryAndTypeInput{
				Category: "DEPLOYMENT",
				Type:     "BASIC",
			},
			CategoryFields: &ChangeTrackingCategoryFieldsInput{
				Deployment: &ChangeTrackingDeploymentFieldsInput{
					Version:   "1.0.0",
					Changelog: "test deployment changelog",
					Commit:    "12345a",
					DeepLink:  "https://example.com/deployment",
				},
			},
		},
	}

	res, err := a.ChangeTrackingCreateEvent(
		input,
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
	)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestChangeTrackingCreateEvent_OlderThan24HoursTimestampError(t *testing.T) {
	t.Parallel()
	now := time.Now()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingCreateEventInput{
		Description: "This is a test change tracking event with timestamp older than 24 hours",
		EntitySearch: ChangeTrackingEntitySearchInput{
			// Query: fmt.Sprintf("name = '%s'", testhelpers.IntegrationTestApplicationEntityNameNew),
			Query: "name = 'My PHP Application'",
		},
		GroupId:          "event-group",
		ShortDescription: "Test event",
		Timestamp: nrtime.EpochMilliseconds(
			time.Date(
				now.Year(),
				now.Month(),
				now.Day()-2,
				now.Hour()-3,
				now.Minute()-30,
				0,
				0,
				time.Local,
			),
		),
		User: "newrelic-go-client",
		CategoryAndTypeData: &ChangeTrackingCategoryRelatedInput{
			Kind: &ChangeTrackingCategoryAndTypeInput{
				Category: "DEPLOYMENT",
				Type:     "BASIC",
			},
			CategoryFields: &ChangeTrackingCategoryFieldsInput{
				Deployment: &ChangeTrackingDeploymentFieldsInput{
					Version:   "1.0.0",
					Changelog: "test deployment changelog",
					Commit:    "12345a",
					DeepLink:  "https://example.com/deployment",
				},
			},
		},
	}

	res, err := a.ChangeTrackingCreateEvent(
		input,
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
	)
	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile("not be more than 24 hours"), err.Error())
	require.Nil(t, res)
}

func TestChangeTrackingCreateEvent_CustomAttributes(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	// Create custom attributes map
	customAttributes := map[string]interface{}{
		"environment":    "staging",
		"region":         "us-east-1",
		"cloud_vendor":   "aws",
		"isProd":         false,
		"instance_count": 3,
		"deploy_time":    10.5,
	}

	input := ChangeTrackingCreateEventInput{
		Description: "This is a test change tracking event with custom attributes",
		EntitySearch: ChangeTrackingEntitySearchInput{
			Query: "name = 'My PHP Application'",
		},
		CustomAttributes: ChangeTrackingRawCustomAttributesMap(customAttributes),
		GroupId:          "custom-attributes-group",
		ShortDescription: "Test event with custom attributes",
		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
		User:             "newrelic-go-client",
		CategoryAndTypeData: &ChangeTrackingCategoryRelatedInput{
			Kind: &ChangeTrackingCategoryAndTypeInput{
				Category: "DEPLOYMENT",
				Type:     "BASIC",
			},
			CategoryFields: &ChangeTrackingCategoryFieldsInput{
				Deployment: &ChangeTrackingDeploymentFieldsInput{
					Version:   "2.0.0",
					Changelog: "test deployment with custom attributes",
					Commit:    "abc123",
					DeepLink:  "https://example.com/deployment/custom",
				},
			},
		},
	}

	res, err := a.ChangeTrackingCreateEvent(
		input,
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
	)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.NotNil(t, res.ChangeTrackingEvent)

	// Type assert to access the fields
	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
		require.NotEmpty(t, event.ChangeTrackingId)
		require.NotNil(t, event.CustomAttributes)
	}
}

func TestChangeTrackingCreateEvent_AllowCustomCategoryType(t *testing.T) {
	t.Parallel()

	a := newIntegrationTestClient(t)

	input := ChangeTrackingCreateEventInput{
		Description: "This is a test change tracking event with ALLOW_CUSTOM_CATEGORY_OR_TYPE flag",
		EntitySearch: ChangeTrackingEntitySearchInput{
			Query: "name = 'My PHP Application'",
		},
		GroupId:          "allow-custom-group",
		ShortDescription: "Test event with allow custom flag",
		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
		User:             "newrelic-go-client",
		CategoryAndTypeData: &ChangeTrackingCategoryRelatedInput{
			Kind: &ChangeTrackingCategoryAndTypeInput{
				Category: "DEPLOYMENT",
				Type:     "BASIC",
			},
			CategoryFields: &ChangeTrackingCategoryFieldsInput{
				Deployment: &ChangeTrackingDeploymentFieldsInput{
					Version:   "3.0.0",
					Changelog: "test deployment with allow custom flag",
					Commit:    "def456",
					DeepLink:  "https://example.com/deployment/allow-custom",
				},
			},
		},
	}

	res, err := a.ChangeTrackingCreateEvent(
		input,
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.ALLOW_CUSTOM_CATEGORY_OR_TYPE}},
	)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.NotNil(t, res.ChangeTrackingEvent)

	// Type assert to access the fields
	if event, ok := res.ChangeTrackingEvent.(*ChangeTrackingEvent); ok {
		require.NotEmpty(t, event.ChangeTrackingId)
		require.Equal(t, "DEPLOYMENT", event.Category)
	}
}

func newIntegrationTestClient(t *testing.T) Changetracking {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
