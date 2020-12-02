// Code generated by tutone: DO NOT EDIT
package nrdb

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
)

// ChartFormatType - Represents all the format types available for static charts.
type ChartFormatType string

var ChartFormatTypeTypes = struct {
	PDF ChartFormatType

	PNG ChartFormatType
}{

	PDF: "PDF",

	PNG: "PNG",
}

// ChartImageType - Represents all the visualization types available for static charts.
type ChartImageType string

var ChartImageTypeTypes = struct {
	APDEX ChartImageType

	AREA ChartImageType

	BAR ChartImageType

	BASELINE ChartImageType

	BILLBOARD ChartImageType

	BULLET ChartImageType

	EVENT_FEED ChartImageType

	FUNNEL ChartImageType

	HEATMAP ChartImageType

	HISTOGRAM ChartImageType

	LINE ChartImageType

	PIE ChartImageType

	SCATTER ChartImageType

	STACKED_HORIZONTAL_BAR ChartImageType

	TABLE ChartImageType

	VERTICAL_BAR ChartImageType
}{

	APDEX: "APDEX",

	AREA: "AREA",

	BAR: "BAR",

	BASELINE: "BASELINE",

	BILLBOARD: "BILLBOARD",

	BULLET: "BULLET",

	EVENT_FEED: "EVENT_FEED",

	FUNNEL: "FUNNEL",

	HEATMAP: "HEATMAP",

	HISTOGRAM: "HISTOGRAM",

	LINE: "LINE",

	PIE: "PIE",

	SCATTER: "SCATTER",

	STACKED_HORIZONTAL_BAR: "STACKED_HORIZONTAL_BAR",

	TABLE: "TABLE",

	VERTICAL_BAR: "VERTICAL_BAR",
}

// EmbeddedChartType - Represents all the visualization types available for embedded charts.
type EmbeddedChartType string

var EmbeddedChartTypeTypes = struct {
	APDEX EmbeddedChartType

	AREA EmbeddedChartType

	BAR EmbeddedChartType

	BASELINE EmbeddedChartType

	BILLBOARD EmbeddedChartType

	BULLET EmbeddedChartType

	EMPTY EmbeddedChartType

	EVENT_FEED EmbeddedChartType

	FUNNEL EmbeddedChartType

	HEATMAP EmbeddedChartType

	HISTOGRAM EmbeddedChartType

	JSON EmbeddedChartType

	LINE EmbeddedChartType

	MARKDOWN EmbeddedChartType

	PIE EmbeddedChartType

	SCATTER EmbeddedChartType

	STACKED_HORIZONTAL_BAR EmbeddedChartType

	TABLE EmbeddedChartType

	TRAFFIC_LIGHT EmbeddedChartType

	VERTICAL_BAR EmbeddedChartType
}{

	APDEX: "APDEX",

	AREA: "AREA",

	BAR: "BAR",

	BASELINE: "BASELINE",

	BILLBOARD: "BILLBOARD",

	BULLET: "BULLET",

	EMPTY: "EMPTY",

	EVENT_FEED: "EVENT_FEED",

	FUNNEL: "FUNNEL",

	HEATMAP: "HEATMAP",

	HISTOGRAM: "HISTOGRAM",

	JSON: "JSON",

	LINE: "LINE",

	MARKDOWN: "MARKDOWN",

	PIE: "PIE",

	SCATTER: "SCATTER",

	STACKED_HORIZONTAL_BAR: "STACKED_HORIZONTAL_BAR",

	TABLE: "TABLE",

	TRAFFIC_LIGHT: "TRAFFIC_LIGHT",

	VERTICAL_BAR: "VERTICAL_BAR",
}

// EventAttributeDefinition - A human-readable definition of an NRDB Event Type Attribute
type EventAttributeDefinition struct {
	// This attribute's category
	Category string `json:"category"`
	// A short description of this attribute
	Definition string `json:"definition"`
	// The New Relic docs page for this attribute
	DocumentationURL string `json:"documentationUrl"`
	// The human-friendly formatted name of the attribute
	Label string `json:"label"`
	// The name of the attribute
	Name string `json:"name"`
}

// EventDefinition - A human-readable definition of an NRDB Event Type
type EventDefinition struct {
	// A list of attribute definitions for this event type
	Attributes []EventAttributeDefinition `json:"attributes"`
	// A short description of this event
	Definition string `json:"definition"`
	// The human-friendly formatted name of the event
	Label string `json:"label"`
	// The name of the event
	Name string `json:"name"`
}

// NRDBMetadata - An object containing metadata about the query and result.
type NRDBMetadata struct {
	// A list of the event types that were queried.
	EventTypes []string `json:"eventTypes"`
	// A list of facets that were queried.
	Facets []string `json:"facets"`
	// Messages from NRDB included with the result.
	Messages []string `json:"messages"`
	// Details about the query time window.
	TimeWindow NRDBMetadataTimeWindow `json:"timeWindow"`
}

// NRDBMetadataTimeWindow - An object representing details about a query's time window.
type NRDBMetadataTimeWindow struct {
	// Timestamp marking the query begin time.
	Begin EpochMilliseconds `json:"begin"`
	// A clause representing the comparison time window.
	CompareWith string `json:"compareWith"`
	// Timestamp marking the query end time.
	End EpochMilliseconds `json:"end"`
	// SINCE clause resulting from the query
	Since string `json:"since"`
	// UNTIL clause resulting from the query
	Until string `json:"until"`
}

// NRDBResultContainer - A data structure that contains the results of the NRDB query along
// with other capabilities that enhance those results.
//
// Direct query results are available through `results`, `totalResult` and
// `otherResult`. The query you made is accessible through `nrql`, along with
// `metadata` about the query itself. Enhanced capabilities include
// `eventDefinitions`, `suggestedFacets` and more.
type NRDBResultContainer struct {
	// In a `COMPARE WITH` query, the `currentResults` contain the results for the current comparison time window.
	CurrentResults []NRDBResult `json:"currentResults"`
	// Generate a publicly sharable Embedded Chart URL for the NRQL query.
	//
	// For more details, see [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/query-nrql-through-new-relic-graphql-api#embeddable-charts).
	EmbeddedChartURL string `json:"embeddedChartUrl"`
	// Retrieve a list of event type definitions, providing descriptions
	// of the event types returned by this query, as well as details
	// of their attributes.
	EventDefinitions []EventDefinition `json:"eventDefinitions"`
	// Metadata about the query and result.
	Metadata NRDBMetadata `json:"metadata"`
	// The [NRQL](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions) query that was executed to yield these results.
	NRQL NRQL `json:"nrql"`
	// In a `FACET` query, the `otherResult` contains the aggregates representing the events _not_
	// contained in an individual `results` facet
	OtherResult NRDBResult `json:"otherResult"`
	// In a `COMPARE WITH` query, the `previousResults` contain the results for the previous comparison time window.
	PreviousResults []NRDBResult `json:"previousResults"`
	// The query results. This is a flat list of objects who's structure matches the query submitted.
	Results []NRDBResult `json:"results"`
	// Generate a publicly sharable static chart URL for the NRQL query.
	StaticChartURL string `json:"staticChartUrl"`
	// Retrieve a list of suggested NRQL facets for this NRDB query, to be used with
	// the `FACET` keyword in NRQL.
	//
	// Results are based on historical query behaviors.
	//
	// If the query already has a `FACET` clause, it will be ignored for the purposes
	// of suggesting facets.
	//
	// For more details, see [our docs](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/nerdgraph-graphiql-nrql-tutorial#suggest-facets).
	SuggestedFacets []NRQLFacetSuggestion `json:"suggestedFacets"`
	// Suggested queries that could help explain an anomaly in your timeseries based on either statistical differences in the data or historical usage.
	//
	// If no `anomalyTimeWindow` is supplied, we will attempt to detect a spike in the NRQL results. If no spike is found, the suggested query results will be empty.
	//
	// Input NRQL must be a TIMESERIES query and must have exactly one result.
	SuggestedQueries SuggestedNRQLQueryResponse `json:"suggestedQueries"`
	// In a `FACET` query, the `totalResult` contains the aggregates representing _all_ the events,
	// whether or not they are contained in an individual `results` facet
	TotalResult NRDBResult `json:"totalResult"`
}

// NRQLFacetSuggestion - A suggested NRQL facet. Facet suggestions may be either a single attribute, or
// a list of attributes in the case of multi-attribute facet suggestions.
type NRQLFacetSuggestion struct {
	// A list of attribute names comprising the suggested facet.
	//
	// Raw attribute names will be returned here. Attribute names may need to be
	// backtick-quoted before inclusion in a NRQL query.
	Attributes []string `json:"attributes"`
	// A modified version of the input NRQL, with a `FACET ...` clause appended.
	// If the original NRQL had a `FACET` clause already, it will be replaced.
	NRQL NRQL `json:"nrql"`
}

// NRQLHistoricalQuery - An NRQL query executed in the past.
type NRQLHistoricalQuery struct {
	// The Account ID queried.
	AccountID int `json:"accountId"`
	// The NRQL query executed.
	NRQL NRQL `json:"nrql"`
	// The time the query was executed.
	Timestamp EpochSeconds `json:"timestamp"`
}

// SuggestedAnomalyBasedNRQLQuery - A query suggestion based on analysis of events within a specific anomalous time
// range vs. nearby events outside of that time range.
type SuggestedAnomalyBasedNRQLQuery struct {
	// Information about the anomaly upon which this suggestion is based
	Anomaly SuggestedNRQLQueryAnomaly `json:"anomaly"`
	// The NRQL string to run for the suggested query
	NRQL string `json:"nrql"`
	// A human-readable title describing what the query shows
	Title string `json:"title"`
}

func (x *SuggestedAnomalyBasedNRQLQuery) ImplementsSuggestedNRQLQuery() {}

// SuggestedHistoryBasedNRQLQuery - query suggestion based on historical query patterns.
type SuggestedHistoryBasedNRQLQuery struct {
	// The NRQL string to run for the suggested query
	NRQL string `json:"nrql"`
	// A human-readable title describing what the query shows
	Title string `json:"title"`
}

func (x *SuggestedHistoryBasedNRQLQuery) ImplementsSuggestedNRQLQuery() {}

// SuggestedNRQLQuery - Interface type representing a query suggestion.
type SuggestedNRQLQuery struct {
	// The NRQL string to run for the suggested query
	NRQL string `json:"nrql"`
	// A human-readable title describing what the query shows
	Title string `json:"title"`
}

func (x *SuggestedNRQLQuery) ImplementsSuggestedNRQLQuery() {}

// SuggestedNRQLQueryAnomaly - Information about the anomaly upon which this analysis was based.
type SuggestedNRQLQueryAnomaly struct {
	// The approximate time range of the anomalous region
	TimeWindow TimeWindow `json:"timeWindow"`
}

// SuggestedNRQLQueryResponse - result type encapsulating suggested queries
type SuggestedNRQLQueryResponse struct {
	// List of suggested queries.
	Suggestions []SuggestedNRQLQueryInterface `json:"suggestions"`
}

// special
func (x *SuggestedNRQLQueryResponse) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	for k, v := range objMap {
		switch k {
		case "suggestions":
			var rawMessageSuggestions []*json.RawMessage
			err = json.Unmarshal(*v, &rawMessageSuggestions)
			if err != nil {
				return err
			}

			for _, m := range rawMessageSuggestions {
				xxx, err := UnmarshalSuggestedNRQLQueryInterface(*m)
				if err != nil {
					return err
				}

				if xxx != nil {
					x.Suggestions = append(x.Suggestions, *xxx)
				}
			}
		}
	}

	return nil
}

// TimeWindow - Represents a time window.
type TimeWindow struct {
	// The end time of the time window the number of milliseconds since the Unix epoch.
	EndTime EpochMilliseconds `json:"endTime"`
	// The start time of the time window the number of milliseconds since the Unix epoch.
	StartTime EpochMilliseconds `json:"startTime"`
}

// TimeWindowInput - Represents a time window input.
type TimeWindowInput struct {
	// The end time of the time window the number of milliseconds since the Unix epoch.
	EndTime EpochMilliseconds `json:"endTime"`
	// The start time of the time window the number of milliseconds since the Unix epoch.
	StartTime EpochMilliseconds `json:"startTime"`
}

// EpochMilliseconds - The `EpochMilliseconds` scalar represents the number of milliseconds since the Unix epoch
type EpochMilliseconds serialization.EpochTime

// EpochSeconds - The `EpochSeconds` scalar represents the number of seconds since the Unix epoch
type EpochSeconds serialization.EpochTime

// NRDBResult - This scalar represents a NRDB Result. It is a `Map` of `String` keys to values.
//
// The shape of these objects reflect the query used to generate them, the contents
// of the objects is not part of the GraphQL schema.
type NRDBResult string

// NRQL - This scalar represents a NRQL query string.
//
// See the [NRQL Docs](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions) for more information about NRQL syntax.
type NRQL string

// SuggestedNRQLQuery - Interface type representing a query suggestion.
type SuggestedNRQLQueryInterface interface {
	ImplementsSuggestedNRQLQuery()
}

//yes
func UnmarshalSuggestedNRQLQueryInterface(b []byte) (*SuggestedNRQLQueryInterface, error) {
	var err error

	var rawMessageSuggestedNRQLQuery map[string]*json.RawMessage
	err = json.Unmarshal(b, &rawMessageSuggestedNRQLQuery)
	if err != nil {
		return nil, err
	}

	var typeName string

	if rawTypeName, ok := rawMessageSuggestedNRQLQuery["__typename"]; ok {
		err = json.Unmarshal(*rawTypeName, &typeName)
		if err != nil {
			return nil, err
		}

		switch typeName {
		case "SuggestedAnomalyBasedNrqlQuery":
			var interfaceType SuggestedAnomalyBasedNRQLQuery
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx SuggestedNRQLQueryInterface = &interfaceType

			return &xxx, nil
		case "SuggestedHistoryBasedNrqlQuery":
			var interfaceType SuggestedHistoryBasedNRQLQuery
			err = json.Unmarshal(b, &interfaceType)
			if err != nil {
				return nil, err
			}

			var xxx SuggestedNRQLQueryInterface = &interfaceType

			return &xxx, nil
		}
	} else {
		keys := []string{}
		for k := range rawMessageSuggestedNRQLQuery {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("interface SuggestedNRQLQuery did not include a __typename field for inspection: %s", keys)
	}

	return nil, fmt.Errorf("interface SuggestedNRQLQuery was not matched against all PossibleTypes: %s", typeName)
}
