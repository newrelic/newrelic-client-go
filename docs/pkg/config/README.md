# config
--
    import "github.com/newrelic/newrelic-client-go/pkg/config"


## Usage

```go
const (
	// US represents New Relic's US-based production deployment.
	US = iota

	// EU represents New Relic's EU-based production deployment.
	EU

	// Staging represents New Relic's US-based staging deployment.
	// This is for internal New Relic use only.
	Staging
)
```

```go
var Region = struct {
	US      RegionType
	EU      RegionType
	Staging RegionType
}{
	US:      US,
	EU:      EU,
	Staging: Staging,
}
```
Region specifies the New Relic environment to target.

#### type RegionType

```go
type RegionType int
```

RegionType represents the members of the Region enumeration.

#### type ReplacementConfig

```go
type ReplacementConfig struct {
	BaseURL       string
	APIKey        string
	Timeout       *time.Duration
	HTTPTransport *http.RoundTripper
	UserAgent     string
	Region        RegionType
}
```

Config contains all the configuration data for the API Client.
