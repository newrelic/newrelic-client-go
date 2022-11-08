// Code generated by tutone: DO NOT EDIT
package servicelevel

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/accounts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

// ServiceLevelEventsQuerySelectFunction - The function to use in the SELECT clause.
type ServiceLevelEventsQuerySelectFunction string

var ServiceLevelEventsQuerySelectFunctionTypes = struct {
	// COUNT.
	COUNT ServiceLevelEventsQuerySelectFunction
	// SUM.
	SUM ServiceLevelEventsQuerySelectFunction
}{
	// COUNT.
	COUNT: "COUNT",
	// SUM.
	SUM: "SUM",
}

// ServiceLevelObjectiveRollingTimeWindowUnit - The rolling time window units.
type ServiceLevelObjectiveRollingTimeWindowUnit string

var ServiceLevelObjectiveRollingTimeWindowUnitTypes = struct {
	// Day.
	DAY ServiceLevelObjectiveRollingTimeWindowUnit
}{
	// Day.
	DAY: "DAY",
}

// Actor - The `Actor` object contains fields that are scoped to the API user's access level.
type Actor struct {
	// Fetch a single entity.
	//
	// For more details on entities, visit our [entity docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/use-new-relic-graphql-api-query-entities).
	Entity EntityInterface `json:"entity,omitempty"`
}

// Entity - The `Entity` interface allows fetching detailed entity information for a single entity.
//
// To understand more about entities and entity types, look at [our docs](https://docs.newrelic.com/docs/what-are-new-relic-entities).
type Entity struct {
	// The New Relic account ID associated with this entity.
	AccountID int `json:"accountId,omitempty"`
	// The entity's domain
	Domain string `json:"domain,omitempty"`
	// The name of this entity.
	Name string `json:"name,omitempty"`
	// The url to the entity.
	Permalink string `json:"permalink,omitempty"`
	// The service level defined for the entity.
	ServiceLevel ServiceLevelDefinition `json:"serviceLevel,omitempty"`
	// The entity's type
	Type string `json:"type,omitempty"`
}

// ServiceLevelDefinition - The service level defined for a specific entity.
type ServiceLevelDefinition struct {
	// The SLIs attached to the entity.
	Indicators []ServiceLevelIndicator `json:"indicators"`
}

// ServiceLevelEvents - The events that define the SLI.
type ServiceLevelEvents struct {
	// The New Relic account to fetch the events from.
	Account accounts.AccountReference `json:"account,omitempty"`
	// The definition of bad events.
	BadEvents *ServiceLevelEventsQuery `json:"badEvents,omitempty"`
	// The definition of good events.
	GoodEvents *ServiceLevelEventsQuery `json:"goodEvents,omitempty"`
	// The definition of valid events.
	ValidEvents *ServiceLevelEventsQuery `json:"validEvents"`
}

// ServiceLevelEventsCreateInput - The events that define the SLI.
type ServiceLevelEventsCreateInput struct {
	// The New Relic account ID where the events are fetched from.
	AccountID int `json:"accountId"`
	// The definition of bad events.
	BadEvents *ServiceLevelEventsQueryCreateInput `json:"badEvents,omitempty"`
	// The definition of good events.
	GoodEvents *ServiceLevelEventsQueryCreateInput `json:"goodEvents,omitempty"`
	// The definition of valid events.
	ValidEvents *ServiceLevelEventsQueryCreateInput `json:"validEvents,omitempty"`
}

// ServiceLevelEventsQuery - The query that represents the events to fetch.
type ServiceLevelEventsQuery struct {
	// The NRDB event to fetch the data from.
	From NRQL `json:"from"`
	// The NRQL SELECT clause to aggregate events.
	Select ServiceLevelEventsQuerySelect `json:"select,omitempty"`
	// The NRQL condition to filter the events.
	Where NRQL `json:"where,omitempty"`
}

// ServiceLevelEventsQueryCreateInput - The query that represents the events to fetch.
type ServiceLevelEventsQueryCreateInput struct {
	// The NRDB event to fetch the data from.
	From NRQL `json:"from"`
	// The NRQL SELECT clause to aggregate events. Default is COUNT(*).
	Select *ServiceLevelEventsQuerySelectCreateInput `json:"select,omitempty"`
	// The NRQL condition to filter the events.
	Where NRQL `json:"where,omitempty"`
}

// ServiceLevelEventsQuerySelect - The resulting NRQL SELECT clause to aggregate events.
type ServiceLevelEventsQuerySelect struct {
	// The event attribute to use in the SELECT clause.
	Attribute string `json:"attribute,omitempty"`
	// The function to use in the SELECT clause.
	Function ServiceLevelEventsQuerySelectFunction `json:"function"`
}

// ServiceLevelEventsQuerySelectCreateInput - The NRQL SELECT clause to aggregate events.
type ServiceLevelEventsQuerySelectCreateInput struct {
	// The event attribute to use in the SELECT clause.
	Attribute string `json:"attribute,omitempty"`
	// The function to use in the SELECT clause.
	Function ServiceLevelEventsQuerySelectFunction `json:"function"`
}

// ServiceLevelEventsQuerySelectUpdateInput - The NRQL SELECT clause to aggregate events.
type ServiceLevelEventsQuerySelectUpdateInput struct {
	// The event attribute to use in the SELECT clause.
	Attribute string `json:"attribute,omitempty"`
	// The function to use in the SELECT clause.
	Function ServiceLevelEventsQuerySelectFunction `json:"function"`
}

// ServiceLevelEventsQueryUpdateInput - The query that represents the events to fetch.
type ServiceLevelEventsQueryUpdateInput struct {
	// The NRDB event to fetch the data from.
	From NRQL `json:"from"`
	// The NRQL SELECT clause to aggregate events. Default is COUNT(*).
	Select *ServiceLevelEventsQuerySelectUpdateInput `json:"select,omitempty"`
	// The NRQL condition to filter the events.
	Where NRQL `json:"where,omitempty"`
}

// ServiceLevelEventsUpdateInput - The events that define the SLI.
type ServiceLevelEventsUpdateInput struct {
	// The definition of bad events.
	BadEvents *ServiceLevelEventsQueryUpdateInput `json:"badEvents,omitempty"`
	// The definition of good events.
	GoodEvents *ServiceLevelEventsQueryUpdateInput `json:"goodEvents,omitempty"`
	// The definition of valid events.
	ValidEvents *ServiceLevelEventsQueryUpdateInput `json:"validEvents,omitempty"`
}

// ServiceLevelIndicator - The definition of the SLI.
type ServiceLevelIndicator struct {
	// The date when the SLI was created represented in the number of milliseconds since the Unix epoch.
	CreatedAt *nrtime.EpochMilliseconds `json:"createdAt"`
	// The user who created the SLI.
	CreatedBy UserReference `json:"createdBy,omitempty"`
	// The description of the SLI.
	Description string `json:"description,omitempty"`
	// The entity which the SLI is attached to.
	EntityGUID common.EntityGUID `json:"entityGuid"`
	// The events that define the SLI.
	Events ServiceLevelEvents `json:"events"`
	// The unique entity identifier of the SLI.
	GUID common.EntityGUID `json:"guid"`
	// The unique identifier of the SLI.
	ID string `json:"id"`
	// The name of the SLI.
	Name string `json:"name"`
	// A list of objective definitions.
	Objectives []ServiceLevelObjective `json:"objectives"`
	// The date when the SLI was last updated represented in the number of milliseconds since the Unix epoch.
	UpdatedAt *nrtime.EpochMilliseconds `json:"updatedAt,omitempty"`
	// The user who last update the SLI.
	UpdatedBy UserReference `json:"updatedBy,omitempty"`
}

// ServiceLevelIndicatorCreateInput - The input object that represents the SLI that will be created.
type ServiceLevelIndicatorCreateInput struct {
	// The description of the SLI.
	Description string `json:"description,omitempty"`
	// The events that define the SLI.
	Events ServiceLevelEventsCreateInput `json:"events,omitempty"`
	// The name of the SLI.
	Name string `json:"name"`
	// A list of objective definitions.
	Objectives []ServiceLevelObjectiveCreateInput `json:"objectives,omitempty"`
}

// ServiceLevelIndicatorResultQueries - The resulting NRQL queries that help consume the metrics of the SLI.
type ServiceLevelIndicatorResultQueries struct {
	// The NRQL query that measures the good events.
	GoodEvents ServiceLevelResultQuery `json:"goodEvents"`
	// The NRQL query that measures the value of the SLI.
	Indicator ServiceLevelResultQuery `json:"indicator"`
	// The NRQL query that measures the valid events.
	ValidEvents ServiceLevelResultQuery `json:"validEvents"`
}

// ServiceLevelIndicatorUpdateInput - The input object that represents the SLI that will be updated.
type ServiceLevelIndicatorUpdateInput struct {
	// The description of the SLI.
	Description string `json:"description,omitempty"`
	// The events that define the SLI.
	Events *ServiceLevelEventsUpdateInput `json:"events,omitempty"`
	// The name of the SLI.
	Name string `json:"name,omitempty"`
	// A list of objective definitions.
	Objectives []ServiceLevelObjectiveUpdateInput `json:"objectives,omitempty"`
}

// ServiceLevelObjective - An objective definition.
type ServiceLevelObjective struct {
	// The description of the SLO.
	Description string `json:"description,omitempty"`
	// The name of the SLO.
	Name string `json:"name,omitempty"`
	// The target percentage of the SLO.
	Target float64 `json:"target"`
	// The time window configuration of the SLO.
	TimeWindow ServiceLevelObjectiveTimeWindow `json:"timeWindow"`
}

// ServiceLevelObjectiveCreateInput - The input object that represents an objective definition.
type ServiceLevelObjectiveCreateInput struct {
	// The description of the SLO.
	Description string `json:"description,omitempty"`
	// The name of the SLO.
	Name string `json:"name,omitempty"`
	// The target percentage of the SLO. Maximum value is 100.
	Target float64 `json:"target"`
	// The time window configuration of the SLO.
	TimeWindow ServiceLevelObjectiveTimeWindowCreateInput `json:"timeWindow,omitempty"`
}

// ServiceLevelObjectiveResultQueries - The resulting NRQL queries that help consume the metrics of the SLO.
type ServiceLevelObjectiveResultQueries struct {
	// The NRQL query that measures the attainment of the SLO target.
	Attainment ServiceLevelResultQuery `json:"attainment"`
}

// ServiceLevelObjectiveRollingTimeWindow - The rolling time window configuration of the SLO.
type ServiceLevelObjectiveRollingTimeWindow struct {
	// The count of time units.
	Count int `json:"count"`
	// The time unit.
	Unit ServiceLevelObjectiveRollingTimeWindowUnit `json:"unit"`
}

// ServiceLevelObjectiveRollingTimeWindowCreateInput - The rolling time window configuration of the SLO.
type ServiceLevelObjectiveRollingTimeWindowCreateInput struct {
	// The count of time units. Accepted values are 1, 7 and 28 days.
	Count int `json:"count"`
	// The time unit.
	Unit ServiceLevelObjectiveRollingTimeWindowUnit `json:"unit"`
}

// ServiceLevelObjectiveRollingTimeWindowUpdateInput - The rolling time window configuration of the SLO.
type ServiceLevelObjectiveRollingTimeWindowUpdateInput struct {
	// The count of time units. Accepted values are 1, 7 and 28 days.
	Count int `json:"count"`
	// The time unit.
	Unit ServiceLevelObjectiveRollingTimeWindowUnit `json:"unit"`
}

// ServiceLevelObjectiveTimeWindow - The time window configuration of the SLO.
type ServiceLevelObjectiveTimeWindow struct {
	// The rolling time window configuration of the SLO.
	Rolling ServiceLevelObjectiveRollingTimeWindow `json:"rolling,omitempty"`
}

// ServiceLevelObjectiveTimeWindowCreateInput - The time window configuration of the SLO.
type ServiceLevelObjectiveTimeWindowCreateInput struct {
	// The rolling time window configuration of the SLO.
	Rolling ServiceLevelObjectiveRollingTimeWindowCreateInput `json:"rolling,omitempty"`
}

// ServiceLevelObjectiveTimeWindowUpdateInput - The time window configuration of the SLO.
type ServiceLevelObjectiveTimeWindowUpdateInput struct {
	// The rolling time window configuration of the SLO.
	Rolling ServiceLevelObjectiveRollingTimeWindowUpdateInput `json:"rolling,omitempty"`
}

// ServiceLevelObjectiveUpdateInput - The input object that represents an objective definition.
type ServiceLevelObjectiveUpdateInput struct {
	// The description of the SLO.
	Description string `json:"description,omitempty"`
	// The name of the SLO.
	Name string `json:"name,omitempty"`
	// The target percentage of the SLO. Maximum value is 100.
	Target float64 `json:"target"`
	// The time window configuration of the SLO.
	TimeWindow ServiceLevelObjectiveTimeWindowUpdateInput `json:"timeWindow,omitempty"`
}

// ServiceLevelResultQuery - A resulting query.
type ServiceLevelResultQuery struct {
	// A NRQL query.
	NRQL NRQL `json:"nrql"`
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

type indicatorsResponse struct {
	Actor Actor `json:"actor"`
}

// Float - The `Float` scalar type represents signed double-precision fractional
// values as specified by
// [IEEE 754](http://en.wikipedia.org/wiki/IEEE_floating_point).
type Float string

// ID - The `ID` scalar type represents a unique identifier, often used to
// refetch an object or as key for a cache. The ID type appears in a JSON
// response as a String; however, it is not intended to be human-readable.
// When expected as an input type, any string (such as `"4"`) or integer
// (such as `4`) input value will be accepted as an ID.
type ID string

// NRQL - This scalar represents a NRQL query string.
//
// See the [NRQL Docs](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions) for more information about NRQL syntax.
type NRQL string
