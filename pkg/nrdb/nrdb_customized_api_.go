// Package nrdb provides a programmatic API for interacting with NRDB, New Relic's Datastore.
// This package is NOT covered by Tutone.

package nrdb

import "context"

// WARNING! The following function, 'Query' is used by newrelic-cli to run pre-install
// validation procedures before the actual installation begins; and is hence, extremely fragile.
// Please do not resort to changing this function unless necessary, such as in the case
// of a deprecation or an end-of-life. Kindly duplicate this function to allow more
// attributes, or carefully modify functions following this function, below.

// Query facilitates making an NRQL query using NerdGraph.
func (n *Nrdb) PerformNRQLQuery(accountID int, query NRQL) (*NRDBResultContainerMultiResultCustomized, error) {
	return n.PerformNRQLQueryWithContext(context.Background(), accountID, query)
}

// WARNING! This function is extremely fragile.
// Please read the note above the function 'Query': refrain from making changes unless extremely necessary.

// QueryWithContext facilitates making a NRQL query.
func (n *Nrdb) PerformNRQLQueryWithContext(ctx context.Context, accountID int, query NRQL) (*NRDBResultContainerMultiResultCustomized, error) {
	respBody := gqlNRQLQueryResponseCustomized{}

	vars := map[string]interface{}{
		"accountId": accountID,
		"query":     query,
	}

	if err := n.client.NerdGraphQueryWithContext(ctx, gqlNrqlQuery, vars, &respBody); err != nil {
		return nil, err
	}

	return &respBody.Actor.Account.NRQL, nil
}

type gqlNRQLQueryResponseCustomized struct {
	Actor struct {
		Account struct {
			NRQL NRDBResultContainerMultiResultCustomized
		}
	}
}
