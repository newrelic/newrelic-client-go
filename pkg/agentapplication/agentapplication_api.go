// Code generated by tutone: DO NOT EDIT
package agentapplication

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

// If you aren't using an auto-instrumenting agent on the backend, use this to set up browser monitoring for an application. For more information on enabling copy/paste, [see our docs](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#copy-paste-app)
func (a *AgentApplication) AgentApplicationCreateBrowser(
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
func (a *AgentApplication) AgentApplicationCreateBrowserWithContext(
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
func (a *AgentApplication) AgentApplicationDelete(
	gUID common.EntityGUID,
) (*AgentApplicationDeleteResult, error) {
	return a.AgentApplicationDeleteWithContext(context.Background(),
		gUID,
	)
}

// Deletes a browser, mobile, or APM application. This isn't allowed if an application is actively reporting data.
func (a *AgentApplication) AgentApplicationDeleteWithContext(
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
func (a *AgentApplication) AgentApplicationEnableApmBrowser(
	gUID common.EntityGUID,
	settings AgentApplicationBrowserSettingsInput,
) (*AgentApplicationEnableBrowserResult, error) {
	return a.AgentApplicationEnableApmBrowserWithContext(context.Background(),
		gUID,
		settings,
	)
}

// Enable browser monitoring for an application monitored by APM. For information about specific APM agents, [see our docs](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#agent-instrumentation)
func (a *AgentApplication) AgentApplicationEnableApmBrowserWithContext(
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
