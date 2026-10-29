# nrqlmacros

The `nrqlmacros` package provides CRUD operations for [NRQL Macros](https://docs.newrelic.com/docs/nrql/nrql-macros/nrql-macros-overview/) via the New Relic Entity Management API.

NRQL macros are reusable expression fragments that can be substituted into NRQL queries using the `$macroName` syntax. They are scoped to an organization and available across all accounts within that organization.

## Install

```go
import "github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
```

## Create a client

```go
package main

import (
    "github.com/newrelic/newrelic-client-go/v2/pkg/config"
    "github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
)

func main() {
    cfg := config.New()
    cfg.PersonalAPIKey = "YOUR_API_KEY"
    client := nrqlmacros.New(cfg)
    _ = client
}
```

## Key types

- `EntityManagementNRQLMacroEntity` — the macro entity returned by read/search operations
- `EntityManagementNRQLMacroEntityCreateInput` — input for creating a macro
- `EntityManagementNRQLMacroEntityUpdateInput` — input for updating a macro (expression, tags, numParameters only; name and scope are immutable)
- `EntityManagementScopedReferenceInput` — scope reference (only `ORGANIZATION` scope is supported)

## Examples

### Create a macro

```go
result, err := client.EntityManagementCreateNRQLMacro(nrqlmacros.EntityManagementNRQLMacroEntityCreateInput{
    Name:       "filterProdTraffic",
    Expression: "environment = 'production' AND appName NOT LIKE '%test%'",
    Scope: nrqlmacros.EntityManagementScopedReferenceInput{
        Type: nrqlmacros.EntityManagementEntityScopeTypes.ORGANIZATION,
        ID:   "YOUR_ORG_GUID",
    },
})
```

Reference in a query: `SELECT count(*) FROM Transaction WHERE $filterProdTraffic`

### Read a macro by ID

```go
entity, err := client.GetEntity("MACRO_ENTITY_ID")
macro, ok := (*entity).(*nrqlmacros.EntityManagementNRQLMacroEntity)
```

### List all macros in the organization

```go
result, err := client.GetEntitySearch("", "type = 'NRQL_MACRO'")
for _, e := range result.Entities {
    if macro, ok := e.(*nrqlmacros.EntityManagementNRQLMacroEntity); ok {
        fmt.Println(macro.Name, macro.Expression)
    }
}
```

### Update a macro's expression

```go
result, err := client.EntityManagementUpdateNRQLMacro(macroID, nrqlmacros.EntityManagementNRQLMacroEntityUpdateInput{
    Expression: "environment = 'staging' AND appName NOT LIKE '%test%'",
})
```

> **Warning:** Updating a macro's expression immediately affects all dashboards, alerts, and queries that reference it.

### Delete a macro

```go
_, err := client.EntityManagementDelete(macroID)
```

## Behavioral notes

- **Scope**: Only `ORGANIZATION` scope is supported. Attempting to use `ACCOUNT` scope returns an error.
- **Immutability**: A macro's `name` and `scope` cannot be changed after creation.
- **Instant propagation**: When a macro's `expression` is updated, all queries, dashboards, and alert conditions that reference it will use the new expression on their next execution.
- **Optimistic locking**: The server tracks a `version` field. Both update and delete accept an optional version parameter for conflict detection (not exposed in this client by default).

## Regenerating this package

```bash
tutone -c .tutone.yml generate --package nrqlmacros --refetch
```

`--refetch` is required. Several types in this package
(`EntityManagementNrqlMacroEntityCreateInput`, `EntityManagementNrqlMacroEntityUpdateInput`,
`EntityManagementScopedReferenceInput`, `EntityManagementTagInput`) do not appear in
tutone's general schema introspection dump but are reachable when queried by name —
`--refetch` triggers those targeted `__type` lookups.

**After every regeneration, two manual steps are required:**

**1. In `types.go` — restore `*int` on `EntityManagementNRQLMacroEntityUpdateInput.NumParameters`:**

Tutone regenerates this field as `int` (the raw schema type). It is marked
`// DO NOT REPLACE` in `types.go` and must be kept as `*int` so callers can explicitly
set `numParameters` to zero without the value being silently dropped by `omitempty`.

```go
// Tutone output (revert this):
NumParameters int `json:"numParameters,omitempty"`

// Correct (restore to this):
NumParameters *int `json:"numParameters,omitempty"`
```

**2. In `nrqlmacros_api.go` — re-comment the `version int` parameters:**

Tutone restores `version int` on `EntityManagementDelete` and
`EntityManagementUpdateNRQLMacro` (and their `WithContext` variants). Comment them out
again to match the pipelinecontrol convention:

```go
// Tutone output (revert this):
func (a *Nrqlmacros) EntityManagementDelete(iD string, version int) ...

// Correct (restore to this):
func (a *Nrqlmacros) EntityManagementDelete(
    iD string,
    // version int,
) ...
```

`nrqlmacros_extra.go` requires no changes after regeneration.

## Links

- [NRQL Macros overview](https://docs.newrelic.com/docs/nrql/nrql-macros/nrql-macros-overview/)
- [What's new: NRQL Macros GA](https://docs.newrelic.com/whats-new/2026/10/whats-new-10-26-nrql-macros-ga/)
