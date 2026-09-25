package scorecards

// This file is hand-written because tutone cannot currently produce valid Go
// for the `[ID]` scalar-list return type of the AddCollectionMembers /
// RemoveCollectionMembers mutations — it emits an invalid field selector
// `resp.[]string`. See NGEP_ANALYSIS.md § "Member operations semantics".
//
// Response shape (as observed live):
//   {"data":{"entityManagementAddCollectionMembers":["id1","id2"]}}
//   Slots align 1:1 with the input `ids`. On per-item failure the slot is
//   `null` and a companion `errors[]` entry with `path: [<mutation>, <index>]`
//   describes it (e.g. COLLECTION_DUPLICATE_MEMBER, Target entity not found).

import "context"

// EntityManagementAddCollectionMembers adds one or more entities to a
// collection. Returns the list of IDs actually added — slots for failures are
// nil pointers, with matching entries in the transport-level error slice.
func (a *Scorecards) EntityManagementAddCollectionMembers(
	collectionID string,
	ids []string,
) ([]*string, error) {
	return a.EntityManagementAddCollectionMembersWithContext(context.Background(), collectionID, ids)
}

// EntityManagementAddCollectionMembersWithContext is the context-aware form of
// EntityManagementAddCollectionMembers.
func (a *Scorecards) EntityManagementAddCollectionMembersWithContext(
	ctx context.Context,
	collectionID string,
	ids []string,
) ([]*string, error) {

	resp := entityManagementCollectionMembersResponse{
		key: "entityManagementAddCollectionMembers",
	}
	vars := map[string]interface{}{
		"collectionId": collectionID,
		"ids":          ids,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, entityManagementAddCollectionMembersMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.decoded(), nil
}

// EntityManagementRemoveCollectionMembers detaches one or more entities from a
// collection. Same slot-aligned semantics as
// EntityManagementAddCollectionMembers.
func (a *Scorecards) EntityManagementRemoveCollectionMembers(
	collectionID string,
	ids []string,
) ([]*string, error) {
	return a.EntityManagementRemoveCollectionMembersWithContext(context.Background(), collectionID, ids)
}

// EntityManagementRemoveCollectionMembersWithContext is the context-aware form
// of EntityManagementRemoveCollectionMembers.
func (a *Scorecards) EntityManagementRemoveCollectionMembersWithContext(
	ctx context.Context,
	collectionID string,
	ids []string,
) ([]*string, error) {

	resp := entityManagementCollectionMembersResponse{
		key: "entityManagementRemoveCollectionMembers",
	}
	vars := map[string]interface{}{
		"collectionId": collectionID,
		"ids":          ids,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, entityManagementRemoveCollectionMembersMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.decoded(), nil
}

// entityManagementCollectionMembersResponse is a small helper that decodes the
// scalar-list mutation response uniformly for Add and Remove.
//
// NerdGraph camel-cases the top-level response field to
// "EntityManagementAddCollectionMembers" / "EntityManagementRemoveCollectionMembers"
// in the client-side unmarshal. We unmarshal to whichever key is populated.
type entityManagementCollectionMembersResponse struct {
	// key is populated by the caller before decoding — Add or Remove.
	key string

	EntityManagementAddCollectionMembers    []*string `json:"EntityManagementAddCollectionMembers"`
	EntityManagementRemoveCollectionMembers []*string `json:"EntityManagementRemoveCollectionMembers"`
}

func (r *entityManagementCollectionMembersResponse) decoded() []*string {
	switch r.key {
	case "entityManagementAddCollectionMembers":
		return r.EntityManagementAddCollectionMembers
	case "entityManagementRemoveCollectionMembers":
		return r.EntityManagementRemoveCollectionMembers
	}
	return nil
}

const entityManagementAddCollectionMembersMutation = `mutation(
	$collectionId: ID!,
	$ids: [ID!]!,
) { entityManagementAddCollectionMembers(
	collectionId: $collectionId,
	ids: $ids,
) }`

const entityManagementRemoveCollectionMembersMutation = `mutation(
	$collectionId: ID!,
	$ids: [ID!]!,
) { entityManagementRemoveCollectionMembers(
	collectionId: $collectionId,
	ids: $ids,
) }`
