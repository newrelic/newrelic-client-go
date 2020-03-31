package nerdgraph
type NerdStorageScope int 

const (
	ACCOUNT = iota
	ACTOR
	ENTITY
)
type NerdStorageScopeInput struct {
	 Id string `json:"id"`

	 Name NerdStorageScope `json:"name"`

}
