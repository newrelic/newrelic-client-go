# synthetics
--
    import "."


## Usage

#### type Monitor

```go
type Monitor struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Frequency    uint           `json:"frequency"`
	URI          string         `json:"uri"`
	Locations    []string       `json:"locations"`
	Status       string         `json:"status"`
	SLAThreshold float64        `json:"slaThreshold"`
	UserID       uint           `json:"userId,omitempty"`
	APIVersion   string         `json:"apiVersion,omitempty"`
	ModifiedAt   time.Time      `json:"modified_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	Options      MonitorOptions `json:"options,omitempty"`
}
```

Monitor represents a New Relic Synthetics monitor.

#### type MonitorOptions

```go
type MonitorOptions struct {
	ValidationString       string `json:"validationString,omitempty"`
	VerifySSL              bool   `json:"verifySSL,omitempty"`
	BypassHEADRequest      bool   `json:"bypassHEADRequest,omitempty"`
	TreatRedirectAsFailure bool   `json:"treatRedirectAsFailure,omitempty"`
}
```


#### type Synthetics

```go
type Synthetics struct {
}
```


#### func  New

```go
func New(config config.Config) Synthetics
```
New is used to create a new Synthetics client instance.

#### func (*Synthetics) ListMonitors

```go
func (s *Synthetics) ListMonitors() ([]Monitor, error)
```
ListMonitors is used to retrieve New Relic Synthetics monitors.
