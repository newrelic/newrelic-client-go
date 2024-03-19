//go:build integration
// +build integration

package nrdb

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationNRDBQuery(t *testing.T) {
	t.Parallel()

	// Request a constant so we can easily validate
	query := "SELECT 1 FROM Transaction"
	client := newNRDBIntegrationTestClient(t)

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		t.Skipf("integration testing requires NEW_RELIC_ACOUNT_ID")
	}

	res, err := client.Query(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, 1, len(res.Results))

	if v, ok := res.Results[0]["constant"]; ok {
		assert.Equal(t, float64(1), v)
	}
}

func TestIntegrationNRDBQueryWithAdditionalOptions(t *testing.T) {
	t.Parallel()

	// Request a constant so we can easily validate
	query := "SELECT 1 from Transaction"
	client := newNRDBIntegrationTestClient(t)

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		t.Skipf("integration testing requires NEW_RELIC_ACOUNT_ID")
	}

	res, err := client.QueryWithAdditionalOptions(
		accountID,
		NRQL(query),
		30,
		false,
	)

	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, 1, len(res.Results))

	if v, ok := res.Results[0]["constant"]; ok {
		assert.Equal(t, float64(1), v)
	}
}

func TestIntegrationNRDBQueryHistoryQuery(t *testing.T) {
	t.Parallel()

	client := newNRDBIntegrationTestClient(t)

	res, err := client.QueryHistory()

	require.NoError(t, err)
	require.NotNil(t, res)

	require.GreaterOrEqual(t, len(*res), 1)
}
