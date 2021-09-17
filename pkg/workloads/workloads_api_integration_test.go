//go:build integration
// +build integration

package workloads

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/common"
	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationWorkloadCreate(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	name := "newrelic-client-go-test-workload-" + mock.RandSeq(5)
	workload := WorkloadCreateInput{
		Name: name,
		ScopeAccounts: &WorkloadScopeAccountsInput{
			AccountIDs: []int{testAccountID},
		},
		EntitySearchQueries: []WorkloadEntitySearchQueryInput{
			{
				Query: "(name like 'tf_test' or id = 'tf_test' or domainId = 'tf_test')",
			},
		},
	}

	created, err := client.WorkloadCreate(testAccountID, workload)
	require.NoError(t, err)
	require.NotNil(t, created)

	assert.Equal(t, name, created.Name)
	assert.NotEmpty(t, created.GUID)
	assert.Equal(t, testAccountID, created.Account.ID)
	assert.Equal(t, 1, len(created.EntitySearchQueries))

	// Wait for indexing to catch up
	time.Sleep(30 * time.Second)

	// Cleanup
	_, err = client.WorkloadDelete(common.EntityGUID(created.GUID))
	require.NoError(t, err)
}

func TestIntegrationWorkloadDuplicate(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	name := "newrelic-client-go-test-workload-" + mock.RandSeq(5)
	workload := WorkloadCreateInput{
		Name: name,
		ScopeAccounts: &WorkloadScopeAccountsInput{
			AccountIDs: []int{testAccountID},
		},
		EntitySearchQueries: []WorkloadEntitySearchQueryInput{
			{
				Query: "(name like 'tf_test' or id = 'tf_test' or domainId = 'tf_test')",
			},
		},
	}

	created, err := client.WorkloadCreate(testAccountID, workload)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.GUID)

	// Wait for indexing to catch up
	time.Sleep(10 * time.Second)

	dup, err := client.WorkloadDuplicate(testAccountID, created.GUID, WorkloadDuplicateInput{
		Name: name + "-duplicate",
	})
	require.NoError(t, err)
	require.NotNil(t, dup)
	require.NotEmpty(t, dup.GUID)

	// Wait for indexing to catch up
	time.Sleep(30 * time.Second)

	// Cleanup
	_, err = client.WorkloadDelete(common.EntityGUID(created.GUID))
	assert.NoError(t, err)
	_, err = client.WorkloadDelete(common.EntityGUID(dup.GUID))
	assert.NoError(t, err)
}

func TestIntegrationWorkloadUpdate(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	client := newIntegrationTestClient(t)

	// Test vars
	name := "newrelic-client-go-test-workload-" + mock.RandSeq(5)
	workload := WorkloadCreateInput{
		Name: name,
		ScopeAccounts: &WorkloadScopeAccountsInput{
			AccountIDs: []int{testAccountID},
		},
		EntitySearchQueries: []WorkloadEntitySearchQueryInput{
			{
				Query: "(name like 'tf_test' or id = 'tf_test' or domainId = 'tf_test')",
			},
		},
	}

	created, err := client.WorkloadCreate(testAccountID, workload)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.GUID)

	// Wait for indexing to catch up
	time.Sleep(30 * time.Second)

	up, err := client.WorkloadUpdate(created.GUID, WorkloadUpdateInput{
		Name: name + "-updated",
	})
	require.NoError(t, err)
	require.NotNil(t, up)
	require.NotEmpty(t, up.GUID)

	// Cleanup
	_, err = client.WorkloadDelete(common.EntityGUID(up.GUID))
	assert.NoError(t, err)
}
