//go:build unit || integration
// +build unit integration

package keytransaction

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func newIntegrationTestClient(t *testing.T) Keytransaction {
	tc := testhelpers.NewIntegrationTestConfig(t)
	return New(tc)
}

func newIntegrationTestClient_Entities(t *testing.T) entities.Entities {
	tc := testhelpers.NewIntegrationTestConfig(t)

	return entities.New(tc)
}

var (
	testKeyTransactionName = fmt.Sprintf(
		"%s-key-transaction",
		testhelpers.GenerateRandomName(10),
	)
)
