package nrqlmacros

import "encoding/json"

// DO NOT DELETE
//
// Tutone generates "// special" UnmarshalJSON for structs that contain a *slice* of an
// interface type (e.g. EntityManagementEntitySearchResult.Entities []EntityManagementEntityInterface),
// but NOT for structs that contain a single interface-typed field.
// EntityManagementActorStitchedFields.Entity is the latter case, so tutone skips it.
// Without this function, GetEntity always returns a JSON decode error because
// encoding/json cannot unmarshal into a bare interface value.
//
// The identical gap exists in pkg/pipelinecontrol/types.go where the same function
// carries the same "// DO NOT DELETE" marker. This is not an include_implementations
// issue — it is a tutone limitation specific to single-value (non-slice) interface fields.
//
// Regenerating pkg/nrqlmacros should leave this file intact.
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
		}
	}
	return nil
}
