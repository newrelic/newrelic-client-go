# NRDB Package

The `nrdb` package provides a programmatic API for interacting with NRDB, New Relic's Datastore, allowing you to execute NRQL queries and process their results in Go applications.

## Overview

This package contains several methods for querying New Relic's NRDB database using NRQL (New Relic Query Language). The core functionality revolves around executing NRQL queries against your New Relic accounts and handling the structured response data.

## Query Functions

### Original Query Function

The original `Query()` function is the standard way to execute NRQL queries:

```go
func (n *Nrdb) Query(accountID int, query NRQL) (*NRDBResultContainer, error)
```

This function processes the query and returns an `NRDBResultContainer` with the results.

#### Usage Example

```go
client := newrelic.New(config)
accountID := 12345678

query := `SELECT count(*) FROM Transaction`
results, err := client.Nrdb.Query(accountID, nrdb.NRQL(query))
if err != nil {
    log.Fatal("Error running query:", err)
}

// Access the results
fmt.Printf("Query results: %v\n", results.Results)
```

### Limitation with FACET + TIMESERIES Queries

The original `Query()` function has a limitation when handling queries that combine both `FACET` and `TIMESERIES` clauses. In these cases, the NerdGraph API may return `otherResult` and `totalResult` as arrays rather than single objects, which conflicts with the `NRDBResultContainer` structure that expects single objects for these fields.

This can lead to JSON unmarshaling errors like:

```
json: cannot unmarshal array into Go struct field .data.Actor.Account.NRQL.otherResult of type nrdb.NRDBResult
```

#### Example of a Problematic Query

```go
// This query will cause unmarshaling errors with the original Query() function
query := `SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago`
```

### Enhanced Query Function

To address this limitation, we've introduced the `PerformNRQLQuery()` function:

```go
func (n *Nrdb) PerformNRQLQuery(accountID int, query NRQL) (*NRDBResultContainerMultiResultCustomized, error)
```

This function uses a customized result container that can handle both single objects and arrays for `otherResult` and `totalResult` fields, making it compatible with all types of NRQL queries, including those with combined `FACET` and `TIMESERIES` clauses.

#### Usage Example

```go
client := newrelic.New(config)
accountID := 12345678

// Using a query that combines FACET and TIMESERIES
query := `SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour SINCE 1 day ago`
results, err := client.Nrdb.PerformNRQLQuery(accountID, nrdb.NRQL(query))
if err != nil {
    log.Fatal("Error running query:", err)
}

// Check if otherResult is an array or single object
if otherResults, isArray := results.OtherResult.([]nrdb.NRDBResult); isArray {
    fmt.Printf("Got %d otherResults\n", len(otherResults))
    // Process array of results
    for i, result := range otherResults {
        fmt.Printf("Other Result %d: %v\n", i, result)
    }
} else if singleResult, isSingle := results.OtherResult.(nrdb.NRDBResult); isSingle {
    fmt.Printf("Single Other Result: %v\n", singleResult)
}

// Similarly for totalResult
if totalResults, isArray := results.TotalResult.([]nrdb.NRDBResult); isArray {
    // Process array of results
} else if singleResult, isSingle := results.TotalResult.(nrdb.NRDBResult); isSingle {
    // Process single result
}
```

## When to Use Each Function

1. **Use `Query()` when**:
   - You're executing simple NRQL queries without combined `FACET` and `TIMESERIES` clauses
   - You need the most stable and long-supported API
   - You're sure your queries will return single objects for `otherResult` and `totalResult`

2. **Use `PerformNRQLQuery()` when**:
   - You're executing queries that combine `FACET` and `TIMESERIES` clauses
   - You're unsure about the structure of the API response
   - You need to handle both array and object responses for `otherResult` and `totalResult`

## Query Types and Expected Response Formats

| Query Type | otherResult Type | totalResult Type | Example |
|------------|-----------------|------------------|---------|
| Simple | Single Object | Single Object | `SELECT count(*) FROM Transaction` |
| FACET only | Single Object | Single Object | `SELECT count(*) FROM Transaction FACET appName` |
| TIMESERIES only | Single Object | Single Object | `SELECT count(*) FROM Transaction TIMESERIES 1 hour` |
| FACET + TIMESERIES | Array | Array | `SELECT count(*) FROM Transaction FACET appName TIMESERIES 1 hour` |

## Extended Functionality

Both sets of functions have variants that provide additional capabilities:

- **WithContext**: Allows passing a context for request cancellation or timeouts
- **WithExtendedResponse**: Returns additional metadata about the query
- **WithAdditionalOptions**: Allows specifying timeout and async parameters

### Example with Extended Response

```go
query := `SELECT count(*) FROM Transaction`
results, err := client.Nrdb.PerformNRQLQueryWithExtendedResponse(accountID, nrdb.NRQL(query))
```

### Example with Additional Options

```go
query := `SELECT count(*) FROM Transaction`
timeout := nrdb.Seconds(30)  // 30 second timeout
async := false               // Synchronous query
results, err := client.Nrdb.QueryWithAdditionalOptions(accountID, nrdb.NRQL(query), timeout, async)
```

## Best Practices

1. **Use the appropriate function** for your query type to avoid unmarshaling errors
2. **Type-check interface{} fields** when using `PerformNRQLQuery()` since `otherResult` and `totalResult` can be either arrays or single objects
3. **Set reasonable timeouts** for queries that may take longer to execute
4. **Consider using async mode** for queries that process large amounts of data

## Troubleshooting

If you encounter JSON unmarshaling errors when using `Query()` with complex queries, switch to `PerformNRQLQuery()` to handle the variable response structure.

If you need to maintain compatibility with code that expects the standard `NRDBResultContainer`, but want to use `PerformNRQLQuery()` for its flexibility, you can convert between types:

```go
customResults, err := client.Nrdb.PerformNRQLQuery(accountID, nrdb.NRQL(query))
if err != nil {
    return nil, err
}

// Convert otherResult if it's an array (take first element)
if otherResults, isArray := customResults.OtherResult.([]nrdb.NRDBResult); isArray && len(otherResults) > 0 {
    // Use the first element as a single result
    singleOtherResult := otherResults[0]
    // Use with code expecting NRDBResult
}
```

## Further Reading

- [New Relic NRQL Documentation](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/nrql-resources/nrql-syntax-components-functions)
- [NerdGraph API Documentation](https://docs.newrelic.com/docs/apis/graphql-api/tutorials/query-nrql-through-new-relic-graphql-api)
