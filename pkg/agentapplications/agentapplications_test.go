package agentapplications

import (
	"fmt"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	testApplicationSettings = AgentApplicationSettingsUpdateInput{
		Alias: func(s string) *string { return &s }("tf_test_updated"),
		ApmConfig: &AgentApplicationSettingsApmConfigInput{
			ApdexTarget:         0,
			UseServerSideConfig: func(b bool) *bool { return &b }(false),
		},
		ThreadProfiler: &AgentApplicationSettingsThreadProfilerInput{
			Enabled: func(b bool) *bool { return &b }(false),
		},
		ErrorCollector: &AgentApplicationSettingsErrorCollectorInput{
			Enabled:              func(b bool) *bool { return &b }(false),
			ExpectedErrorCodes:   nil,
			ExpectedErrorClasses: nil,
			IgnoredErrorCodes:    nil,
			IgnoredErrorClasses:  nil,
		},
		TransactionTracer: &AgentApplicationSettingsTransactionTracerInput{
			Enabled:                   func(b bool) *bool { return &b }(false),
			TransactionThresholdValue: func(f float64) *float64 { return &f }(0),
			TransactionThresholdType:  "",
			RecordSql:                 "",
			LogSql:                    func(b bool) *bool { return &b }(false),
			StackTraceThreshold:       func(f float64) *float64 { return &f }(0),
			ExplainEnabled:            func(b bool) *bool { return &b }(false),
			ExplainThresholdValue:     func(f float64) *float64 { return &f }(0),
			ExplainThresholdType:      "",
		},
		TracerType: &AgentApplicationSettingsTracerTypeInput{"NONE"},
	}

	testApplicationJson = `{
		  "apmSettings": {
			"Alias": "tf_test_updated",
			"ApmConfig": {
			  "ApdexTarget": 0,
			  "UseServerSideConfig": false
			},
			"ThreadProfilerEnabled": false,
			"ErrorCollector": {
			  "Enabled": false,
			  "ExpectedErrorCodes": [],
			  "ExpectedErrorClasses": [],
			  "IgnoredErrorCodes": [],
			  "IgnoredErrorClasses": []
			},
			"TransactionTracing": {
			  "Enabled": false,
			  "TransactionThresholdValue": 0,
			  "TransactionThresholdType": "off",
			  "RecordSql": "off",
			  "LogSql": false,
			  "StackTraceThresholdValue": 0,
			  "ExplainQueryPlanEnabled": false,
			  "ExplainQueryPlanThresholdValue": 0,
			  "ExplainQueryPlanThresholdType": "off"
			},
			"TracerType": "NONE"
		  }
	}`
)

func TestGetApmApplicationDetails(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"application": %s}`, testApplicationJson)
	client := newMockResponse(t, responseJSON, http.StatusOK)

	actual, err := client.GetEntity(mock.IntegrationTestApplicationEntityGUIDNew)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func TestUpdateApmApplicationDetails(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{ "application": %s}`, testApplicationJson)
	client := newMockResponseApm(t, responseJSON, http.StatusOK)

	actual, err := client.AgentApplicationSettingsUpdate(mock.IntegrationTestApplicationEntityGUIDNew, testApplicationSettings)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
