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

func TestSyntheticsMonitors_Basic(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}
	location := SyntheticsLocationsInput{
		Public: "AP_SOUTH_1",
	}
	tags := SyntheticsTag{
		Key:    "Name",
		Values: "testSimpleBrowserMonitor",
	}

	//Create synthetic script api Monitor
	simpleBrowserMonitor := SyntheticsCreateSimpleBrowserMonitorInput{
		Locations: location,
		Name:      "testSimpleMonitor",
		Period:    "EVERY_5_MINUTES",
		Status:    "ENABLED",
		Tags:      tags,
		Uri:       "https://www.one.newrelic.com",
	}
	a := newIntegrationTestClient(t)

	//create a simple browser monitor
	createSimpleBrowserMonitor, err := a.SyntheticsCreateSimpleBrowserMonitor(testAccountID, simpleBrowserMonitor)

	require.NoError(t, err)
	require.NotNil(t, createSimpleBrowserMonitor)

	simpleBrowserMonitor := SyntheticsCreateSimpleBrowserMonitorInput{
		Locations: location,
		Name:      "testSimpleMonitor-updated",
		Period:    "EVERY_5_MINUTES",
		Status:    "ENABLED",
		Tags:      tags,
		Uri:       "https://www.one.newrelic.com",
	}
	//update a simple browser monitor
	updateSimpleBrowerMonitor, err := a.SyntheticsUpdateSimpleBrowserMonitor(createSimpleBrowserMonitor.Monitor.GUID, simpleBrowserMonitor)

	require.NoError(t, err)
	require.NotNil(t, updateSimpleBrowerMonitor)

	//delete a simple browser monitor
	deleteSimpleBrowserMonitor, err := a.SyntheticsDeleteMonitor(createSimpleBrowserMonitor.Monitor.GUID)

	require.NotNil(t, deleteSimpleBrowserMonitor)

	simpleMonitor := SyntheticsCreateSimpleMonitorInput{
		Locations: location,
		Name:      "testSimpleMonitor",
		Period:    "EVERY_5_MINUTES",
		Status:    "ENABLED",
		Tags:      tags,
		Uri:       "https://www.one.newrelic.com",
	}

	//create a simple monitor
	createSimpleMonitor, err := a.SyntheticsCreateSimpleMonitor(testAccountID, simpleMonitor)

	require.NoError(t, err)
	require.NotNil(t, createSimpleMonitor)

	simpleMonitor := SyntheticsCreateSimpleMonitorInput{
		Locations: location,
		Name:      "testSimpleMonitor-updated",
		Period:    "EVERY_5_MINUTES",
		Status:    "ENABLED",
		Tags:      tags,
		Uri:       "https://www.one.newrelic.com",
	}
	//update a simple monitor
	updateSimpleMonitor, err := a.SyntheticsUpdateSimpleMonitor(createSimpleMonitor.Monitor.GUID, simpleBrowserMonitor)

	require.NoError(t, err)
	require.NotNil(t, updateSimpleMonitor)

	//delete a simple monitor
	deleteSimpleMonitor, err := a.SyntheticsDeleteMonitor(createSimpleMonitor)

	require.Nil(t, deleteSimpleMonitor)
}

func newIntegrationTestClient(t *testing.T) Synthetics {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
