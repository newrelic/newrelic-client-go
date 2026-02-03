# Service Name Auto-Detection

## Overview

The newrelic-client-go now automatically detects the calling service name for internal telemetry purposes. This helps New Relic track which applications and services are using the client library, without requiring users to manually configure service names.

## How It Works

When a client is initialized without an explicit `ServiceName` configured, the library attempts to auto-detect the calling service using the following priority:

1. **Go Module Path** (Primary) - Reads the main module path from Go's build info
   - Example: `github.com/acme-corp/my-monitoring-service`
   - Works with any module path (GitHub, GitLab, corporate internal repos, etc.)
   - Only works when the application is built properly (`go build`, `go install`)

2. **Binary Name** (Fallback) - Uses the executable's filename
   - Example: `my-monitoring-service`
   - Works when module info isn't available (e.g., `go run`)

3. **Process Name** (Last Resort) - Uses the process name from `os.Args[0]`
   - Used when other methods fail

If all detection methods fail, the library falls back to the default: `"newrelic-client-go"`

## Header Format

The detected service name is included in the `NewRelic-Requesting-Services` header:

```
NewRelic-Requesting-Services: <detected-name>|newrelic-client-go
```

Examples:
- With auto-detection: `github.com/acme-corp/monitoring-app|newrelic-client-go`
- With custom name: `my-custom-app|newrelic-client-go`
- Detection failed: `newrelic-client-go`
- With env var: `github-action|detected-name|newrelic-client-go`

## Customer Impact

### No Breaking Changes
- Existing behavior is preserved for customers who explicitly set `ServiceName`
- No new configuration required
- No errors thrown if detection fails
- Completely transparent to end users

### Backward Compatibility
- Custom service names continue to work exactly as before
- `NEW_RELIC_SERVICE_NAME` environment variable still prepends as expected
- Terraform provider detection logic remains intact

## Implementation Details

### Key Files Modified
1. `internal/http/client.go` - Modified `NewClient()` to call auto-detection
2. `internal/http/module_detector.go` - New file with detection logic
3. `internal/http/module_detector_test.go` - Comprehensive test coverage

### Safety Features
- All detection functions are panic-safe (use defer/recover)
- Detection failures never block client creation
- Test binaries are filtered out (`.test`, temp directories, etc.)
- The library's own module path is filtered out

### Test Coverage
- Unit tests for all detection functions
- Integration tests for `NewClient()` behavior
- Edge case handling (empty detection, test binaries, env variables)
- All existing tests continue to pass

## When Detection Works

✅ **Works automatically:**
- Applications built with `go build` or `go build .`
- Applications installed with `go install`
- Docker container builds
- CI/CD pipeline builds
- Any module path (GitHub, GitLab, Bitbucket, corporate internal, custom domains)

⚠️ **Limited detection:**
- `go run main.go` - Falls back to binary name (e.g., "main" or temp binary name)
- Development environments without proper builds

❌ **No detection (uses default):**
- Test binaries (automatically filtered out)
- Builds within the newrelic-client-go repo itself (for library development)

## Examples

### Example 1: Corporate Internal Service
```go
// Module: go.internal.acme-corp.com/observability/monitoring-service

import "github.com/newrelic/newrelic-client-go/v2/newrelic"

client, _ := newrelic.New(newrelic.ConfigPersonalAPIKey("NRAK-..."))
// Header: go.internal.acme-corp.com/observability/monitoring-service|newrelic-client-go
```

### Example 2: Custom Service Name (Explicit)
```go
client, _ := newrelic.New(
    newrelic.ConfigPersonalAPIKey("NRAK-..."),
    newrelic.ConfigServiceName("my-custom-app"),
)
// Header: my-custom-app|newrelic-client-go
```

### Example 3: With Environment Variable
```bash
export NEW_RELIC_SERVICE_NAME="github-action"
```
```go
client, _ := newrelic.New(newrelic.ConfigPersonalAPIKey("NRAK-..."))
// Header: github-action|github.com/user/repo|newrelic-client-go
```

### Example 4: Terraform Provider (Special Handling)
```go
// Module: github.com/newrelic/terraform-provider-newrelic

// Auto-detected, and special header is set:
// X-Query-Source-Capability-Id: TERRAFORM
```

## Technical Notes

### Why Module Path Detection Works

Go's `debug.ReadBuildInfo()` returns the **main module** that was built, not the library's module. When a customer imports newrelic-client-go as a dependency and builds their application, `ReadBuildInfo()` returns the customer's module information, which we use for tracking.

### Filtering Logic

The implementation filters out unhelpful identifiers:
- Test binaries (`.test` suffix, `_` prefix, `go-build` in path)
- Generic names (`main`)
- The library's own module path
- Temporary directories (`/tmp/`, `/T/`)

This ensures we only report meaningful service identifiers.

## Testing

Run tests with:
```bash
go test -v -tags unit github.com/newrelic/newrelic-client-go/v2/internal/http
```

All existing tests pass, plus new tests verify:
- Detection functions work correctly
- No panics occur
- Proper fallback behavior
- Service name formatting is correct
- Environment variables work as expected
