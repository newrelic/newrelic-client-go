package entities

import (
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// Type represents a New Relic One entity type (short)
type Type string

var (
	// Types specifies the possible types for a New Relic One entity.
	Types = struct {
		Application Type
		Dashboard   Type
		Host        Type
		Monitor     Type
		Workload    Type
	}{
		Application: "APPLICATION",
		Dashboard:   "DASHBOARD",
		Host:        "HOST",
		Monitor:     "MONITOR",
		Workload:    "WORKLOAD",
	}
)

// EntityDomainType represents a New Relic One entity domain.
type EntityDomainType string

var (
	// EntityDomains specifies the possible domains for a New Relic One entity.
	EntityDomains = struct {
		APM            EntityDomainType
		Browser        EntityDomainType
		Infrastructure EntityDomainType
		Mobile         EntityDomainType
		Nr1            EntityDomainType
		Synthetics     EntityDomainType
		Visualization  EntityDomainType
	}{
		APM:            "APM",
		Browser:        "BROWSER",
		Infrastructure: "INFRA",
		Mobile:         "MOBILE",
		Nr1:            "NR1",
		Synthetics:     "SYNTH",
		Visualization:  "VIZ",
	}
)

// EntityAlertSeverityType represents a New Relic One entity alert severity.
type EntityAlertSeverityType string

var (
	// EntityAlertSeverities specifies the possible alert severities for a New Relic One entity.
	EntityAlertSeverities = struct {
		Critical      EntityAlertSeverityType
		NotAlerting   EntityAlertSeverityType
		NotConfigured EntityAlertSeverityType
		Warning       EntityAlertSeverityType
	}{
		Critical:      "CRITICAL",
		NotAlerting:   "NOT_ALERTING",
		NotConfigured: "NOT_CONFIGURED",
		Warning:       "WARNING",
	}
)

// SearchEntitiesParams represents a set of search parameters for retrieving New Relic One entities.
type SearchEntitiesParams struct {
	AlertSeverity                 EntityAlertSeverityType `json:"alertSeverity,omitempty"`
	Domain                        EntityDomainType        `json:"domain,omitempty"`
	InfrastructureIntegrationType string                  `json:"infrastructureIntegrationType,omitempty"`
	Name                          string                  `json:"name,omitempty"`
	Reporting                     *bool                   `json:"reporting,omitempty"`
	Tags                          *TagValue               `json:"tags,omitempty"`
	Type                          EntityType              `json:"type,omitempty"`
}

// SearchEntities searches New Relic One entities based on the provided search parameters.
func (e *Entities) SearchEntities(params SearchEntitiesParams) ([]*Entity, error) {
	entities := []*Entity{}
	var nextCursor *string

	for ok := true; ok; ok = nextCursor != nil {
		resp := searchEntitiesResponse{}
		vars := map[string]interface{}{
			"queryBuilder": params,
			"cursor":       nextCursor,
		}

		if err := e.client.NerdGraphQuery(searchEntitiesQuery, vars, &resp); err != nil {
			return nil, err
		}

		entities = append(entities, resp.Actor.EntitySearch.Results.Entities...)

		nextCursor = resp.Actor.EntitySearch.Results.NextCursor
	}

	return entities, nil
}

// GetEntities retrieves a set of New Relic One entities by their entity guids.
func (e *Entities) GetEntities(guids []string) ([]*Entity, error) {
	resp := getEntitiesResponse{}
	vars := map[string]interface{}{
		"guids": guids,
	}

	if err := e.client.NerdGraphQuery(getEntitiesQuery, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.Actor.Entities) == 0 {
		return nil, errors.NewNotFound("")
	}

	return resp.Actor.Entities, nil
}

// GetEntity retrieve a set of New Relic One entities by their entity guids.
func (e *Entities) GetEntity(guid string) (*Entity, error) {
	resp := getEntityResponse{}
	vars := map[string]interface{}{
		"guid": guid,
	}

	if err := e.client.NerdGraphQuery(getEntityQuery, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Actor.Entity == nil {
		return nil, errors.NewNotFound("")
	}

	return resp.Actor.Entity, nil
}

const (
	// graphqlEntityStructFields is the set of fields that we want returned on entity queries,
	// and should map back directly to the Entity struct
	graphqlEntityStructFields = `
					accountId
					domain
					entityType
					guid
					name
					permalink
					reporting
					type
`
	graphqlEntityStructTagsFields = `
					tagsWithMetadata {
						key
						values {
							mutable
							value
						}
					}
					tags {
						key
						values
					}
`

	graphqlApmApplicationEntityFields = `
					... on ApmApplicationEntity {
						applicationId
						alertSeverity
						language
						runningAgentVersions {
							maxVersion
							minVersion
						}
						settings {
							apdexTarget
							serverSideConfig
						}
					}`

	graphqlApmApplicationEntityOutlineFields = `
					... on ApmApplicationEntityOutline {
						applicationId
						alertSeverity
						language
						runningAgentVersions {
							maxVersion
							minVersion
						}
						settings {
							apdexTarget
							serverSideConfig
						}
					}`

	graphqlBrowserApplicationEntityFields = `
		... on BrowserApplicationEntity {
			alertSeverity
			applicationId
			servingApmApplicationId
	}`

	graphqlBrowserApplicationEntityOutlineFields = `
		... on BrowserApplicationEntityOutline {
			alertSeverity
			applicationId
			servingApmApplicationId
	}`

	graphqlMobileApplicationEntityFields = `
		... on MobileApplicationEntity {
			alertSeverity
			applicationId
	}`

	graphqlMobileApplicationEntityOutlineFields = `
		... on MobileApplicationEntityOutline {
			alertSeverity
			applicationId
	}`

	getEntitiesQuery = `query($guids: [String!]!) { actor { entities(guids: $guids)  {` +
		graphqlEntityStructFields +
		graphqlEntityStructTagsFields +
		graphqlApmApplicationEntityFields +
		graphqlBrowserApplicationEntityFields +
		graphqlMobileApplicationEntityFields +
		` } } }`

	getEntityQuery = `query($guid: String!) { actor { entity(guid: $guid)  {` +
		graphqlEntityStructFields +
		graphqlEntityStructTagsFields +
		graphqlApmApplicationEntityFields +
		graphqlBrowserApplicationEntityFields +
		graphqlMobileApplicationEntityFields +
		` } } }`

	searchEntitiesQuery = `
		query($queryBuilder: EntitySearchQueryBuilder, $cursor: String) {
			actor {
				entitySearch(queryBuilder: $queryBuilder)  {
					results(cursor: $cursor) {
						nextCursor
						entities {` +
		graphqlEntityStructFields +
		graphqlApmApplicationEntityOutlineFields +
		graphqlBrowserApplicationEntityOutlineFields +
		graphqlMobileApplicationEntityOutlineFields +
		` } } } } }`
)

type searchEntitiesResponse struct {
	Actor struct {
		EntitySearch struct {
			Results struct {
				NextCursor *string
				Entities   []*Entity
			}
		}
	}
}

type getEntitiesResponse struct {
	Actor struct {
		Entities []*Entity
	}
}

type getEntityResponse struct {
	Actor struct {
		Entity *Entity
	}
}
