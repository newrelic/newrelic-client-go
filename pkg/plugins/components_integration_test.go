// +build integration

package plugins

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationComponents(t *testing.T) {
	t.Parallel()

	tc := mock.NewIntegrationTestConfig(t)

	api := New(tc)
	a, err := api.ListComponents(nil)

	require.NoError(t, err)
	require.NotNil(t, a)

	if len(a) < 1 {
		t.Skipf("Skipping `GetComponent` integration test due to zero plugins found")
	}

	c, err := api.GetComponent(a[0].ID)

	require.NoError(t, err)
	require.NotNil(t, c)

	m, err := api.ListComponentMetrics(c.ID, nil)

	require.NoError(t, err)
	require.NotNil(t, m)

	if len(m) < 1 {
		t.Skipf("Skipping `GetComponentMetricData` integration test due to zero plugin metrics found")
	}
	params := GetComponentMetricDataParams{
		Names: []string{m[0].Name},
	}
	_, err = api.GetComponentMetricData(a[0].ID, &params)

	require.NoError(t, err)
}
