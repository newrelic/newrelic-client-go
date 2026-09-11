# notebooks

The `notebooks` package manages New Relic Notebook entities from Go.

The Notebooks backend is split across two services and this package hides that split behind a single `Notebooks` client:

| Operation | Backend | Method |
|---|---|---|
| Create a notebook (with initial content) | Blob Storage REST | `CreateNotebook` |
| Read notebook content | Blob Storage REST | `GetNotebookContent` |
| Overwrite content (new revision) | Blob Storage REST | `UpdateNotebookContent` |
| Rename a notebook | Blob Storage REST | `RenameNotebook` |
| Delete a notebook | Blob Storage REST | `DeleteNotebook` |
| Read notebook metadata (name, tags, scope, version) | NerdGraph | `GetNotebook` |
| List notebooks | NerdGraph | `SearchNotebooks` |

Rename shares the update endpoint via a `NewRelic-Entity` header; `RenameNotebook` takes the current content because the Blob API has no rename-only path. The Blob Storage API lives at `blob-api.service.newrelic.com/v1/e` (regional variants in `pkg/region/region_constants.go`); NerdGraph is `api.newrelic.com/graphql`. Both use the same User API key.

## Usage

```go
import (
    "github.com/newrelic/newrelic-client-go/v2/pkg/config"
    "github.com/newrelic/newrelic-client-go/v2/pkg/notebooks"
)

cfg := config.New()
cfg.PersonalAPIKey = "YOUR_API_KEY"
client := notebooks.New(cfg)

orgID := "your-organization-uuid"

// Create.
created, _ := client.CreateNotebook(orgID, "checkout investigation", map[string]interface{}{
    "version": "1",
    "blocks":  []interface{}{},
})

// Update content.
_, _ = client.UpdateNotebookContent(orgID, created.EntityGUID, newBody)

// Rename (must re-POST the current content).
_, _ = client.RenameNotebook(orgID, created.EntityGUID, "resolved", newBody)

// Read metadata via NerdGraph.
nb, _ := client.GetNotebook(created.EntityGUID)
fmt.Println(nb.Name, "version", nb.Metadata.Version)

// List (predicate uses UPPER_SNAKE runtime type).
result, _ := client.SearchNotebooks("", "type = 'NOTEBOOK'")
for _, nb := range result.Notebooks { fmt.Println(nb.ID, nb.Name) }

// Delete.
_ = client.DeleteNotebook(orgID, created.EntityGUID)
```

## Content format

The Blob API stores whatever versioned JSON you POST. `CreateNotebook`, `UpdateNotebookContent`, and `RenameNotebook` accept any JSON-serialisable value (struct, map, `json.RawMessage`). For the block/widget shapes the New Relic UI understands, see the [Blob Storage API for notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/blob-storage-api-for-notebooks/) and [Visualizations in notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/visualizations-in-notebooks/) in the public documentation.

## Regenerating with Tutone

This package requires a version of Tutone that supports the `include_implementations` key on `EntityManagementEntity`. Vanilla Tutone does not filter interface implementations at the query level, which pulls unrelated subtypes into the generated files. To regenerate, check out a Tutone build that includes `include_implementations` support and run:

```
/path/to/tutone -c .tutone.yml generate --package notebooks --refetch
```

## Integration tests

```
NEW_RELIC_FLEET_TEST_API_KEY=NRAK-...
NEW_RELIC_FLEET_TEST_ORGANIZATION_ID=<uuid>   # optional; defaults to the known org UUID
go test -tags integration ./pkg/notebooks/...
```

Every mutation is followed by an independent verification call against the platform - Blob API GET for content, NerdGraph `entity(id:)` / `entitySearch(query:)` for metadata - so an incorrect implementation cannot pass by returning fake success envelopes. Cleanup retries a few times with backoff and logs if a notebook is leaked.

## Unit tests

```
go test -tags unit ./pkg/notebooks/...
```

Fully covers the hand-written code with a mock httptest server; no credentials or network required.
