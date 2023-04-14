//go:build integration
// +build integration

package synthetics

import (
	"testing"

	"github.com/stretchr/testify/require"

	"fmt"
	"os"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var tv bool = true

func TestSyntheticsSecureCredential_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	createResp, err := a.SyntheticsCreateSecureCredential(testAccountID, "test secure credential", "TEST", "secure value")
	require.NoError(t, err)
	require.NotNil(t, createResp)

	updateResp, err := a.SyntheticsUpdateSecureCredential(testAccountID, "test secure credential", "TEST", "new secure value")
	require.NoError(t, err)
	require.NotNil(t, updateResp)

	deleteResp, err := a.SyntheticsDeleteSecureCredential(testAccountID, "TEST")
	require.Equal(t, "", deleteResp.Key)
}

func TestSyntheticsSecureCredential_Error(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	createResp, respErr := a.SyntheticsCreateSecureCredential(testAccountID, "test secure credential", "TEST-BAD-KEY", "secure value")
	require.NoError(t, respErr)
	require.Greater(t, len(createResp.Errors), 0)
}

func TestSyntheticsSimpleBrowserMonitor_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := mock.RandSeq(5)

	simpleBrowserMonitorInput := SyntheticsCreateSimpleBrowserMonitorInput{
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
		Runtime: &SyntheticsRuntimeInput{
			RuntimeType:        "CHROME_BROWSER",
			RuntimeTypeVersion: SemVer("100"),
			ScriptLanguage:     "JAVASCRIPT",
		},
		AdvancedOptions: SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
			ResponseValidationText:             "SUCCESS",
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitor",
					Value: "synthetics",
				},
			},
			UseTlsValidation: &tv,
		},
	}

	createSimpleBrowserMonitor, err := a.SyntheticsCreateSimpleBrowserMonitor(testAccountID, simpleBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, createSimpleBrowserMonitor)
	require.Equal(t, 0, len(createSimpleBrowserMonitor.Errors))

	simpleBrowserMonitorInputUpdated := SyntheticsUpdateSimpleBrowserMonitorInput{
		AdvancedOptions: SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitor",
					Value: "Synthetics",
				},
			},
			EnableScreenshotOnFailureAndScript: &tv,
			ResponseValidationText:             "Success",
			UseTlsValidation:                   &tv,
		},
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
		Runtime: &SyntheticsRuntimeInput{
			RuntimeType:        "CHROME_BROWSER",
			RuntimeTypeVersion: SemVer("100"),
			ScriptLanguage:     "JAVASCRIPT",
		},
	}

	updateSimpleBrowserMonitor, err := a.SyntheticsUpdateSimpleBrowserMonitor(createSimpleBrowserMonitor.Monitor.GUID, simpleBrowserMonitorInputUpdated)
	require.NoError(t, err)
	require.NotNil(t, updateSimpleBrowserMonitor)
	require.Equal(t, 0, len(updateSimpleBrowserMonitor.Errors))

	deleteSimpleBrowserMonitor, err := a.SyntheticsDeleteMonitor(createSimpleBrowserMonitor.Monitor.GUID)
	require.NotNil(t, deleteSimpleBrowserMonitor)
	require.NoError(t, err)
}

func TestSyntheticsSimpleMonitor_Basic(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := mock.RandSeq(5)

	simpleMonitorInput := SyntheticsCreateSimpleMonitorInput{
		AdvancedOptions: SyntheticsSimpleMonitorAdvancedOptionsInput{
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitor",
					Value: "Synthetics",
				},
			},
			ResponseValidationText:  "Success",
			RedirectIsFailure:       &tv,
			ShouldBypassHeadRequest: &tv,
			UseTlsValidation:        &tv,
		},
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
	}

	createSimpleMonitor, err := a.SyntheticsCreateSimpleMonitor(testAccountID, simpleMonitorInput)

	require.NoError(t, err)
	require.NotNil(t, createSimpleMonitor)
	require.Equal(t, 0, len(createSimpleMonitor.Errors))

	simpleMonitorInputUpdated := SyntheticsUpdateSimpleMonitorInput{
		AdvancedOptions: SyntheticsSimpleMonitorAdvancedOptionsInput{
			CustomHeaders: []SyntheticsCustomHeaderInput{
				{
					Name:  "Monitors",
					Value: "Synthetics",
				},
			},
			ResponseValidationText:  "Success",
			RedirectIsFailure:       &tv,
			ShouldBypassHeadRequest: &tv,
			UseTlsValidation:        &tv,
		},
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
		Uri: "https://www.one.newrelic.com",
	}

	updateSimpleMonitor, err := a.SyntheticsUpdateSimpleMonitor(createSimpleMonitor.Monitor.GUID, simpleMonitorInputUpdated)
	require.NoError(t, err)
	require.NotNil(t, updateSimpleMonitor)
	require.Equal(t, 0, len(updateSimpleMonitor.Errors))

	deleteSimpleMonitor, err := a.SyntheticsDeleteMonitor(createSimpleMonitor.Monitor.GUID)
	require.NotNil(t, deleteSimpleMonitor)
	require.NoError(t, err)
}

func TestSyntheticsScriptApiMonitor_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := mock.RandSeq(5)

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

	scriptApiMonitorInput := SyntheticsCreateScriptAPIMonitorInput{
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
		Runtime: &SyntheticsRuntimeInput{
			RuntimeTypeVersion: SemVer("16.10"),
			RuntimeType:        "NODE_API",
			ScriptLanguage:     "JAVASCRIPT",
		},
	}

	createScriptApiMonitor, err := a.SyntheticsCreateScriptAPIMonitor(testAccountID, scriptApiMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, createScriptApiMonitor)
	require.Equal(t, 0, len(createScriptApiMonitor.Errors))

	updatedScriptApiMonitorInput := SyntheticsUpdateScriptAPIMonitorInput{
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
		Runtime: &SyntheticsRuntimeInput{
			RuntimeTypeVersion: SemVer("16.10"),
			RuntimeType:        "NODE_API",
			ScriptLanguage:     "JAVASCRIPT",
		},
	}

	updateScriptApiMonitor, err := a.SyntheticsUpdateScriptAPIMonitor(createScriptApiMonitor.Monitor.GUID, updatedScriptApiMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, updateScriptApiMonitor)
	require.Equal(t, 0, len(updateScriptApiMonitor.Errors))

	deleteScriptApiMonitor, err := a.SyntheticsDeleteMonitor(createScriptApiMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteScriptApiMonitor)
}

func TestSyntheticsScriptApiMonitorLegacy_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := mock.RandSeq(5)

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

	scriptApiMonitorInput := SyntheticsCreateScriptAPIMonitorInput{
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
	}

	createScriptApiMonitor, err := a.SyntheticsCreateScriptAPIMonitor(testAccountID, scriptApiMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, createScriptApiMonitor)
	require.Equal(t, 0, len(createScriptApiMonitor.Errors), createScriptApiMonitor.Errors)

	updatedScriptApiMonitorInput := SyntheticsUpdateScriptAPIMonitorInput{
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
	}

	updateScriptApiMonitor, err := a.SyntheticsUpdateScriptAPIMonitor(createScriptApiMonitor.Monitor.GUID, updatedScriptApiMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, updateScriptApiMonitor)
	require.Equal(t, 0, len(updateScriptApiMonitor.Errors))

	deleteScriptApiMonitor, err := a.SyntheticsDeleteMonitor(createScriptApiMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteScriptApiMonitor)
}

func TestSyntheticsScriptBrowserMonitor_Basic(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := mock.RandSeq(5)

	scriptBrowserMonitorInput := SyntheticsCreateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Runtime: &SyntheticsRuntimeInput{
			RuntimeTypeVersion: "100",
			RuntimeType:        "CHROME_BROWSER",
			ScriptLanguage:     "JAVASCRIPT",
		},
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
				},
			},
		},
		Script: "var assert = require('assert');\n\n$browser.get('https://one.newrelic.com')",
	}

	createScriptBrowserMonitor, err := a.SyntheticsCreateScriptBrowserMonitor(testAccountID, scriptBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, createScriptBrowserMonitor)
	require.Equal(t, 0, len(createScriptBrowserMonitor.Errors))

	updatedScriptBrowserMonitorInput := SyntheticsUpdateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Runtime: &SyntheticsRuntimeInput{
			RuntimeTypeVersion: "100",
			RuntimeType:        "CHROME_BROWSER",
			ScriptLanguage:     "JAVASCRIPT",
		},
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"script_browser_pizza",
				},
			},
		},
		Script: "var assert = require('assert');\n\n$browser.get('https://one.newrelic.com')",
	}

	updateScriptBrowserMonitor, err := a.SyntheticsUpdateScriptBrowserMonitor(createScriptBrowserMonitor.Monitor.GUID, updatedScriptBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, updateScriptBrowserMonitor)
	require.Equal(t, 0, len(updateScriptBrowserMonitor.Errors))

	deleteScriptBrowserMonitor, err := a.SyntheticsDeleteMonitor(createScriptBrowserMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteScriptBrowserMonitor)
}

func TestSyntheticsScriptBrowserMonitor_DeviceEmulation(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := mock.RandSeq(5)

	scriptBrowserMonitorInput := SyntheticsCreateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
			DeviceEmulation: &SyntheticsDeviceEmulationInput{
				DeviceOrientation: SyntheticsDeviceOrientationTypes.PORTRAIT,
				DeviceType:        SyntheticsDeviceTypeTypes.MOBILE,
			},
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{"AP_SOUTH_1"},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_12_HOURS,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Runtime: &SyntheticsRuntimeInput{
			RuntimeTypeVersion: "100",
			RuntimeType:        "CHROME_BROWSER",
			ScriptLanguage:     "JAVASCRIPT",
		},
		Script: "console.log('test')",
	}

	createScriptBrowserMonitor, err := a.SyntheticsCreateScriptBrowserMonitor(testAccountID, scriptBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, createScriptBrowserMonitor)
	require.Equal(t, 0, len(createScriptBrowserMonitor.Errors))

	updatedScriptBrowserMonitorInput := SyntheticsUpdateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
			// Test changing device emulation options
			DeviceEmulation: &SyntheticsDeviceEmulationInput{
				DeviceOrientation: SyntheticsDeviceOrientationTypes.LANDSCAPE,
				DeviceType:        SyntheticsDeviceTypeTypes.TABLET,
			},
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{"AP_SOUTH_1"},
		},
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriodTypes.EVERY_12_HOURS,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Runtime: &SyntheticsRuntimeInput{
			RuntimeTypeVersion: "100",
			RuntimeType:        "CHROME_BROWSER",
			ScriptLanguage:     "JAVASCRIPT",
		},
		Script: "console.log('test)",
	}

	updateScriptBrowserMonitor, err := a.SyntheticsUpdateScriptBrowserMonitor(createScriptBrowserMonitor.Monitor.GUID, updatedScriptBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, updateScriptBrowserMonitor)
	require.Equal(t, 0, len(updateScriptBrowserMonitor.Errors))

	deleteScriptBrowserMonitor, err := a.SyntheticsDeleteMonitor(createScriptBrowserMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteScriptBrowserMonitor)
}

func TestSyntheticsScriptBrowserMonitor_LegacyRuntime(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := mock.RandSeq(5)

	scriptBrowserMonitorInput := SyntheticsCreateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Private: []SyntheticsPrivateLocationInput{
				{
					GUID: "MzgwNjUyNnxTWU5USHxQUklWQVRFX0xPQ0FUSU9OfGNhNmZmNTY3LTJlZWItNGNkNi04ODhhLTAxMTFjMGMzMTBjNA",
				},
			},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Script: "console.log('test')",
		// Note: Omitting `Runtime` defaults to the legacy runtime
	}

	createScriptBrowserMonitor, err := a.SyntheticsCreateScriptBrowserMonitor(testAccountID, scriptBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, createScriptBrowserMonitor)
	require.Equal(t, 0, len(createScriptBrowserMonitor.Errors))

	updatedScriptBrowserMonitorInput := SyntheticsUpdateScriptBrowserMonitorInput{
		AdvancedOptions: SyntheticsScriptBrowserMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &tv,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Private: []SyntheticsPrivateLocationInput{
				{
					GUID: "MzgwNjUyNnxTWU5USHxQUklWQVRFX0xPQ0FUSU9OfGNhNmZmNTY3LTJlZWItNGNkNi04ODhhLTAxMTFjMGMzMTBjNA",
				},
			},
		},
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Script: "var assert = require('assert');\n\n$browser.get('https://one.newrelic.com')",
	}

	updateScriptBrowserMonitor, err := a.SyntheticsUpdateScriptBrowserMonitor(createScriptBrowserMonitor.Monitor.GUID, updatedScriptBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, updateScriptBrowserMonitor)
	require.Equal(t, 0, len(updateScriptBrowserMonitor.Errors))

	deleteScriptBrowserMonitor, err := a.SyntheticsDeleteMonitor(createScriptBrowserMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deleteScriptBrowserMonitor)
}

func TestSyntheticsPrivateLocation_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	createResp, err := a.SyntheticsCreatePrivateLocation(testAccountID, "test secure credential", "TEST", true)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	updateResp, err := a.SyntheticsUpdatePrivateLocation("test secure credential", createResp.GUID, true)
	require.NoError(t, err)
	require.NotNil(t, updateResp)

	purgeresp, err := a.SyntheticsPurgePrivateLocationQueue(createResp.GUID)
	require.NotNil(t, purgeresp)

	deleteResp, err := a.SyntheticsDeletePrivateLocation(createResp.GUID)
	require.NotNil(t, deleteResp)
}

func TestSyntheticsBrokenLinksMonitor_Basic(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := fmt.Sprintf("client-integration-test-%s", mock.RandSeq(5))
	monitorInput := SyntheticsCreateBrokenLinksMonitorInput{
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.DISABLED,
		Locations: SyntheticsLocationsInput{
			Public: []string{"AP_SOUTH_1"},
		},
		Tags: []SyntheticsTag{
			{
				Key:    "coconut",
				Values: []string{"avocado"},
			},
		},
		Uri: "https://www.google.com",
	}

	createdMonitor, err := a.SyntheticsCreateBrokenLinksMonitor(testAccountID, monitorInput)
	require.NoError(t, err)
	require.NotNil(t, createdMonitor)
	require.Equal(t, 0, len(createdMonitor.Errors))

	monitorNameUpdate := fmt.Sprintf("%s-updated", monitorName)
	monitorUpdateInput := SyntheticsUpdateBrokenLinksMonitorInput{
		Name:      fmt.Sprintf("%s-updated", monitorName),
		Period:    monitorInput.Period,
		Status:    monitorInput.Status,
		Locations: monitorInput.Locations,
		Tags:      monitorInput.Tags,
		Uri:       fmt.Sprintf("%s?updated=true", monitorInput.Uri),
	}

	updatedMonitor, err := a.SyntheticsUpdateBrokenLinksMonitor(createdMonitor.Monitor.GUID, monitorUpdateInput)
	require.NoError(t, err)
	require.NotNil(t, updatedMonitor.Monitor)
	require.Equal(t, 0, len(updatedMonitor.Errors))
	require.Equal(t, monitorNameUpdate, updatedMonitor.Monitor.Name)
	require.Equal(t, "https://www.google.com?updated=true", updatedMonitor.Monitor.Uri)
	require.Equal(t, createdMonitor.Monitor.GUID, updatedMonitor.Monitor.GUID)

	deletedMonitor, err := a.SyntheticsDeleteMonitor(createdMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deletedMonitor)
	require.Equal(t, createdMonitor.Monitor.GUID, deletedMonitor.DeletedGUID)
}

func TestSyntheticsCertCheckMonitor_Basic(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := fmt.Sprintf("client-integration-test-%s", mock.RandSeq(5))
	monitorInput := SyntheticsCreateCertCheckMonitorInput{
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
		Status: SyntheticsMonitorStatusTypes.DISABLED,
		Locations: SyntheticsLocationsInput{
			Public: []string{"AP_SOUTH_1"},
		},
		Tags: []SyntheticsTag{
			{
				Key:    "coconut",
				Values: []string{"avocado"},
			},
		},
		Domain:                            "https://www.google.com",
		NumberDaysToFailBeforeCertExpires: 1,
	}

	createdMonitor, err := a.SyntheticsCreateCertCheckMonitor(testAccountID, monitorInput)
	require.NoError(t, err)
	require.NotNil(t, createdMonitor)
	require.Equal(t, 0, len(createdMonitor.Errors))

	monitorNameUpdate := fmt.Sprintf("%s-updated", monitorName)
	monitorUpdateInput := SyntheticsUpdateCertCheckMonitorInput{
		Name:                              fmt.Sprintf("%s-updated", monitorName),
		Period:                            monitorInput.Period,
		Status:                            monitorInput.Status,
		Locations:                         monitorInput.Locations,
		Tags:                              monitorInput.Tags,
		Domain:                            fmt.Sprintf("%s?updated=true", monitorInput.Domain),
		NumberDaysToFailBeforeCertExpires: 2,
	}

	updatedMonitor, err := a.SyntheticsUpdateCertCheckMonitor(createdMonitor.Monitor.GUID, monitorUpdateInput)
	require.NoError(t, err)
	require.NotNil(t, updatedMonitor.Monitor)
	require.Equal(t, 0, len(updatedMonitor.Errors))
	require.Equal(t, monitorNameUpdate, updatedMonitor.Monitor.Name)
	require.Equal(t, "https://www.google.com?updated=true", updatedMonitor.Monitor.Domain)
	require.Equal(t, 2, updatedMonitor.Monitor.NumberDaysToFailBeforeCertExpires)
	require.Equal(t, createdMonitor.Monitor.GUID, updatedMonitor.Monitor.GUID)

	deletedMonitor, err := a.SyntheticsDeleteMonitor(createdMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deletedMonitor)
	require.Equal(t, createdMonitor.Monitor.GUID, deletedMonitor.DeletedGUID)
}

func TestSyntheticsStepMonitor_Basic(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := fmt.Sprintf("client-integration-test-%s", mock.RandSeq(5))
	enableScreenshotOnFailureAndScript := true
	monitorInput := SyntheticsCreateStepMonitorInput{
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_DAY,
		Status: SyntheticsMonitorStatusTypes.DISABLED,
		AdvancedOptions: SyntheticsStepMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &enableScreenshotOnFailureAndScript,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{"AP_SOUTH_1"},
		},
		Tags: []SyntheticsTag{
			{
				Key:    "step",
				Values: []string{"monitor"},
			},
		},
		Steps: []SyntheticsStepInput{
			{
				Ordinal: 0,
				Type:    SyntheticsStepTypeTypes.NAVIGATE,
				Values:  []string{"https://one.newrelic.com"},
			},
			{
				Ordinal: 1,
				Type:    SyntheticsStepTypeTypes.ASSERT_TITLE,
				Values:  []string{"%=", "New Relic"}, // %= is used for "contains" logic
			},
		},
	}

	createdMonitor, err := a.SyntheticsCreateStepMonitor(testAccountID, monitorInput)
	require.NoError(t, err)
	require.NotNil(t, createdMonitor)
	require.Equal(t, 0, len(createdMonitor.Errors))
	require.Equal(t, 2, len(createdMonitor.Monitor.Steps))

	monitorNameUpdate := fmt.Sprintf("%s-updated", monitorName)
	monitorUpdateInput := SyntheticsUpdateStepMonitorInput{
		Name: fmt.Sprintf("%s-updated", monitorName),
		Steps: []SyntheticsStepInput{
			{
				Ordinal: 0,
				Type:    SyntheticsStepTypeTypes.NAVIGATE,
				Values:  []string{"https://one.newrelic.com"},
			},
			{
				Ordinal: 1,
				Type:    SyntheticsStepTypeTypes.ASSERT_TITLE,
				Values:  []string{"%=", "New Relic"}, // %= is used for "contains" logic
			},
			{
				Ordinal: 2,
				Type:    SyntheticsStepTypeTypes.ASSERT_ELEMENT,
				Values:  []string{"h2.NewDesign", "present", "true"},
			},
		},
	}

	updatedMonitor, err := a.SyntheticsUpdateStepMonitor(createdMonitor.Monitor.GUID, monitorUpdateInput)
	require.NoError(t, err)
	require.NotNil(t, updatedMonitor.Monitor)
	require.Equal(t, 0, len(updatedMonitor.Errors))
	require.Equal(t, monitorNameUpdate, updatedMonitor.Monitor.Name)
	require.Equal(t, 3, len(updatedMonitor.Monitor.Steps))

	deletedMonitor, err := a.SyntheticsDeleteMonitor(createdMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deletedMonitor)
	require.Equal(t, createdMonitor.Monitor.GUID, deletedMonitor.DeletedGUID)
}

func TestSyntheticsStepMonitor_GetSteps(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := fmt.Sprintf("client-integration-test-%s", mock.RandSeq(5))
	enableScreenshotOnFailureAndScript := true
	monitorInput := SyntheticsCreateStepMonitorInput{
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_DAY,
		Status: SyntheticsMonitorStatusTypes.DISABLED,
		AdvancedOptions: SyntheticsStepMonitorAdvancedOptionsInput{
			EnableScreenshotOnFailureAndScript: &enableScreenshotOnFailureAndScript,
		},
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{"AP_SOUTH_1"},
		},
		Tags: []SyntheticsTag{
			{
				Key:    "step",
				Values: []string{"monitor"},
			},
		},
		Steps: []SyntheticsStepInput{
			{
				Ordinal: 0,
				Type:    SyntheticsStepTypeTypes.NAVIGATE,
				Values:  []string{"https://one.newrelic.com"},
			},
			{
				Ordinal: 1,
				Type:    SyntheticsStepTypeTypes.ASSERT_TITLE,
				Values:  []string{"%=", "New Relic"}, // %= is used for "contains" logic
			},
		},
	}

	createdMonitor, err := a.SyntheticsCreateStepMonitor(testAccountID, monitorInput)
	require.NoError(t, err)
	require.NotNil(t, createdMonitor)

	// Test the `steps` query endpoint
	steps, err := a.GetSteps(testAccountID, createdMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, steps)
	require.Equal(t, 2, len(*steps))

	deletedMonitor, err := a.SyntheticsDeleteMonitor(createdMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deletedMonitor)
	require.Equal(t, createdMonitor.Monitor.GUID, deletedMonitor.DeletedGUID)
}

func TestSyntheticsStepMonitor_GetScript(t *testing.T) {
	t.Parallel()
	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	monitorName := fmt.Sprintf("client-integration-test-%s", mock.RandSeq(5))
	monitorInput := SyntheticsCreateScriptBrowserMonitorInput{
		Locations: SyntheticsScriptedMonitorLocationsInput{
			Public: []string{"AP_SOUTH_1"},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriodTypes.EVERY_HOUR,
		Status: SyntheticsMonitorStatusTypes.ENABLED,
		Runtime: &SyntheticsRuntimeInput{
			RuntimeTypeVersion: "100",
			RuntimeType:        "CHROME_BROWSER",
			ScriptLanguage:     "JAVASCRIPT",
		},
		Script: "var assert = require('assert');\n\n$browser.get('https://api.newrelic.com')",
	}

	createdMonitor, err := a.SyntheticsCreateScriptBrowserMonitor(testAccountID, monitorInput)
	require.NoError(t, err)
	require.NotNil(t, createdMonitor)

	// Test the `steps` query endpoint
	script, err := a.GetScript(testAccountID, createdMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, script)
	require.NotEmpty(t, script.Text)

	deletedMonitor, err := a.SyntheticsDeleteMonitor(createdMonitor.Monitor.GUID)
	require.NoError(t, err)
	require.NotNil(t, deletedMonitor)
	require.Equal(t, createdMonitor.Monitor.GUID, deletedMonitor.DeletedGUID)
}

func newIntegrationTestClient(t *testing.T) Synthetics {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
