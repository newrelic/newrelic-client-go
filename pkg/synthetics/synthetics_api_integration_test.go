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

	monitorName := mock.RandSeq(5)

	////Simple Browser monitor
	//Input for simple browser monitor
	simpleBrowserMonitorInput := SyntheticsCreateSimpleBrowserMonitorInput{
		Locations: SyntheticsLocationsInput{
			Public: []string{
				"AP_SOUTH_1",
			},
		},
		Name:   monitorName,
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
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
		Runtime: SyntheticsRuntimeInput{
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

	//Test to create simple browser monitor
	createSimpleBrowserMonitor, err := a.SyntheticsCreateSimpleBrowserMonitor(testAccountID, simpleBrowserMonitorInput)

	require.NoError(t, err)
	require.NotNil(t, createSimpleBrowserMonitor)
	require.Equal(t, 0, len(createSimpleBrowserMonitor.Errors))

	//Input for simple browser monitor for updating
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
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
		Runtime: SyntheticsRuntimeInput{
			RuntimeType:        "CHROME_BROWSER",
			RuntimeTypeVersion: SemVer("100"),
			ScriptLanguage:     "JAVASCRIPT",
		},
	}

	//Test to update simple browser monitor
	updateSimpleBrowserMonitor, err := a.SyntheticsUpdateSimpleBrowserMonitor(createSimpleBrowserMonitor.Monitor.GUID, simpleBrowserMonitorInputUpdated)
	require.NoError(t, err)
	require.NotNil(t, updateSimpleBrowserMonitor)
	require.Equal(t, 0, len(updateSimpleBrowserMonitor.Errors))

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

	monitorName := mock.RandSeq(5)

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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
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

	//Test to create simple monitor
	createSimpleMonitor, err := a.SyntheticsCreateSimpleMonitor(testAccountID, simpleMonitorInput)

	require.NoError(t, err)
	require.NotNil(t, createSimpleMonitor)
	require.Equal(t, 0, len(createSimpleMonitor.Errors))

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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
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

	//Test to update simple monitor
	updateSimpleMonitor, err := a.SyntheticsUpdateSimpleMonitor(createSimpleMonitor.Monitor.GUID, simpleMonitorInputUpdated)
	require.NoError(t, err)
	require.NotNil(t, updateSimpleMonitor)
	require.Equal(t, 0, len(updateSimpleMonitor.Errors))

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

	monitorName := mock.RandSeq(5)

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
		Name:   monitorName,
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
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
		Name:   monitorName + "-updated",
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Script: apiScript,
		Tags: []SyntheticsTag{
			{
				Key: "pineapple",
				Values: []string{
					"pizza",
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
	deleteScriptApiMonitor, err := a.SyntheticsDeleteMonitor(createScriptApiMonitor.Monitor.GUID)
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

	monitorName := mock.RandSeq(5)

	//Input to create script browser monitor
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Runtime: SyntheticsRuntimeInput{
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

	//test to create script browser monitor
	createScriptBrowserMonitor, err := a.SyntheticsCreateScriptBrowserMonitor(testAccountID, scriptBrowserMonitorInput)
	require.NoError(t, err)
	require.NotNil(t, createScriptBrowserMonitor)
	require.Equal(t, 0, len(createScriptBrowserMonitor.Errors))

	//Input to update script browser monitor
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Runtime: SyntheticsRuntimeInput{
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

// Integration testing for private location
func TestSyntheticsPrivateLocation_Basic(t *testing.T) {
	t.Parallel()

	testAccountID, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

	a := newIntegrationTestClient(t)

	// Test to Create private location
	createResp, err := a.SyntheticsCreatePrivateLocation(testAccountID, "test secure credential", "TEST", true)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Test to update private location
	updateResp, err := a.SyntheticsUpdatePrivateLocation("test secure credential", createResp.GUID, true)
	require.NoError(t, err)
	require.NotNil(t, updateResp)

	// Test to purge private location queue
	purgeresp, err := a.SyntheticsPurgePrivateLocationQueue(createResp.GUID)
	require.NotNil(t, purgeresp)

	// Test to delete private location
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.DISABLED),
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.DISABLED),
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_DAY),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.DISABLED),
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_DAY),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.DISABLED),
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
		Period: SyntheticsMonitorPeriod(SyntheticsMonitorPeriodTypes.EVERY_HOUR),
		Status: SyntheticsMonitorStatus(SyntheticsMonitorStatusTypes.ENABLED),
		Runtime: SyntheticsRuntimeInput{
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
