# NGEP Analysis — Teams, Scorecards, Collections

Live-API exploration performed against `https://api.newrelic.com/graphql` using an ORG-scoped personal API key. This document captures the shape, semantics, side-effects, cascade behaviour, and quirks that shape how the `pkg/scorecards` package (which covers **all three** related resources) is designed.

---

## 1. Entity type inventory

All these types implement the `EntityManagementEntity` interface.

| GraphQL type | `type` field value | Scope | Auto-created? |
|---|---|---|---|
| `EntityManagementTeamEntity` | `TEAM` | `ORGANIZATION` only | No |
| `EntityManagementTeamsHierarchyLevelEntity` | `TEAMS_HIERARCHY_LEVEL` | `ORGANIZATION` | No (RBAC-gated mutation) |
| `EntityManagementTeamsOrganizationSettingsEntity` | `TEAMS_ORGANIZATION_SETTINGS` | `ORGANIZATION` | No (RBAC-gated mutation) |
| `EntityManagementScorecardEntity` | `SCORECARD` | `ORGANIZATION` only | No |
| `EntityManagementScorecardRuleEntity` | `SCORECARD_RULE` | `ORGANIZATION` | No |
| `EntityManagementCollectionEntity` | `COLLECTION` | `ORGANIZATION` (ACCOUNT allowed at schema but RBAC-restricted) | **Yes**, from Team + Scorecard creation |
| `EntityManagementUserEntity` | `USER` | `ORGANIZATION` | Managed by identity subsystem |

Entity IDs decode from base64 to `<accountId>|NGEP|<TYPE>|<uuid>`, where `accountId` is the org's hub account (`4829673` in this org). Do not hand-craft or split these IDs — treat them as opaque strings.

---

## 2. Relationship mind-map

```
                        ┌────────────────────────────────────────┐
                        │            ORGANIZATION                │
                        │  (scope.type=ORGANIZATION, id=orgUUID) │
                        └─────────────────┬──────────────────────┘
                                          │ scope
                ┌─────────────────────────┼─────────────────────────┐
                │                         │                         │
        ┌───────▼─────────┐    ┌──────────▼─────────┐   ┌───────────▼──────────┐
        │      TEAM       │    │     SCORECARD      │   │      COLLECTION      │
        │  (parentId → )  │    │                    │   │  (user-created,      │
        │        │        │    │                    │   │   heterogeneous)     │
        │        ▼        │    │                    │   └───────┬──────────────┘
        │   TEAM (parent) │    │                    │           │ contains
        └──┬──────────┬───┘    └──────┬─────────────┘           │
           │ auto     │ auto          │ auto                    ▼
           │          │               │                    any EntityManagementEntity
           ▼          ▼               ▼                    (USER, APM APP, TEAM, RULE, …)
    ┌──────────┐ ┌──────────┐   ┌──────────┐
    │COLLECTION│ │COLLECTION│   │COLLECTION│
    │membership│ │ownership │   │  rules   │
    └────┬─────┘ └──────────┘   └────┬─────┘
         │ contains                  │ contains
         ▼                           ▼
    USER, APM APP, etc.        SCORECARD_RULE   ──────► standalone lifecycle
                                                        (survives scorecard delete)
    TEAM.managers ⊂ TEAM.membership.members
```

Every parent→collection edge is exposed by:
- Forward: parent-object field (`Team.membership`, `Team.ownership`, `Scorecard.rules`) — value is the Collection sub-object.
- Reverse: `actor.entityManagement.collectionsContainingEntity(entityId: <memberId>)` returns `[{collection, parentInfo{parentEntity, parentField}}]`. `parentField` names the parent's field (`"membership"`, `"ownership"`, `"rules"`, or empty for user-created collections).

---

## 3. Auto-created (backing) Collections

Certain resources implicitly create Collections at creation time. **They are managed** — you can't delete them directly:

| Parent | Field name | Auto-collection name | Contains |
|---|---|---|---|
| `TEAM` | `membership` | `"<team-name> Collection"` | Users |
| `TEAM` | `ownership` | `"<team-name> Collection"` | Any resource the team owns |
| `SCORECARD` | `rules` | `"<scorecard-name> Collection"` | ScorecardRule entities |

`entityManagementDelete(id: <backing-collection-id>)` → `Cannot delete Collection entity with existing parent`.

Their IDs are **stable** across parent updates (verified: renamed scorecard, `rules.id` unchanged).

**Cascade on parent delete:** deleting the parent (Team or Scorecard) deletes the backing Collections. The **members** of those Collections survive (ScorecardRules become orphans; Users remain).

---

## 4. Mutation matrix

| Resource | Create | Update | Delete | Members |
|---|---|---|---|---|
| Team | `entityManagementCreateTeam` | `entityManagementUpdateTeam(id, teamEntity, version)` | `entityManagementDelete(id, version)` | via `membership` / `ownership` auto-collections + `entityManagementAddCollectionMembers` |
| TeamsHierarchyLevel | `entityManagementCreateTeamsHierarchyLevel` | `entityManagementUpdateTeamsHierarchyLevel` | `entityManagementDelete` | n/a |
| TeamsOrganizationSettings | `entityManagementCreateTeamsOrganizationSettings` | `entityManagementUpdateTeamsOrganizationSettings` | `entityManagementDelete` | n/a |
| Scorecard | `entityManagementCreateScorecard` | `entityManagementUpdateScorecard` | `entityManagementDelete` | via `rules` auto-collection |
| ScorecardRule | `entityManagementCreateScorecardRule` | `entityManagementUpdateScorecardRule` | `entityManagementDelete` | n/a (is a member) |
| Collection (user) | `entityManagementCreateCollection` | `entityManagementUpdateCollection` | `entityManagementDelete` | `entityManagementAddCollectionMembers` / `entityManagementRemoveCollectionMembers` |

**All update mutations are PATCH** — omitted fields are preserved. Explicit `null` clears the value; explicit empty list `[]` sets to empty list. Every mutation increments `metadata.version` (occasionally by 2, suggesting internal retries).

---

## 5. Field-level semantics — verified live

### 5.1 Team (`EntityManagementTeamEntityUpdateInput`)

| Field | Type | Patch behaviour | Notes |
|---|---|---|---|
| `name` | `String` | update if present | Non-null on create; `String` (nullable) on update |
| `description` | `String` | update if present | |
| `aliases` | `[String!]` | list replaces on update; `[]` clears to empty; `null` clears to null | Read-back preserves the choice — `[]` != `null` |
| `managers` | `[ID!]` | list replaces | **Every ID must be a pre-existing member** of the team's `membership` Collection. Use the **NGEP-encoded user ID** (`NDgyOTY3M3xOR0VQfFVTRVJ…`), NOT the raw int userId. |
| `resources` | `[EntityManagementTeamResourceUpdateInput!]` | list replaces; each item needs `content!`+`type!` | Not additive |
| `tags` | `[EntityManagementTagInput!]` | list replaces; `[]` and `null` both accepted and distinguishable | |
| `parentId` | `ID` | update if present; `null` unsets | Cycle-checked (`A cannot be its own parent`); invalid ID rejected |
| `externalIntegration` | `EntityManagementTeamExternalIntegrationUpdateInput` | update if present; `null` unsets | **GITHUB_TEAM** and **SERVICENOW_TEAM** cannot be set manually (`Access denied`). **IAM_GROUP** works. **Only one integration type per org.** |
| `hierarchyLevelId` | not on update input | — | Set only via hierarchy mutations (RBAC-gated) |

### 5.2 Scorecard (`EntityManagementScorecardEntityUpdateInput`)

| Field | Type | Behaviour |
|---|---|---|
| `name` | `String` | update |
| `description` | `String` | update |
| `progressLevels` | `[EntityManagementProgressLevelDefinitionUpdateInput!]` | **Currently blocked by a backend bug** — any change fails with `Unknown type 'EntityManagementScorecardRuleEntity'`. Empty list rejected: `must contain at least 1 level`. Effectively **immutable via update today**. |
| `tags` | `[EntityManagementTagInput!]` | list replaces; `[]` and `null` distinguishable |

`Scorecard.rules.id` is **stable** across updates.

### 5.3 ScorecardRule (`EntityManagementScorecardRuleEntityUpdateInput`)

| Field | Behaviour |
|---|---|
| `name`, `description` | update |
| `enabled` | update; toggle safely |
| `impactWeight` | update; `null` clears |
| `progressLevel` | must be an `id` from the parent scorecard's progressLevels; validated against **the currently-attached parent's** progressLevels (`Progress level 'x' does not exist in parent scorecard '…'`) |
| `runInterval` | Enum: **`[60, 360, 720, 1440, 4320]`** — anything else rejected |
| `nrqlEngine` | NRQL must contain an `if(...)` alias `'score'`, and must be `FACET`-ed. Validated at update time. |
| `schedule` | **`schedule.enabled is not supported for scorecard rules`** — schema exposes it but backend rejects. Treat schedule as always-null for rules. |
| `tags` | list replaces |

### 5.4 Collection (`EntityManagementCollectionEntityUpdateInput`)

Only `name` and `tags` are mutable. All member manipulation is via the dedicated `AddCollectionMembers`/`RemoveCollectionMembers` mutations.

---

## 6. Member operations semantics

`entityManagementAddCollectionMembers(collectionId, ids)` and `RemoveCollectionMembers(collectionId, ids)` both return `[ID]` — one slot per input `ids` element.

| Outcome | Slot | Companion error |
|---|---|---|
| Successful add / remove | the ID itself | — |
| Target entity doesn't exist | `null` | `"<id>: Target entity not found"` |
| Duplicate add | `null` | `"<id>: Entity already belongs to collection <collectionId>"` (`COLLECTION_DUPLICATE_MEMBER`) |
| Non-member remove | (untested — assume `null` + descriptive error) | — |

Both queries **preserve order** in the response slot alignment. A partial-success response has non-null slots for the successes and null slots + `errors[]` entries (by path index) for the failures.

Members can be heterogeneous (`USER` + `APM-APPLICATION` co-existed in the same collection). Members not represented in Go by an explicit interface impl are typed as `EntityManagementGenericEntity` in the collectionElements response.

---

## 7. Reads

- `actor.entityManagement.entity(id)` — direct read; use inline fragments per implementation.
- `actor.entityManagement.entitySearch(query, cursor)` — search. **Only accepts `type = 'X'` as a predicate** (single equality). `AND`, `LIKE`, or extra clauses return `INVALID_INPUT: Query should contain only one entityType filter`. Pagination via `cursor`+`nextCursor`.
- `actor.entityManagement.collectionElements(filter: {collectionId: {eq: $id}}, cursor, first)` → `{items: [EntityManagementEntity!], nextCursor}`.
- `actor.entityManagement.collectionsContainingEntity(entityId)` → `[{collection, parentInfo{parentEntity, parentField}}]`. **Best tool for walking the graph in reverse.**

---

## 8. Scope semantics

Enum `EntityManagementEntityScope`: `ACCOUNT`, `ORGANIZATION`.

| Resource | Accepted scope in this environment |
|---|---|
| Team | `ORGANIZATION` only (`ACCOUNT` → `Scope not supported: ACCOUNT`) |
| Scorecard | `ORGANIZATION` |
| ScorecardRule | `ORGANIZATION` |
| Collection | `ORGANIZATION` succeeded; `ACCOUNT` returned `ACCESS_DENIED` (per-resource RBAC) |
| TeamsHierarchyLevel, TeamsOrganizationSettings | `ORGANIZATION` |

---

## 9. RBAC-gated behaviour observed

- Hierarchy mutations (`Create/UpdateTeamsHierarchyLevel`, `Create/UpdateTeamsOrganizationSettings`) return `Cannot query field "…" on type "RootMutationType"` from a personal-API-key user — the field is **visible via introspection but not exposed at parse time** for insufficient RBAC. Requires a workspace-level admin permission.
- Externals: `GITHUB_TEAM`, `SERVICENOW_TEAM` externalIntegrations reject with `Access denied. Integration with … cannot be setup manually`. Only `IAM_GROUP` is user-settable.
- Only **one** integration type per org (`Only one integration type is allowed per organization, and IAM_GROUP already exists`).

---

## 10. Backend flakes worth handling in downstream (Terraform)

- **Ghost NOT_FOUND on freshly-created entities**: shortly after create, `entity(id:)` / update / delete / add-members may return `: Entity not found.` (leading colon, empty prefix). `entitySearch` still returns the entity. This looks like a read-through-cache inconsistency in the NGEP resolver. Distinguishable from **real** NOT_FOUND by the leading colon + empty prefix in the error `message` (real deletions carry the ID in the prefix: `<id>: Entity not found.`).
- **Progress levels update backend bug**: `Unknown type 'EntityManagementScorecardRuleEntity'` during the internal deletion-validation query. Treat progressLevels as immutable via `updateScorecard` until backend is patched.
- **Rate-limit-esque `Access restricted` on rapid Collection creates** after several creations in a session.
- **Version numbers occasionally skip** (create → v1, update → v2, another update → v5). Not deterministic; do not rely on strict monotonic +1.

Recommended handling at the Terraform/provider layer:
- Retry `NOT_FOUND` with empty-prefix message on any post-create operation, capped at ~30s.
- On real NOT_FOUND (id-prefixed message), do not retry — surface as-is.
- For scorecards, warn users up-front that `progressLevels` requires recreate.

---

## 11. Consolidated package design decision

All three resource families (Teams, Scorecards, Collections) share the same NGEP interface and the same `entityManagementDelete` / `entity` / `entitySearch` endpoints, and they interoperate (Teams have backing Collections; Scorecards have a backing Collection of Rules; Collections are the join layer). Duplicating types across three packages would fight both Tutone's interface-generation and NGEP's own type hierarchy.

Decision: single package `pkg/scorecards` that covers `Team`, `TeamsHierarchyLevel`, `TeamsOrganizationSettings`, `Scorecard`, `ScorecardRule`, `Collection`. The name is a placeholder — can be renamed to `pkg/entitymanagement` or split later if the taxonomy needs it.

`include_implementations` in `.tutone.yml` for this package covers all seven (six above + `EntityManagementUserEntity` so member reads type correctly) plus a graceful fallback for `EntityManagementGenericEntity`.
