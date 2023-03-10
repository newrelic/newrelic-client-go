package accountmanagement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCreateAccountResponseJSON = `{
    "accountManagementCreateAccount": {
      "managedAccount": {
        "id": 3833407,
        "name": "test sub account",
        "regionCode": "us01"
      }
    }
  }`
	testUpdateAccountResponseJSON = `{
    "accountManagementUpdateAccount": {
      "managedAccount": {
        "id": 3833407,
        "name": "test sub account",
        "regionCode": "us01"
      }
    }
  }`
)

func TestUpdateAccount(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testUpdateAccountResponseJSON)
	accountManagement := newMockResponse(t, respJSON, http.StatusCreated)

	updateAccountInput := AccountManagementUpdateInput{
		Name: "test sub account",
		ID:   3833407,
	}
	managedAccount := AccountManagementManagedAccount{
		Name:       updateAccountInput.Name,
		RegionCode: "us01",
		ID:         updateAccountInput.ID,
	}
	expected := &AccountManagementUpdateResponse{
		ManagedAccount: managedAccount,
	}

	actual, err := accountManagement.AccountManagementUpdateAccount(updateAccountInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateAccount(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testCreateAccountResponseJSON)
	accountManagement := newMockResponse(t, respJSON, http.StatusCreated)

	createAccountInput := AccountManagementCreateInput{
		Name:       "test sub account",
		RegionCode: "us01",
	}
	managedAccount := AccountManagementManagedAccount{
		Name:       createAccountInput.Name,
		RegionCode: createAccountInput.RegionCode,
		ID:         3833407,
	}
	expected := &AccountManagementCreateResponse{
		ManagedAccount: managedAccount,
	}

	actual, err := accountManagement.AccountManagementCreateAccount(createAccountInput)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
