//go:build integration
// +build integration

package nrdb

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

// nolint
func newNrdbIntegrationTestClient(t *testing.T) Nrdb {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}

func TestIntegrationNrdbQuery(t *testing.T) {
	t.Parallel()

	// Request a constant so we can easily validate
	query := "SELECT 1 FROM Transaction"
	client := newNrdbIntegrationTestClient(t)

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		t.Skipf("integration testing requires NEW_RELIC_ACOUNT_ID")
	}

	res, err := client.Query(accountID, NRQL(query))

	require.NoError(t, err)
	require.NotNil(t, res)

	fmt.Printf("%+v\n", res)
	require.Equal(t, 1, len(res.Results))

	if v, ok := res.Results[0]["constant"]; ok {
		assert.Equal(t, float64(1), v)
	}
}

func TestIntegrationNrdbQueryHistoryQuery(t *testing.T) {
	t.Parallel()

	client := newNrdbIntegrationTestClient(t)

	res, err := client.QueryHistory()

	require.NoError(t, err)
	require.NotNil(t, res)

	require.GreaterOrEqual(t, len(*res), 1)
}
