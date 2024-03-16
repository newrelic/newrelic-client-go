//go:build unit

package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindTagByKey(t *testing.T) {
	t.Parallel()

	entityTags := []EntityTag{
		{
			Key:    "test",
			Values: []string{"someTag"},
		},
	}

	result := FindTagByKey(entityTags, "test")
	require.Equal(t, []string{"someTag"}, result)
}

func TestFindTagByKeyNotFound(t *testing.T) {
	t.Parallel()

	entityTags := []EntityTag{
		{
			Key:    "test",
			Values: []string{"someTag"},
		},
	}

	result := FindTagByKey(entityTags, "notFound")
	require.Empty(t, result)
}

func TestBuildTagsQueryFragment_SingleTag(t *testing.T) {
	t.Parallel()

	expected := "tags.`tagKey` = 'tagValue'"

	tags := []map[string]string{
		map[string]string{
			"key":   "tagKey",
			"value": "tagValue",
		},
	}

	result := BuildTagsQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildTagsQueryFragment_MultipleTags(t *testing.T) {
	t.Parallel()

	expected := "tags.`tagKey` = 'tagValue' AND tags.`tagKey2` = 'tagValue2' AND tags.`tagKey3` = 'tagValue3'"

	tags := []map[string]string{
		map[string]string{
			"key":   "tagKey",
			"value": "tagValue",
		},
		map[string]string{
			"key":   "tagKey2",
			"value": "tagValue2",
		},
		map[string]string{
			"key":   "tagKey3",
			"value": "tagValue3",
		},
	}

	result := BuildTagsQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildTagsQueryFragment_EmptyTags(t *testing.T) {
	t.Parallel()

	expected := ""
	tags := []map[string]string{}

	result := BuildTagsQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildEntitySearchQuery(t *testing.T) {
	t.Parallel()

	tags := []map[string]string{}

	// Name only
	expected := "name = 'Dummy App'"
	result := BuildEntitySearchQuery("Dummy App", "", "", tags)
	require.Equal(t, expected, result)

	// Name & Domain
	expected = "name = 'Dummy App' AND domain = 'APM'"
	result = BuildEntitySearchQuery("Dummy App", "APM", "", tags)
	require.Equal(t, expected, result)

	// Name, domain, and type
	expected = "name = 'Dummy App' AND domain = 'APM' AND type = 'APPLICATION'"
	result = BuildEntitySearchQuery("Dummy App", "APM", "APPLICATION", tags)
	require.Equal(t, expected, result)

	// Name, domain, type, and tags
	expected = "name = 'Dummy App' AND domain = 'APM' AND type = 'APPLICATION' AND tags.`tagKey` = 'tagValue' AND tags.`tagKey2` = 'tagValue2'"
	tags = []map[string]string{
		map[string]string{
			"key":   "tagKey",
			"value": "tagValue",
		},
		map[string]string{
			"key":   "tagKey2",
			"value": "tagValue2",
		},
	}
	result = BuildEntitySearchQuery("Dummy App", "APM", "APPLICATION", tags)
	require.Equal(t, expected, result)
}
