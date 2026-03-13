package notifications

// EntityScopeTypeInput - Scope type for destinations
type EntityScopeTypeInput string

var EntityScopeTypeInputTypes = struct {
	// Organization scope type
	ORGANIZATION EntityScopeTypeInput
	// Account scope type
	ACCOUNT EntityScopeTypeInput
}{
	ORGANIZATION: "ORGANIZATION",
	ACCOUNT:      "ACCOUNT",
}

// EntityScopeInput - Scope input for destinations
type EntityScopeInput struct {
	// id - Organization UUID for ORGANIZATION scope, Account ID for ACCOUNT scope
	ID string `json:"id"`
	// type - Scope type (ORGANIZATION or ACCOUNT)
	Type EntityScopeTypeInput `json:"type"`
}

// EntityScope - Scope response from API
type EntityScope struct {
	// id
	ID string `json:"id,omitempty"`
	// type
	Type EntityScopeTypeInput `json:"type,omitempty"`
}
