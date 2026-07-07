package notebooks

import (
	"context"
	"fmt"
)

// GetNotebook fetches metadata for a single notebook by entity GUID and
// returns a typed *EntityManagementNotebookEntity, saving callers the type
// assertion that the tutone-generated GetEntity requires (its return type is
// the EntityManagementEntity interface).
//
// If the notebook doesn't exist the server returns an "Entity not found"
// error and this method surfaces it as a wrapped error.
func (a *Notebooks) GetNotebook(entityGUID string) (*EntityManagementNotebookEntity, error) {
	return a.GetNotebookWithContext(context.Background(), entityGUID)
}

// GetNotebookWithContext is the context-aware variant of GetNotebook.
func (a *Notebooks) GetNotebookWithContext(ctx context.Context, entityGUID string) (*EntityManagementNotebookEntity, error) {
	if entityGUID == "" {
		return nil, fmt.Errorf("notebooks: entity GUID is required")
	}

	var resp getNotebookResponse
	vars := map[string]interface{}{"id": entityGUID}
	if err := a.client.NerdGraphQueryWithContext(ctx, getNotebookQuery, vars, &resp); err != nil {
		return nil, err
	}

	nb := resp.Actor.EntityManagement.Entity
	if nb == nil || nb.ID == "" {
		return nil, fmt.Errorf("notebooks: entity %s not found", entityGUID)
	}
	return nb, nil
}

// SearchNotebooks lists notebooks that match the given NerdGraph
// entityManagement.entitySearch predicate. Callers usually pass
// "type = 'NOTEBOOK'" (optionally combined with scope.id filters).
// Pass an empty cursor for the first page; subsequent pages come from
// NotebookSearchResult.NextCursor.
//
// Returns a pre-filtered slice of *EntityManagementNotebookEntity instead of
// the interface-typed entities the tutone-generated GetEntitySearch returns.
func (a *Notebooks) SearchNotebooks(cursor string, query string) (*NotebookSearchResult, error) {
	return a.SearchNotebooksWithContext(context.Background(), cursor, query)
}

// SearchNotebooksWithContext is the context-aware variant of SearchNotebooks.
func (a *Notebooks) SearchNotebooksWithContext(ctx context.Context, cursor string, query string) (*NotebookSearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("notebooks: search query is required (typically \"type = 'NOTEBOOK'\")")
	}

	var resp searchNotebooksResponse
	vars := map[string]interface{}{"query": query}
	if cursor != "" {
		vars["cursor"] = cursor
	}
	if err := a.client.NerdGraphQueryWithContext(ctx, searchNotebooksQuery, vars, &resp); err != nil {
		return nil, err
	}

	result := &NotebookSearchResult{
		NextCursor: resp.Actor.EntityManagement.EntitySearch.NextCursor,
	}
	for _, e := range resp.Actor.EntityManagement.EntitySearch.Entities {
		if e != nil {
			result.Notebooks = append(result.Notebooks, e)
		}
	}
	return result, nil
}

// NotebookSearchResult is the paginated view returned by SearchNotebooks. Only
// NotebookEntity results are included; the underlying entitySearch call can
// return other entity types but this method filters them out client-side.
type NotebookSearchResult struct {
	Notebooks  []*EntityManagementNotebookEntity
	NextCursor string
}

// getNotebookResponse is the shape of the raw GraphQL response for
// getNotebookQuery. Kept package-private to prevent accidental use from
// outside the client.
type getNotebookResponse struct {
	Actor struct {
		EntityManagement struct {
			Entity *EntityManagementNotebookEntity `json:"entity"`
		} `json:"entityManagement"`
	} `json:"actor"`
}

type searchNotebooksResponse struct {
	Actor struct {
		EntityManagement struct {
			EntitySearch struct {
				Entities   []*EntityManagementNotebookEntity `json:"entities"`
				NextCursor string                            `json:"nextCursor"`
			} `json:"entitySearch"`
		} `json:"entityManagement"`
	} `json:"actor"`
}

// getNotebookQuery selects only fields defined on EntityManagementNotebookEntity
// plus the entity interface's common fields. Fields on Blob (`content`) that
// were observed at runtime not to be queryable - storageIdentifier, checksum,
// checksumAlgorithm - are deliberately omitted so this query cannot regress
// with future Blob schema changes.
const getNotebookQuery = `query($id: ID!) {
	actor {
		entityManagement {
			entity(id: $id) {
				__typename
				... on EntityManagementNotebookEntity {
					id
					name
					type
					scope { id type }
					tags { key values }
					content { id contentType }
					metadata {
						version
						createdAt
						updatedAt
						createdBy { __typename id type }
						updatedBy { __typename id type }
					}
				}
			}
		}
	}
}`

const searchNotebooksQuery = `query($query: String!, $cursor: String) {
	actor {
		entityManagement {
			entitySearch(query: $query, cursor: $cursor) {
				entities {
					__typename
					... on EntityManagementNotebookEntity {
						id
						name
						type
						scope { id type }
						tags { key values }
						metadata {
							version
							createdAt
							updatedAt
						}
					}
				}
				nextCursor
			}
		}
	}
}`
