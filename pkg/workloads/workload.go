package workloads

import (
	"time"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// Workload represents a New Relic One workload.
type Workload struct {
	Account             *AccountReference            `json:"account,omitempty"`
	CreatedAt           *serialization.EpochTime     `json:"created_at,omitempty"`
	CreatedBy           *UserReference               `json:"created_by,omitempty"`
	Entities            []*WorkloadEntityRef         `json:"entities,omitempty"`
	EntitySearchQueries []*WorkloadEntitySearchQuery `json:"entitySearchQueries,omitempty"`
	EntitySearchQuery   string                       `json:"entitySearchQuery,omitempty"`
	GUID                *string                      `json:"guid,omitempty"`
	ID                  *int                         `json:"id,omitempty"`
	Name                *string                      `json:"name,omitempty"`
	Permalink           *string                      `json:"permalink,omitempty"`
	ScopeAccounts       *WorkloadScopeAccounts       `json:"scopeAccounts,omitempty"`
	UpdatedAt           *serialization.EpochTime     `json:"updated_at,omitempty"`
}

// AccountReference represents the account this workload belongs to.
type AccountReference struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// WorkloadEntityRef represents an entity referenced by this workload.
type WorkloadEntityRef struct {
	GUID string `json:"id,omitempty"`
}

// WorkloadEntitySearchQuery represents an entity search used by this workload.
type WorkloadEntitySearchQuery struct {
	CreatedAt *time.Time               `json:"createdAt,omitempty"`
	CreatedBy *UserReference           `json:"createdBy,omitempty"`
	ID        *int                     `json:"id,omitempty"`
	Name      *string                  `json:"name,omitempty"`
	Query     *string                  `json:"query,omitempty"`
	UpdatedAt *serialization.EpochTime `json:"updatedAt,omitempty"`
}

// UserReference represents a user referenced by this workload's search query.
type UserReference struct {
	Email    string `json:"email,omitempty"`
	Gravatar string `json:"gravatar,omitempty"`
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
}

// WorkloadScopeAccounts represents the accounts used to scope this workload.
type WorkloadScopeAccounts struct {
	AccountIDs []*int `json:"accountIds,omitempty"`
}

// ListWorkloads retrieves a set of New Relic One workloads by their account ID.
func (e *Workloads) ListWorkloads(accountID int) ([]*Workload, error) {
	resp := listWorkloadsResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
	}

	if err := e.client.Query(listWorkloadsQuery, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.Actor.Account.Workload.Collections) == 0 {
		return nil, errors.NewNotFound("")
	}

	return resp.Actor.Account.Workload.Collections, nil
}

// GetWorkload retrieve a New Relic One workload by its ID.
func (e *Workloads) GetWorkload(accountID int, workloadID int) (*Workload, error) {
	resp := getWorkloadResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"id":        workloadID,
	}

	if err := e.client.Query(getWorkloadQuery, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Actor.Account.Workload.Collection.ID == nil {
		return nil, errors.NewNotFound("")
	}

	return &resp.Actor.Account.Workload.Collection, nil
}

const (
	// graphqlWorkloadStructFields is the set of fields that we want returned on workload queries,
	// and should map back directly to the Workload struct
	graphqlEntityStructFields = `
			account {
				id
				name
			}
			createdBy {
				email
				gravatar
				id
				name
			}
			createdAt
			entities {
				guid
			}
			entitySearchQueries {
				createdAt
				createdBy {
					email
					gravatar
					id
					name
				}
				name
				id
				query
				updatedAt
			}
			entitySearchQuery
			guid
			id
			permalink
			name
			scopeAccounts {
				accountIds
			}
			updatedAt
`

	getWorkloadQuery = `query($id: Int!, $accountId: Int!) { actor { account(id: $accountId) { workload { collection(id: $id)  {` +
		graphqlEntityStructFields +
		` } } } } }`

	listWorkloadsQuery = `query($accountId: Int!) { actor { account(id: $accountId) { workload { collections {` +
		graphqlEntityStructFields +
		` } } } } }`
)

type listWorkloadsResponse struct {
	Actor struct {
		Account struct {
			Workload struct {
				Collections []*Workload
			}
		}
	}
}

type getWorkloadResponse struct {
	Actor struct {
		Account struct {
			Workload struct {
				Collection Workload
			}
		}
	}
}
