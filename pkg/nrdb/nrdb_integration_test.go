// +build integration

package nrdb

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationNrdbQuery(t *testing.T) {
	t.Parallel()

	var query Nrql
	query = "SELECT count(*) FROM Transaction"

	client := newNrdbIntegrationTestClient(t)

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		t.Skipf("integration testing requires NEW_RELIC_ACOUNT_ID")
	}

	actual, err := client.Query(accountID, query)

	require.NoError(t, err)
	require.NotNil(t, actual)
}

// nolint
func newNrdbIntegrationTestClient(t *testing.T) Nrdb {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
