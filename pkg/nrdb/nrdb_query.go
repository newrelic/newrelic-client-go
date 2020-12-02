// Package nrdb provides a programmatic API for interacting with NRDB, New Relic's Datastore
package nrdb

import "context"

func (n *Nrdb) Query(accountID int, query Nrql) (*NrdbResultContainer, error) {
	return n.QueryWithContext(context.Background(), accountID, query)
}

// QueryWithContext facilitates making a NRQL query.
func (n *Nrdb) QueryWithContext(ctx context.Context, accountID int, query Nrql) (*NrdbResultContainer, error) {
	respBody := gqlNrglQueryResponse{}

	vars := map[string]interface{}{
		"accountId": accountID,
		"query":     query,
	}

	if err := n.client.NerdGraphQueryWithContext(ctx, gqlNrqlQuery, vars, &respBody); err != nil {
		return nil, err
	}

	return &respBody.Actor.Account.Nrql, nil
}

func (n *Nrdb) QueryHistory() (*[]NrqlHistoricalQuery, error) {
	return n.QueryHistoryWithContext(context.Background())
}

func (n *Nrdb) QueryHistoryWithContext(ctx context.Context) (*[]NrqlHistoricalQuery, error) {
	respBody := gqlNrglQueryHistoryResponse{}
	vars := map[string]interface{}{}

	if err := n.client.NerdGraphQueryWithContext(ctx, gqlNrqlQueryHistoryQuery, vars, &respBody); err != nil {
		return nil, err
	}

	return &respBody.Actor.NrqlQueryHistory, nil
}

const (
	gqlNrqlQueryHistoryQuery = `{ actor { nrqlQueryHistory { accountId nrql timestamp } } }`

	gqlNrqlQuery = `query($query: Nrql!, $accountId: Int!) { actor { account(id: $accountId) { nrql(query: $query) {
    currentResults otherResult previousResults results totalResult
    metadata { eventTypes facets messages timeWindow { begin compareWith end since until } }
  } } } }`
)

type gqlNrglQueryResponse struct {
	Actor struct {
		Account struct {
			Nrql NrdbResultContainer
		}
	}
}

type gqlNrglQueryHistoryResponse struct {
	Actor struct {
		NrqlQueryHistory []NrqlHistoricalQuery
	}
}
