package notebooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// NotebookMutationResponse is the envelope returned by the Blob Storage API
// on create and update. The entityGuid is the stable identifier for the
// notebook (use it in NerdGraph metadata queries); blobId identifies the
// specific version of the content just written. blobVersionEntity is
// documented but currently observed to be null - leaving the field so future
// server behavior is preserved automatically.
type NotebookMutationResponse struct {
	EntityGUID        string               `json:"entityGuid,omitempty"`
	BlobID            string               `json:"blobId,omitempty"`
	BlobVersionEntity *NotebookBlobVersion `json:"blobVersionEntity,omitempty"`
}

// NotebookBlobVersion mirrors the doc's schema for the (currently null)
// blobVersionEntity subfield.
type NotebookBlobVersion struct {
	EntityGUID string `json:"entityGuid,omitempty"`
	Version    int    `json:"version,omitempty"`
}

// NotebookContent is a helper wrapper around any JSON-serialisable value
// representing the notebook body. Callers can pass a struct, a map, or
// json.RawMessage - whatever the API accepts. The Blob API does not enforce
// a fixed schema; it stores whatever versioned JSON you send. See the docs
// draft for currently supported block/widget shapes.
type NotebookContent = interface{}

// CreateNotebook creates a new notebook inside an organization. The notebook
// is created with the given name and initial content and returns the
// entityGuid that identifies it in every subsequent call.
//
// The Blob Storage API requires content on create - pass an empty document
// like map[string]interface{}{"version": "1", "blocks": []interface{}{}}
// if you want to seed a blank notebook.
func (a *Notebooks) CreateNotebook(
	organizationID string,
	name string,
	content NotebookContent,
) (*NotebookMutationResponse, error) {
	return a.CreateNotebookWithContext(context.Background(), organizationID, name, content)
}

// CreateNotebookWithContext is the context-aware variant of CreateNotebook.
func (a *Notebooks) CreateNotebookWithContext(
	ctx context.Context,
	organizationID string,
	name string,
	content NotebookContent,
) (*NotebookMutationResponse, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("notebooks: organization ID is required")
	}
	if name == "" {
		return nil, fmt.Errorf("notebooks: notebook name is required")
	}
	if content == nil {
		return nil, fmt.Errorf("notebooks: notebook content is required (pass an empty document to seed a blank notebook)")
	}

	resp := NotebookMutationResponse{}
	_, err := a.client.PostWithContext(
		ctx,
		a.notebooksURL(organizationID, ""),
		blobRequestHeaders(map[string]string{newRelicEntityHeader: encodeEntityHeader(name)}),
		content,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateNotebookContent overwrites the notebook's content with a new version.
// Every call produces a new revision (metadata.version increments on the
// entity). The previous version is retained by the server for a short
// retention window; consult the docs for the current TTL.
func (a *Notebooks) UpdateNotebookContent(
	organizationID string,
	entityGUID string,
	content NotebookContent,
) (*NotebookMutationResponse, error) {
	return a.UpdateNotebookContentWithContext(context.Background(), organizationID, entityGUID, content)
}

// UpdateNotebookContentWithContext is the context-aware variant.
func (a *Notebooks) UpdateNotebookContentWithContext(
	ctx context.Context,
	organizationID string,
	entityGUID string,
	content NotebookContent,
) (*NotebookMutationResponse, error) {
	return a.postNotebookContent(ctx, organizationID, entityGUID, content, "")
}

// RenameNotebook renames the notebook and rewrites its content in a single
// atomic call. The Blob Storage API's rename channel is the NewRelic-Entity
// header on the update endpoint - there is no rename-only path, so this
// method requires the current content to be POSTed back. Fetch it first with
// GetNotebookContent if the caller only wants to change the name.
func (a *Notebooks) RenameNotebook(
	organizationID string,
	entityGUID string,
	newName string,
	content NotebookContent,
) (*NotebookMutationResponse, error) {
	return a.RenameNotebookWithContext(context.Background(), organizationID, entityGUID, newName, content)
}

// RenameNotebookWithContext is the context-aware variant.
func (a *Notebooks) RenameNotebookWithContext(
	ctx context.Context,
	organizationID string,
	entityGUID string,
	newName string,
	content NotebookContent,
) (*NotebookMutationResponse, error) {
	if newName == "" {
		return nil, fmt.Errorf("notebooks: new name is required for rename")
	}
	return a.postNotebookContent(ctx, organizationID, entityGUID, content, newName)
}

func (a *Notebooks) postNotebookContent(
	ctx context.Context,
	organizationID string,
	entityGUID string,
	content NotebookContent,
	renameTo string,
) (*NotebookMutationResponse, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("notebooks: organization ID is required")
	}
	if entityGUID == "" {
		return nil, fmt.Errorf("notebooks: entity GUID is required")
	}
	if content == nil {
		return nil, fmt.Errorf("notebooks: notebook content is required")
	}

	// The Blob Storage API returns text/plain (not JSON) for error responses,
	// so we bypass the internal JSON client and make a raw HTTP request. This
	// lets us read the actual error body and surface it to the caller rather
	// than returning an opaque "400 response returned" with no context.
	body, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("notebooks: marshal content: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.notebooksURL(organizationID, entityGUID), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("notebooks: build POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", a.config.PersonalAPIKey)
	if a.config.UserAgent != "" {
		req.Header.Set("User-Agent", a.config.UserAgent)
	}
	if renameTo != "" {
		req.Header.Set(newRelicEntityHeader, encodeEntityHeader(renameTo))
	}

	httpClient := &http.Client{}
	if a.config.Timeout != nil {
		httpClient.Timeout = *a.config.Timeout
	}
	if a.config.HTTPTransport != nil {
		httpClient.Transport = a.config.HTTPTransport
	}

	httpResp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notebooks: POST notebook content: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("notebooks: read POST response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(respBody))
		// Translate the most actionable Blob API error messages into explicit
		// Terraform-friendly descriptions so users know exactly what to fix.
		switch {
		case httpResp.StatusCode == 400 && strings.Contains(msg, "unique constraint"):
			return nil, fmt.Errorf(
				"notebooks: a notebook with this title already exists in the organization — "+
					"choose a unique title (API: %s)", msg)
		case httpResp.StatusCode == 403:
			return nil, fmt.Errorf(
				"notebooks: permission denied — ensure the API key has write access to notebooks (API: %s)", msg)
		case httpResp.StatusCode == 404:
			return nil, fmt.Errorf("notebooks: notebook not found: %s", entityGUID)
		default:
			return nil, fmt.Errorf("notebooks: POST notebook content: unexpected status %d: %s",
				httpResp.StatusCode, msg)
		}
	}

	var result NotebookMutationResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("notebooks: decode POST response: %w", err)
	}
	return &result, nil
}

// GetNotebookContent fetches the current content of the notebook and returns
// it as a raw JSON message. The Blob Storage API returns whatever versioned
// JSON was last written - this package doesn't presume a schema. Callers can
// json.Unmarshal into their own type or leave it as a RawMessage for
// pass-through use cases.
func (a *Notebooks) GetNotebookContent(
	organizationID string,
	entityGUID string,
) (json.RawMessage, error) {
	return a.GetNotebookContentWithContext(context.Background(), organizationID, entityGUID)
}

// GetNotebookContentWithContext is the context-aware variant.
func (a *Notebooks) GetNotebookContentWithContext(
	ctx context.Context,
	organizationID string,
	entityGUID string,
) (json.RawMessage, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("notebooks: organization ID is required")
	}
	if entityGUID == "" {
		return nil, fmt.Errorf("notebooks: entity GUID is required")
	}

	// The Blob API returns the notebook body verbatim as the response body -
	// no JSON envelope, so we bypass the standard client's decode-into-struct
	// path and read the raw bytes ourselves (same pattern as
	// FleetControlGetConfiguration).
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.notebooksURL(organizationID, entityGUID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Api-Key", a.config.PersonalAPIKey)
	if a.config.UserAgent != "" {
		req.Header.Set("User-Agent", a.config.UserAgent)
	}

	httpClient := &http.Client{}
	if a.config.Timeout != nil {
		httpClient.Timeout = *a.config.Timeout
	}
	if a.config.HTTPTransport != nil {
		httpClient.Transport = a.config.HTTPTransport
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notebooks: GET notebook content: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("notebooks: read notebook content: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("notebooks: notebook not found: %s", entityGUID)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("notebooks: GET notebook content: unexpected status %d: %s", resp.StatusCode, string(body))
	}
	return json.RawMessage(body), nil
}

// DeleteNotebook removes the notebook and all its content versions. The
// Blob Storage API's DELETE endpoint returns HTTP 200 with an empty body on
// success (contrary to the doc claim of 204).
func (a *Notebooks) DeleteNotebook(
	organizationID string,
	entityGUID string,
) error {
	return a.DeleteNotebookWithContext(context.Background(), organizationID, entityGUID)
}

// DeleteNotebookWithContext is the context-aware variant.
func (a *Notebooks) DeleteNotebookWithContext(
	ctx context.Context,
	organizationID string,
	entityGUID string,
) error {
	if organizationID == "" {
		return fmt.Errorf("notebooks: organization ID is required")
	}
	if entityGUID == "" {
		return fmt.Errorf("notebooks: entity GUID is required")
	}

	_, err := a.client.DeleteWithContext(ctx, a.notebooksURL(organizationID, entityGUID), nil, nil)
	return err
}
