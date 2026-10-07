//go:build unit
// +build unit

package entities

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestAddCursorArg_QueryStrings verifies the derived *WithCursor query strings
// declare the $cursor variable and forward it into results(cursor: $cursor).
// If tutone regeneration reshapes either source query, the addCursorArg helper
// panics at package init; this test additionally guards the resulting text.
func TestAddCursorArg_QueryStrings(t *testing.T) {
	t.Parallel()

	for name, q := range map[string]string{
		"getEntitySearchByQueryWithCursor": getEntitySearchByQueryWithCursor,
		"getEntitySearchWithCursor":        getEntitySearchWithCursor,
	} {
		require.Contains(t, q, "$cursor: String,", "%s missing $cursor variable declaration", name)
		require.Contains(t, q, "results(cursor: $cursor) {", "%s does not forward cursor to results()", name)
		// The pre-existing top-level `results {` must be replaced, not duplicated.
		require.NotContains(t, q, "\n\tresults {\n", "%s still contains un-cursor'd results { block", name)
	}
}

// captureRequestServer is a testing mock that captures each inbound GraphQL
// request body and returns a canned response. It's local to this file because
// the shared testhelpers.NewMockServer discards the request body.
type captureRequestServer struct {
	*httptest.Server
	requests [][]byte
}

func newCaptureRequestServer(t *testing.T, response string) *captureRequestServer {
	t.Helper()
	crs := &captureRequestServer{}
	crs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		crs.requests = append(crs.requests, body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(response))
		require.NoError(t, err)
	}))
	return crs
}

const entitySearchCursorResponse = `{
	"data": {
		"actor": {
			"entitySearch": {
				"count": 2,
				"query": "domain = 'APM'",
				"results": {
					"entities": [
						{
							"__typename": "ApmApplicationEntityOutline",
							"guid": "MTIzfEFQTXxBUFBMSUNBVElPTnwxMTE=",
							"name": "app-1"
						},
						{
							"__typename": "ApmApplicationEntityOutline",
							"guid": "MTIzfEFQTXxBUFBMSUNBVElPTnwyMjI=",
							"name": "app-2"
						}
					],
					"nextCursor": "PAGE-2-CURSOR"
				}
			}
		}
	}
}`

// TestGetEntitySearchByQueryWithCursor reproduces Scott's Slack question:
// prior to this method the caller had no way to pass a nextCursor value back
// in for pagination. This test asserts (1) the returned NextCursor is parsed,
// (2) an empty cursor is sent as `null`, and (3) a non-empty cursor is
// forwarded on the wire.
func TestGetEntitySearchByQueryWithCursor(t *testing.T) {
	t.Parallel()

	crs := newCaptureRequestServer(t, entitySearchCursorResponse)
	defer crs.Close()

	client := New(mock.NewTestConfig(t, crs.Server))

	// First page — empty cursor should serialize to `null`.
	got, err := client.GetEntitySearchByQueryWithCursor(
		EntitySearchOptions{}, "domain = 'APM'", nil, "",
	)
	require.NoError(t, err)
	require.Equal(t, "PAGE-2-CURSOR", got.Results.NextCursor)

	// Second page — feed the cursor back in.
	_, err = client.GetEntitySearchByQueryWithCursor(
		EntitySearchOptions{}, "domain = 'APM'", nil, "PAGE-2-CURSOR",
	)
	require.NoError(t, err)

	require.Len(t, crs.requests, 2)
	require.Contains(t, string(crs.requests[0]), `"cursor":null`, "empty cursor should be sent as null")
	require.Contains(t, string(crs.requests[1]), `"cursor":"PAGE-2-CURSOR"`, "cursor value should round-trip")

	// Sanity: the outbound query must declare and forward $cursor.
	for _, body := range crs.requests {
		require.True(t, strings.Contains(string(body), "$cursor: String"), "outbound query missing $cursor declaration")
		require.True(t, strings.Contains(string(body), "results(cursor: $cursor)"), "outbound query does not forward cursor to results()")
	}
}

// TestGetEntitySearchWithCursor mirrors TestGetEntitySearchByQueryWithCursor
// for the queryBuilder variant.
func TestGetEntitySearchWithCursor(t *testing.T) {
	t.Parallel()

	crs := newCaptureRequestServer(t, entitySearchCursorResponse)
	defer crs.Close()

	client := New(mock.NewTestConfig(t, crs.Server))

	qb := EntitySearchQueryBuilder{Domain: EntitySearchQueryBuilderDomainTypes.APM}

	got, err := client.GetEntitySearchWithCursor(
		EntitySearchOptions{}, "", qb, nil, nil, "",
	)
	require.NoError(t, err)
	require.Equal(t, "PAGE-2-CURSOR", got.Results.NextCursor)

	_, err = client.GetEntitySearchWithCursor(
		EntitySearchOptions{}, "", qb, nil, nil, "PAGE-2-CURSOR",
	)
	require.NoError(t, err)

	require.Len(t, crs.requests, 2)
	require.Contains(t, string(crs.requests[0]), `"cursor":null`)
	require.Contains(t, string(crs.requests[1]), `"cursor":"PAGE-2-CURSOR"`)
}
