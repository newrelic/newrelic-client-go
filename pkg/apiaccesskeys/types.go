package apiaccesskeys

// APIAccessCreateKeysInput represents a list of the configurations for each key you want to create.
type APIAccessCreateKeysInput struct {
	Ingest []APIAccessCreateIngestKeyInput `json:"ingest,omitempty"`
	User   []APIAccessCreateUserKeyInput   `json:"user,omitempty"`
}

// APIAccessCreateIngestKeyInput represents the input for any ingest keys you want to create.
// Each ingest key must have a type that communicates what kind of data it is for.
// You can optionally add a name or notes to your key, which can be updated later.
type APIAccessCreateIngestKeyInput struct {
	AccountID  int    `json:"accountId"`
	IngestType string `json:"ingestType"`
	Name       string `json:"name"`
	Notes      string `json:"notes"`
}

// APIAccessCreateUserKeyInput represents a request to create user keys. You can optionally add a name or notes to your key,
// which can be updated later.
type APIAccessCreateUserKeyInput struct {
	AccountID int    `json:"accountId"`
	UserID    int    `json:"userId"`
	Name      string `json:"name"`
	Notes     string `json:"notes"`
}

// APIAccessDeleteInput represents a list of each key id that you want to delete.
type APIAccessDeleteInput struct {
	IngestKeyIds []string `json:"ingestKeyIds,omitempty"`
	UserKeyIds   []string `json:"userKeyIds,omitempty"`
}

// APIAccessGetInput represents a single key by ID and type that you want to retrieve from the API.
type APIAccessGetInput struct {
	ID      string `json:"id"`
	KeyType string `json:"keyType"`
}

// APIAccessUpdateInput represents the input needed to update an API access key.
type APIAccessUpdateInput struct {
	Ingest []APIAccessUpdateKeyInput `json:"ingest,omitempty"`
	User   []APIAccessUpdateKeyInput `json:"user,omitempty"`
}

// APIAccessUpdateKeyInput represents the individual fields required to update an API access key.
type APIAccessUpdateKeyInput struct {
	ID    string `json:"keyId"`
	Name  string `json:"name"`
	Notes string `json:"notes"`
}
