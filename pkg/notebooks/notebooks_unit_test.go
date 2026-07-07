//go:build unit
// +build unit

// Unit tests for the notebooks package. These do not touch the platform:
// every network call is served by an httptest.Server so the tests are
// fully deterministic and runnable in any environment. Together with the
// integration suite they push the reachable branches of the hand-written
// code (notebooks.go, notebooks_blob.go, notebooks_metadata.go) to full
// coverage.

package notebooks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
)

// -----------------------------------------------------------------------------
// Test scaffolding: an httptest server that answers both the Blob Storage API
// path (/organizations/.../Notebooks[/id]) and the NerdGraph endpoint (/graphql)
// and a config wired to point at that server for both URL surfaces.
// -----------------------------------------------------------------------------

type mockPlatform struct {
	server *httptest.Server
	// handler is set by each test with the desired behaviour. Wrapped in a
	// closure so a test can rebind the handler without swapping the server.
	handler http.Handler
}

func newMockPlatform(t *testing.T) *mockPlatform {
	t.Helper()

	mp := &mockPlatform{}
	mp.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mp.handler == nil {
			http.Error(w, "no handler set", http.StatusInternalServerError)
			return
		}
		mp.handler.ServeHTTP(w, r)
	}))
	t.Cleanup(mp.server.Close)
	return mp
}

// newNotebooksClient returns a Notebooks client whose blob-service and
// NerdGraph endpoints are both routed to mp.server. Uses a fixed API key so
// tests don't require environment variables.
func (mp *mockPlatform) newNotebooksClient(t *testing.T) Notebooks {
	t.Helper()

	cfg := config.New()
	cfg.PersonalAPIKey = "NRAK-TEST-KEY"
	cfg.LogLevel = "error"
	cfg.Region().SetBlobServiceBaseURL(mp.server.URL)
	cfg.Region().SetNerdGraphBaseURL(mp.server.URL + "/graphql")
	return New(cfg)
}

// -----------------------------------------------------------------------------
// Validation-branch coverage: every public method's argument checks.
// -----------------------------------------------------------------------------

// TestValidation_CreateNotebook exercises every early-return in
// CreateNotebookWithContext. No network happens.
func TestValidation_CreateNotebook(t *testing.T) {
	client := New(config.New())

	_, err := client.CreateNotebook("", "name", blankContent())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "organization ID")

	_, err = client.CreateNotebook("org", "", blankContent())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")

	_, err = client.CreateNotebook("org", "name", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "content")
}

// TestValidation_UpdateNotebookContent covers postNotebookContent's early
// returns when called via the update path.
func TestValidation_UpdateNotebookContent(t *testing.T) {
	client := New(config.New())

	_, err := client.UpdateNotebookContent("", "guid", blankContent())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "organization ID")

	_, err = client.UpdateNotebookContent("org", "", blankContent())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "entity GUID")

	_, err = client.UpdateNotebookContent("org", "guid", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "content")
}

// TestValidation_RenameNotebook covers RenameNotebookWithContext's early
// returns, including the delegation to postNotebookContent's checks.
func TestValidation_RenameNotebook(t *testing.T) {
	client := New(config.New())

	_, err := client.RenameNotebook("org", "guid", "", blankContent())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "new name")

	_, err = client.RenameNotebook("", "guid", "new", blankContent())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "organization ID")

	_, err = client.RenameNotebook("org", "", "new", blankContent())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "entity GUID")

	_, err = client.RenameNotebook("org", "guid", "new", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "content")
}

// TestValidation_GetNotebookContent covers GetNotebookContentWithContext's
// argument checks.
func TestValidation_GetNotebookContent(t *testing.T) {
	client := New(config.New())

	_, err := client.GetNotebookContent("", "guid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "organization ID")

	_, err = client.GetNotebookContent("org", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "entity GUID")
}

// TestValidation_DeleteNotebook covers DeleteNotebookWithContext's argument
// checks.
func TestValidation_DeleteNotebook(t *testing.T) {
	client := New(config.New())

	err := client.DeleteNotebook("", "guid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "organization ID")

	err = client.DeleteNotebook("org", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "entity GUID")
}

// TestValidation_GetNotebook and _SearchNotebooks cover the NerdGraph
// metadata methods' argument checks.
func TestValidation_GetNotebook(t *testing.T) {
	client := New(config.New())

	_, err := client.GetNotebook("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "entity GUID")
}

func TestValidation_SearchNotebooks(t *testing.T) {
	client := New(config.New())

	_, err := client.SearchNotebooks("", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "query is required")
}

// -----------------------------------------------------------------------------
// Happy-path coverage via mock httptest server: exercises the network paths
// end to end without needing a real platform account.
// -----------------------------------------------------------------------------

// TestBlobAPI_HappyPath drives create -> get -> update -> rename -> delete
// against a mock server. Verifies that the request path, method, and
// NewRelic-Entity header are set exactly as the Blob Storage API expects.
func TestBlobAPI_HappyPath(t *testing.T) {
	mp := newMockPlatform(t)

	const (
		orgID  = "test-org"
		guid   = "test-guid-abc"
		blobA  = "blob-v1"
		blobB  = "blob-v2"
		nbName = "unit-test-notebook"
	)

	// Recorded server-side state; each request updates it so the test can
	// assert on the observed sequence of calls.
	var (
		posts     int
		lastEntityHeader string
		lastPath   string
		lastMethod string
	)

	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastMethod = r.Method
		lastPath = r.URL.Path

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/organizations/"+orgID+"/Notebooks":
			posts++
			lastEntityHeader = r.Header.Get("NewRelic-Entity")
			require.NotEmpty(t, lastEntityHeader, "create should set NewRelic-Entity header")
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			require.Equal(t, "NRAK-TEST-KEY", r.Header.Get("Api-Key"))
			respond(t, w, http.StatusOK, map[string]interface{}{
				"entityGuid":        guid,
				"blobId":            blobA,
				"blobVersionEntity": nil,
			})

		case r.Method == http.MethodGet && r.URL.Path == "/organizations/"+orgID+"/Notebooks/"+guid:
			// Return the last content the client posted. For simplicity in
			// this test we just echo a known payload.
			respond(t, w, http.StatusOK, blankContent())

		case r.Method == http.MethodPost && r.URL.Path == "/organizations/"+orgID+"/Notebooks/"+guid:
			posts++
			lastEntityHeader = r.Header.Get("NewRelic-Entity")
			respond(t, w, http.StatusOK, map[string]interface{}{
				"entityGuid": guid,
				"blobId":     blobB,
			})

		case r.Method == http.MethodDelete && r.URL.Path == "/organizations/"+orgID+"/Notebooks/"+guid:
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "unexpected: "+r.Method+" "+r.URL.Path, http.StatusNotFound)
		}
	})

	client := mp.newNotebooksClient(t)

	// Cleanup discipline mirrors the integration test: if a test failure
	// leaves a live notebook (impossible against a mock, but keeps the
	// pattern consistent so operators of this test suite recognise the
	// shape), the defer tries to delete it and logs on failure.
	var createdGUID string
	defer func() {
		if createdGUID == "" {
			return
		}
		if err := client.DeleteNotebook(orgID, createdGUID); err != nil {
			t.Logf("cleanup: DeleteNotebook(%s) failed: %v", createdGUID, err)
		}
	}()

	// Create.
	created, err := client.CreateNotebook(orgID, nbName, blankContent())
	require.NoError(t, err)
	require.Equal(t, guid, created.EntityGUID)
	require.Equal(t, blobA, created.BlobID)
	require.Equal(t, `{"name":"`+nbName+`"}`, lastEntityHeader, "NewRelic-Entity should contain the name JSON")
	createdGUID = created.EntityGUID

	// Read content back.
	body, err := client.GetNotebookContent(orgID, guid)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &got))
	require.Equal(t, "1", got["version"])

	// Update content: NewRelic-Entity header should be absent.
	lastEntityHeader = "sentinel"
	updated, err := client.UpdateNotebookContent(orgID, guid, blankContent())
	require.NoError(t, err)
	require.Equal(t, blobB, updated.BlobID)
	require.Equal(t, "", lastEntityHeader, "update content must not send NewRelic-Entity")

	// Rename: NewRelic-Entity should carry the new name.
	renameResp, err := client.RenameNotebook(orgID, guid, "renamed", blankContent())
	require.NoError(t, err)
	require.Equal(t, blobB, renameResp.BlobID)
	require.Equal(t, `{"name":"renamed"}`, lastEntityHeader)

	// Delete: verifies path + method observed by the server.
	require.NoError(t, client.DeleteNotebook(orgID, guid))
	require.Equal(t, http.MethodDelete, lastMethod)
	require.Equal(t, "/organizations/"+orgID+"/Notebooks/"+guid, lastPath)
	createdGUID = "" // cleanup no longer needed

	require.Equal(t, 3, posts, "expected 3 POSTs (create, update, rename)")
}

// TestGetNotebookContent_NotFound verifies the 404 error surface is a
// well-formed 'not found' message the caller can distinguish from other
// failure modes.
func TestGetNotebookContent_NotFound(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, "Blob not found.")
	})
	client := mp.newNotebooksClient(t)

	_, err := client.GetNotebookContent("org", "guid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Contains(t, err.Error(), "guid", "error should include the requested entity GUID")
}

// TestGetNotebookContent_ServerError covers the non-2xx-non-404 branch.
func TestGetNotebookContent_ServerError(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "boom")
	})
	client := mp.newNotebooksClient(t)

	_, err := client.GetNotebookContent("org", "guid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status 500")
	assert.Contains(t, err.Error(), "boom", "error should include the response body for diagnostics")
}

// TestGetNotebookContent_TransportError covers the http.Client.Do error path
// by pointing the client at a server we've already closed.
func TestGetNotebookContent_TransportError(t *testing.T) {
	cfg := config.New()
	cfg.PersonalAPIKey = "NRAK-TEST"
	// A URL that cannot resolve produces a Do-time error rather than an HTTP
	// status. Using an unbound port on localhost is deterministic.
	cfg.Region().SetBlobServiceBaseURL("http://127.0.0.1:1")
	client := New(cfg)

	_, err := client.GetNotebookContent("org", "guid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "GET notebook content")
}

// TestGetNotebookContent_InvalidURL exercises the http.NewRequestWithContext
// failure branch by pointing the client at a URL containing a control byte.
// Go's url.Parse rejects it before the request is dispatched.
func TestGetNotebookContent_InvalidURL(t *testing.T) {
	cfg := config.New()
	cfg.PersonalAPIKey = "NRAK-TEST"
	// A NUL byte in the host part triggers url.Parse failure inside
	// http.NewRequestWithContext.
	cfg.Region().SetBlobServiceBaseURL("http://example\x7f.com")
	client := New(cfg)

	_, err := client.GetNotebookContent("org", "guid")
	require.Error(t, err)
}

// TestGetNotebookContent_UsesConfiguredHTTPTuning verifies that a User-Agent,
// custom Timeout, and custom Transport on the config are honoured when the
// content GET builds its request. Necessary to cover the tuning branches in
// GetNotebookContentWithContext.
func TestGetNotebookContent_UsesConfiguredHTTPTuning(t *testing.T) {
	mp := newMockPlatform(t)
	var observedUA string
	transportInvoked := false
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		observedUA = r.Header.Get("User-Agent")
		respond(t, w, http.StatusOK, blankContent())
	})

	cfg := config.New()
	cfg.PersonalAPIKey = "NRAK-TEST"
	cfg.UserAgent = "notebook-client-unit-test"
	timeout := 30 * time.Second
	cfg.Timeout = &timeout
	// Wrap the default transport so we can detect that our transport pointer
	// was actually used.
	cfg.HTTPTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		transportInvoked = true
		return http.DefaultTransport.RoundTrip(req)
	})
	cfg.Region().SetBlobServiceBaseURL(mp.server.URL)
	client := New(cfg)

	body, err := client.GetNotebookContent("org", "guid")
	require.NoError(t, err)
	require.NotEmpty(t, body)
	assert.Equal(t, "notebook-client-unit-test", observedUA, "custom User-Agent should be sent")
	assert.True(t, transportInvoked, "custom HTTPTransport should be used")
}

// roundTripFunc adapts a function to http.RoundTripper for the tuning test.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

// TestCreateNotebook_ServerError verifies the PostWithContext error path in
// CreateNotebookWithContext propagates cleanly.
func TestCreateNotebook_ServerError(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	client := mp.newNotebooksClient(t)

	_, err := client.CreateNotebook("org", "name", blankContent())
	require.Error(t, err)
}

// TestUpdateNotebookContent_ServerError verifies the PostWithContext error
// path in postNotebookContent propagates cleanly.
func TestUpdateNotebookContent_ServerError(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	client := mp.newNotebooksClient(t)

	_, err := client.UpdateNotebookContent("org", "guid", blankContent())
	require.Error(t, err)
}

// TestGetNotebook_ServerError verifies GetNotebookWithContext surfaces a
// NerdGraph transport-level error rather than returning a partial result.
func TestGetNotebook_ServerError(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "graphql exploded", http.StatusInternalServerError)
	})
	client := mp.newNotebooksClient(t)

	_, err := client.GetNotebook("some-guid")
	require.Error(t, err)
}

// TestSearchNotebooks_ServerError verifies SearchNotebooksWithContext surfaces
// a NerdGraph transport-level error.
func TestSearchNotebooks_ServerError(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "graphql exploded", http.StatusInternalServerError)
	})
	client := mp.newNotebooksClient(t)

	_, err := client.SearchNotebooks("", "type = 'NOTEBOOK'")
	require.Error(t, err)
}

// TestSearchNotebooks_EmptyCursorAndNilEntity covers the empty-cursor branch
// (the `if cursor != ""` false path) plus the nil-entity filter in the result
// loop.
func TestSearchNotebooks_EmptyCursorAndNilEntity(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The client should NOT send a cursor variable when cursor is empty.
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.NotContains(t, string(body), `"cursor"`, "empty cursor should be omitted from GraphQL variables")

		respond(t, w, http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"actor": map[string]interface{}{
					"entityManagement": map[string]interface{}{
						"entitySearch": map[string]interface{}{
							"entities": []interface{}{
								nil, // filtered out client-side
								map[string]interface{}{
									"__typename": "EntityManagementNotebookEntity",
									"id":         "kept",
								},
							},
							"nextCursor": "",
						},
					},
				},
			},
		})
	})
	client := mp.newNotebooksClient(t)

	result, err := client.SearchNotebooks("", "type = 'NOTEBOOK'")
	require.NoError(t, err)
	require.Len(t, result.Notebooks, 1, "nil entity in the response list should be filtered out")
	assert.Equal(t, "kept", result.Notebooks[0].ID)
}

// TestEncodeEntityHeader confirms the helper produces valid JSON for the
// characters that would trip a naive fmt.Sprintf approach (quotes,
// backslashes, control chars) so the callers can safely rely on it.
func TestEncodeEntityHeader(t *testing.T) {
	cases := map[string]string{
		"simple":      `{"name":"simple"}`,
		`with"quote`:  `{"name":"with\"quote"}`,
		"with\ttab":   `{"name":"with\ttab"}`,
		"with\\slash": `{"name":"with\\slash"}`,
	}
	for input, want := range cases {
		require.Equal(t, want, encodeEntityHeader(input), "input %q", input)
	}
}

// TestEncodeEntityHeader_MarshalPanic exercises the defensive panic guarding
// against a json.Marshal contract violation. json.Marshal on a
// map[string]string cannot fail in practice, but the panic must remain in
// place to make the runtime state visible if that contract were ever broken.
// We reach it here by swapping the package-level marshaler for a stub that
// deliberately errors.
func TestEncodeEntityHeader_MarshalPanic(t *testing.T) {
	orig := entityHeaderMarshaler
	entityHeaderMarshaler = func(interface{}) ([]byte, error) {
		return nil, assertErr("forced marshal failure")
	}
	defer func() { entityHeaderMarshaler = orig }()

	require.PanicsWithValue(t,
		"notebooks: unreachable json.Marshal failure encoding entity header: forced marshal failure",
		func() { _ = encodeEntityHeader("anything") },
	)
}

// assertErr is a minimal error type used by the panic test to give the panic
// message a stable string form.
type assertErr string

func (e assertErr) Error() string { return string(e) }

// TestGetNotebookContent_BodyReadError installs a transport whose response
// body errors mid-read. That exercises the io.ReadAll failure branch in
// GetNotebookContentWithContext which no ordinary httptest server can reach.
func TestGetNotebookContent_BodyReadError(t *testing.T) {
	cfg := config.New()
	cfg.PersonalAPIKey = "NRAK-TEST"
	cfg.Region().SetBlobServiceBaseURL("http://blob.invalid")
	cfg.HTTPTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(errorReader{}),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})
	client := New(cfg)

	_, err := client.GetNotebookContent("org", "guid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read notebook content")
}

// errorReader is an io.Reader that always returns an error, used to trigger
// the io.ReadAll error branch above.
type errorReader struct{}

func (errorReader) Read(_ []byte) (int, error) { return 0, io.ErrUnexpectedEOF }

// TestGetNotebookContent_ContextCancelled verifies the context path in
// GetNotebookContentWithContext short-circuits when the caller cancels.
func TestGetNotebookContent_ContextCancelled(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should never be reached because the context is cancelled before Do.
		t.Errorf("unexpected server invocation")
	})
	client := mp.newNotebooksClient(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.GetNotebookContentWithContext(ctx, "org", "guid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

// -----------------------------------------------------------------------------
// NerdGraph metadata query coverage: exercise decoding + not-found surface.
// -----------------------------------------------------------------------------

// TestGetNotebook_HappyPath asserts we decode the NerdGraph shape correctly
// and return a *EntityManagementNotebookEntity ready for use.
func TestGetNotebook_HappyPath(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/graphql", r.URL.Path)
		respond(t, w, http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"actor": map[string]interface{}{
					"entityManagement": map[string]interface{}{
						"entity": map[string]interface{}{
							"__typename": "EntityManagementNotebookEntity",
							"id":         "guid-1",
							"name":       "test",
							"type":       "NOTEBOOK",
							"scope":      map[string]interface{}{"id": "org", "type": "ORGANIZATION"},
							"tags":       []interface{}{},
							"metadata":   map[string]interface{}{"version": 1},
						},
					},
				},
			},
		})
	})
	client := mp.newNotebooksClient(t)

	nb, err := client.GetNotebook("guid-1")
	require.NoError(t, err)
	require.NotNil(t, nb)
	assert.Equal(t, "guid-1", nb.ID)
	assert.Equal(t, "test", nb.Name)
	assert.Equal(t, "NOTEBOOK", nb.Type)
	assert.Equal(t, 1, nb.Metadata.Version)
}

// TestGetNotebook_MissingEntity covers the "entity is null" branch which the
// server returns for unknown IDs when the query goes through cleanly (as
// opposed to erroring out).
func TestGetNotebook_MissingEntity(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond(t, w, http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"actor": map[string]interface{}{
					"entityManagement": map[string]interface{}{"entity": nil},
				},
			},
		})
	})
	client := mp.newNotebooksClient(t)

	_, err := client.GetNotebook("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestSearchNotebooks_HappyPath verifies decoding + filtering plus cursor
// propagation (the with-cursor variable branch).
func TestSearchNotebooks_HappyPath(t *testing.T) {
	mp := newMockPlatform(t)
	mp.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The client sends the cursor as a GraphQL variable when non-empty;
		// echo the request payload back so we can assert on it.
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), `"cursor":"page2"`, "cursor should be forwarded when set")

		respond(t, w, http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"actor": map[string]interface{}{
					"entityManagement": map[string]interface{}{
						"entitySearch": map[string]interface{}{
							"entities": []interface{}{
								map[string]interface{}{
									"__typename": "EntityManagementNotebookEntity",
									"id":         "guid-1",
									"name":       "n1",
								},
								map[string]interface{}{
									"__typename": "EntityManagementNotebookEntity",
									"id":         "guid-2",
									"name":       "n2",
								},
							},
							"nextCursor": "page3",
						},
					},
				},
			},
		})
	})
	client := mp.newNotebooksClient(t)

	result, err := client.SearchNotebooks("page2", "type = 'NOTEBOOK'")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Notebooks, 2)
	assert.Equal(t, "guid-1", result.Notebooks[0].ID)
	assert.Equal(t, "page3", result.NextCursor)
}

// -----------------------------------------------------------------------------
// Small helpers.
// -----------------------------------------------------------------------------

func blankContent() map[string]interface{} {
	return map[string]interface{}{"version": "1", "blocks": []interface{}{}}
}

// respond serialises v as JSON and writes it with the given status code.
func respond(t *testing.T, w http.ResponseWriter, status int, v interface{}) {
	t.Helper()

	buf, err := json.Marshal(v)
	require.NoError(t, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(buf)
	require.NoError(t, err)
}

