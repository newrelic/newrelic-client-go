// Code generated by tutone: DO NOT EDIT
package agentapplications

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

// If you aren't using an auto-instrumenting agent on the backend, use this to set up browser monitoring for an application. For more information on enabling copy/paste, [see our docs](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#copy-paste-app)
func (a *AgentApplications) AgentApplicationCreateBrowser(
	accountID int,
	name string,
	settings AgentApplicationBrowserSettingsInput,
) (*AgentApplicationCreateBrowserResult, error) {
	return a.AgentApplicationCreateBrowserWithContext(context.Background(),
		accountID,
		name,
		settings,
	)
}

// If you aren't using an auto-instrumenting agent on the backend, use this to set up browser monitoring for an application. For more information on enabling copy/paste, [see our docs](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#copy-paste-app)
func (a *AgentApplications) AgentApplicationCreateBrowserWithContext(
	ctx context.Context,
	accountID int,
	name string,
	settings AgentApplicationBrowserSettingsInput,
) (*AgentApplicationCreateBrowserResult, error) {

	resp := AgentApplicationCreateBrowserQueryResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"name":      name,
		"settings":  settings,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AgentApplicationCreateBrowserMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AgentApplicationCreateBrowserResult, nil
}

type AgentApplicationCreateBrowserQueryResponse struct {
	AgentApplicationCreateBrowserResult AgentApplicationCreateBrowserResult `json:"AgentApplicationCreateBrowser"`
}

const AgentApplicationCreateBrowserMutation = `mutation(
	$accountId: Int!,
	$name: String!,
	$settings: AgentApplicationBrowserSettingsInput,
) { agentApplicationCreateBrowser(
	accountId: $accountId,
	name: $name,
	settings: $settings,
) {
	guid
	name
	settings {
		cookiesEnabled
		distributedTracingEnabled
		loaderScript
		loaderType
	}
} }`

// Deletes a browser, mobile, or APM application. This isn't allowed if an application is actively reporting data.
func (a *AgentApplications) AgentApplicationDelete(
	gUID common.EntityGUID,
) (*AgentApplicationDeleteResult, error) {
	return a.AgentApplicationDeleteWithContext(context.Background(),
		gUID,
	)
}

// Deletes a browser, mobile, or APM application. This isn't allowed if an application is actively reporting data.
func (a *AgentApplications) AgentApplicationDeleteWithContext(
	ctx context.Context,
	gUID common.EntityGUID,
) (*AgentApplicationDeleteResult, error) {

	resp := AgentApplicationDeleteQueryResponse{}
	vars := map[string]interface{}{
		"guid": gUID,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AgentApplicationDeleteMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AgentApplicationDeleteResult, nil
}

type AgentApplicationDeleteQueryResponse struct {
	AgentApplicationDeleteResult AgentApplicationDeleteResult `json:"AgentApplicationDelete"`
}

const AgentApplicationDeleteMutation = `mutation(
	$guid: EntityGuid!,
) { agentApplicationDelete(
	guid: $guid,
) {
	success
} }`

// Enable browser monitoring for an application monitored by APM. For information about specific APM agents, [see our docs](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#agent-instrumentation)
func (a *AgentApplications) AgentApplicationEnableApmBrowser(
	gUID common.EntityGUID,
	settings AgentApplicationBrowserSettingsInput,
) (*AgentApplicationEnableBrowserResult, error) {
	return a.AgentApplicationEnableApmBrowserWithContext(context.Background(),
		gUID,
		settings,
	)
}

// Enable browser monitoring for an application monitored by APM. For information about specific APM agents, [see our docs](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#agent-instrumentation)
func (a *AgentApplications) AgentApplicationEnableApmBrowserWithContext(
	ctx context.Context,
	gUID common.EntityGUID,
	settings AgentApplicationBrowserSettingsInput,
) (*AgentApplicationEnableBrowserResult, error) {

	resp := AgentApplicationEnableApmBrowserQueryResponse{}
	vars := map[string]interface{}{
		"guid":     gUID,
		"settings": settings,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AgentApplicationEnableApmBrowserMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AgentApplicationEnableBrowserResult, nil
}

type AgentApplicationEnableApmBrowserQueryResponse struct {
	AgentApplicationEnableBrowserResult AgentApplicationEnableBrowserResult `json:"AgentApplicationEnableApmBrowser"`
}

const AgentApplicationEnableApmBrowserMutation = `mutation(
	$guid: EntityGuid!,
	$settings: AgentApplicationBrowserSettingsInput,
) { agentApplicationEnableApmBrowser(
	guid: $guid,
	settings: $settings,
) {
	name
	settings {
		cookiesEnabled
		distributedTracingEnabled
		loaderType
	}
} }`

// Update configuration for APM applications. Includes thresholds for how often to record Transaction traces, SQL traces, enabling Distributed traces, ignoring certain error classes. This is the main mutation that powers the Application > Settings page in APM.
func (a *AgentApplications) AgentApplicationSettingsUpdate(
	gUID common.EntityGUID,
	settings AgentApplicationSettingsUpdateInput,
) (*AgentApplicationSettingsUpdateResult, error) {
	return a.AgentApplicationSettingsUpdateWithContext(context.Background(),
		gUID,
		settings,
	)
}

// Update configuration for APM applications. Includes thresholds for how often to record Transaction traces, SQL traces, enabling Distributed traces, ignoring certain error classes. This is the main mutation that powers the Application > Settings page in APM.
func (a *AgentApplications) AgentApplicationSettingsUpdateWithContext(
	ctx context.Context,
	gUID common.EntityGUID,
	settings AgentApplicationSettingsUpdateInput,
) (*AgentApplicationSettingsUpdateResult, error) {

	resp := AgentApplicationSettingsUpdateQueryResponse{}
	vars := map[string]interface{}{
		"guid":     gUID,
		"settings": settings,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, AgentApplicationSettingsUpdateMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AgentApplicationSettingsUpdateResult, nil
}

type AgentApplicationSettingsUpdateQueryResponse struct {
	AgentApplicationSettingsUpdateResult AgentApplicationSettingsUpdateResult `json:"AgentApplicationSettingsUpdate"`
}

const AgentApplicationSettingsUpdateMutation = `mutation(
	$guid: EntityGuid!,
	$settings: AgentApplicationSettingsUpdateInput!,
) { agentApplicationSettingsUpdate(
	guid: $guid,
	settings: $settings,
) {
	alias
	apmSettings {
		apmConfig {
			apdexTarget
			useServerSideConfig
		}
		errorCollector {
			enabled
			expectedErrorClasses
			expectedErrorCodes
			ignoredErrorClasses
			ignoredErrorCodes
		}
		jfr {
			enabled
		}
		originalName
		slowSql {
			enabled
		}
		threadProfiler {
			enabled
		}
		tracerType
		transactionTracer {
			captureMemcacheKeys
			enabled
			explainEnabled
			explainThresholdType
			explainThresholdValue
			logSql
			recordSql
			stackTraceThreshold
			transactionThresholdType
			transactionThresholdValue
		}
	}
	browserProperties {
		jsLoaderScript
	}
	browserSettings {
		browserConfig {
			apdexTarget
		}
		browserMonitoring {
			ajax {
				denyList
			}
			distributedTracing {
				allowedOrigins
				corsEnabled
				corsUseNewrelicHeader
				corsUseTracecontextHeaders
				enabled
				excludeNewrelicHeader
			}
			loader
			privacy {
				cookiesEnabled
			}
		}
	}
	errors {
		description
		errorClass
		field
	}
	guid
	mobileSettings {
		networkSettings {
			aliases {
				alias
				hosts
			}
			filterMode
			hideList
			ignoredStatusCodeRules {
				hosts
				statusCodes
			}
			showList
		}
		useCrashReports
	}
	name
} }`
