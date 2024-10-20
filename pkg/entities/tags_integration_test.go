//go:build integration
// +build integration

package entities

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationListTags(t *testing.T) {
	t.Parallel()

	var (
		// GUID of Dummy App
		testGUID = common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew)
	)

	client := newIntegrationTestClient(t)

	actual, err := client.ListTags(testGUID)
	require.NoError(t, err)
	require.Greater(t, len(actual), 0)

	actual, err = client.ListAllTags(testGUID)
	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

func TestIntegrationTaggingAddTagsToEntityAndGetTags(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUIDNew)
	)

	client := newIntegrationTestClient(t)

	tags := []TaggingTagInput{
		{
			Key:    "test",
			Values: []string{"value"},
		},
	}
	result, err := client.TaggingAddTagsToEntity(testGUID, tags)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 0, len(result.Errors))

	actual, err := client.GetTagsForEntity(testGUID)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)

	actual, err = client.GetTagsForEntityMutable(testGUID)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)

}

func TestIntegrationTaggingReplaceTagsOnEntity(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)
	)

	client := newIntegrationTestClient(t)

	tags := []TaggingTagInput{
		{
			Key:    "test",
			Values: []string{"value"},
		},
	}
	result, err := client.TaggingReplaceTagsOnEntity(testGUID, tags)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 0, len(result.Errors))

}

func TestIntegrationDeleteTags(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)
	)

	client := newIntegrationTestClient(t)

	tagKeys := []string{"test"}
	result, err := client.TaggingDeleteTagFromEntity(testGUID, tagKeys)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 0, len(result.Errors))
}

func TestIntegrationDeleteTagValues(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)
	)

	client := newIntegrationTestClient(t)

	tagValues := []TaggingTagValueInput{
		{
			Key:   "test",
			Value: "value",
		},
	}
	result, err := client.TaggingDeleteTagValuesFromEntity(testGUID, tagValues)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 0, len(result.Errors))
}

func TestIntegrationEntityTagsReservedKeysMutation(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID(testhelpers.IntegrationTestApplicationEntityGUID)
	)

	client := newIntegrationTestClient(t)

	// Test: To add a reserved key(immutable key)
	tags := []TaggingTagInput{
		{
			Key:    "account",
			Values: []string{"Random-name"},
		},
	}
	result, err := client.TaggingAddTagsToEntity(testGUID, tags)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Greater(t, len(result.Errors), 0)
	message := result.Errors[0].Message
	match, er := regexp.MatchString("reserved", message)
	require.NoError(t, er)
	require.True(t, match)

	// Test: To update a reserved key(immutable key)
	result, err = client.TaggingReplaceTagsOnEntity(testGUID, tags)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Greater(t, len(result.Errors), 0)
	message = result.Errors[0].Message
	match, er = regexp.MatchString("reserved", message)
	require.NoError(t, er)
	require.True(t, match)
}
