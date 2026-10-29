package nrqlmacros

import "encoding/json"

// This file supplies types that tutone's generator misses because the
// EntityManagementNrqlMacroEntityCreateInput and EntityManagementNrqlMacroEntityUpdateInput
// types are absent from tutone's schema introspection dump. The NerdGraph mutations are
// real and tested; the input types here are derived directly from live API responses.
//
// Regenerating pkg/nrqlmacros should leave this file intact.

// EntityManagementNRQLMacroEntityCreateInput is the input for creating a NRQL macro.
type EntityManagementNRQLMacroEntityCreateInput struct {
	// The NRQL expression fragment that this macro expands to. Required.
	Expression string `json:"expression"`
	// A unique, human-readable name for the macro. Immutable after creation. Required.
	Name string `json:"name"`
	// Number of positional parameters ($1, $2, ...) used in the expression. Optional.
	NumParameters int `json:"numParameters,omitempty"`
	// The scope of the macro. Only ORGANIZATION scope is currently supported. Required.
	Scope EntityManagementScopedReferenceInput `json:"scope"`
	// Optional key-value tags to associate with the macro.
	Tags []EntityManagementTagInput `json:"tags,omitempty"`
}

// EntityManagementNRQLMacroEntityUpdateInput is the input for updating a NRQL macro.
// Only expression, numParameters, and tags are mutable; name and scope cannot be changed.
type EntityManagementNRQLMacroEntityUpdateInput struct {
	// The updated NRQL expression fragment. Optional.
	Expression string `json:"expression,omitempty"`
	// Updated parameter count. Optional.
	NumParameters int `json:"numParameters,omitempty"`
	// Updated tags. Optional.
	Tags []EntityManagementTagInput `json:"tags,omitempty"`
}

// EntityManagementNRQLMacroEntityCreateResult holds the result of a create mutation.
type EntityManagementNRQLMacroEntityCreateResult struct {
	Entity EntityManagementNRQLMacroEntity `json:"entity"`
}

// EntityManagementNRQLMacroEntityUpdateResult holds the result of an update mutation.
type EntityManagementNRQLMacroEntityUpdateResult struct {
	Entity EntityManagementNRQLMacroEntity `json:"entity"`
}

// EntityManagementScopedReferenceInput is the input for specifying the scope of an entity.
type EntityManagementScopedReferenceInput struct {
	// The ID of the scope target (e.g. organization GUID).
	ID string `json:"id"`
	// The type of the scope (ACCOUNT or ORGANIZATION).
	Type EntityManagementEntityScope `json:"type"`
}

// EntityManagementTagInput represents a single key-value tag.
type EntityManagementTagInput struct {
	// The tag key.
	Key string `json:"key"`
	// The tag values (one or more).
	Values []string `json:"values"`
}

// GetID satisfies the EntityManagementEntityInterface for NRQLMacroEntity.
func (x *EntityManagementNRQLMacroEntity) GetID() string { return x.ID }

// UnmarshalJSON provides custom JSON unmarshaling for EntityManagementActorStitchedFields
// because the Entity field is an interface. Tutone does not generate this for the
// nrqlmacros package; without it GetEntity returns a decode error.
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
