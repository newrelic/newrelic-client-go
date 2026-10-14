# servicearchintelligence

The `servicearchintelligence` package covers four related NGEP (New Relic Entity Platform) resource families that all live under `entityManagement`:

- **Team** (`TEAM`)
- **Scorecard** (`SCORECARD`)
- **ScorecardRule** (`SCORECARD_RULE`)
- **Collection** (`COLLECTION`) — includes the auto-created backing collections attached to Teams (`membership`, `ownership`) and Scorecards (`rules`)

All resources share the same GraphQL interface, the same `entityManagement*` mutations, and the same read path.

## Install

```go
import "github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
```

## Create a client

```go
package main

import (
    "github.com/newrelic/newrelic-client-go/v2/pkg/config"
    "github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
)

func main() {
    cfg := config.New()
    cfg.PersonalAPIKey = "YOUR_API_KEY"
    client := servicearchintelligence.New(cfg)
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
res, err := client.EntityManagementCreateTeam(servicearchintelligence.EntityManagementTeamEntityCreateInput{
    Name:        "platform",
    Description: "Platform engineering",
    Scope: servicearchintelligence.EntityManagementScopedReferenceInput{
        ID:   orgID,
        Type: servicearchintelligence.EntityManagementEntityScopeTypes.ORGANIZATION,
    },
    Aliases: []string{"platform-alias"},
    Tags:    []servicearchintelligence.EntityManagementTagInput{{Key: "team", Values: []string{"platform"}}},
})
```

Creating a Team **auto-creates two Collections** returned as `res.Entity.Membership` and `res.Entity.Ownership`. Their IDs are stable across future Team updates.

### Field-by-field update

The Update input is a JSON-Patch-style partial update — fields you omit are left alone. Explicit `nil`/`[]` for list-typed fields clears them.

```go
_, err := client.EntityManagementUpdateTeam(res.Entity.ID, servicearchintelligence.EntityManagementTeamEntityUpdateInput{
    Description: "Now with 20% more platform",
})
```

Notable constraints:

- `managers` requires **NGEP-encoded user IDs** (`NDgy…VVNFUnw…`), not raw int userIds — and the user must already be in the team's `membership` Collection.
- `externalIntegration.type = "GITHUB_TEAM"` / `"SERVICENOW_TEAM"` are managed by integrations and cannot be set through this API. `IAM_GROUP` is user-settable but only one integration type per org.
- `parentId` is validated: no self-parenting, must point at an existing Team.

### Add users to the team

```go
_, err := client.EntityManagementAddCollectionMembers(res.Entity.Membership.ID, []string{userNGEPID})
```

## Scorecards

```go
sc, err := client.EntityManagementCreateScorecard(servicearchintelligence.EntityManagementScorecardEntityCreateInput{
    Name:  "engineering-best-practices",
    Scope: servicearchintelligence.EntityManagementScopedReferenceInput{ID: orgID, Type: servicearchintelligence.EntityManagementEntityScopeTypes.ORGANIZATION},
    ProgressLevels: []servicearchintelligence.EntityManagementProgressLevelDefinitionCreateInput{
        {ID: "red",   Name: "Red",   HexColorCode: "#FF0000"},
        {ID: "green", Name: "Green", HexColorCode: "#00CC00"},
    },
})
// sc.Entity.Rules.ID is the auto-created rules Collection.
```

Progress levels can also be added, changed, or removed in-place after creation via `EntityManagementUpdateScorecard`.

## ScorecardRules

Rules are standalone entities — they exist even without being attached to a Scorecard, and detaching a rule does NOT delete it.

```go
r, err := client.EntityManagementCreateScorecardRule(servicearchintelligence.EntityManagementScorecardRuleEntityCreateInput{
    Name:    "APM apps have alerts",
    Enabled: true,
    Scope:   servicearchintelligence.EntityManagementScopedReferenceInput{ID: orgID, Type: servicearchintelligence.EntityManagementEntityScopeTypes.ORGANIZATION},
    NRQLEngine: &servicearchintelligence.EntityManagementNRQLRuleEngineCreateInput{
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
col, err := client.EntityManagementCreateCollection(servicearchintelligence.EntityManagementCollectionEntityCreateInput{
    Name:  "my-favourite-apps",
    Scope: servicearchintelligence.EntityManagementScopedReferenceInput{ID: orgID, Type: servicearchintelligence.EntityManagementEntityScopeTypes.ORGANIZATION},
})

_, err = client.EntityManagementAddCollectionMembers(col.Entity.ID, []string{apmEntityGUID, userNGEPID})
```

### Read the members

```go
res, err := client.GetCollectionElements(
    "",
    servicearchintelligence.EntityManagementCollectionElementsFilter{
        CollectionID: servicearchintelligence.EntityManagementCollectionIdFilterArgument{Eq: col.Entity.ID},
    },
    25, // page size
)
for _, item := range res.Items {
    switch e := item.(type) {
    case *servicearchintelligence.EntityManagementUserEntity:
        fmt.Println("user:", e.Name, e.UserId)
    default:
        // Unknown implementations deserialise as *EntityManagementGenericEntity.
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

This package relies on the tutone fork with `include_implementations` (see [`.tutone.yml`](../../.tutone.yml)). Two hand-written files are kept alongside the generated code:

- `scorecards_collection_members.go` — hand-writes the `Add/RemoveCollectionMembers` mutations because tutone cannot yet emit valid Go for their `[ID]` scalar-list return.
- `scorecards_types_extra.go` — adds the missing `EntityManagementCollectionElementsFilter` and `EntityManagementCollectionIdFilterArgument` input types, plus an `UnmarshalJSON` for `EntityManagementActorStitchedFields` that dispatches the `entity` interface field via `__typename`.

After every `tutone generate --packages scorecards`, strip the `version` argument from all generated `Update*` and `Delete` methods. NGEP treats `version: 0` as "wrong version" and returns `Concurrent modification of entity`. Remove `version int` parameters, the corresponding `"version": version` variables, and the `$version: Int` / `version: $version` selections from the mutation strings.

Do not commit the vendored `tutone` binary or the `schema.json` file — they are gitignored and only used at regen time.

## Links

- [NGEP scorecards tutorial (docs)](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-scorecards-tutorial/)
- [NGEP scorecards custom-roles tutorial (docs)](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-scorecards-custom-tutorial/)
