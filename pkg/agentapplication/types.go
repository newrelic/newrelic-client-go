// Code generated by tutone: DO NOT EDIT
package agentapplication

import "github.com/newrelic/newrelic-client-go/v2/pkg/common"

// AgentApplicationBrowserLoader - Determines which browser loader will be configured. There are three browser loader types. They are Pro+SPA, Pro, and Lite.
// See [documentation](https://docs.newrelic.com/docs/browser/browser-monitoring/installation/install-browser-monitoring-agent/#agent-types) for further information.
type AgentApplicationBrowserLoader string

var AgentApplicationBrowserLoaderTypes = struct {
	// Use PRO instead.
	FULL AgentApplicationBrowserLoader
	// Lite: Gives you information about some basic page load timing and browser user information. Lacks the Browser Pro features and SPA features.
	LITE AgentApplicationBrowserLoader
	// Don't use an agent.
	NONE AgentApplicationBrowserLoader
	// Pro: Gives you access to the Browser Pro features. Lacks the functionality designed for single page app monitoring.
	PRO AgentApplicationBrowserLoader
	// Pro+SPA: This is the default installed agent when you enable browser monitoring. Gives you access to all of the Browser Pro features and to Single Page App (SPA) monitoring. Provides detailed page timing data and the most up-to-date New Relic features, including distributed tracing, for all types of applications.
	SPA AgentApplicationBrowserLoader
}{
	// Use PRO instead.
	FULL: "FULL",
	// Lite: Gives you information about some basic page load timing and browser user information. Lacks the Browser Pro features and SPA features.
	LITE: "LITE",
	// Don't use an agent.
	NONE: "NONE",
	// Pro: Gives you access to the Browser Pro features. Lacks the functionality designed for single page app monitoring.
	PRO: "PRO",
	// Pro+SPA: This is the default installed agent when you enable browser monitoring. Gives you access to all of the Browser Pro features and to Single Page App (SPA) monitoring. Provides detailed page timing data and the most up-to-date New Relic features, including distributed tracing, for all types of applications.
	SPA: "SPA",
}

// AgentApplicationBrowserSettings - The settings of a browser application. Includes loader script.
type AgentApplicationBrowserSettings struct {
	// Configure cookies. The default is enabled: true.
	CookiesEnabled bool `json:"cookiesEnabled"`
	// Configure distributed tracing in browser apps. The default is enabled: true.
	DistributedTracingEnabled bool `json:"distributedTracingEnabled"`
	// The snippet of JavaScript used to copy/paste into your JavaScript app if you aren't using an auto-instrumenting agent on the backend. Note that the resulting snippet will be a JSON string that will need to be parsed before using in your browser application.
	LoaderScript string `json:"loaderScript,omitempty"`
	// Determines which browser loader will be configured. The default is "SPA".
	LoaderType AgentApplicationBrowserLoader `json:"loaderType"`
}

// AgentApplicationBrowserSettingsInput - Configure additional browser settings here.
type AgentApplicationBrowserSettingsInput struct {
	// Configure cookies. The default is enabled: true.
	CookiesEnabled bool `json:"cookiesEnabled,omitempty"`
	// Configure distributed tracing in browser apps. The default is enabled: true.
	DistributedTracingEnabled bool `json:"distributedTracingEnabled,omitempty"`
	// Determines which browser loader is configured. The default is "SPA".
	LoaderType AgentApplicationBrowserLoader `json:"loaderType,omitempty"`
}

// AgentApplicationCreateBrowserResult - The result of creating a browser application.
type AgentApplicationCreateBrowserResult struct {
	// The GUID for the affected Entity.
	GUID common.EntityGUID `json:"guid"`
	// The name of the application.
	Name string `json:"name"`
	// Fields related to browser settings.
	Settings AgentApplicationBrowserSettings `json:"settings,omitempty"`
}

// AgentApplicationDeleteResult - The result of deleting an application.
type AgentApplicationDeleteResult struct {
	// Did the delete succeed?
	Success bool `json:"success"`
}
