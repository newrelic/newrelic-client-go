//go:build integration
// +build integration

package changetracking

import (
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

	var customAttributes = map[string]string{
		"test":  "123",
		"test2": "456",
	}

	input := ChangeTrackingDeploymentInput{
		Changelog:        "test",
		Commit:           "12345a",
		CustomAttributes: &customAttributes,
		DeepLink:         "newrelic-client-go",
		DeploymentType:   ChangeTrackingDeploymentTypeTypes.BASIC,
		Description:      "This is a test description",
		EntityGUID:       common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID),
		GroupId:          "deployment",
		Timestamp:        nrtime.EpochMilliseconds(time.Now()),
		User:             "newrelic-go-client",
		Version:          "0.0.1",
	}

	res, err := a.ChangeTrackingCreateDeployment(input)
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

	res, err := a.ChangeTrackingCreateDeployment(input)
	require.Error(t, err)
	require.Nil(t, res)
}

func newIntegrationTestClient(t *testing.T) Changetracking {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return New(tc)
}
