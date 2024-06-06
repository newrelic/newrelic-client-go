//go:build integration
// +build integration

package changetracking

import (
	"encoding/json"
	"fmt"
	"log"
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

func TestChangeTrackingCreateDeployment_CustomAttributes(t *testing.T) {
	skipMsg := fmt.Sprintf("Skipping %s until custom attributes are out of limited preview.", t.Name())
	t.Skip(skipMsg)
	t.Parallel()

	a := newIntegrationTestClient(t)

	var customAttributes = `{"a":"1","b":"two","c":"1.5","d":"true"}`
	attrs := make(map[string]interface{})
	err := json.Unmarshal([]byte(customAttributes), &attrs)
	if err != nil {
		log.Fatal(err)
	}

	input := ChangeTrackingDeploymentInput{
		Changelog:        "test",
		Commit:           "12345a",
		CustomAttributes: ChangeTrackingRawCustomAttributesMap(attrs),
		DeepLink:         "newrelic-client-go",
		DeploymentType:   ChangeTrackingDeploymentTypeTypes.BASIC,
		Description:      "This is a test description",
		EntityGUID:       common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew),
		GroupId:          "deployment",
		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
		User:             "newrelic-go-client",
		Version:          "0.0.1",
	}

	res, err := a.ChangeTrackingCreateDeployment(
		ChangeTrackingDataHandlingRules{ValidationFlags: []ChangeTrackingValidationFlag{ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
		input,
	)
	require.NoError(t, err)

	require.NotNil(t, res)
	require.Equal(t, res.EntityGUID, input.EntityGUID)
}

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

func newIntegrationTestClient(t *testing.T) Changetracking {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
