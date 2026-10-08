package servicearchintelligence

// This file provides a type-agnostic way to read collection membership.
//
// The standard GetCollectionElements/pageCollectionItems path requires every
// entity type that may appear in a collection to be explicitly registered in
// UnmarshalEntityManagementEntityInterface. That creates a maintenance burden:
// whenever NGEP adds a new entity type, both the go-client unmarshal switch and
// any provider-side switch must be updated, or members of that type are silently
// dropped.
//
// CollectionMemberIDsWithContext bypasses typed unmarshal entirely. It issues a
// targeted NerdGraph query that selects only `id` and `tags` on each element,
// decodes directly into a plain struct, and pages automatically. The result is
// stable across all current and future entity types.

import (
	"context"
	"fmt"
)

// CollectionMember holds the minimal fields needed by the provider for
// non-authoritative ownership split logic: the entity GUID and its tags.
type CollectionMember struct {
	ID   string                `json:"id"`
	Tags []EntityManagementTag `json:"tags"`
}

// GetCollectionMemberIDs is the convenience wrapper for
// GetCollectionMemberIDsWithContext.
func (a *Scorecards) GetCollectionMemberIDs(collectionID string) ([]CollectionMember, error) {
	return a.GetCollectionMemberIDsWithContext(context.Background(), collectionID)
}

// GetCollectionMemberIDsWithContext pages through a collection and returns a
// slice of CollectionMember for every entity in it, regardless of entity type.
// It is the preferred function for enumerating team ownership collections because
// it does not require knowledge of each member's concrete entity type.
func (a *Scorecards) GetCollectionMemberIDsWithContext(ctx context.Context, collectionID string) ([]CollectionMember, error) {
	var all []CollectionMember
	cursor := ""

	for {
		var resp collectionMemberIDsResponse
		vars := map[string]interface{}{
			"id":     collectionID,
			"cursor": cursor,
		}

		if err := a.client.NerdGraphQueryWithContext(ctx, collectionMemberIDsQuery, vars, &resp); err != nil {
			return nil, fmt.Errorf("fetching collection member IDs for %s: %w", collectionID, err)
		}

		elems := resp.Actor.EntityManagement.CollectionElements
		all = append(all, elems.Entities...)

		if elems.NextCursor == nil || *elems.NextCursor == "" {
			break
		}
		cursor = *elems.NextCursor
	}

	return all, nil
}

// collectionMemberIDsResponse is the minimal response shape for the
// collectionMemberIDsQuery below — only id and tags are selected.
type collectionMemberIDsResponse struct {
	Actor struct {
		EntityManagement struct {
			CollectionElements struct {
				Entities   []CollectionMember `json:"entities"`
				NextCursor *string            `json:"nextCursor"`
			} `json:"collectionElements"`
		} `json:"entityManagement"`
	} `json:"actor"`
}

// collectionMemberIDsQuery selects only id and tags — no __typename, no entity
// type dispatch. This query works for any entity type that may appear in a
// collection, including types not yet registered in the go-client.
const collectionMemberIDsQuery = `query(
	$id: ID!,
	$cursor: String,
) {
  actor {
    entityManagement {
      collectionElements(id: $id, cursor: $cursor) {
        entities {
          id
          tags { key values }
        }
        nextCursor
      }
    }
  }
}`
