//go:build integration
// +build integration

package nrqlmacros

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationNrqlmacros_NrqlMacro_CRUD(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	orgID, err := testhelpers.GetTestOrganizationID()
	if err != nil || orgID == "" {
		t.Skipf("skipping: no organization ID available (%v)", err)
	}

	macroName := fmt.Sprintf("nr-test-macro-%s", testhelpers.RandSeq(6))
	createInput := EntityManagementNRQLMacroEntityCreateInput{
		Name:       macroName,
		Expression: "environment = 'production' AND appName NOT LIKE '%test%'",
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ORGANIZATION,
			ID:   orgID,
		},
		Tags: []EntityManagementTagInput{
			{Key: "purpose", Values: []string{"integration-test"}},
		},
	}

	// 1. Create
	createResult, err := client.EntityManagementCreateNRQLMacro(createInput)
	require.NoError(t, err)
	require.NotNil(t, createResult)
	require.NotEmpty(t, createResult.Entity.ID)
	require.Equal(t, macroName, createResult.Entity.Name)
	require.Equal(t, createInput.Expression, createResult.Entity.Expression)

	macroID := createResult.Entity.ID
	defer func() {
		_, deleteErr := client.EntityManagementDelete(macroID)
		if deleteErr != nil {
			t.Logf("warning: failed to delete macro %s: %v", macroID, deleteErr)
		}
	}()

	// 2. Read via GetEntity
	getResult, err := client.GetEntity(macroID)
	require.NoError(t, err)
	require.NotNil(t, getResult)

	macroEntity, ok := (*getResult).(*EntityManagementNRQLMacroEntity)
	require.True(t, ok, "expected *EntityManagementNRQLMacroEntity, got %T", *getResult)
	require.Equal(t, macroID, macroEntity.ID)
	require.Equal(t, macroName, macroEntity.Name)
	require.Equal(t, createInput.Expression, macroEntity.Expression)

	// 3. Read via GetEntitySearch
	searchResult, err := client.GetEntitySearch("", "type = 'NRQL_MACRO'")
	require.NoError(t, err)
	require.NotNil(t, searchResult)

	foundInSearch := false
	for _, e := range searchResult.Entities {
		if me, ok := e.(*EntityManagementNRQLMacroEntity); ok {
			if me.ID == macroID {
				foundInSearch = true
				break
			}
		}
	}
	require.True(t, foundInSearch, "created macro not found in entity search results")

	// 4. Update expression
	updatedExpression := "environment = 'staging' AND appName NOT LIKE '%test%'"
	updateInput := EntityManagementNRQLMacroEntityUpdateInput{
		Expression: updatedExpression,
	}
	updateResult, err := client.EntityManagementUpdateNRQLMacro(macroID, updateInput)
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, updatedExpression, updateResult.Entity.Expression)
	require.Equal(t, macroName, updateResult.Entity.Name, "name must not change on update")
	require.Equal(t, createResult.Entity.Metadata.Version+1, updateResult.Entity.Metadata.Version)

	// 5. Read back to verify update
	getAfterUpdate, err := client.GetEntity(macroID)
	require.NoError(t, err)
	updatedEntity, ok := (*getAfterUpdate).(*EntityManagementNRQLMacroEntity)
	require.True(t, ok)
	require.Equal(t, updatedExpression, updatedEntity.Expression)

	// 6. Delete is handled by the deferred function above
}

func TestIntegrationNrqlmacros_NrqlMacro_AccountScopeFails(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	accountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skipf("skipping: no account ID available (%v)", err)
	}

	createInput := EntityManagementNRQLMacroEntityCreateInput{
		Name:       fmt.Sprintf("nr-test-macro-acct-%s", testhelpers.RandSeq(6)),
		Expression: "host = 'test-host'",
		Scope: EntityManagementScopedReferenceInput{
			Type: EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   fmt.Sprintf("%d", accountID),
		},
	}

	// ACCOUNT scope is not supported for NRQL macros
	_, err = client.EntityManagementCreateNRQLMacro(createInput)
	require.Error(t, err, "expected error when creating macro with ACCOUNT scope")
}
