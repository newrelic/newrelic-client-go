//go:build integration
// +build integration

// These integration tests exercise the full Notebook lifecycle against a real
// New Relic account. They are the primary correctness signal for the package,
// since the two-API architecture (Blob Storage for content + rename, NerdGraph
// for metadata) is only observable when both endpoints are reachable.
//
// Credentials: uses the shared Fleet test account (NEW_RELIC_FLEET_TEST_API_KEY
// + NEW_RELIC_FLEET_TEST_ORGANIZATION_ID) which is the only account currently
// entitled for Notebooks. If those env vars are absent the tests skip.
//
// Every mutation in these tests is followed by an independent verification
// call against the platform - either a Blob API GET or a NerdGraph
// entity(id:)/entitySearch. Assertions on the mutation's response envelope
// alone would not prove the resource actually landed on the server.

package notebooks

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// blankNotebookContent is the minimum payload the Blob API accepts on create.
// A blank notebook renders empty in the UI and gives the test a stable
// baseline to update from.
func blankNotebookContent() map[string]interface{} {
	return map[string]interface{}{
		"version": "1",
		"blocks":  []interface{}{},
	}
}

// markdownWidgetContent produces a small notebook body containing one
// markdown widget. Used to prove the content actually changed on update.
func markdownWidgetContent(text string) map[string]interface{} {
	return map[string]interface{}{
		"version": "1",
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "widget",
				"content": map[string]interface{}{
					"type": "visualization",
					"id":   "viz.markdown",
					"props": map[string]interface{}{
						"text": text,
					},
				},
			},
		},
	}
}

// TestIntegrationNotebookLifecycle drives create -> read -> search -> update
// content -> rename -> delete, verifying platform state after each mutation
// rather than trusting the response envelope. Every ID used is captured from
// a live API call - no hardcoded GUIDs.
func TestIntegrationNotebookLifecycle(t *testing.T) {
	t.Parallel()

	if _, err := mock.GetFleetTestAccountID(); err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)
	orgID := testOrganizationID

	notebookName := fmt.Sprintf("integration-test-notebook-%d", time.Now().UnixNano())
	initialContent := blankNotebookContent()

	// Cleanup discipline: whatever path the test takes (fatal failure, panic,
	// early return) this defer runs and tries to remove any live notebook so
	// integration accounts never accumulate orphaned resources. If create
	// never happened entityGUID is empty and we skip; if delete already ran
	// deleted is true and we skip.
	//
	// Delete is retried a few times with backoff because the Blob API has
	// been observed to return transient 5xx during high load and a leaked
	// notebook is worse than a slow teardown. Any final failure is logged
	// (not fatal - the test's own assertion is the primary signal) so CI
	// operators can spot leaks even when the run appears green.
	var entityGUID string
	deleted := false
	defer func() {
		if deleted || entityGUID == "" {
			return
		}
		cleanupNotebook(t, client, orgID, entityGUID)
	}()

	// --- Create -----------------------------------------------------------
	createResp, err := client.CreateNotebook(orgID, notebookName, initialContent)
	require.NoError(t, err, "create notebook")
	require.NotNil(t, createResp)
	require.NotEmpty(t, createResp.EntityGUID, "expected entityGuid on create response")
	require.NotEmpty(t, createResp.BlobID, "expected blobId on create response")
	entityGUID = createResp.EntityGUID
	firstBlobID := createResp.BlobID

	// Verify the notebook exists on the platform by fetching its metadata.
	notebook := requireNotebookMetadata(t, client, entityGUID)
	require.Equal(t, notebookName, notebook.Name, "notebook name persisted incorrectly")
	require.Equal(t, "NOTEBOOK", notebook.Type, "expected type = NOTEBOOK on the server")
	require.Equal(t, orgID, notebook.Scope.ID, "notebook should be scoped to the requesting organization")
	require.Equal(t, "ORGANIZATION", string(notebook.Scope.Type), "Blob API always creates ORGANIZATION-scoped notebooks")
	require.GreaterOrEqual(t, notebook.Metadata.Version, 1, "metadata.version should be >= 1 after create")
	initialVersion := notebook.Metadata.Version

	// Verify the initial content round-trips through the Blob API.
	assertContentMatches(t, client, orgID, entityGUID, initialContent)

	// --- List / search verification --------------------------------------
	// NerdGraph's entitySearch can lag behind Blob API writes by seconds; poll
	// briefly so the test is stable without being flaky.
	requireNotebookAppearsInSearch(t, client, entityGUID, notebookName)

	// --- Update content ---------------------------------------------------
	updatedContent := markdownWidgetContent("# Updated by TestIntegrationNotebookLifecycle")
	updateResp, err := client.UpdateNotebookContent(orgID, entityGUID, updatedContent)
	require.NoError(t, err, "update notebook content")
	require.NotNil(t, updateResp)
	require.Equal(t, entityGUID, updateResp.EntityGUID, "update response should reference the same entity")
	require.NotEmpty(t, updateResp.BlobID, "update should return a new blobId")
	require.NotEqual(t, firstBlobID, updateResp.BlobID, "blobId should rotate on content update")

	// Content was actually rewritten on the server.
	assertContentMatches(t, client, orgID, entityGUID, updatedContent)

	// Version incremented in metadata. The server bumps once per content POST
	// so we expect strictly greater than the initial version.
	afterUpdate := requireNotebookMetadata(t, client, entityGUID)
	require.Greater(t, afterUpdate.Metadata.Version, initialVersion, "metadata.version should increment after content update")
	require.Equal(t, notebookName, afterUpdate.Name, "name should not change on a content-only update")

	// --- Rename -----------------------------------------------------------
	// Rename is atomic with a content re-POST (Blob API has no name-only path).
	renamedName := fmt.Sprintf("%s-renamed", notebookName)
	renameResp, err := client.RenameNotebook(orgID, entityGUID, renamedName, updatedContent)
	require.NoError(t, err, "rename notebook")
	require.NotNil(t, renameResp)
	require.Equal(t, entityGUID, renameResp.EntityGUID)

	// The rename is only visible via NerdGraph - Blob GET returns content, not
	// name. Poll briefly to allow propagation.
	requireEventually(t, "rename to propagate", func() error {
		nb := requireNotebookMetadata(t, client, entityGUID)
		if nb.Name != renamedName {
			return fmt.Errorf("notebook name is still %q, want %q", nb.Name, renamedName)
		}
		return nil
	})

	// --- Delete -----------------------------------------------------------
	require.NoError(t, client.DeleteNotebook(orgID, entityGUID), "delete notebook")
	deleted = true

	// Blob GET should return not-found. The client currently surfaces this as
	// a wrapped error containing the entity GUID.
	_, err = client.GetNotebookContent(orgID, entityGUID)
	require.Error(t, err, "GET content after delete should error")
	require.Contains(t, err.Error(), "not found", "expected a 'not found' style error, got: %v", err)

	// NerdGraph should also stop returning the entity. NerdGraph propagation
	// after a delete has been observed to take a few seconds, so poll.
	requireEventually(t, "delete to propagate to NerdGraph", func() error {
		nb, err := client.GetNotebook(entityGUID)
		if err != nil {
			// "Entity not found" surfaces as an error on the client side -
			// that's the state we actually want.
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if nb == nil {
			return nil
		}
		return fmt.Errorf("GetNotebook still returns %s after delete", nb.ID)
	})
}

// TestIntegrationNotebookRejectsMissingArgs is a fast unit-shaped sanity check
// for the argument validation in the Blob API surface. It doesn't hit the
// network - the client short-circuits on empty inputs - so it stays under the
// integration tag purely to co-locate with the rest of the suite.
func TestIntegrationNotebookRejectsMissingArgs(t *testing.T) {
	t.Parallel()

	if _, err := mock.GetFleetTestAccountID(); err != nil {
		t.Skipf("%s", err)
	}
	client := newIntegrationTestClient(t)

	_, err := client.CreateNotebook("", "name", blankNotebookContent())
	assert.Error(t, err, "empty organization ID should fail")

	_, err = client.CreateNotebook(testOrganizationID, "", blankNotebookContent())
	assert.Error(t, err, "empty name should fail")

	_, err = client.CreateNotebook(testOrganizationID, "name", nil)
	assert.Error(t, err, "nil content should fail")

	_, err = client.UpdateNotebookContent(testOrganizationID, "", blankNotebookContent())
	assert.Error(t, err, "empty entity GUID should fail on update")

	_, err = client.RenameNotebook(testOrganizationID, "some-guid", "", blankNotebookContent())
	assert.Error(t, err, "empty new name should fail on rename")

	assert.Error(t, client.DeleteNotebook("", "some-guid"), "empty organization ID should fail on delete")
	assert.Error(t, client.DeleteNotebook(testOrganizationID, ""), "empty entity GUID should fail on delete")
}

// requireNotebookMetadata fetches the notebook via the narrow GetNotebook
// query. Every place that fetches metadata should go through this so the
// error surface is consistent. We deliberately don't use the generated
// GetEntity - that query expands fragments for every EntityManagement
// implementation and breaks on unrelated schema drift.
func requireNotebookMetadata(t *testing.T, client Notebooks, entityGUID string) *EntityManagementNotebookEntity {
	t.Helper()

	nb, err := client.GetNotebook(entityGUID)
	require.NoError(t, err, "GetNotebook(%s)", entityGUID)
	require.NotNil(t, nb, "GetNotebook returned nil entity")
	return nb
}

// requireNotebookAppearsInSearch polls entitySearch and asserts the notebook
// eventually shows up. NerdGraph's projection lags Blob API writes so we
// retry briefly instead of asserting on the first attempt.
func requireNotebookAppearsInSearch(t *testing.T, client Notebooks, entityGUID, expectedName string) {
	t.Helper()

	requireEventually(t, "notebook to appear in entitySearch", func() error {
		result, err := client.SearchNotebooks("", "type = 'NOTEBOOK'")
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("nil search result")
		}
		for _, nb := range result.Notebooks {
			if nb.ID == entityGUID {
				if nb.Name != expectedName {
					return fmt.Errorf("found notebook by id but name = %q, want %q", nb.Name, expectedName)
				}
				return nil
			}
		}
		return fmt.Errorf("notebook %s not present in entitySearch results", entityGUID)
	})
}

// assertContentMatches downloads the current notebook content and structurally
// compares it to the expected value. Passes both through json.Marshal +
// Unmarshal to normalise map/interface types before reflect.DeepEqual.
func assertContentMatches(t *testing.T, client Notebooks, orgID, entityGUID string, expected map[string]interface{}) {
	t.Helper()

	raw, err := client.GetNotebookContent(orgID, entityGUID)
	require.NoError(t, err, "GET notebook content")

	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &got), "unmarshal notebook body")

	// Round-trip expected through JSON so both sides use identical concrete
	// types (json.Unmarshal yields float64 for numbers, etc.).
	normalisedBytes, err := json.Marshal(expected)
	require.NoError(t, err)
	var normalisedExpected map[string]interface{}
	require.NoError(t, json.Unmarshal(normalisedBytes, &normalisedExpected))

	if !reflect.DeepEqual(got, normalisedExpected) {
		t.Fatalf("notebook content mismatch\n  got:      %s\n  expected: %s", string(raw), string(normalisedBytes))
	}
}

// requireEventually retries fn every 500ms for up to 15s, failing the test if
// the condition never holds. Keeps the test stable without turning transient
// server-side propagation into false flakiness.
func requireEventually(t *testing.T, describe string, fn func() error) {
	t.Helper()

	deadline := time.Now().Add(15 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		lastErr = fn()
		if lastErr == nil {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s: %v", describe, lastErr)
}

// cleanupNotebook removes a notebook created earlier in the test. It retries
// a few times with backoff so transient Blob API errors don't leak the
// resource, and after final failure it verifies via NerdGraph whether the
// notebook is still present so the log line clearly says "leaked" vs
// "already gone". Called from a defer, so it must never call t.Fatal (that
// would mask the test's real failure); instead it uses t.Logf and
// t.Errorf so CI surfaces leaks without hiding earlier assertions.
func cleanupNotebook(t *testing.T, client Notebooks, orgID, entityGUID string) {
	t.Helper()

	const maxAttempts = 4
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		lastErr = client.DeleteNotebook(orgID, entityGUID)
		if lastErr == nil {
			t.Logf("cleanup: DeleteNotebook(%s) succeeded on attempt %d", entityGUID, attempt)
			return
		}
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}

	// Delete kept failing. Confirm whether the notebook still exists so
	// the operator knows whether this is a true leak or a delete API bug
	// against an already-gone resource.
	if _, err := client.GetNotebookContent(orgID, entityGUID); err != nil && strings.Contains(err.Error(), "not found") {
		t.Logf("cleanup: DeleteNotebook(%s) returned %v but notebook is already gone", entityGUID, lastErr)
		return
	}
	t.Errorf("cleanup: LEAKED notebook %s in org %s (last delete error: %v). Please delete manually.", entityGUID, orgID, lastErr)
}
