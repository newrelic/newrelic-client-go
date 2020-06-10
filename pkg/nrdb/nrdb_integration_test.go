// +build integration

package nrdb

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationNrdbQuery(t *testing.T) {
	t.Parallel()

	// Request a constant so we can easily validate
	query := "SELECT 1 FROM Transaction"
	client := newNrdbIntegrationTestClient(t)

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		t.Skipf("integration testing requires NEW_RELIC_ACOUNT_ID")
	}

	actual, err := client.Query(accountID, Nrql(query))

	require.NoError(t, err)
	require.NotNil(t, actual)

	require.Equal(t, 1, len(actual.Results))

	if v, ok := actual.Results[0]["constant"]; ok {
		assert.Equal(t, float64(1), v)
	}
}

// nolint
func newNrdbIntegrationTestClient(t *testing.T) Nrdb {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
