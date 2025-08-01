// Code generated by tutone: DO NOT EDIT
package dashboards

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

// DashboardCreateErrorType - Expected error types that can be returned by create operation.
type DashboardCreateErrorType string

var DashboardCreateErrorTypeTypes = struct {
	// Invalid input error.
	INVALID_INPUT DashboardCreateErrorType
}{
	// Invalid input error.
	INVALID_INPUT: "INVALID_INPUT",
}

// DashboardDeleteErrorType - Expected error types that can be returned by delete operation.
type DashboardDeleteErrorType string

var DashboardDeleteErrorTypeTypes = struct {
	// Dashboard not found in the system.
	DASHBOARD_NOT_FOUND DashboardDeleteErrorType
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION DashboardDeleteErrorType
}{
	// Dashboard not found in the system.
	DASHBOARD_NOT_FOUND: "DASHBOARD_NOT_FOUND",
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION: "FORBIDDEN_OPERATION",
}

// DashboardDeleteResultStatus - Result status of delete operation.
type DashboardDeleteResultStatus string

var DashboardDeleteResultStatusTypes = struct {
	// FAILURE.
	FAILURE DashboardDeleteResultStatus
	// SUCCESS.
	SUCCESS DashboardDeleteResultStatus
}{
	// FAILURE.
	FAILURE: "FAILURE",
	// SUCCESS.
	SUCCESS: "SUCCESS",
}

// DashboardLiveURLType - Live URL type.
type DashboardLiveURLType string

var DashboardLiveURLTypeTypes = struct {
	// Dashboard.
	DASHBOARD DashboardLiveURLType
	// Widget.
	WIDGET DashboardLiveURLType
}{
	// Dashboard.
	DASHBOARD: "DASHBOARD",
	// Widget.
	WIDGET: "WIDGET",
}

// DashboardUpdateErrorType - Expected error types that can be returned by update operation.
type DashboardUpdateErrorType string

var DashboardUpdateErrorTypeTypes = struct {
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION DashboardUpdateErrorType
	// Invalid input error.
	INVALID_INPUT DashboardUpdateErrorType
}{
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION: "FORBIDDEN_OPERATION",
	// Invalid input error.
	INVALID_INPUT: "INVALID_INPUT",
}

// DashboardUpdatePageErrorType - Expected error types that can be returned by updatePage operation.
type DashboardUpdatePageErrorType string

var DashboardUpdatePageErrorTypeTypes = struct {
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION DashboardUpdatePageErrorType
	// Invalid input error.
	INVALID_INPUT DashboardUpdatePageErrorType
	// Page not found in the system.
	PAGE_NOT_FOUND DashboardUpdatePageErrorType
}{
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION: "FORBIDDEN_OPERATION",
	// Invalid input error.
	INVALID_INPUT: "INVALID_INPUT",
	// Page not found in the system.
	PAGE_NOT_FOUND: "PAGE_NOT_FOUND",
}

// DashboardUpdateWidgetsInPageErrorType - Expected error types that can be returned by updateWidgetsInPage operation.
type DashboardUpdateWidgetsInPageErrorType string

var DashboardUpdateWidgetsInPageErrorTypeTypes = struct {
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION DashboardUpdateWidgetsInPageErrorType
	// Invalid input error.
	INVALID_INPUT DashboardUpdateWidgetsInPageErrorType
	// Page not found in the system.
	PAGE_NOT_FOUND DashboardUpdateWidgetsInPageErrorType
	// Widget not found in the system.
	WIDGET_NOT_FOUND DashboardUpdateWidgetsInPageErrorType
}{
	// User is not allowed to execute the operation.
	FORBIDDEN_OPERATION: "FORBIDDEN_OPERATION",
	// Invalid input error.
	INVALID_INPUT: "INVALID_INPUT",
	// Page not found in the system.
	PAGE_NOT_FOUND: "PAGE_NOT_FOUND",
	// Widget not found in the system.
	WIDGET_NOT_FOUND: "WIDGET_NOT_FOUND",
}

// DashboardVariableReplacementStrategy - Possible strategies when replacing variables in a NRQL query.
type DashboardVariableReplacementStrategy string

var DashboardVariableReplacementStrategyTypes = struct {
	// Replace the variable based on its automatically-inferred type.
	DEFAULT DashboardVariableReplacementStrategy
	// Replace the variable value as an identifier.
	IDENTIFIER DashboardVariableReplacementStrategy
	// Replace the variable value as a number.
	NUMBER DashboardVariableReplacementStrategy
	// Replace the variable value as a string.
	STRING DashboardVariableReplacementStrategy
}{
	// Replace the variable based on its automatically-inferred type.
	DEFAULT: "DEFAULT",
	// Replace the variable value as an identifier.
	IDENTIFIER: "IDENTIFIER",
	// Replace the variable value as a number.
	NUMBER: "NUMBER",
	// Replace the variable value as a string.
	STRING: "STRING",
}

// DashboardVariableType - Indicates where a variable's possible values may come from.
type DashboardVariableType string

var DashboardVariableTypeTypes = struct {
	// Value comes from an enumerated list of possible values.
	ENUM DashboardVariableType
	// Value comes from the results of a NRQL query.
	NRQL DashboardVariableType
	// Dashboard user can supply an arbitrary string value to variable.
	STRING DashboardVariableType
}{
	// Value comes from an enumerated list of possible values.
	ENUM: "ENUM",
	// Value comes from the results of a NRQL query.
	NRQL: "NRQL",
	// Dashboard user can supply an arbitrary string value to variable.
	STRING: "STRING",
}

// DashboardAreaWidgetConfigurationInput - Configuration for visualization type 'viz.area'
type DashboardAreaWidgetConfigurationInput struct {
	// NRQL queries.
	NRQLQueries []DashboardWidgetNRQLQueryInput `json:"nrqlQueries,omitempty"`
}

// DashboardBarWidgetConfigurationInput - Configuration for visualization type 'viz.bar'. Learn more about [bar](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#bar) widget.
type DashboardBarWidgetConfigurationInput struct {
	// NRQL queries.
	NRQLQueries []DashboardWidgetNRQLQueryInput `json:"nrqlQueries,omitempty"`
}

// DashboardBillboardWidgetConfigurationInput - Configuration for visualization type 'viz.billboard'. Learn more about [billboard](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#billboard) widget.
type DashboardBillboardWidgetConfigurationInput struct {
	// NRQL queries.
	NRQLQueries []DashboardWidgetNRQLQueryInput `json:"nrqlQueries,omitempty"`
	// Array of thresholds to categorize the results of the query in different groups.
	Thresholds []DashboardBillboardWidgetThresholdInput `json:"thresholds,omitempty"`
}

// DashboardCreateError - Expected errors that can be returned by create operation.
type DashboardCreateError struct {
	// Error description.
	Description string `json:"description,omitempty"`
	// Error type.
	Type DashboardCreateErrorType `json:"type"`
}

// DashboardCreateResult - Result of create operation.
type DashboardCreateResult struct {
	// Dashboard creation result.
	EntityResult DashboardEntityResult `json:"entityResult,omitempty"`
	// Expected errors while processing request.
	Errors []DashboardCreateError `json:"errors,omitempty"`
}

// DashboardDeleteError - Expected error types that can be returned by delete operation.
type DashboardDeleteError struct {
	// Error description.
	Description string `json:"description,omitempty"`
	// Error type.
	Type DashboardDeleteErrorType `json:"type"`
}

// DashboardDeleteResult - Result of delete operation.
type DashboardDeleteResult struct {
	// Expected errors while processing request.
	Errors []DashboardDeleteError `json:"errors,omitempty"`
	// The status of the attempted delete.
	Status DashboardDeleteResultStatus `json:"status,omitempty"`
}

// DashboardEntityResult - Public schema - `DashboardEntity` result representation for mutations. It's a subset of the `DashboardEntity` that inherits from the Entity type, but a complete different type.
type DashboardEntityResult struct {
	// The New Relic account where the dashboard is created.
	AccountID int `json:"accountId,omitempty"`
	// Dashboard creation timestamp.
	CreatedAt nrtime.DateTime `json:"createdAt,omitempty"`
	// Brief text describing the dashboard.
	Description string `json:"description,omitempty"`
	// Unique entity identifier.
	GUID common.EntityGUID `json:"guid,omitempty"`
	// The name of the dashboard.
	Name string `json:"name,omitempty"`
	// Information of the user that owns the dashboard.
	Owner entities.DashboardOwnerInfo `json:"owner,omitempty"`
	// A nested block of all pages belonging to the dashboard.
	Pages []entities.DashboardPage `json:"pages,omitempty"`
	// Dashboard permissions configuration.
	Permissions entities.DashboardPermissions `json:"permissions,omitempty"`
	// Dashboard update timestamp.
	UpdatedAt nrtime.DateTime `json:"updatedAt,omitempty"`
	// Dashboard-local variable definitions.
	Variables []entities.DashboardVariable `json:"variables,omitempty"`
}

// DashboardInput - Dashboard input.
type DashboardInput struct {
	// Brief text describing the dashboard.
	Description string `json:"description,omitempty"`
	// The name of the dashboard.
	Name string `json:"name"`
	// A nested block of all pages belonging to the dashboard.
	Pages []DashboardPageInput `json:"pages,omitempty"`
	// Permissions to set level of visibility & editing.
	Permissions entities.DashboardPermissions `json:"permissions"`
	// Dashboard-local variable definitions.
	Variables []DashboardVariableInput `json:"variables,omitempty"`
}

// DashboardLineWidgetConfigurationInput - Configuration for visualization type 'viz.line'. Learn more about [line](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#line) widget.
type DashboardLineWidgetConfigurationInput struct {
	// NRQL queries.
	NRQLQueries []DashboardWidgetNRQLQueryInput `json:"nrqlQueries,omitempty"`
}

// DashboardLiveURL - Live URL.
type DashboardLiveURL struct {
	// Creation date.
	CreatedAt nrtime.EpochMilliseconds `json:"createdAt,omitempty"`
	// Title that describes the source entity that is accessible through the public live URL.
	Title string `json:"title,omitempty"`
	// Live URL type.
	Type DashboardLiveURLType `json:"type,omitempty"`
	// Public URL.
	URL string `json:"url"`
	// The unique identifier of the public live URL.
	Uuid string `json:"uuid"`
}

// DashboardMarkdownWidgetConfigurationInput - Configuration for visualization type 'viz.markdown'. Learn more about [markdown](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#markdown) widget.
type DashboardMarkdownWidgetConfigurationInput struct {
	// Markdown content of the widget.
	Text string `json:"text"`
}

// DashboardPageInput - Page input.
type DashboardPageInput struct {
	// Brief text describing the page.
	Description string `json:"description,omitempty"`
	// Unique entity identifier of the Page to be updated. When null, it means a new Page will be created.
	GUID common.EntityGUID `json:"guid,omitempty"`
	// The name of the page.
	Name string `json:"name"`

	// NOTE: The JSON description of the following attribute, "Widgets" has been modified manually
	// (removal of "omitempty") to facilitate creating pages with no widgets (empty pages).
	// Please DO NOT regenerate/modify this attribute and its datatype via Tutone (which would add "omitempty" back).

	// A nested block of all widgets belonging to the page.
	Widgets []DashboardWidgetInput `json:"widgets"`
}

// DashboardPieWidgetConfigurationInput - Configuration for visualization type 'viz.pie'.  Learn more about [pie](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#pie) widget.
type DashboardPieWidgetConfigurationInput struct {
	// NRQL queries.
	NRQLQueries []DashboardWidgetNRQLQueryInput `json:"nrqlQueries,omitempty"`
}

// DashboardSnapshotURLInput - Parameters that affect the data and the rendering of the dashboards returned by the snapshot url mutation.
type DashboardSnapshotURLInput struct {
	// Period of time from which the data to be displayed on the dashboard will be obtained.
	TimeWindow DashboardSnapshotURLTimeWindowInput `json:"timeWindow,omitempty"`
}

// DashboardSnapshotURLTimeWindowInput - Period of time from which the data to be displayed on the dashboard will be obtained.
type DashboardSnapshotURLTimeWindowInput struct {
	// The starting time of the time window. If specified, an endTime or a duration must also be specified.
	BeginTime nrtime.EpochMilliseconds `json:"beginTime,omitempty"`
	// The duration of the time window.
	Duration nrtime.Milliseconds `json:"duration,omitempty"`
	// The end time of the time window. If specified, a beginTime or a duration must also be specified.
	EndTime nrtime.EpochMilliseconds `json:"endTime,omitempty"`
}

// DashboardTableWidgetConfigurationInput - Configuration for visualization type 'viz.table'.  Learn more about [table](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#table) widget.
type DashboardTableWidgetConfigurationInput struct {
	// NRQL queries.
	NRQLQueries []DashboardWidgetNRQLQueryInput `json:"nrqlQueries,omitempty"`
}

// DashboardUpdateError - Expected errors that can be returned by update operation.
type DashboardUpdateError struct {
	// Error description.
	Description string `json:"description,omitempty"`
	// Error type.
	Type DashboardUpdateErrorType `json:"type"`
}

// DashboardUpdatePageError - Expected errors that can be returned by updatePage operation.
type DashboardUpdatePageError struct {
	// Error description.
	Description string `json:"description,omitempty"`
	// Error type.
	Type DashboardUpdatePageErrorType `json:"type"`
}

// DashboardUpdatePageInput - Page input used when updating an individual page.
type DashboardUpdatePageInput struct {
	// Page description.
	Description string `json:"description,omitempty"`
	// Page name.
	Name string `json:"name"`

	// NOTE: The JSON description of the following attribute, "Widgets" has been modified manually
	// (removal of "omitempty") to facilitate creating pages with no widgets (empty pages).
	// Please DO NOT regenerate/modify this attribute and its datatype via Tutone (which would add "omitempty" back).

	// Page widgets.
	Widgets []DashboardWidgetInput `json:"widgets"`
}

// DashboardUpdatePageResult - Result of updatePage operation.
type DashboardUpdatePageResult struct {
	// Expected errors while processing request. No errors means successful request.
	Errors []DashboardUpdatePageError `json:"errors,omitempty"`
}

// DashboardUpdateResult - Result of update operation.
type DashboardUpdateResult struct {
	// Dashboard update result.
	EntityResult DashboardEntityResult `json:"entityResult,omitempty"`
	// Expected errors while processing request.
	Errors []DashboardUpdateError `json:"errors,omitempty"`
}

// DashboardUpdateWidgetInput - Input type used when updating widgets.
type DashboardUpdateWidgetInput struct {
	// Typed widgets are area, bar, billboard, line, markdown, pie, and table. Check our [docs](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#widget-typed) for more info.
	Configuration DashboardWidgetConfigurationInput `json:"configuration,omitempty"`
	// ID of the widget to be updated.
	ID string `json:"id"`
	// The widget's position and size in the dashboard.
	Layout DashboardWidgetLayoutInput `json:"layout,omitempty"`
	// Entities related to the widget. Currently only supports one Dashboard entity guid, but may allow other cases in the future.
	LinkedEntityGUIDs []common.EntityGUID `json:"linkedEntityGuids"`
	// Untyped widgets are all other widgets, such as bullet, histogram, inventory, etc. Check our [docs](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#widget-untyped) for more info.
	RawConfiguration entities.DashboardWidgetRawConfiguration `json:"rawConfiguration,omitempty"`
	// A title for the widget.
	Title string `json:"title,omitempty"`
	// Specifies how this widget will be visualized. If null, the WidgetConfigurationInput will be used to determine the visualization.
	Visualization DashboardWidgetVisualizationInput `json:"visualization,omitempty"`
}

// DashboardUpdateWidgetsInPageError - Expected errors that can be returned by updateWidgetsInPage operation.
type DashboardUpdateWidgetsInPageError struct {
	// Error description.
	Description string `json:"description,omitempty"`
	// Error type.
	Type DashboardUpdateWidgetsInPageErrorType `json:"type"`
}

// DashboardUpdateWidgetsInPageResult - Result of updateWidgetsInPage operation.
type DashboardUpdateWidgetsInPageResult struct {
	// Expected errors while processing request. No errors means successful request.
	Errors []DashboardUpdateWidgetsInPageError `json:"errors,omitempty"`
}

// DashboardVariable - Definition of a variable that is local to this dashboard. Variables are placeholders for dynamic values in widget NRQLs.
type DashboardVariable struct {
	// [DEPRECATED] Default value for this variable. The actual value to be used will depend on the type.
	DefaultValue DashboardVariableDefaultValue `json:"defaultValue,omitempty"`
	// Default values for this variable. The actual value to be used will depend on the type.
	DefaultValues []DashboardVariableDefaultItem `json:"defaultValues,omitempty"`
	// Indicates whether this variable supports multiple selection or not. Only applies to variables of type NRQL or ENUM.
	IsMultiSelection bool `json:"isMultiSelection,omitempty"`
	// List of possible values for variables of type ENUM.
	Items []DashboardVariableEnumItem `json:"items,omitempty"`
	// Configuration for variables of type NRQL.
	NRQLQuery DashboardVariableNRQLQuery `json:"nrqlQuery,omitempty"`
	// Variable identifier.
	Name string `json:"name,omitempty"`
	// Options applied to the variable.
	Options DashboardVariableOptions `json:"options,omitempty"`
	// Indicates the strategy to apply when replacing a variable in a NRQL query.
	ReplacementStrategy DashboardVariableReplacementStrategy `json:"replacementStrategy,omitempty"`
	// Human-friendly display string for this variable.
	Title string `json:"title,omitempty"`
	// Specifies the data type of the variable and where its possible values may come from.
	Type DashboardVariableType `json:"type,omitempty"`
}

// DashboardVariableDefaultItem - Represents a possible default value item.
type DashboardVariableDefaultItem struct {
	// The value of this default item.
	Value DashboardVariableDefaultValue `json:"value,omitempty"`
}

// DashboardVariableDefaultItemInput - Represents a possible default value item.
type DashboardVariableDefaultItemInput struct {
	// The value of this default item.
	Value DashboardVariableDefaultValueInput `json:"value,omitempty"`
}

// DashboardVariableDefaultValue - Specifies a default value for variables.
type DashboardVariableDefaultValue struct {
	// Default string value.
	String string `json:"string,omitempty"`
}

// DashboardVariableDefaultValueInput - Specifies a default value for variables.
type DashboardVariableDefaultValueInput struct {
	// Default string value.
	String string `json:"string,omitempty"`
}

// DashboardVariableEnumItem - Represents a possible value for a variable of type ENUM.
type DashboardVariableEnumItem struct {
	// A human-friendly display string for this value.
	Title string `json:"title,omitempty"`
	// A possible variable value.
	Value string `json:"value,omitempty"`
}

// DashboardVariableEnumItemInput - Input type that represents a possible value for a variable of type ENUM.
type DashboardVariableEnumItemInput struct {
	// A human-friendly display string for this value.
	Title string `json:"title,omitempty"`
	// A possible variable value
	Value string `json:"value"`
}

// DashboardVariableInput - Definition of a variable that is local to this dashboard. Variables are placeholders for dynamic values in widget NRQLs.
type DashboardVariableInput struct {
	// [DEPRECATED] Default value for this variable. The actual value to be used will depend on the type.
	DefaultValue *DashboardVariableDefaultValueInput `json:"defaultValue,omitempty"`
	// Default values for this variable. The actual value to be used will depend on the type.
	DefaultValues *[]DashboardVariableDefaultItemInput `json:"defaultValues,omitempty"`
	// Indicates whether this variable supports multiple selection or not. Only applies to variables of type NRQL or ENUM.
	IsMultiSelection bool `json:"isMultiSelection,omitempty"`
	// List of possible values for variables of type ENUM
	Items []DashboardVariableEnumItemInput `json:"items,omitempty"`
	// Configuration for variables of type NRQL.
	NRQLQuery *DashboardVariableNRQLQueryInput `json:"nrqlQuery,omitempty"`
	// Variable identifier.
	Name string `json:"name"`
	// Options applied to the variable
	Options *DashboardVariableOptionsInput `json:"options,omitempty"`
	// Indicates the strategy to apply when replacing a variable in a NRQL query.
	ReplacementStrategy DashboardVariableReplacementStrategy `json:"replacementStrategy,omitempty"`
	// Human-friendly display string for this variable.
	Title string `json:"title,omitempty"`
	// Specifies the data type of the variable and where its possible values may come from.
	Type DashboardVariableType `json:"type"`
}

// DashboardVariableNRQLQuery - Configuration for variables of type NRQL.
type DashboardVariableNRQLQuery struct {
	// New Relic account ID(s) to issue the query against.
	AccountIDs []int `json:"accountIds,omitempty"`
	// NRQL formatted query.
	Query nrdb.NRQL `json:"query"`
}

// DashboardVariableNRQLQueryInput - Configuration for variables of type NRQL.
type DashboardVariableNRQLQueryInput struct {
	// New Relic account ID(s) to issue the query against.
	AccountIDs []int `json:"accountIds"`
	// NRQL formatted query.
	Query nrdb.NRQL `json:"query"`
}

// DashboardVariableOptions - Options applied to the variable.
type DashboardVariableOptions struct {
	// With this turned on, query condition defined with the variable will not be included in the query.
	Excluded bool `json:"excluded,omitempty"`
	// Only applies to variables of type NRQL. With this turned on, the time range for the NRQL query will override the time picker on dashboards and other pages. Turn this off to use the time picker as normal.
	IgnoreTimeRange bool `json:"ignoreTimeRange,omitempty"`
}

// DashboardVariableOptionsInput - Options applied to the variable
type DashboardVariableOptionsInput struct {
	// With this turned on, query condition defined with the variable will not be included in the query.
	Excluded *bool `json:"excluded,omitempty"`
	// Only applies to variables of type NRQL. With this turned on, the time range for the NRQL query will override the time picker on dashboards and other pages. Turn this off to use the time picker as normal.
	IgnoreTimeRange *bool `json:"ignoreTimeRange,omitempty"`
}

// DashboardWidgetConfigurationInput - Typed configuration for known visualizations. At most one may be populated.
type DashboardWidgetConfigurationInput struct {
	// Configuration for visualization type 'viz.area'
	Area *DashboardAreaWidgetConfigurationInput `json:"area,omitempty"`
	// Configuration for visualization type 'viz.bar'
	Bar *DashboardBarWidgetConfigurationInput `json:"bar,omitempty"`
	// Configuration for visualization type 'viz.billboard'
	Billboard *DashboardBillboardWidgetConfigurationInput `json:"billboard,omitempty"`
	// Configuration for visualization type 'viz.line'
	Line *DashboardLineWidgetConfigurationInput `json:"line,omitempty"`
	// Configuration for visualization type 'viz.markdown'
	Markdown *DashboardMarkdownWidgetConfigurationInput `json:"markdown,omitempty"`
	// Configuration for visualization type 'viz.pie'
	Pie *DashboardPieWidgetConfigurationInput `json:"pie,omitempty"`
	// Configuration for visualization type 'viz.table'
	Table *DashboardTableWidgetConfigurationInput `json:"table,omitempty"`
}

// DashboardWidgetInput - Widget input.
type DashboardWidgetInput struct {
	// Typed widgets are area, bar, billboard, line, markdown, pie, and table. Check our [docs](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#widget-typed) for more info.
	Configuration DashboardWidgetConfigurationInput `json:"configuration,omitempty"`
	// ID of the widget. If null, a new widget will be created and added to a dashboard.
	ID string `json:"id,omitempty"`
	// The widget's position and size in the dashboard.
	Layout DashboardWidgetLayoutInput `json:"layout,omitempty"`
	// Entities related to the widget. Currently only supports one Dashboard entity guid, but may allow other cases in the future.
	LinkedEntityGUIDs []common.EntityGUID `json:"linkedEntityGuids"`
	// Untyped widgets are all other widgets, such as bullet, histogram, inventory, etc. Check our [docs](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#widget-untyped) for more info.
	RawConfiguration entities.DashboardWidgetRawConfiguration `json:"rawConfiguration,omitempty"`
	// A title for the widget.
	Title string `json:"title,omitempty"`
	// Specifies how this widget will be visualized. If null, the WidgetConfigurationInput will be used to determine the visualization.
	Visualization DashboardWidgetVisualizationInput `json:"visualization,omitempty"`
}

// DashboardWidgetLayoutInput - Widget layout input.
type DashboardWidgetLayoutInput struct {
	// Column position of widget from top left, starting at 1.
	Column int `json:"column,omitempty"`
	// Height of the widget. Valid values are 1 to 12 inclusive. Defaults to 3.
	Height int `json:"height,omitempty"`
	// Row position of widget from top left, starting at 1.
	Row int `json:"row,omitempty"`
	// Width of the widget. Valid values are 1 to 12 inclusive. Defaults to 4.
	Width int `json:"width,omitempty"`
}

// DashboardWidgetNRQLQueryInput - NRQL query used by a widget.
type DashboardWidgetNRQLQueryInput struct {
	// New Relic account ID to issue the query against.
	AccountID int `json:"accountId,omitempty"`
	// New Relic account IDs to issue the query against.
	AccountIDS []int `json:"accountIds,omitempty"`
	// NRQL formatted query.
	Query nrdb.NRQL `json:"query"`
}

// DashboardWidgetVisualizationInput - Visualization configuration.
type DashboardWidgetVisualizationInput struct {
	// This field can either have a known type like `viz.area` or `<nerdpack-id>.<visualization-id>` in the case of [custom visualizations](https://developer.newrelic.com/explore-docs/custom-viz/build-visualization/). Check out [docs](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/#widget-schema) for more info.
	ID string `json:"id,omitempty"`
}
