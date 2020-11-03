// +build unit

package infrastructure

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testLinkCloudAccountResponseJSON = `{
		"data": {
			"linkedAccounts": [
				{
					"id": 17327,
					"name": "unit-test-linked-account1",
					"authLabel": "unit-test-auth-label1"
				},
				{
					"id": 28934,
					"name": "unit-test-linked-account2",
					"authLabel": "unit-test-auth-label2"
				}
			]
		}
	}`
	testUnlinkCloudAccountResponseJSON = `{
		"data": {
			"unlinkedAccounts": [
				{
					"id": 17327,
					"name": "unit-test-linked-account1",
				},
				{
					"id": 28934,
					"name": "unit-test-linked-account2",
				}
			]
		}
	}`
	testErrorResponseJSON = `{
		"errors": [
			{
				"message": "Could not link cloud account"
			}
		]
	}`
)

func TestInfrastructure_LinkCloudAccount_Success(t *testing.T) {
	t.Parallel()
	infrastructure := newMockResponse(t, testLinkCloudAccountResponseJSON, http.StatusOK)

	expected := []LinkedCloudAccount{
		{
			ID:        17327,
			Name:      "unit-test-linked-account1",
			AuthLabel: "unit-test-auth-label1",
		},
		{
			ID:        28934,
			Name:      "unit-test-linked-account2",
			AuthLabel: "unit-test-auth-label2",
		},
	}

	actual, err := infrastructure.LinkCloudAccount(LinkCloudAccountInput{AccountID: 3242})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestInfrastructure_LinkCloudAccount_Error(t *testing.T) {
	t.Parallel()
	infrastructure := newMockResponse(t, testErrorResponseJSON, http.StatusOK)
	_, err := infrastructure.LinkCloudAccount(LinkCloudAccountInput{AccountID: 3242})

	assert.Error(t, err)
}

func TestInfrastructure_UnlinkCloudAccount_Success(t *testing.T) {
	t.Parallel()
	infrastructure := newMockResponse(t, testLinkCloudAccountResponseJSON, http.StatusOK)

	expected := []UnlinkedCloudAccount{
		{
			ID:   17327,
			Name: "unit-test-linked-account1",
		},
		{
			ID:   28934,
			Name: "unit-test-linked-account2",
		},
	}

	actual, err := infrastructure.UnlinkCloudAccount(3242, UnlinkCloudAccountInput{
		Accounts: []LinkedCloudAccountRef{
			{LinkedAccountId: 17327},
			{LinkedAccountId: 28934},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestInfrastructure_UnlinkCloudAccount_Error(t *testing.T) {
	t.Parallel()
	infrastructure := newMockResponse(t, testErrorResponseJSON, http.StatusOK)
	_, err := infrastructure.UnlinkCloudAccount(3242, UnlinkCloudAccountInput{})

	assert.Error(t, err)
}
