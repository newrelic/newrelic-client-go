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

func TestBuildTagsNrqlQueryFragment_SingleTag(t *testing.T) {
	t.Parallel()

	expected := "tags.`tagKey` = 'tagValue'"

	tags := []map[string]string{
		map[string]string{
			"key":   "tagKey",
			"value": "tagValue",
		},
	}

	result := BuildTagsNrqlQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildTagsNrqlQueryFragment_MultipleTags(t *testing.T) {
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

	result := BuildTagsNrqlQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildTagsNrqlQueryFragment_EmptyTags(t *testing.T) {
	t.Parallel()

	expected := ""
	tags := []map[string]string{}

	result := BuildTagsNrqlQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildEntitySearchNrqlQuery(t *testing.T) {
	t.Parallel()

	// Name only
	expected := "name LIKE 'Dummy App'"
	searchParams := EntitySearchParams{
		Name: "Dummy App",
	}
	result := BuildEntitySearchNrqlQuery(searchParams)
	require.Equal(t, expected, result)

	// Is reporting = true
	isReporting := true
	expected = "reporting = 'true'"
	searchParams = EntitySearchParams{
		IsReporting: &isReporting,
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Equal(t, expected, result)

	// Is reporting = false
	isReporting = false
	expected = "reporting = 'false'"
	searchParams = EntitySearchParams{
		IsReporting: &isReporting,
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Equal(t, expected, result)

	// Is reporting = true and type
	isReporting = true
	searchParams = EntitySearchParams{
		Type:        "APPLICATION",
		IsReporting: &isReporting,
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Contains(t, result, "reporting = 'true'")
	require.Contains(t, result, "type = 'APPLICATION'")

	// Is reporting = false and type
	isReporting = false
	searchParams = EntitySearchParams{
		Type:        "APPLICATION",
		IsReporting: &isReporting,
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Contains(t, result, "reporting = 'false'")
	require.Contains(t, result, "type = 'APPLICATION'")

	// Case-sensitive search (applies to `name` only)
	expected = "name = 'Dummy App'"
	searchParams = EntitySearchParams{
		Name:            "Dummy App",
		IsCaseSensitive: true,
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Equal(t, expected, result)

	// Name & Domain
	searchParams = EntitySearchParams{
		Name:   "Dummy App",
		Domain: "APM",
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Contains(t, result, "name LIKE 'Dummy App'")
	require.Contains(t, result, "domain = 'APM'")

	// Name, domain, and type
	searchParams = EntitySearchParams{
		Name:   "Dummy App",
		Domain: "APM",
		Type:   "APPLICATION",
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Contains(t, result, "name LIKE 'Dummy App'")
	require.Contains(t, result, "domain = 'APM'")
	require.Contains(t, result, "type = 'APPLICATION'")

	// Name, domain, type, and tags
	searchParams = EntitySearchParams{
		Name:   "Dummy App",
		Domain: "APM",
		Type:   "APPLICATION",
		Tags: []map[string]string{
			map[string]string{
				"key":   "tagKey",
				"value": "tagValue",
			},
			map[string]string{
				"key":   "tagKey2",
				"value": "tagValue2",
			},
		},
	}
	result = BuildEntitySearchNrqlQuery(searchParams)
	require.Contains(t, result, "name LIKE 'Dummy App'")
	require.Contains(t, result, "domain = 'APM'")
	require.Contains(t, result, "type = 'APPLICATION'")
	require.Contains(t, result, " AND tags.`tagKey` = 'tagValue' AND tags.`tagKey2` = 'tagValue2'")
}
