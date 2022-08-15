// DashboardBillboardWidgetThresholdInput - Billboard widget threshold input.
package dashboards

import "github.com/newrelic/newrelic-client-go/pkg/entities"

type DashboardBillboardWidgetThresholdInput struct {
	// alert severity.
	AlertSeverity entities.DashboardAlertSeverity `json:"alertSeverity,omitempty"`
	// value.
	Value *float64 `json:"value,omitempty"`
}

type RawConfiguration struct {
	// Used by all widgets
	NRQLQueries     []DashboardWidgetNRQLQueryInput  `json:"nrqlQueries,omitempty"`
	PlatformOptions *RawConfigurationPlatformOptions `json:"platformOptions,omitempty"`

	// Used by viz.bullet
	Limit float64 `json:"limit,omitempty"`

	// Used by viz.markdown
	Text string `json:"text,omitempty"`

	// Used by viz.billboard
	Thresholds []DashboardBillboardWidgetThresholdInput `json:"thresholds,omitempty"`
}

type RawConfigurationPlatformOptions struct {
	IgnoreTimeRange bool `json:"ignoreTimeRange,omitempty"`
}
