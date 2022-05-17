//go:build integration
// +build integration

package synthetics

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

func TestSyntheticsSecureCredential_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	// Create a secure credential
	createResp, err := a.SyntheticsCreateSecureCredential(testAccountID, "test secure credential", "TEST", "secure value")

	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Update secure credential
	updateResp, err := a.SyntheticsUpdateSecureCredential(testAccountID, "test secure credential", "TEST", "new secure value")

	require.NoError(t, err)
	require.NotNil(t, updateResp)

	// Delete secure credential
	deleteResp, err := a.SyntheticsDeleteSecureCredential(testAccountID, "TEST")

	require.Nil(t, deleteResp)
}

// Integration testing for private location

func TestSyntheticsPrivateLocation_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	// Create a private location
	createResp, err := a.SyntheticsCreatePrivateLocation(testAccountID, "test secure credential", "TEST", true)

	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Update private location
	updateResp, err := a.SyntheticsUpdatePrivateLocation("test secure credential", createResp.GUID, true)

	require.NoError(t, err)
	require.NotNil(t, updateResp)

	// Delete private location
	deleteResp, err := a.SyntheticsDeletePrivateLocation(createResp.GUID)

	require.Nil(t, deleteResp)

	// Purge private location queue
	purgeresp, err := a.SyntheticsPurgePrivateLocationQueue(createResp.GUID)

	require.NotNil(t, purgeresp)
}

func newIntegrationTestClient(t *testing.T) Synthetics {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
