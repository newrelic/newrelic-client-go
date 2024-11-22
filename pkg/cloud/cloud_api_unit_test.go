package cloud

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	testCreateAzureMonitorIntegration = `{
				"__typename": "CloudAzureMonitorIntegration",
				"createdAt": 1682411205,
				"enabled": true,
				"excludeTags": ["env:staging", "env:testing"],
				"id": 1709478,
				"includeTags": ["env:production"],
				"inventoryPollingInterval": null,
				"metricsPollingInterval": 1200,
				"name": "Azure Monitor metrics",
				"nrAccountId": ` + nrAccountID + `,
				"resourceGroups": ["resource_groups"],
				"resourceTypes": ["microsoft.datashare/accounts"],
				"updatedAt": 1682413262
				}`
	testDeleteAzureMonitorDisabledIntegration = `{
				"__typename": "CloudAzureMonitorIntegration",
				"createdAt": 1682437459,
				"enabled": true,
				"excludeTags": ["env:staging"],
				"id": 1709859,
				"includeTags": ["env:production", "env:testing"],
				"inventoryPollingInterval": null,
				"metricsPollingInterval": 1200,
				"name": "Azure Monitor metrics",
				"nrAccountId": ` + nrAccountID + `,
				"resourceGroups": ["resource_groups"],
				"resourceTypes": ["microsoft.datashare/accounts"],
				"updatedAt": 1682437515
			}`

	testDeleteAzureMonitor = `{
	"data": {
		"cloudDisableIntegration": {
			"disabledIntegrations": [` + testDeleteAzureMonitorDisabledIntegration + `],
			"errors": []
		}
	}
  }`
	testCreateAzureMonitor = `{
	"data": {
		"cloudConfigureIntegration": {
			"errors": [],
			"integrations": [` + testCreateAzureMonitorIntegration + `]
		}
	}
  }`
	linkedAccountID = fmt.Sprintf("%06d", rand.Int63n(1e6))
	nrAccountID     = fmt.Sprintf("%06d", rand.Int63n(1e6))

	testUpdateAzureLinkAccount = `
{
  "data": {
    "cloudUpdateAccount": {
      "linkedAccounts": [
        {
          "authLabel": "36840357-ac3e-4273-94f0-eccg108ff0e9",
          "disabled": false,
          "externalId": "agjs-dha57-687hag-shgafshd-f79hh",
          "id": ` + linkedAccountID + `,
          "name": "TEST_AZURE_ACCOUNT_UPDATED",
          "nrAccountId": ` + nrAccountID + `,
          "updatedAt": 1729674748
        }
      ]
    }
  }
}
	`

	testCreateFOSSALinkAccount = `
{
  "data": {
    "cloudLinkAccount": {
      "linkedAccounts": [
        {
          "id": ` + linkedAccountID + `,
          "name": "TEST_FOSSA_ACCOUNT",
          "nrAccountId": ` + nrAccountID + `,
        }
      ]
    }
  }
}
	`
)

// Unit Test to test the creation of an Azure Monitor.
// Applies to update too, as create and update use the same mutation 'CloudConfigureIntegration'.
func TestUnitCreateAzureMonitor(t *testing.T) {
	t.Parallel()
	createAzureMonitorResponse := newMockResponse(t, testCreateAzureMonitor, http.StatusCreated)
	linkedAccountIDAsInt, _ := strconv.Atoi(linkedAccountID)
	createAzureMonitorInput := CloudIntegrationsInput{
		Azure: CloudAzureIntegrationsInput{
			AzureMonitor: []CloudAzureMonitorIntegrationInput{{
				LinkedAccountId:        linkedAccountIDAsInt,
				Enabled:                true,
				ExcludeTags:            []string{"env:staging", "env:testing"},
				IncludeTags:            []string{"env:production"},
				MetricsPollingInterval: 1200,
				ResourceTypes:          []string{"microsoft.datashare/accounts"},
				ResourceGroups:         []string{"resource_groups"},
			}},
		},
	}

	NRAccountIDInt, _ := strconv.Atoi(nrAccountID)
	actual, err := createAzureMonitorResponse.CloudConfigureIntegration(NRAccountIDInt, createAzureMonitorInput)

	responseJSON, _ := json.Marshal(actual.Integrations[0])
	responseJSONAsString := string(responseJSON)
	objActual, objExpected, objError := unmarshalAzureCloudIntegrationJSON(responseJSONAsString, testCreateAzureMonitorIntegration)

	assert.NoError(t, err)
	assert.NoError(t, objError)
	assert.NotNil(t, actual)
	assert.Equal(t, objActual, objExpected)

}

// Unit Test to test the deletion of an Azure Monitor.
func TestUnitDeleteAzureMonitor(t *testing.T) {
	t.Parallel()
	deleteAzureMonitorResponse := newMockResponse(t, testDeleteAzureMonitor, http.StatusCreated)
	linkedAccountIDAsInt, _ := strconv.Atoi(linkedAccountID)

	deleteAzureMonitorInput := CloudDisableIntegrationsInput{
		Azure: CloudAzureDisableIntegrationsInput{
			AzureMonitor: []CloudDisableAccountIntegrationInput{{
				LinkedAccountId: linkedAccountIDAsInt,
			}},
		},
	}

	NRAccountIDInt, _ := strconv.Atoi(nrAccountID)
	actual, err := deleteAzureMonitorResponse.CloudDisableIntegration(NRAccountIDInt, deleteAzureMonitorInput)

	responseJSON, _ := json.Marshal(actual.DisabledIntegrations[0])
	responseJSONAsString := string(responseJSON)

	objActual, objExpected, objError := unmarshalAzureCloudIntegrationJSON(responseJSONAsString, testDeleteAzureMonitorDisabledIntegration)

	assert.NoError(t, err)
	assert.NoError(t, objError)
	assert.NotNil(t, actual)
	assert.Equal(t, objActual, objExpected)

}

func unmarshalAzureCloudIntegrationJSON(actualJSONString string, expectedJSONString string) (CloudAzureMonitorIntegration, CloudAzureMonitorIntegration, error) {
	var actual, expected CloudAzureMonitorIntegration

	errActual := json.Unmarshal([]byte(actualJSONString), &actual)
	errExpected := json.Unmarshal([]byte(expectedJSONString), &expected)

	if errActual != nil {
		return actual, expected, errActual
	}

	if errExpected != nil {
		return actual, expected, errExpected
	}
	return actual, expected, nil
}

func newMockResponse(t *testing.T, mockJSONResponse string, statusCode int) Cloud {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)
	tc := mock.NewTestConfig(t, ts)

	return New(tc)
}

// unit test to test azure link account update mutation
func TestUnitAzureLinkAccountUpdate(t *testing.T) {
	t.Parallel()
	azureUpdateAccountResponse := newMockResponse(t, testUpdateAzureLinkAccount, http.StatusOK)
	NRAccountIDInt, _ := strconv.Atoi(nrAccountID)
	linkedAccountIDInt, _ := strconv.Atoi(linkedAccountID)
	disabled := false

	updateAccountInput := CloudUpdateCloudAccountsInput{
		Azure: []CloudAzureUpdateAccountInput{{
			ApplicationID:   "36840357-ac3e-4273-94f0-eccg108ff0e9",
			ClientSecret:    "gdsajysgda676t5ahgsdhafsdga67as",
			Disabled:        &disabled,
			LinkedAccountId: linkedAccountIDInt,
			Name:            "TEST_AZURE_ACCOUNT-UPDATED",
			SubscriptionId:  "agjs-dha57-687hag-shgafshd-f79hh",
			TenantId:        "ajkhsdjkas676hjgasdhjga687yhhj",
		}},
	}

	actual, err := azureUpdateAccountResponse.CloudUpdateAccount(NRAccountIDInt, updateAccountInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, linkedAccountID, strconv.Itoa(actual.LinkedAccounts[0].ID))
	assert.Equal(t, "TEST_AZURE_ACCOUNT_UPDATED", actual.LinkedAccounts[0].Name)
	assert.Equal(t, false, actual.LinkedAccounts[0].Disabled)
	assert.Equal(t, nrAccountID, strconv.Itoa(actual.LinkedAccounts[0].NrAccountId))

}

func TestUnitFOSSALinkAccountCreate(t *testing.T) {
	t.Parallel()
	createFOSSALinkAccountResponse := newMockResponse(t, testCreateFOSSALinkAccount, http.StatusOK)
	NRAccountIDInt, _ := strconv.Atoi(nrAccountID)

	createAccountInput := CloudLinkCloudAccountsInput{
		Fossa: []CloudFossaLinkAccountInput{{
			APIKey: "f3ab70a3d29028c94ab8057f1b30838b",
			//ExternalId: "agjs-dha57-687hag-shgafshd-f79hh",
			Name: "TEST_FOSSA_ACCOUNT",
		}},
	}

	actual, err := createFOSSALinkAccountResponse.CloudLinkAccount(NRAccountIDInt, createAccountInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, linkedAccountID, strconv.Itoa(actual.LinkedAccounts[0].ID))
	assert.Equal(t, "TEST_FOSSA_ACCOUNT", actual.LinkedAccounts[0].Name)
	assert.Equal(t, "agjs-dha57-687hag-shgafshd-f79hh", actual.LinkedAccounts[0].ExternalId)
}
