package notebooks

import (
	"encoding/json"
	"fmt"
)

// customHeadersKey is the sentinel field name that the internal HTTP client
// looks for in the queryParams map to know it should convert the entries into
// request headers. It's a private convention of internal/http/client.go -
// keeping the string in one place here so nobody has to remember it.
const customHeadersKey = "x-newrelic-client-go-custom-headers"

// newRelicEntityHeader is the name of the Blob Storage API's entity-metadata
// header. On create it carries `{"name": "<notebook name>"}`; on update it can
// be re-supplied to rename the notebook atomically alongside the content
// change.
const newRelicEntityHeader = "NewRelic-Entity"

// notebooksURL builds a Blob Storage API URL for the notebooks resource. If
// entityGUID is empty, returns the collection URL used for creation.
func (a *Notebooks) notebooksURL(organizationID, entityGUID string) string {
	if entityGUID == "" {
		return a.config.Region().BlobServiceURL(fmt.Sprintf("/organizations/%s/Notebooks", organizationID))
	}
	return a.config.Region().BlobServiceURL(fmt.Sprintf("/organizations/%s/Notebooks/%s", organizationID, entityGUID))
}

// blobRequestHeaders returns the map shape the internal HTTP client expects
// when the caller wants to add custom headers to a request. Returns an
// untyped nil when no headers were supplied - the return type must be
// interface{} rather than map[string]interface{} because a typed-nil map
// would satisfy != nil at the queryParams interface{} check upstream, then
// fall through to go-querystring's Values() which rejects maps.
func blobRequestHeaders(headers map[string]string) interface{} {
	if len(headers) == 0 {
		return nil
	}
	return map[string]interface{}{customHeadersKey: headers}
}

// entityHeaderMarshaler is the JSON marshaler used by encodeEntityHeader.
// It is a package-level variable so unit tests can substitute a failing
// marshaler and exercise the defensive panic below. In production it is
// always json.Marshal.
var entityHeaderMarshaler = json.Marshal

// encodeEntityHeader serialises the NewRelic-Entity header value. json.Marshal
// on a map[string]string is guaranteed by the encoding/json contract not to
// fail, so callers can consume the result directly without an error path.
// The panic is documented as unreachable and exists only to satisfy the
// json.Marshal signature; if it fires the runtime is in an inconsistent state
// and there is nothing meaningful the caller could do to recover.
func encodeEntityHeader(name string) string {
	b, err := entityHeaderMarshaler(map[string]string{"name": name})
	if err != nil {
		panic(fmt.Sprintf("notebooks: unreachable json.Marshal failure encoding entity header: %v", err))
	}
	return string(b)
}
