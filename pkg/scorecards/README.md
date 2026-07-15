# scorecards

The `scorecards` package covers four related NGEP (New Relic Entity Platform) resource families that all live under `entityManagement`:

- **Team** (`TEAM`)
- **Scorecard** (`SCORECARD`)
- **ScorecardRule** (`SCORECARD_RULE`)
- **Collection** (`COLLECTION`) — includes the auto-created backing collections attached to Teams (`membership`, `ownership`) and Scorecards (`rules`)

The name of the package is provisional; it groups everything into one place because the resources share the same GraphQL interface, the same `entityManagement*` mutations, and the same reads. See [`NGEP_ANALYSIS.md`](./NGEP_ANALYSIS.md) for the API relationships, side-effect map, and the transient-error taxonomy this package works around.

## Install

```go
import "github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
```

## Create a client

```go
package main

import (
    "github.com/newrelic/newrelic-client-go/v2/pkg/config"
    "github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

func main() {
    cfg := config.New()
    cfg.PersonalAPIKey = "YOUR_API_KEY"
    // cfg.Region = "EU" // optional; default is US
    client := scorecards.New(cfg)
    _ = client
}
```

## Scope

All entities exposed by this package are `ORGANIZATION`-scoped. You'll need your organization's UUID — fetch it once with:

```graphql
{ actor { organization { id name } } }
```

`Collection` also accepts `ACCOUNT` scope in the schema, but RBAC on many orgs rejects it; prefer `ORGANIZATION` unless you know account-scoped collections are enabled for you.

## Teams

### Create

```go
res, err := client.EntityManagementCreateTeam(scorecards.EntityManagementTeamEntityCreateInput{
    Name:        "platform",
    Description: "Platform engineering",
    Scope: scorecards.EntityManagementScopedReferenceInput{
        ID:   orgID,
        Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION,
    },
    Aliases: []string{"platform-alias"},
    Tags:    []scorecards.EntityManagementTagInput{{Key: "team", Values: []string{"platform"}}},
})
```

Creating a Team **auto-creates two Collections** returned as `res.Entity.Membership` and `res.Entity.Ownership`. Their IDs are stable across future Team updates.

### Field-by-field update

The Update input is a JSON-Patch-style partial update — fields you omit are left alone. Explicit `nil`/`[]` for list-typed fields clears them (and the two are distinguishable server-side).

```go
_, err := client.EntityManagementUpdateTeam(res.Entity.ID, scorecards.EntityManagementTeamEntityUpdateInput{
    Description: "Now with 20% more platform",
})
```

Notable constraints — see `NGEP_ANALYSIS.md § 5.1` for the full matrix:

- `managers` requires **NGEP-encoded user IDs** (`NDgy…VVNFUnw…`), not raw int userIds — and the user must already be in the team's `membership` Collection.
- `externalIntegration.type = "GITHUB_TEAM"` / `"SERVICENOW_TEAM"` are managed by integrations and cannot be set through this API. `IAM_GROUP` is user-settable but only one integration type per org.
- `parentId` is validated: no self-parenting, must point at an existing Team.

### Add users to the team

```go
_, err := client.EntityManagementAddCollectionMembers(res.Entity.Membership.ID, []string{userNGEPID})
```

## Scorecards

```go
sc, err := client.EntityManagementCreateScorecard(scorecards.EntityManagementScorecardEntityCreateInput{
    Name:  "engineering-best-practices",
    Scope: scorecards.EntityManagementScopedReferenceInput{ID: orgID, Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION},
    ProgressLevels: []scorecards.EntityManagementProgressLevelDefinitionCreateInput{
        {ID: "red",   Name: "Red",   HexColorCode: "#FF0000"},
        {ID: "green", Name: "Green", HexColorCode: "#00CC00"},
    },
})
// sc.Entity.Rules.ID is the auto-created rules Collection.
```

Note: at the time of writing, changing `progressLevels` through `EntityManagementUpdateScorecard` triggers a backend bug (`Unknown type 'EntityManagementScorecardRuleEntity'`). Treat progress levels as **create-time-only** until fixed.

## ScorecardRules

Rules are standalone entities — they exist even without being attached to a Scorecard, and detaching a rule does NOT delete it.

```go
r, err := client.EntityManagementCreateScorecardRule(scorecards.EntityManagementScorecardRuleEntityCreateInput{
    Name:    "APM apps have alerts",
    Enabled: true,
    Scope:   scorecards.EntityManagementScopedReferenceInput{ID: orgID, Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION},
    NRQLEngine: &scorecards.EntityManagementNRQLRuleEngineCreateInput{
        Accounts: []int{accountID},
        Query:    "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago",
    },
})

// Attach the rule to a scorecard's rules Collection.
_, err = client.EntityManagementAddCollectionMembers(sc.Entity.Rules.ID, []string{r.Entity.ID})
```

NRQL constraints (per the current backend):

- Must be `FACET`-ed.
- The SELECT clause must alias to `'score'` and be either `if(...)` or `latest(...)`.
- `runInterval` (minutes) must be one of `{60, 360, 720, 1440, 4320}`.
- `schedule` is declared on the schema but rejected at runtime for rules — do not set it.

## Collections (standalone)

Collections created directly (not the ones auto-generated for Teams/Scorecards) can hold **heterogeneous** members — Users, APM apps, Teams, other Collections. Deleting a Collection detaches its members silently; they survive.

```go
col, err := client.EntityManagementCreateCollection(scorecards.EntityManagementCollectionEntityCreateInput{
    Name:  "my-favourite-apps",
    Scope: scorecards.EntityManagementScopedReferenceInput{ID: orgID, Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION},
})

// Response is slot-aligned with your input `ids`. Each slot is either the
// added id or nil on failure; failures are surfaced through the returned
// error (transport error) with per-index paths in `errors[]`.
_, err = client.EntityManagementAddCollectionMembers(col.Entity.ID, []string{apmEntityGUID, userNGEPID})
```

### Read the members

```go
res, err := client.GetCollectionElements(
    "",
    scorecards.EntityManagementCollectionElementsFilter{
        CollectionID: scorecards.EntityManagementCollectionIdFilterArgument{Eq: col.Entity.ID},
    },
    25, // page size
)
for _, item := range res.Items {
    switch e := item.(type) {
    case *scorecards.EntityManagementUserEntity:
        fmt.Println("user:", e.Name, e.UserId)
    default:
        // Anything that isn't a known implementation (e.g. APM applications)
        // deserialises as *EntityManagementGenericEntity.
    }
}
```

### Reverse-lookup

```go
containers, err := client.GetCollectionsContainingEntity(userNGEPID)
for _, c := range *containers {
    fmt.Println(c.Collection.Name, "→ parent field:", c.ParentInfo.ParentField)
    // ParentField is "membership" / "ownership" / "rules" / "" (user-created)
}
```

## Deletes and cascade

`entityManagementDelete(id)` is used for every entity type in this package.

- Deleting a **Team** removes both auto-created backing Collections. Members (Users) survive.
- Deleting a **Scorecard** removes its `rules` backing Collection. Attached ScorecardRules become **orphans** — they are not cascaded. Delete them yourself.
- Deleting a **standalone Collection** silently detaches its members.
- The **backing** membership/ownership/rules Collections cannot be deleted directly (`Cannot delete Collection entity with existing parent`) — delete the parent instead.

## Regenerating this package via tutone

This package relies on the tutone fork with `include_implementations` (see [`.tutone.yml`](../../.tutone.yml)). A few hand-patches are applied after every `tutone generate --packages scorecards`:

1. **Strip the `version` argument** from every generated `Update*` and `Delete` method. NGEP treats `version: 0` as "wrong version" and returns `Concurrent modification of entity`. The regen recipe strips `version int` parameters, the corresponding `"version": version` variables, and the `$version: Int` / `version: $version` selections from the mutation string. See the small Python sed at the top of `Makefile`-adjacent regen scripts, or reuse the block in `NGEP_ANALYSIS.md § 11`.
2. **Two hand-written files** are kept alongside the generated code:
   - `scorecards_collection_members.go` — hand-writes the `Add/RemoveCollectionMembers` mutations because tutone cannot yet emit valid Go for their `[ID]` scalar-list return.
   - `scorecards_types_extra.go` — adds the missing `EntityManagementCollectionElementsFilter` and `EntityManagementCollectionIdFilterArgument` input types plus an `UnmarshalJSON` for `EntityManagementActorStitchedFields` (dispatches the `entity` interface field via `__typename`).

Do not commit the vendored `tutone` binary or the `schema.json` file that tutone writes to the worktree root — they are gitignored and only used at regen time.

## Links

- [NGEP scorecards tutorial (docs)](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-scorecards-tutorial/)
- [NGEP scorecards custom-roles tutorial (docs)](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-scorecards-custom-tutorial/)
