//go:build integration
// +build integration

package synthetics

import (
	"github.com/stretchr/testify/require"
	"testing"

	"fmt"
	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
	"os"
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

//Test simple browser monitor
func TestSyntheticsSimpleBrowserMonitor_Basic(t *testing.T) {

	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	////Simple Browser monitor
	//Input for simple browser monitor
	simpleBrowserMonitorInput := SyntheticsCreateSimpleBrowserMonitorInput{
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testSimpleBrowserMonitor",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "Name",
				Values: []string{
					"synthetics",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
		Runtime: SyntheticsRuntimeInput{
			RuntimeType:        "CHROME_BROWSER",
			RuntimeTypeVersion: SemVer("100"),
			ScriptLanguage:     "JAVASCRIPT",
		},
		AdvancedOptions: SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: true,
			ResponseValidationText:             "SUCCESS",
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitor",
					Value: "synthetics",
				},
			},
			UseTlsValidation: true,
		},
	}

	//Test to create simple browser monitor
	createSimpleBrowserMonitor, err := a.SyntheticsCreateSimpleBrowserMonitor(testAccountID, simpleBrowserMonitorInput)

	require.Equal(t, 0, len(createSimpleBrowserMonitor.Errors))
	require.NoError(t, err)
	require.NotNil(t, createSimpleBrowserMonitor)

	//Input for simple browser monitor for updating
	simpleBrowserMonitorInputUpdated := SyntheticsUpdateSimpleBrowserMonitorInput{
		AdvancedOptions: SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitor",
					Value: "Synthetics",
				},
			},
			EnableScreenshotOnFailureAndScript: true,
			ResponseValidationText:             "Success",
			UseTlsValidation:                   true,
		},
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testSimpleBrowserMonitorUpdated",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "Name",
				Values: []string{
					"synthetics",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
		Runtime: SyntheticsRuntimeInput{
			RuntimeType:        "CHROME_BROWSER",
			RuntimeTypeVersion: SemVer("100"),
			ScriptLanguage:     "JAVASCRIPT",
		},
	}

	//Test to update simple browser monitor
	updateSimpleBrowserMonitor, err := a.SyntheticsUpdateSimpleBrowserMonitor(createSimpleBrowserMonitor.Monitor.GUID, simpleBrowserMonitorInputUpdated)

	require.Equal(t, 0, len(updateSimpleBrowserMonitor.Errors))
	require.NoError(t, err)
	require.NotNil(t, updateSimpleBrowserMonitor)

	//Test to delete a simple browser monitor
	deleteSimpleBrowserMonitor, err := a.SyntheticsDeleteMonitor(createSimpleBrowserMonitor.Monitor.GUID)

	require.NotNil(t, deleteSimpleBrowserMonitor)
	require.NoError(t, err)
}

//TestSyntheticsSimpleMonitor_Basic function to test simple monitor
func TestSyntheticsSimpleMonitor_Basic(t *testing.T) {

	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	////simple monitor
	//Input for creating a simple monitor
	simpleMonitorInput := SyntheticsCreateSimpleMonitorInput{
		AdvancedOptions: SyntheticsSimpleMonitorAdvancedOptionsInput{
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitor",
					Value: "Synthetics",
				},
			},
			ResponseValidationText:  "Success",
			RedirectIsFailure:       true,
			ShouldBypassHeadRequest: true,
			UseTlsValidation:        true,
		},
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testSimpleMonitor",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "Name",
				Values: []string{
					"Synthetics",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
	}

	//Test to create simple monitor
	createSimpleMonitor, err := a.SyntheticsCreateSimpleMonitor(testAccountID, simpleMonitorInput)

	require.Equal(t, 0, len(createSimpleMonitor.Errors))
	require.NoError(t, err)
	require.NotNil(t, createSimpleMonitor)

	//Input to update simple monitor
	simpleMonitorInputUpdated := SyntheticsUpdateSimpleMonitorInput{
		AdvancedOptions: SyntheticsSimpleMonitorAdvancedOptionsInput{
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitors",
					Value: "Synthetics",
				},
			},
			ResponseValidationText:  "Success",
			RedirectIsFailure:       true,
			ShouldBypassHeadRequest: true,
			UseTlsValidation:        true,
		},
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testSimpleMonitorUpdated",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "Name",
				Values: []string{
					"Synthetics",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
	}

	//Test to update simple monitor
	updateSimpleMonitor, err := a.SyntheticsUpdateSimpleMonitor(createSimpleMonitor.Monitor.GUID, simpleMonitorInputUpdated)

	require.Equal(t, 0, len(updateSimpleMonitor.Errors))
	require.NoError(t, err)
	require.NotNil(t, updateSimpleMonitor)

	//Test to delete simple monitor
	deleteSimpleMonitor, err := a.SyntheticsDeleteMonitor(createSimpleMonitor.Monitor.GUID)

	require.NotNil(t, deleteSimpleMonitor)
	require.NoError(t, err)
}

//TestSyntheticsScriptApiMonitor_Basic to test the script api monitor
func TestSyntheticsScriptApiMonitor_Basic(t *testing.T) {

	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	////Scripted API monitor
	apiScript := fmt.Sprintf(`
		const myAccountId = '%s';
		const myAPIKey = '%s';
		const options = {
		// Define endpoint URI, https://api.eu.newrelic.com/graphql for EU accounts
		uri: 'https://api.newrelic.com/graphql',
		headers: {
		'API-key': myAPIKey,
		'Content-Type': 'application/json',
		},
		body: JSON.stringify({
		query: "
		query getNrqlResults($accountId: Int!, $nrql: Nrql!) {
		actor {
		account(id: $accountId) {
		nrql(query: $nrql) {results}}}}",
		variables: {accountId: Number(myAccountId),nrql: 'SELECT average(duration) FROM Transaction'}})};
		
		// Define expected results using callback function
		function callback(err, response, body) {
		// Log JSON results from endpoint to Synthetics console
		console.log(body);
		console.log('Script execution completed');
		}
		
		// Make POST request, passing in options and callback
		$http.post(options, callback);
		`, os.Getenv("NEW_RELIC_ACCOUNT_ID"), os.Getenv("NEW_RELIC_API_KEY"))

	//input for script api monitor
	scriptApiMonitorInput := SyntheticsCreateScriptAPIMonitorInput{
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testScriptApiMonitor",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "Name",
				Values: []string{
					"ScriptAPIMonitor",
				},
			},
		},
		Runtime: SyntheticsRuntimeInput{
			RuntimeTypeVersion: SemVer("16.10"),
			RuntimeType:        "NODE_API",
			ScriptLanguage:     "JAVASCRIPT",
		},
	}

	//Test to Create scripted api monitor
	createScriptApiMonitor, err := a.SyntheticsCreateScriptAPIMonitor(testAccountID, scriptApiMonitorInput)

	require.NoError(t, err)
	require.NotNil(t, createScriptApiMonitor)
	require.Equal(t, 0, len(createScriptApiMonitor.Errors))

	//input to update script api monitor
	updatedScriptApiMonitorInput := SyntheticsUpdateScriptAPIMonitorInput{
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testScriptApiMonitorUpdated",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "Name",
				Values: []string{
					"ScriptAPIMonitor",
				},
			},
		},
		Runtime: SyntheticsRuntimeInput{
			RuntimeTypeVersion: SemVer("16.10"),
			RuntimeType:        "NODE_API",
			ScriptLanguage:     "JAVASCRIPT",
		},
	}

	//Test to update scripted api monitor
	updateScriptApiMonitor, err := a.SyntheticsUpdateScriptAPIMonitor(createScriptApiMonitor.Monitor.GUID, updatedScriptApiMonitorInput)

	require.NoError(t, err)
	require.NotNil(t, updateScriptApiMonitor)
	require.Equal(t, 0, len(updateScriptApiMonitor.Errors))

	//Test to delete scripted api monitor
	deleteScriptApiMonitor, err := a.SyntheticsDeleteMonitor(updateScriptApiMonitor.Monitor.GUID)

	require.NoError(t, err)
	require.NotNil(t, deleteScriptApiMonitor)
}

//TestSyntheticsScriptBrowserMonitor_Basic function to test script browser monitor
func TestSyntheticsScriptBrowserMonitor_Basic(t *testing.T) {

	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	//Input to create script browser monitor
	scriptBrowserMonitorInput := SyntheticsCreateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: true,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testScriptBrowserScript",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Runtime: SyntheticsRuntimeInput{
			RuntimeTypeVersion: "100",
			RuntimeType:        "CHROME_BROWSER",
			ScriptLanguage:     "JAVASCRIPT",
		},
		Tags: []SyntheticsTag{
			{
				Key: "NAME",
				Values: []string{
					"scriptBrowserScript",
				},
			},
		},
		Script: "var assert = require('assert');\n\n$browser.get('https://one.newrelic.com')",
	}

	//test to create script browser monitor
	createScriptBrowserMonitor, err := a.SyntheticsCreateScriptBrowserMonitor(testAccountID, scriptBrowserMonitorInput)

	require.NoError(t, err)
	require.NotNil(t, createScriptBrowserMonitor)
	require.Equal(t, 0, len(createScriptBrowserMonitor.Errors))

	//Input to update script browser monitor
	updatedScriptBrowserMonitorInput := SyntheticsUpdateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: true,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   "testScriptBrowserScriptUpdated",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Runtime: SyntheticsRuntimeInput{
			RuntimeTypeVersion: "100",
			RuntimeType:        "CHROME_BROWSER",
			ScriptLanguage:     "JAVASCRIPT",
		},
		Tags: []SyntheticsTag{
			{
				Key: "NAME",
				Values: []string{
					"script_browser_script",
				},
			},
		},
		Script: "var assert = require('assert');\n\n$browser.get('https://one.newrelic.com')",
	}
	//test to update script browser monitor
	updateScriptBrowserMonitor, err := a.SyntheticsUpdateScriptBrowserMonitor(createScriptBrowserMonitor.Monitor.GUID, updatedScriptBrowserMonitorInput)

	require.NoError(t, err)
	require.NotNil(t, updateScriptBrowserMonitor)
	require.Equal(t, 0, len(updateScriptBrowserMonitor.Errors))

	//test to delete script browser monitor
	deleteScriptBrowserMonitor, err := a.SyntheticsDeleteMonitor(createScriptBrowserMonitor.Monitor.GUID)

	require.NoError(t, err)
	require.NotNil(t, deleteScriptBrowserMonitor)
}

func newIntegrationTestClient(t *testing.T) Synthetics {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
