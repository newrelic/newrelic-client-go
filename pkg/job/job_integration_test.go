//go:build integration
// +build integration

package job

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetJobs(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	resp, err := client.JobGetOrganizationCreateResults(
		jobId,
	)
	require.NoError(t, err)
	require.Equal(t, 0, len(resp.Items))
}
