//go:build integration
// +build integration

package entities

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/common"
)

func TestIntegrationListTags(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID("MjUwODI1OXxBUE18QVBQTElDQVRJT058MjA0MjYxMzY4")
	)

	client := newIntegrationTestClient(t)

	actual, err := client.ListTags(testGUID)
	require.NoError(t, err)
	require.Greater(t, len(actual), 0)

	actual, err = client.ListAllTags(testGUID)
	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

func TestIntegrationGetTagsForEntity(t *testing.T) {
	t.Parallel()

	var (
		// GUID of Dummy App
		testGUID = common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")
	)

	client := newIntegrationTestClient(t)

	actual, err := client.GetTagsForEntity(testGUID)
	require.NoError(t, err)
	require.Greater(t, len(actual), 0)

	actual, err = client.GetTagsForEntityMutable(testGUID)
	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

func TestIntegrationTaggingAddTagsToEntity(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")
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
}

func TestIntegrationTaggingReplaceTagsOnEntity(t *testing.T) {
	t.Parallel()

	var (
		testGUID = common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")
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
		testGUID = common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")
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
		testGUID = common.EntityGUID("MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1")
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
