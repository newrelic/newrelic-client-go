package dashboards

import (
	"context"
	"fmt"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// ListDashboardsParams represents a set of filters to be
// used when querying New Relic dashboards.
type ListDashboardsParams struct {
	Category      string     `url:"filter[category],omitempty"`
	CreatedAfter  *time.Time `url:"filter[created_after],omitempty"`
	CreatedBefore *time.Time `url:"filter[created_before],omitempty"`
	Page          int        `url:"page,omitempty"`
	PerPage       int        `url:"per_page,omitempty"`
	Sort          string     `url:"sort,omitempty"`
	Title         string     `url:"filter[title],omitempty"`
	UpdatedAfter  *time.Time `url:"filter[updated_after],omitempty"`
	UpdatedBefore *time.Time `url:"filter[updated_before],omitempty"`
}

// ListDashboards is used to retrieve New Relic dashboards.
func (d *Dashboards) ListDashboards(params *ListDashboardsParams) ([]*Dashboard, error) {
	dashboard := []*Dashboard{}
	nextURL := d.config.Region().RestURL("dashboards.json")

	for nextURL != "" {
		response := dashboardsResponse{}
		resp, err := d.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		dashboard = append(dashboard, response.Dashboards...)

		paging := d.pager.Parse(resp)
		nextURL = paging.Next
	}

	return dashboard, nil
}

// GetDashboardEntity is used to retrieve a single New Relic One Dashboard
func (d *Dashboards) GetDashboardEntity(gUID entities.EntityGUID) (*entities.DashboardEntity, error) {
	return d.GetDashboardEntityWithContext(context.Background(), gUID)
}

// GetDashboardEntityWithContext is used to retrieve a single New Relic One Dashboard
func (d *Dashboards) GetDashboardEntityWithContext(ctx context.Context, gUID entities.EntityGUID) (*entities.DashboardEntity, error) {
	resp := struct {
		Actor entities.Actor `json:"actor"`
	}{}
	vars := map[string]interface{}{
		"guid": gUID,
	}

	if err := d.client.NerdGraphQueryWithContext(ctx, getDashboardEntityQuery, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Actor.Entity == nil {
		return nil, errors.NewNotFound("entity not found. GUID: '" + string(gUID) + "'")
	}

	return resp.Actor.Entity.(*entities.DashboardEntity), nil
}

// getDashboardEntityQuery is not auto-generated as tutone does not currently support
// generation of queries that return a specific interface.
const getDashboardEntityQuery = `query ($guid: EntityGuid!) {
  actor {
    entity(guid: $guid) {
      guid
      ... on DashboardEntity {
        __typename
        accountId
        createdAt
        dashboardParentGuid
        description
        indexedAt
        name
        owner { email userId }
        pages {
          createdAt
          description
          guid
          name
          owner { email userId }
          updatedAt
          widgets {
            rawConfiguration
            configuration {
              area { nrqlQueries { accountId query } }
              bar { nrqlQueries { accountId query } }
              billboard { nrqlQueries { accountId query } thresholds { alertSeverity value } }
              line { nrqlQueries { accountId query } }
              markdown { text }
              pie { nrqlQueries { accountId query } }
              table { nrqlQueries { accountId query } }
            }
            layout { column height row width }
            title
            visualization { id }
            id
            linkedEntities {
              __typename
              guid
              name
              accountId
              tags { key values }
              ... on DashboardEntityOutline {
                dashboardParentGuid
              }
            }
          }
        }
        permalink
        permissions
        tags { key values }
        tagsWithMetadata { key values { mutable value } }
        updatedAt
      }
    }
  }
}`

// GetDashboard is used to retrieve a single New Relic dashboard.
// Deprecated: Use GetDashboardEntity instead
func (d *Dashboards) GetDashboard(dashboardID int) (*Dashboard, error) {
	response := dashboardResponse{}
	url := fmt.Sprintf("/dashboards/%d.json", dashboardID)

	_, err := d.client.Get(d.config.Region().RestURL(url), nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

// CreateDashboard is used to create a New Relic dashboard.
// Deprecated: Use DashboardCreate instead
func (d *Dashboards) CreateDashboard(dashboard Dashboard) (*Dashboard, error) {
	response := dashboardResponse{}
	reqBody := dashboardRequest{
		Dashboard: dashboard,
	}
	_, err := d.client.Post(d.config.Region().RestURL("dashboards.json"), nil, &reqBody, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

// UpdateDashboard is used to update a New Relic dashboard.
// Deprecated: Use DashboardUpdate instead
func (d *Dashboards) UpdateDashboard(dashboard Dashboard) (*Dashboard, error) {
	response := dashboardResponse{}
	url := fmt.Sprintf("/dashboards/%d.json", dashboard.ID)
	reqBody := dashboardRequest{
		Dashboard: dashboard,
	}

	_, err := d.client.Put(d.config.Region().RestURL(url), nil, &reqBody, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

// DeleteDashboard is used to delete a New Relic dashboard.
// Deprecated: Use DashboardDelete instead
func (d *Dashboards) DeleteDashboard(dashboardID int) (*Dashboard, error) {
	response := dashboardResponse{}
	url := fmt.Sprintf("/dashboards/%d.json", dashboardID)

	_, err := d.client.Delete(d.config.Region().RestURL(url), nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

type dashboardsResponse struct {
	Dashboards []*Dashboard `json:"dashboards,omitempty"`
}

type dashboardResponse struct {
	Dashboard Dashboard `json:"dashboard,omitempty"`
}

type dashboardRequest struct {
	Dashboard Dashboard `json:"dashboard"`
}
