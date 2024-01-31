//go:build integration
// +build integration

package usermanagement

import (
	"testing"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestIntegrationCreateWorkflow(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("%s", err)
	}

}
