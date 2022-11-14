//go:build integration
// +build integration

package logconfigurations

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationObfuscationExpression(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	var (
		rand            = mock.RandSeq(5)
		testName        = "testName_" + rand
		testDescription = "testDescription_" + rand
		testCreateInput = LogConfigurationsCreateObfuscationExpressionInput{

			Description: testDescription,
			Name:        testName,
			Regex:       "(^http.*)",
		}
	)

	client := newIntegrationTestClient(t)

	// Test: Create
	created, err := client.LogConfigurationsCreateObfuscationExpression(testAccountID, testCreateInput)

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created)

	// Test: Delete
	testDeleteInput := created.ID
	deleted, err := client.LogConfigurationsDeleteObfuscationExpression(testAccountID, testDeleteInput)

	require.NoError(t, err)
	require.NotNil(t, deleted)
	require.NotEmpty(t, deleted)
}

func newIntegrationTestClient(t *testing.T) Logconfigurations {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
