package scorecards

// This file supplies bits that tutone's generator misses today:
//   (a) the EntityManagementCollectionElementsFilter and CollectionIdFilterArgument
//       input types (not emitted for reasons unknown, though referenced by the
//       generated GetCollectionElements method), and
//   (b) an UnmarshalJSON on EntityManagementActorStitchedFields, which is
//       needed because .Entity is an interface and Go's default json decoder
//       can't dispatch it. Compare pipelinecontrol/types.go where tutone did
//       emit this method.
// Regenerating pkg/scorecards should leave this file intact.

import "encoding/json"

// EntityManagementCollectionElementsFilter is the required filter object
// for the `collectionElements` query. Only `collectionId` (equality) is
// currently accepted server-side.
type EntityManagementCollectionElementsFilter struct {
	CollectionID EntityManagementCollectionIdFilterArgument `json:"collectionId"`
}

// EntityManagementCollectionIdFilterArgument is the sole predicate available
// under EntityManagementCollectionElementsFilter — an equality match on the
// parent collection's ID.
type EntityManagementCollectionIdFilterArgument struct {
	Eq string `json:"eq"`
}

// UnmarshalJSON dispatches the polymorphic `entity` field on
// EntityManagementActorStitchedFields to the correct concrete implementation
// based on `__typename`. The other three fields are decoded normally.
func (x *EntityManagementActorStitchedFields) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	if err := json.Unmarshal(b, &objMap); err != nil {
		return err
	}

	for k, v := range objMap {
		if v == nil {
			continue
		}
		switch k {
		case "entity":
			xxx, err := UnmarshalEntityManagementEntityInterface(*v)
			if err != nil {
				return err
			}
			if xxx != nil {
				x.Entity = *xxx
			}
		case "entitySearch":
			if err := json.Unmarshal(*v, &x.EntitySearch); err != nil {
				return err
			}
		case "collectionElements":
			if err := json.Unmarshal(*v, &x.CollectionElements); err != nil {
				return err
			}
		case "collectionsContainingEntity":
			if err := json.Unmarshal(*v, &x.CollectionsContainingEntity); err != nil {
				return err
			}
		}
	}
	return nil
}
