// Code generated by tutone: DO NOT EDIT
package logconfigurations

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

// LogConfigurationsObfuscationMethod - Methods for replacing obfuscated values.
type LogConfigurationsObfuscationMethod string

var LogConfigurationsObfuscationMethodTypes = struct {
	// Replace the matched data with a SHA256 hash.
	HASH_SHA256 LogConfigurationsObfuscationMethod
	// Replace the matched data with a static value.
	MASK LogConfigurationsObfuscationMethod
}{
	// Replace the matched data with a SHA256 hash.
	HASH_SHA256: "HASH_SHA256",
	// Replace the matched data with a static value.
	MASK: "MASK",
}

// LogConfigurationsParsingRuleMutationErrorType - Expected default error types as result of mutating an existing parsing rule.
type LogConfigurationsParsingRuleMutationErrorType string

var LogConfigurationsParsingRuleMutationErrorTypeTypes = struct {
	// Invalid Grok
	INVALID_GROK LogConfigurationsParsingRuleMutationErrorType
	// Number format error. ID should be convertible to int.
	INVALID_ID LogConfigurationsParsingRuleMutationErrorType
	// Invalid NRQL
	INVALID_NRQL LogConfigurationsParsingRuleMutationErrorType
	// Couldn't find the specified parsing rule.
	NOT_FOUND LogConfigurationsParsingRuleMutationErrorType
}{
	// Invalid Grok
	INVALID_GROK: "INVALID_GROK",
	// Number format error. ID should be convertible to int.
	INVALID_ID: "INVALID_ID",
	// Invalid NRQL
	INVALID_NRQL: "INVALID_NRQL",
	// Couldn't find the specified parsing rule.
	NOT_FOUND: "NOT_FOUND",
}

// Account - The `Account` object provides general data about the account, as well as
// being the entry point into more detailed data about a single account.
//
// Account configuration data is queried through this object, as well as
// telemetry data that is specific to a single account.
type Account struct {
	//
	ID int `json:"id,omitempty"`
	//
	LicenseKey string `json:"licenseKey,omitempty"`
	// This field provides access to LogConfigurations data.
	LogConfigurations LogConfigurationsAccountStitchedFields `json:"logConfigurations,omitempty"`
	//
	Name string `json:"name,omitempty"`
}

// Actor - The `Actor` object contains fields that are scoped to the API user's access level.
type Actor struct {
	// The `account` field is the entry point into data that is scoped to a single account.
	Account Account `json:"account,omitempty"`
}

// LogConfigurationsAccountStitchedFields -
type LogConfigurationsAccountStitchedFields struct {
	// Look up for all obfuscation expressions for a given account
	ObfuscationExpressions []LogConfigurationsObfuscationExpression `json:"obfuscationExpressions"`
	// Look up for all obfuscation rules for a given account.
	ObfuscationRules []LogConfigurationsObfuscationRule `json:"obfuscationRules"`
	// Look up for all parsing rules for a given account.
	ParsingRules []*LogConfigurationsParsingRule `json:"parsingRules"`
	// Test a Grok pattern against a list of log lines.
	TestGrok []LogConfigurationsGrokTestResult `json:"testGrok"`
}

// LogConfigurationsCreateObfuscationActionInput - Input for creating an obfuscation action on a rule being created.
type LogConfigurationsCreateObfuscationActionInput struct {
	// Attribute names for action. An empty list applies the action to all the attributes.
	Attributes []string `json:"attributes"`
	// Expression Id for action.
	ExpressionId string `json:"expressionId"`
	// Obfuscation method to use.
	Method LogConfigurationsObfuscationMethod `json:"method"`
}

// LogConfigurationsCreateObfuscationExpressionInput - Input for creating an obfuscation expression.
type LogConfigurationsCreateObfuscationExpressionInput struct {
	// Description of expression.
	Description string `json:"description,omitempty"`
	// Name of expression.
	Name string `json:"name"`
	// Regex of expression.
	Regex string `json:"regex"`
}

// LogConfigurationsCreateObfuscationRuleInput - Input for creating an obfuscation rule.
type LogConfigurationsCreateObfuscationRuleInput struct {
	// Actions for the rule. The actions will be applied in the order specified by this list.
	Actions []LogConfigurationsCreateObfuscationActionInput `json:"actions,omitempty"`
	// Description of rule.
	Description string `json:"description,omitempty"`
	// Whether the rule should be applied or not to incoming data.
	Enabled bool `json:"enabled"`
	// NRQL for determining whether a given log record should have obfuscation actions applied.
	Filter NRQL `json:"filter"`
	// Name of rule.
	Name string `json:"name"`
}

// LogConfigurationsCreateParsingRuleResponse - The result after creating a new parsing rule.
type LogConfigurationsCreateParsingRuleResponse struct {
	// List of errors, if any.
	Errors []LogConfigurationsParsingRuleMutationError `json:"errors,omitempty"`
	// The created parsing rule.
	Rule *LogConfigurationsParsingRule `json:"rule,omitempty"`
}

// LogConfigurationsDeleteParsingRuleResponse - The result after deleting a parsing rule.
type LogConfigurationsDeleteParsingRuleResponse struct {
	// List of errors, if any.
	Errors []LogConfigurationsParsingRuleMutationError `json:"errors,omitempty"`
}

// LogConfigurationsGrokTestExtractedAttribute - An attribute that was extracted from a Grok test.
type LogConfigurationsGrokTestExtractedAttribute struct {
	// The attribute name.
	Name string `json:"name"`
	// A string representation of the extracted value (which might not be a String).
	Value string `json:"value"`
}

// LogConfigurationsGrokTestResult - The result of testing Grok on a log line.
type LogConfigurationsGrokTestResult struct {
	// Any attributes that were extracted.
	Attributes []LogConfigurationsGrokTestExtractedAttribute `json:"attributes"`
	// The log line that was tested against.
	LogLine string `json:"logLine"`
	// Whether the Grok pattern matched.
	Matched bool `json:"matched"`
}

// LogConfigurationsObfuscationAction - Application of an obfuscation expression with specific a replacement method.
type LogConfigurationsObfuscationAction struct {
	// Log record attributes to apply this expression to. An empty list applies the action to all the attributes.
	Attributes []string `json:"attributes"`
	// Obfuscation expression applied by this action.
	Expression LogConfigurationsObfuscationExpression `json:"expression"`
	// The id of the obfuscation action.
	ID string `json:"id"`
	// How to obfuscate matches for the applied expression.
	Method LogConfigurationsObfuscationMethod `json:"method"`
}

// LogConfigurationsObfuscationExpression - Reusable obfuscation expression.
type LogConfigurationsObfuscationExpression struct {
	// Identifies the date and time when the expression was created.
	CreatedAt nrtime.DateTime `json:"createdAt"`
	// Identifies the user who has created the expression.
	CreatedBy UserReference `json:"createdBy,omitempty"`
	// Description of the expression.
	Description string `json:"description,omitempty"`
	// The id of the obfuscation expression.
	ID string `json:"id"`
	// Name of the expression.
	Name string `json:"name"`
	// Regular expression for this obfuscation expression. Capture groups will be obscured on matching.
	Regex string `json:"regex"`
	// Identifies the date and time when the expression was last updated.
	UpdatedAt nrtime.DateTime `json:"updatedAt"`
	// Identifies the user who has last updated the expression.
	UpdatedBy UserReference `json:"updatedBy,omitempty"`
}

// LogConfigurationsObfuscationRule - Rule for identifying a set of log data to apply specific obfuscation actions to.
type LogConfigurationsObfuscationRule struct {
	// Obfuscation actions to take if a record passes the matching criteria.
	Actions []LogConfigurationsObfuscationAction `json:"actions"`
	// Identifies the date and time when the rule was created.
	CreatedAt nrtime.DateTime `json:"createdAt"`
	// Identifies the user who has created the rule.
	CreatedBy UserReference `json:"createdBy,omitempty"`
	// Description of the obfuscation rule.
	Description string `json:"description,omitempty"`
	// Whether the rule should be applied to incoming logs
	Enabled bool `json:"enabled"`
	// NRQL filter to determine if a log record should have obfuscation actions applied.
	Filter NRQL `json:"filter"`
	// The id of the obfuscation rule.
	ID string `json:"id"`
	// Name of the obfuscation rule.
	Name string `json:"name"`
	// Identifies the date and time when the rule was last updated.
	UpdatedAt nrtime.DateTime `json:"updatedAt"`
	// Identifies the user who has last updated the rule.
	UpdatedBy UserReference `json:"updatedBy,omitempty"`
}

// LogConfigurationsParsingRule - A parsing rule for an account.
type LogConfigurationsParsingRule struct {
	// The account id associated with the rule.
	AccountID int `json:"accountId"`
	// The parsing rule will apply to value of this attribute.
	Attribute string `json:"attribute"`
	// Identifies the user who has created the rule.
	CreatedBy UserReference `json:"createdBy,omitempty"`
	// Whether or not this rule is deleted.
	Deleted bool `json:"deleted"`
	// A description of what this parsing rule represents.
	Description string `json:"description"`
	// Whether or not this rule is enabled.
	Enabled bool `json:"enabled"`
	// The Grok of what to parse.
	Grok string `json:"grok"`
	// Unique parsing rule identifier.
	ID string `json:"id"`
	// The Lucene to match events to the parsing rule.
	Lucene string `json:"lucene"`
	// The NRQL to match events to the parsing rule.
	NRQL NRQL `json:"nrql"`
	// Identifies the date and time when the rule was last updated.
	UpdatedAt nrtime.DateTime `json:"updatedAt,omitempty"`
	// Identifies the user who has last updated the rule.
	UpdatedBy UserReference `json:"updatedBy,omitempty"`
}

// LogConfigurationsParsingRuleConfiguration - A new parsing rule.
type LogConfigurationsParsingRuleConfiguration struct {
	// The parsing rule will apply to value of this attribute. If field is not provided, value will default to message.
	Attribute string `json:"attribute,omitempty"`
	// A description of what this parsing rule represents.
	Description string `json:"description"`
	// Whether or not this rule is enabled.
	Enabled bool `json:"enabled"`
	// The Grok of what to parse.
	Grok string `json:"grok"`
	// The Lucene to match events to the parsing rule.
	Lucene string `json:"lucene"`
	// The NRQL to match events to the parsing rule.
	NRQL NRQL `json:"nrql"`
}

// LogConfigurationsParsingRuleMutationError - Expected errors as a result of mutating a parsing rule.
type LogConfigurationsParsingRuleMutationError struct {
	// The message with the error cause.
	Message string `json:"message,omitempty"`
	// Type of error.
	Type LogConfigurationsParsingRuleMutationErrorType `json:"type,omitempty"`
}

// LogConfigurationsUpdateObfuscationActionInput - Input for creating an obfuscation action on a rule being updated.
type LogConfigurationsUpdateObfuscationActionInput struct {
	// Attribute names for action. An empty list applies the action to all the attributes.
	Attributes []string `json:"attributes"`
	// Expression Id for action.
	ExpressionId string `json:"expressionId"`
	// Obfuscation method to use.
	Method LogConfigurationsObfuscationMethod `json:"method"`
}

// LogConfigurationsUpdateObfuscationExpressionInput - Input for updating an obfuscation expression.
// Null fields are left untouched by mutation.
type LogConfigurationsUpdateObfuscationExpressionInput struct {
	// Description of expression.
	Description string `json:"description,omitempty"`
	// Expression Id.
	ID string `json:"id"`
	// Name of expression.
	Name string `json:"name,omitempty"`
	// Regex of expression.
	Regex string `json:"regex,omitempty"`
}

// LogConfigurationsUpdateObfuscationRuleInput - Input for updating an obfuscation rule.
// Null fields are left untouched by mutation.
type LogConfigurationsUpdateObfuscationRuleInput struct {
	// Actions for the rule. When non-null, this list of actions is used to replace
	// the existing list of actions of the rule. The actions will be applied in the
	// order specified by this list.
	Actions []LogConfigurationsUpdateObfuscationActionInput `json:"actions,omitempty"`
	// Description of rule.
	Description string `json:"description,omitempty"`
	// Whether the rule should be applied or not to incoming data.
	Enabled bool `json:"enabled,omitempty"`
	// NRQL for determining whether a given log record should have obfuscation actions applied.
	Filter NRQL `json:"filter,omitempty"`
	// Rule Id.
	ID string `json:"id"`
	// Name of rule.
	Name string `json:"name,omitempty"`
}

// LogConfigurationsUpdateParsingRuleResponse - The result after updating a parsing rule.
type LogConfigurationsUpdateParsingRuleResponse struct {
	// List of errors, if any.
	Errors []LogConfigurationsParsingRuleMutationError `json:"errors,omitempty"`
	// The updated parsing rule.
	Rule *LogConfigurationsParsingRule `json:"rule,omitempty"`
}

// UserReference - The `UserReference` object provides basic identifying information about the user.
type UserReference struct {
	//
	Email string `json:"email,omitempty"`
	//
	Gravatar string `json:"gravatar,omitempty"`
	//
	ID int `json:"id,omitempty"`
	//
	Name string `json:"name,omitempty"`
}

type obfuscationExpressionsResponse struct {
	Actor Actor `json:"actor"`
}

type obfuscationRulesResponse struct {
	Actor Actor `json:"actor"`
}

type parsingRulesResponse struct {
	Actor Actor `json:"actor"`
}

type testGrokResponse struct {
	Actor Actor `json:"actor"`
}

// NRQL - This scalar represents a NRQL query string.
//
// See the [NRQL Docs](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions) for more information about NRQL syntax.
type NRQL string
