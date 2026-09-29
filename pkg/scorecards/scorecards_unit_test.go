//go:build unit
// +build unit

package scorecards

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Encoded NGEP IDs used across the mock fixtures. They are opaque bytes for
// the client, so anything base64-shaped is fine.
const (
	testOrgID           = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	testTeamID          = "MXxOR0VQfFRFQU18ZmFrZS10ZWFtLWlkLTAwMDAwMDAwMDAwMTIz"
	testTeamMembershipC = "MXxOR0VQfENPTExFQ1RJT058ZmFrZS10ZWFtLW1lbWJlcnNoaXAxMjM"
	testScorecardID     = "MXxOR0VQfFNDT1JFQ0FSRHxmYWtlLXNjLTAwMDAwMDAwMDAwMTIz"
	testScorecardRulesC = "MXxOR0VQfENPTExFQ1RJT058ZmFrZS1zYy1ydWxlczEyMw"
	testRuleID          = "MXxOR0VQfFNDT1JFQ0FSRF9SVUxFfGZha2UtcnVsZS0wMDAwMDAwMDAwMDEyMw"
	testCollectionID    = "MXxOR0VQfENPTExFQ1RJT058ZmFrZS1jb2xsZWN0aW9uLTAwMDAxMjM"
	testUserNGEPID      = "MXxOR0VQfFVTRVJ8ZmFrZS11c2VyLTAwMDAwMDAwMDAwMTIz"
	testApmID           = "MXxBUE18QVBQTElDQVRJT058ZmFrZS1hcHAtMDAwMDEyMw"
)

var (
	testTeamCreateResp = `{
  "data": {
    "entityManagementCreateTeam": {
      "entity": {
        "aliases": ["alpha"],
        "description": "A test team",
        "id": "` + testTeamID + `",
        "managers": null,
        "membership": {"id": "` + testTeamMembershipC + `", "name": "Team A Collection", "type": "COLLECTION"},
        "metadata": {"createdAt": 1783677600000, "updatedAt": 1783677600000, "version": 1},
        "name": "team-a",
        "ownership": {"id": "` + testTeamMembershipC + `_o", "name": "Team A Collection", "type": "COLLECTION"},
        "resources": null,
        "scope": {"id": "` + testOrgID + `", "type": "ORGANIZATION"},
        "tags": [{"key": "purpose", "values": ["testing"]}],
        "type": "TEAM"
      }
    }
  }
}`

	testTeamUpdateResp = `{
  "data": {
    "entityManagementUpdateTeam": {
      "entity": {
        "aliases": ["alpha", "beta"],
        "description": "A test team renamed",
        "id": "` + testTeamID + `",
        "managers": null,
        "membership": {"id": "` + testTeamMembershipC + `", "name": "Team A Collection", "type": "COLLECTION"},
        "metadata": {"createdAt": 1783677600000, "updatedAt": 1783677600000, "version": 2},
        "name": "team-a-renamed",
        "ownership": {"id": "` + testTeamMembershipC + `_o", "name": "Team A Collection", "type": "COLLECTION"},
        "resources": null,
        "scope": {"id": "` + testOrgID + `", "type": "ORGANIZATION"},
        "tags": null,
        "type": "TEAM"
      }
    }
  }
}`

	testTeamGetResp = `{
  "data": {
    "actor": {
      "entityManagement": {
        "entity": {
          "__typename": "EntityManagementTeamEntity",
          "aliases": ["alpha"],
          "description": "A test team",
          "id": "` + testTeamID + `",
          "membership": {"id": "` + testTeamMembershipC + `", "name": "Team A Collection", "type": "COLLECTION"},
          "metadata": {"createdAt": 1783677600000, "updatedAt": 1783677600000, "version": 1},
          "name": "team-a",
          "ownership": {"id": "` + testTeamMembershipC + `_o", "name": "Team A Collection", "type": "COLLECTION"},
          "scope": {"id": "` + testOrgID + `", "type": "ORGANIZATION"},
          "tags": [{"key": "purpose", "values": ["testing"]}],
          "type": "TEAM"
        }
      }
    }
  }
}`

	testScorecardCreateResp = `{
  "data": {
    "entityManagementCreateScorecard": {
      "entity": {
        "description": "A scorecard",
        "id": "` + testScorecardID + `",
        "metadata": {"createdAt": 1783677600000, "updatedAt": 1783677600000, "version": 1},
        "name": "sc-a",
        "progressLevels": [
          {"id": "red", "name": "Red", "description": "Needs work", "hexColorCode": "#FF0000"},
          {"id": "green", "name": "Green", "description": "Good", "hexColorCode": "#00CC00"}
        ],
        "rules": {"id": "` + testScorecardRulesC + `", "name": "sc-a Collection", "type": "COLLECTION"},
        "scope": {"id": "` + testOrgID + `", "type": "ORGANIZATION"},
        "tags": null,
        "type": "SCORECARD"
      }
    }
  }
}`

	testScorecardRuleCreateResp = `{
  "data": {
    "entityManagementCreateScorecardRule": {
      "entity": {
        "description": null,
        "enabled": true,
        "id": "` + testRuleID + `",
        "impactWeight": null,
        "metadata": {"createdAt": 1783677600000, "updatedAt": 1783677600000, "version": 1},
        "name": "rule-a",
        "nrqlEngine": {"accounts": [12345], "joinAccounts": null, "query": "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity FACET id LIMIT MAX SINCE 1 day ago"},
        "progressLevel": null,
        "runInterval": null,
        "schedule": null,
        "scope": {"id": "` + testOrgID + `", "type": "ORGANIZATION"},
        "tags": null,
        "type": "SCORECARD_RULE"
      }
    }
  }
}`

	testCollectionCreateResp = `{
  "data": {
    "entityManagementCreateCollection": {
      "entity": {
        "id": "` + testCollectionID + `",
        "metadata": {"createdAt": 1783677600000, "updatedAt": 1783677600000, "version": 1},
        "name": "coll-a",
        "scope": {"id": "` + testOrgID + `", "type": "ORGANIZATION"},
        "tags": null,
        "type": "COLLECTION"
      }
    }
  }
}`

	testAddMembersResp = `{
  "data": {
    "entityManagementAddCollectionMembers": ["` + testUserNGEPID + `", "` + testApmID + `"]
  }
}`

	testAddMembersPartialResp = `{
  "data": {
    "entityManagementAddCollectionMembers": ["` + testUserNGEPID + `", null]
  },
  "errors": [
    {"message": "` + testApmID + `: Entity already belongs to collection ` + testCollectionID + `", "path": ["entityManagementAddCollectionMembers", 1]}
  ]
}`

	testRemoveMembersResp = `{
  "data": {
    "entityManagementRemoveCollectionMembers": ["` + testUserNGEPID + `"]
  }
}`

	testCollectionElementsResp = `{
  "data": {
    "actor": {
      "entityManagement": {
        "collectionElements": {
          "items": [
            {
              "__typename": "EntityManagementUserEntity",
              "id": "` + testUserNGEPID + `",
              "name": "user1",
              "type": "USER"
            },
            {
              "__typename": "EntityManagementGenericEntity",
              "id": "` + testApmID + `",
              "name": "Dummy App",
              "type": "APM-APPLICATION"
            }
          ],
          "nextCursor": null
        }
      }
    }
  }
}`

	testCollectionsContainingResp = `{
  "data": {
    "actor": {
      "entityManagement": {
        "collectionsContainingEntity": [
          {
            "collection": {"id": "` + testTeamMembershipC + `", "name": "team-a Collection", "type": "COLLECTION"},
            "parentInfo": {
              "parentEntity": {
                "__typename": "EntityManagementTeamEntity",
                "id": "` + testTeamID + `",
                "name": "team-a",
                "type": "TEAM"
              },
              "parentField": "membership"
            }
          }
        ]
      }
    }
  }
}`

	testDeleteResp = `{
  "data": {
    "entityManagementDelete": {
      "id": "` + testTeamID + `"
    }
  }
}`
)

// ---- Teams ----

func TestUnitEntityManagement_CreateTeam(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testTeamCreateResp, http.StatusOK)

	in := EntityManagementTeamEntityCreateInput{
		Name:        "team-a",
		Description: "A test team",
		Aliases:     []string{"alpha"},
		Scope:       EntityManagementScopedReferenceInput{ID: testOrgID, Type: EntityManagementEntityScopeTypes.ORGANIZATION},
	}

	res, err := c.EntityManagementCreateTeam(in)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, testTeamID, res.Entity.ID)
	require.Equal(t, "team-a", res.Entity.Name)
	require.Equal(t, testTeamMembershipC, res.Entity.Membership.ID)
	require.Equal(t, "COLLECTION", res.Entity.Membership.Type)
	require.Equal(t, 1, res.Entity.Metadata.Version)
}

func TestUnitEntityManagement_UpdateTeam(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testTeamUpdateResp, http.StatusOK)

	in := EntityManagementTeamEntityUpdateInput{
		Name:        "team-a-renamed",
		Description: "A test team renamed",
		Aliases:     []string{"alpha", "beta"},
	}
	res, err := c.EntityManagementUpdateTeam(testTeamID, in)
	require.NoError(t, err)
	require.Equal(t, "team-a-renamed", res.Entity.Name)
	require.ElementsMatch(t, []string{"alpha", "beta"}, res.Entity.Aliases)
	require.Equal(t, 2, res.Entity.Metadata.Version)
}

func TestUnitEntityManagement_GetTeamEntity(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testTeamGetResp, http.StatusOK)

	got, err := c.GetEntity(testTeamID)
	require.NoError(t, err)
	require.NotNil(t, got)

	team, ok := (*got).(*EntityManagementTeamEntity)
	require.True(t, ok, "expected *EntityManagementTeamEntity, got %T", *got)
	require.Equal(t, testTeamID, team.ID)
	require.Equal(t, "team-a", team.Name)
	require.Equal(t, testTeamMembershipC, team.Membership.ID)
}

// ---- Scorecards ----

func TestUnitEntityManagement_CreateScorecard(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testScorecardCreateResp, http.StatusOK)

	in := EntityManagementScorecardEntityCreateInput{
		Name:        "sc-a",
		Description: "A scorecard",
		Scope:       EntityManagementScopedReferenceInput{ID: testOrgID, Type: EntityManagementEntityScopeTypes.ORGANIZATION},
		ProgressLevels: []EntityManagementProgressLevelDefinitionCreateInput{
			{ID: "red", Name: "Red", Description: "Needs work", HexColorCode: "#FF0000"},
			{ID: "green", Name: "Green", Description: "Good", HexColorCode: "#00CC00"},
		},
	}
	res, err := c.EntityManagementCreateScorecard(in)
	require.NoError(t, err)
	require.Equal(t, testScorecardID, res.Entity.ID)
	require.Equal(t, testScorecardRulesC, res.Entity.Rules.ID, "auto-created rules collection expected on Scorecard")
	require.Len(t, res.Entity.ProgressLevels, 2)
}

func TestUnitEntityManagement_CreateScorecardRule(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testScorecardRuleCreateResp, http.StatusOK)

	in := EntityManagementScorecardRuleEntityCreateInput{
		Name:    "rule-a",
		Enabled: true,
		Scope:   EntityManagementScopedReferenceInput{ID: testOrgID, Type: EntityManagementEntityScopeTypes.ORGANIZATION},
		NRQLEngine: &EntityManagementNRQLRuleEngineCreateInput{
			Accounts: []int{12345},
			Query:    "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity FACET id LIMIT MAX SINCE 1 day ago",
		},
	}
	res, err := c.EntityManagementCreateScorecardRule(in)
	require.NoError(t, err)
	require.Equal(t, testRuleID, res.Entity.ID)
	require.True(t, res.Entity.Enabled)
	require.Equal(t, "SCORECARD_RULE", res.Entity.Type)
}

// ---- Collections + member ops ----

func TestUnitEntityManagement_CreateCollection(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testCollectionCreateResp, http.StatusOK)

	in := EntityManagementCollectionEntityCreateInput{
		Name:  "coll-a",
		Scope: EntityManagementScopedReferenceInput{ID: testOrgID, Type: EntityManagementEntityScopeTypes.ORGANIZATION},
	}
	res, err := c.EntityManagementCreateCollection(in)
	require.NoError(t, err)
	require.Equal(t, testCollectionID, res.Entity.ID)
	require.Equal(t, "COLLECTION", res.Entity.Type)
}

func TestUnitEntityManagement_AddCollectionMembers_All(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testAddMembersResp, http.StatusOK)

	res, err := c.EntityManagementAddCollectionMembers(testCollectionID, []string{testUserNGEPID, testApmID})
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.NotNil(t, res[0])
	require.Equal(t, testUserNGEPID, *res[0])
	require.NotNil(t, res[1])
	require.Equal(t, testApmID, *res[1])
}

func TestUnitEntityManagement_AddCollectionMembers_Partial(t *testing.T) {
	t.Parallel()
	// Server responds 200 with a partial success — data slot 1 is null, errors[] describes it.
	// The NerdGraph client currently surfaces GraphQL errors as a Go error; we
	// only care that this response shape doesn't panic and that the slot
	// alignment is preserved when errors also come through the transport.
	c := newMockClient(t, testAddMembersPartialResp, http.StatusOK)

	res, err := c.EntityManagementAddCollectionMembers(testCollectionID, []string{testUserNGEPID, testApmID})
	if err == nil {
		// If the client happens to tolerate partial errors, slot alignment must still hold.
		require.Len(t, res, 2)
		require.NotNil(t, res[0])
		require.Equal(t, testUserNGEPID, *res[0])
		require.Nil(t, res[1])
	} else {
		require.Contains(t, err.Error(), "Entity already belongs to collection")
	}
}

func TestUnitEntityManagement_RemoveCollectionMembers(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testRemoveMembersResp, http.StatusOK)

	res, err := c.EntityManagementRemoveCollectionMembers(testCollectionID, []string{testUserNGEPID})
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.NotNil(t, res[0])
	require.Equal(t, testUserNGEPID, *res[0])
}

func TestUnitEntityManagement_GetCollectionElements(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testCollectionElementsResp, http.StatusOK)

	res, err := c.GetCollectionElements(
		"",
		EntityManagementCollectionElementsFilter{CollectionID: EntityManagementCollectionIdFilterArgument{Eq: testTeamMembershipC}},
		25,
	)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Items, 2)

	// The first item should decode as a UserEntity via the interface
	u, ok := res.Items[0].(*EntityManagementUserEntity)
	require.True(t, ok, "expected *EntityManagementUserEntity, got %T", res.Items[0])
	require.Equal(t, testUserNGEPID, u.ID)

	// The second should fall back to GenericEntity
	g, ok := res.Items[1].(*EntityManagementGenericEntity)
	require.True(t, ok, "expected *EntityManagementGenericEntity, got %T", res.Items[1])
	require.Equal(t, testApmID, g.ID)
}

func TestUnitEntityManagement_GetCollectionsContainingEntity(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testCollectionsContainingResp, http.StatusOK)

	res, err := c.GetCollectionsContainingEntity(testUserNGEPID)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, *res, 1)
	require.Equal(t, testTeamMembershipC, (*res)[0].Collection.ID)
	require.Equal(t, "membership", (*res)[0].ParentInfo.ParentField)

	team, ok := (*res)[0].ParentInfo.ParentEntity.(*EntityManagementTeamEntity)
	require.True(t, ok, "expected parent to be *EntityManagementTeamEntity, got %T", (*res)[0].ParentInfo.ParentEntity)
	require.Equal(t, testTeamID, team.ID)
	require.Equal(t, "team-a", team.Name)
}

// ---- Delete (shared) ----

func TestUnitEntityManagement_Delete(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testDeleteResp, http.StatusOK)

	res, err := c.EntityManagementDelete(testTeamID)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, testTeamID, res.ID)
}

func TestUnitEntityManagement_DeleteWithContext(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, testDeleteResp, http.StatusOK)

	res, err := c.EntityManagementDeleteWithContext(context.Background(), testTeamID)
	require.NoError(t, err)
	require.Equal(t, testTeamID, res.ID)
}

func TestUnitEntityManagement_GetEntity_Error(t *testing.T) {
	t.Parallel()
	c := newMockClient(t, `{"errors":[{"message":"Not Found"}]}`, http.StatusNotFound)

	res, err := c.GetEntity("nonexistent")
	require.Error(t, err)
	require.Nil(t, res)
}
