//go:build unit
// +build unit

package entities

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
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
		{
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
		{
			"key":   "tagKey",
			"value": "tagValue",
		},
		{
			"key":   "tagKey2",
			"value": "tagValue2",
		},
		{
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
			{
				"key":   "tagKey",
				"value": "tagValue",
			},
			{
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

// TestEntityAlertViolationInt_LargeID verifies that violationIds exceeding the int32 range
// (as returned by the Alerts backend for 16-digit IDs) unmarshal without overflow or error.
func TestEntityAlertViolationInt_LargeID(t *testing.T) {
	t.Parallel()

	// These are real violation IDs from the bug report (NR-561784) that exceed int32 max (2,147,483,647).
	cases := []struct {
		json string
		want EntityAlertViolationInt
	}{
		{`{"violationId": 1773950568354017}`, EntityAlertViolationInt(1773950568354017)},
		{`{"violationId": 1777567939793015}`, EntityAlertViolationInt(1777567939793015)},
		{`{"violationId": 2147483647}`, EntityAlertViolationInt(2147483647)}, // int32 max — still works
		{`{"violationId": 2147483648}`, EntityAlertViolationInt(2147483648)}, // one above int32 max
	}

	for _, tc := range cases {
		var v EntityAlertViolation
		err := json.Unmarshal([]byte(tc.json), &v)
		require.NoError(t, err, "unmarshal failed for input %s", tc.json)
		require.Equal(t, tc.want, v.ViolationId)
	}
}

// TestGetEntityWithContext_LargeViolationId simulates the full Terraform read path for
// newrelic_synthetics_monitor: a mock NerdGraph server returns a SyntheticMonitorEntity
// with a 16-digit violationId, and GetEntityWithContext must parse it without error.
// This reproduces the crash reported in NR-561784 on 32-bit platforms where
// EntityAlertViolationInt was int (32-bit) instead of int64.
func TestGetEntityWithContext_LargeViolationId(t *testing.T) {
	t.Parallel()

	// Real violation IDs from the bug report that overflow int32.
	const violationID int64 = 1773950568354017

	mockResponse := `{
		"data": {
			"actor": {
				"entity": {
					"__typename": "SyntheticMonitorEntity",
					"guid": "MTAwMjMzMjUwOHxTWU5USHxNT05JVE9SfDEyMzQ1Njc4",
					"name": "test-monitor",
					"recentAlertViolations": [
						{
							"violationId": 1773950568354017,
							"label": "Synthetics monitor test-monitor",
							"level": "CRITICAL",
							"alertSeverity": "CRITICAL"
						}
					]
				}
			}
		}
	}`

	ts := mock.NewMockServer(t, mockResponse, http.StatusOK)
	defer ts.Close()

	cfg := mock.NewTestConfig(t, ts)
	client := New(cfg)

	entity, err := client.GetEntityWithContext(context.Background(), "MTAwMjMzMjUwOHxTWU5USHxNT05JVE9SfDEyMzQ1Njc4")
	require.NoError(t, err)
	require.NotNil(t, entity)

	synthEntity, ok := (*entity).(*SyntheticMonitorEntity)
	require.True(t, ok, "expected *SyntheticMonitorEntity, got %T", *entity)
	require.Len(t, synthEntity.RecentAlertViolations, 1)
	require.Equal(t, EntityAlertViolationInt(violationID), synthEntity.RecentAlertViolations[0].ViolationId)
}
