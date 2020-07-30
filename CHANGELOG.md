<a name="v0.34.0"></a>
## [v0.34.0] - 2020-07-30
### Bug Fixes
- **alerts:** infra condition threshold value should be pointer to support zero-value thresholds
- **alerts:** always send infra condition description even if its an empty string
- **region:** make ConfigRegion case-insensitive

### Documentation Updates
- update README.md

### Features
- **graphql:** capture validation error output from response

<a name="v0.33.2"></a>
## [v0.33.2] - 2020-07-27
### Bug Fixes
- **alerts:** infra condition threshold value should be pointer to support zero-value thresholds

<a name="v0.33.1"></a>
## [v0.33.1] - 2020-07-24
### Bug Fixes
- **alerts:** always send infra condition description even if its an empty string

<a name="v0.33.0"></a>
## [v0.33.0] - 2020-07-23
### Bug Fixes
- fix http client compression
- **build:** run generate-tutone once per make command

### Features
- add a custom event resource
- **alerts:** add description field to infra alert condition

<a name="v0.32.1"></a>
## [v0.32.1] - 2020-07-17
### Bug Fixes
- **accounts:** add accounts API to client
- **nrql_conditions:** proper zero-value marshalling for threshold

<a name="v0.32.0"></a>
## [v0.32.0] - 2020-07-10
### Features
- **accounts:** add an accounts resource
- **synthetics:** add a monitor locations resource

<a name="v0.31.3"></a>
## [v0.31.3] - 2020-07-02
### Bug Fixes
- **alerts:** allow a value of 0 for NRQL condition thresholds

<a name="v0.31.2"></a>
## [v0.31.2] - 2020-07-01
### Bug Fixes
- **alerts:** better nerdgraph operator usage

<a name="v0.31.1"></a>
## [v0.31.1] - 2020-06-26
### Bug Fixes
- **alerts:** drop incorrect MonitorID flag from MultiLocationSyntheticsCondition

<a name="v0.31.0"></a>
## [v0.31.0] - 2020-06-18
### Bug Fixes
- add goreleaser back
- remove goreleaser

### Features
- **entities:** surface underlying application IDs for mobile application entities
- **eventstometrics:** add EventsToMetrics

<a name="v0.30.2"></a>
## [v0.30.2] - 2020-06-15
### Bug Fixes
- **nrdb:** Add String() to Epoch* types
- **serialization:** set tz UTC on MarshalJSON
- **serialization:** Fix nanosecond handling, set tz UTC, add EpochTime.String()

<a name="v0.30.1"></a>
## [v0.30.1] - 2020-06-12
### Bug Fixes
- **deps:** revert goreleaser v0.138.0 (causing import issues)

<a name="v0.30.0"></a>
## [v0.30.0] - 2020-06-12
### Features
- **nrdb:** Fetch nrql query history

<a name="v0.29.1"></a>
## [v0.29.1] - 2020-06-10
### Bug Fixes
- **nrdb:** Correctly unwrap the graphql context for Query, better testing

<a name="v0.29.0"></a>
## [v0.29.0] - 2020-06-10
### Bug Fixes
- **nerdgraph:** Add omitempty for yaml output
- **typegen:** Remove some overrides so types are generated without modifications

### Features
- **nrdb:** Add nrdb.Query(acct, NRQL) via NerdGraph
- **typegen:** Add imports to typegen

<a name="v0.28.1"></a>
## [v0.28.1] - 2020-06-03
### Bug Fixes
- **alerts:** add missing Outlier type to condition NrqlConditionTypes enum

<a name="v0.28.0"></a>
## [v0.28.0] - 2020-06-03
### Bug Fixes
- **alerts:** rename NrqlConditionTerms to singular for accuracy

### Features
- **alerts:** add DeleteNrqlConditionMutation as a proxy to DeleteConditionMutation
- **alerts:** add UpdateNrqlConditionOutlierMutation method for updating outlier NRQL alert conditions
- **alerts:** add CreateNrqlConditionOutlierMutation method for creating outlier NRQL alert conditions

### Refactor
- **nrql_conditions:** conditionID should be a string for consistency in ng implementation

<a name="v0.27.1"></a>
## [v0.27.1] - 2020-05-29
### Bug Fixes
- **serialization:** EpochTime handling of Unix timestamp with milliseconds

### Refactor
- **alerts:** update NG-based condition IDs to string type

<a name="v0.27.0"></a>
## [v0.27.0] - 2020-05-28
### Refactor
- **alerts:** update NG-based policy and condition IDs to string types

<a name="v0.26.0"></a>
## [v0.26.0] - 2020-05-27
### Bug Fixes
- **http:** Use default User-Agent header if none provided
- **http:** Refactor config.Compression and use it in http

### Features
- **http:** Enable compression handling for requests, consolidate POST methods
- **region:** Add Insights insert API URLs

<a name="v0.25.1"></a>
## [v0.25.1] - 2020-05-20
### Bug Fixes
- **alerts:** fix json tag for NrqlConditionInput.ValueFunction

<a name="v0.25.0"></a>
## [v0.25.0] - 2020-05-20
### Features
- enable APIKS auth for alerts and plugins packages
- **alerts:** add nerdgraph-based alert condition deletion

<a name="v0.24.1"></a>
## [v0.24.1] - 2020-05-19
### Bug Fixes
- **apm:** don't return empty zero values for floats for MetricTimesliceValues which can be misleading
- **apm:** adjust MetricDataParams json tags to support query params as arrays

<a name="v0.24.0"></a>
## [v0.24.0] - 2020-05-15
### Features
- **edge:** add trace observer resource

<a name="v0.23.4"></a>
## [v0.23.4] - 2020-05-13
### Bug Fixes
- **alerts:** allow blank runbook URL to be sent

<a name="v0.23.3"></a>
## [v0.23.3] - 2020-05-12
### Bug Fixes
- **workloads:** fix query type for entity guid

<a name="v0.23.2"></a>
## [v0.23.2] - 2020-05-11
### Bug Fixes
- **alerts:** Updating incorrect AlertEvents params
- **region:** gracefully fall back to default region

<a name="v0.23.1"></a>
## [v0.23.1] - 2020-05-04
### Bug Fixes
- **alerts:** Alerts paging was incorrectly generating URLs

<a name="v0.23.0"></a>
## [v0.23.0] - 2020-05-01
### Bug Fixes
- **build:** Github actions `make lint` for commit messages need full history
- **workloads:** Removing deprecated field `name` from `entitySearchQuery`

### Documentation Updates
- **typegen:** Add a README for typegen

### Features
- **events:** add alert events

### Refactor
- **build:** Split up github actions a bit
- **typegen:** Split/refactor much of typegen internals

<a name="v0.22.0"></a>
## [v0.22.0] - 2020-04-23
### Bug Fixes
- **alerts:** include "equal" operator for NRQL condition terms

### Features
- **dashboards:** add grid_column_count field

### Refactor
- **workloads:** query with GUID per upstream API change

<a name="v0.21.1"></a>
## [v0.21.1] - 2020-04-15
### Bug Fixes
- **alerts:** return a NotFound error when policies are not found in NerdGraph

### Refactor
- **workloads:** remove deprecated fields

<a name="v0.21.0"></a>
## [v0.21.0] - 2020-04-06
### Bug Fixes
- **build:** goreleaser now supports libraries, remove hack in config, skip build

### Features
- **typegen:** handle scalar types

### Refactor
- **alerts:** Update go:generate for types signature
- **newrelic:** Continue to fail, do not log invalid region
- **region:** Better error types/more tests
- **region:** Split parsing / fetching of region data so Parse can be reused
- **typegen:** Filter internal context off of generated descriptions if present
- **typegen:** Break out some structs
- **typegen:** Move all schema generation stuff into typegen

<a name="v0.20.1"></a>
## [v0.20.1] - 2020-04-01
### Bug Fixes
- **alerts:** use pointers for the result struct

### Refactor
- **nerdstorage:** check scope ID for zero values during nerdstorage operations

<a name="v0.20.0"></a>
## [v0.20.0] - 2020-03-31
### Bug Fixes
- **alerts:** include missing types to generate
- **nrql_conditions:** fix spelling errors, fix types
- **region:** Better URL building
- **typegen:** ensure we also generate non-input fields
- **typegen:** ensure handling of slice/LIST types
- **typegen:** default to resolving nested types

### Documentation Updates
- **README:** Update example in README, closes [#225](https://github.com/newrelic/newrelic-client-go/issues/225)
- **nerdstorage:** add examples

### Features
- **alerts:** implement NerdGraph policy search
- **internal:** add additional error context to graphQLError
- **nerdstorage:** add a nerdstorage package and resource
- **region:** Add a region package
- **typegen:** handle types of Kind OBJECT

### Refactor
- **alerts:** fix alerts tests
- **alerts:** Move FQDN/URL creation into package, out of http client for REST
- **apm:** Update apm tests
- **apm:** Move FQDN/URL creation into package, out of http client for REST
- **config:** Migrate Region to pkg/region
- **config:** Add new func for config
- **config:** Remove unused config fields
- **dashboards:** Update dashboard tests
- **dashboards:** Move FQDN/URL creation into package, out of http client for REST
- **entities:** Update entities tests
- **http:** introduce a request-scoped API for NerdGraph queries
- **http:** Remove assumption that we are talking to a REST endpoint
- **http:** Move HTTP client to use new region format
- **nerdgraph:** Update nerdgraph tests
- **plugins:** Update plugin tests
- **plugins:** Move FQDN/URL creation into package, out of http client for REST
- **region:** Change access to config.Region to ensure it exists
- **synthetics:** Update synthetics tests
- **synthetics:** Move FQDN/URL creation into package, out of http client for REST
- **typegen:** Convert to using go generate to run typegen, `make generate` to test
- **workloads:** Update workloads tests

<a name="v0.19.0"></a>
## [v0.19.0] - 2020-03-25
### Bug Fixes
- **alerts:** policy update response test
- **workloads:** remove nullable struct fields unless necessary

### Features
- **alerts:** add search method for NRQL conditions
- **alerts:** add get method for query NRQL conditions
- **alerts:** add update methods for baseline and static NRQL conditions
- **alerts:** add create methods for baseline and static NRQL conditions
- **nerdgraph:** begin generating structs from schema

### Refactor
- **alerts:** consolidate Nrql condition structs for better reusability

<a name="v0.18.0"></a>
## [v0.18.0] - 2020-03-20
### Bug Fixes
- **workloads:** fix some bugs in the workloads implementation
- **workloads:** export the workloads API via the newrelic package

### Features
- **alerts:** implement muting rules

<a name="v0.17.1"></a>
## [v0.17.1] - 2020-03-18
### Bug Fixes
- **alerts:** add custom unmarshaling for ConditionTerm
- **workloads:** use epoch time for EntitySearchQuery.CreatedAt

<a name="v0.17.0"></a>
## [v0.17.0] - 2020-03-17
### Bug Fixes
- **workloads:** map non-nullable fields to structs correctly

### Documentation Updates
- **alerts:** add package-level documentation and examples
- **apm:** add package-level documentation and examples
- **client:** add synopses for all packages
- **config:** add package-level documentation
- **dashboards:** add package-level documentation and examples
- **entities:** add package-level documentation and examples
- **errors:** update package-level documentation
- **infrastructure:** add package-level documentation
- **nerdgraph:** add package-level documentation and examples
- **newrelic:** add package-level documentation and examples
- **newrelic:** use single-letter vars for receivers
- **plugins:** add package-level documentation and examples
- **synthetics:** add package-level documentation and examples

### Features
- **alerts:** implement graphql policy methods
- **workloads:** add update operation, rework integration test scenario
- **workloads:** add delete and duplicate mutations
- **workloads:** add a workload create operation
- **workloads:** add a workloads resource, list and get methods

### Refactor
- **alerts:** Fix lint issue
- **alerts:** add types for fields with known values
- **apm:** Move Application REST implementation, use interface
- **dashboards:** add types for fields with known values
- **http:** Move NewRequest, have it follow New* func format
- **http:** Move graphql code out to file
- **http:** Make all fields private, add some setters/getters, more tests
- **http:** Consolidate GraphQL client, rename to http.Client
- **http:** Move GraphQL into http.NewRelicClient as Query()

<a name="v0.16.0"></a>
## [v0.16.0] - 2020-03-11
### Bug Fixes
- **build:** Force pull tags after each checkout

### Documentation Updates
- update community support information

### Features
- **entities:** Add some more details from BrowserApplicationEntity
- **entities:** Return more data on ApmApplicationEntity, and be consistent in what we return between fetch and search

### Refactor
- **entities:** Change Entity.Type type... Add more to the ENUMs

<a name="v0.15.0"></a>
## [v0.15.0] - 2020-03-09
### Bug Fixes
- **apm:** remove unused field
- **build:** Remove working dir config for CircleCI
- **http:** allow overriding of service name

### Refactor
- **alert_conditions:** remove transient PolicyID from struct for consistency with API response
- **alerts:** use consistent types for incident timestamp fields
- **build:** Make the build system consistent with other projects
- **nrql_conditions:** remove transient PolicyID from struct for consistency with API response
- **plugins_conditions:** remove transient PolicyID from struct for consistency with API response

<a name="v0.14.0"></a>
## [v0.14.0] - 2020-03-05
### Features
- **newrelic:** add types for fields with well known values

<a name="v0.13.0"></a>
## [v0.13.0] - 2020-03-03
### Bug Fixes
- **entities:** include applicationId for ApmApplicationEntity results
- **entities:** Make ApplicationID optional in results, omit if not returned by the API
- **http:** create a new errorValue for every request

### Refactor
- **apm:** refactor deployments resource to use new auth strategy
- **http:** refactor client to a request-scoped config context

<a name="v0.12.0"></a>
## [v0.12.0] - 2020-02-28
### Bug Fixes
- **docs:** Fix the release badge

### Features
- **nerdgraph:** implement ability to make raw graphql query

### Refactor
- **alerts:** Move structs into implementing files
- **apm:** Move structs into implementing files
- **config:** BREAKING CHANGE: Change environment vars and rename APIKey to AdminApiKey
- **dashboards:** Move structs into implementing files
- **entities:** Move structs into implementing files
- **synthetics:** Move structs into implementing files

<a name="v0.11.0"></a>
## [v0.11.0] - 2020-02-27
### Features
- **http:** allow personal API keys to be used for alerts and APM resources

### Refactor
- **http:** refactor authentication out of http client

<a name="v0.10.1"></a>
## [v0.10.1] - 2020-02-20
### Bug Fixes
- **entities:** tags filter needs to use type TagValue in graphql query
- **newrelic:** Add option to set ServiceName in Config

<a name="v0.10.0"></a>
## [v0.10.0] - 2020-02-19
### Features
- **ci:** add release make target
- **ci:** the beginnings of some release automation
- **synthetics:** add secure credentials resource
- **synthetics:** implement label monitor support

<a name="v0.9.0"></a>
## [v0.9.0] - 2020-02-05
### Bug Fixes
- allow string representations of JSON for alert channel webhook and payload
- **http:** Clear client responses between pages

### Features
- **alerts:** Implement multi-location synthetics conditions
- **http:** add trace logging with additional request info

<a name="v0.8.0"></a>
## [v0.8.0] - 2020-01-29
### Bug Fixes
- **alerts:** ensure multiple channels can be added via /alerts_policy_channel.json endpoint ([#114](https://github.com/newrelic/newrelic-client-go/issues/114))

### Features
- **apm:** Add support application metric names and data

<a name="v0.7.1"></a>
## [v0.7.1] - 2020-01-24
### Bug Fixes
- **alerts:** handle more complex JSON structures in headers and/or payload
- **logging:** use global methods for the default logger rather than a logrus instance

### Refactor
- **entities:** rename SearchEntities params struct per convention
- **newrelic:** remove reference to pointer for http transport config

<a name="v0.7.0"></a>
## [v0.7.0] - 2020-01-23
### Features
- **newrelic:** add ConfigOptions for logging
- **newrelic:** add the ability to configure base URLs per API

### Refactor
- **newrelic:** incorporate code review feedback

<a name="v0.6.0"></a>
## [v0.6.0] - 2020-01-22
### Features
- **alerts:** add GetSyntheticsCondition method ([#105](https://github.com/newrelic/newrelic-client-go/issues/105))

<a name="v0.5.1"></a>
## [v0.5.1] - 2020-01-21
### Bug Fixes
- **alerts:** custom unmarshal of channel configuration Headers and Payload fields ([#102](https://github.com/newrelic/newrelic-client-go/issues/102))

<a name="v0.5.0"></a>
## [v0.5.0] - 2020-01-16
### Documentation Updates
- **newrelic:** update API key configuration documentation

### Refactor
- **newrelic:** validate that at least one API key is provided

<a name="v0.4.0"></a>
## [v0.4.0] - 2020-01-15
### Bug Fixes
- retry HTTP requests on 429 status codes

### Features
- **entities:** add entities search and entity tagging

### Refactor
- update test helpers to use new mock server, consistent patterns in tests

<a name="v0.3.0"></a>
## [v0.3.0] - 2020-01-13
### Bug Fixes
- make use of ErrorNotFound type for Get methods that are based on List methods
- add policy ID to alert condition

### Documentation Updates
- update example
- **build:** Update README for commit message format
- **changelog:** Add auto-generation of CHANGELOG from git comments via `make changelog`

### Features
- add top-level logging package for convenience
- add option for JSON logging and fail gracefully when log level cannot be parsed
- introduce logging
- update monitor scripts with return design pattern, update tests

### Refactor
- update alerts incidents to follow return design pattern, parallelize and use assert lib in alert incidents tests
- update ListDashboards to return array of pointers, update Dashboard test to use assert
- update ListApplications to return array of pointers, update tests to use assert
- update alert channels to return array of pointers, update tests to use assert
- update synthetics conditions to return array of pointers
- use require lib for dashboards integration tests
- refactor to package-based types files
- remove config pointer references
- remove unnecessary else
- create a logger instance per package
- move logging config code into logging package
- use centralized test helpers and remove old ones
- rescope vars for integration tests to avoid variable name conflicts
- remove redundant 'alert' from file names
- remove redundant 'Alert' from naming convention
- update monitors to use return design pattern where applicable, update tests
- incorporate code review feedback
- consistent use of pointers for &reqBody structs
- **alerts:** Spike example of changes to the mock setup
- **alerts:** Update mock server format, continue to have pkg helper
- **config:** Change Region to a string, then parse with region package
- **newrelic:** Extract config setting to opts ... format
- **region:** Move region out of config into package, add Parse(string)

<a name="v0.2.0"></a>
## [v0.2.0] - 2020-01-08
### Documentation Updates
- update readme example

<a name="v0.1.0"></a>
## v0.1.0 - 2020-01-07
### Bug Fixes
- rename variables to fix redeclared error
- update unit tests to use new method sigs
- fix monitor ID type and GetMonitor URL
- http client needs to handle other 'success' response status codes such as 201
- add godoc as a dep, and a warning about GOPATH and godoc
- fix paging bug for v2 API
- **lint:** formatting fixes for linter

### Documentation Updates
- add alerts package docs
- temporarily checking in broken import paths in generated markdown docs
- add inline documentation
- add badges to README
- fill in missing inline documentation
- document some methods

### Features
- add DeletePluginCondition
- add CreatePluginCondition
- add UpdatePluginCondition
- add GetPluginCondition
- add ListPluginsConditions
- encode monitor script text
- add ability to use 'detailed' query param in ListPlugins method
- add GetPlugin
- add ListPlugins
- publicly expose error types
- finish components endpoints
- add Components
- add internal utils package, move IntArrayToString() util to new home
- add integration tests for key transactions
- add query param filters for ListKeyTransactions
- add GetKeyTransaction
- add ListKeyTransactions
- add DeleteLabel
- add CreateLabel
- add ListLabels, add GetLabel
- add DeleteDeployment
- add CreateDeployment
- add ListDeployments
- centralize apm test helpers
- add DeleteNrqlAlertCondition
- add UpdateNrqlAlertCondition
- add CreateNrqlAlertCondition
- add GetNrqlAlertCondition
- add ListNrqlAlertConditions
- add UpdateAlertPolicy
- add DeleteAlertCondition
- add CreateAlertCondition
- add GetAlertCondition
- add ListAlertConditions
- get infra condition integration tests passing
- add InfrastructureConditions
- add MonitorScripts
- add MonitorScript
- add DeleteAlertPolicyChannel, update unit tests, add integration test (might need to remove this)
- add alert policy channels
- add synthetics alert conditions
- add synthetics alert conditions
- add GetAlertChannel method
- add CreateAlertChannel, ListAlertChannels, DeleteAlertChannel
- add DeleteMonitor
- add UpdateMonitor
- add CreateMonitor
- add dashboards
- add DeleteAlertPolicy method
- add UpdateAlertPolicy method
- add CreateAlertPolicy method
- add GetAlertPolicy method
- add ListAlertPolicies method
- alerts package
- create remaining CRUD methods for application resource
- add new dependency-free client implementation
- add version.go per auto-versioning docs
- add ListAlertConditions for infrastructure
- add infra namespace
- add catchall newrelic package
- add New Relic environment enum
- maximize page size for ListMonitors
- add ListMonitors method for Synthetics monitors
- add application filtering for ListApplications
- get TestListApplications passing

### Refactor
- updates per code review
- use proper noun Plugins in naming convention
- update key txns to use new query string parsing mechanism
- simplify integration test scenarios for components
- move components to the plugins package
- move query string parsing to an external package
- represent query params as a struct rather than a map
- return slices of pointers instead of slices of structs
- simplify parameter handling logic
- optimize IntArrayToString() per review, add test cases
- add integration tests, update unit tests, links should be a pointer for omission
- Makefile cleanup
- optimize pushing to array of pointers
- refactor synthetics conditions to established patterns
- refactor alerts package to established patterns
- refactor synthetics package to established patterns
- update local var names for consistency
- update Epoch to EpochTime
- remove redundant 'Alert' from naming convention
- remove pointer from AlertChannelConfiguration
- utilize testify assert library, other minor refactors
- refactor unit tests to use testify assertions
- add concrete types for field with known possible values
- use Epoch type for date types instead of int64
- consolidate request body structs into one alertPolicyRequestBody
- no pointers for param fields
- integrate new http client
- simplify HTTP method signatures
- add the remaining HTTP methods
- rename the new client types
- remove the old resty-based client
- put new client in place for all resources
- make ListApplications use the new client
- move version into its own internal package for now
- incorporate linter suggestions
- clean up the configuration API for NewRelicClient
- restructuring project files
- extract cross cutting concern for apm resources
- extract paging implementation
- rename packages for clarity, promote Config to the public package

[Unreleased]: https://github.com/newrelic/newrelic-client-go/compare/v0.34.0...HEAD
[v0.34.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.33.2...v0.34.0
[v0.33.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.33.1...v0.33.2
[v0.33.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.33.0...v0.33.1
[v0.33.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.32.1...v0.33.0
[v0.32.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.32.0...v0.32.1
[v0.32.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.3...v0.32.0
[v0.31.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.2...v0.31.3
[v0.31.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.1...v0.31.2
[v0.31.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.31.0...v0.31.1
[v0.31.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.30.2...v0.31.0
[v0.30.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.30.1...v0.30.2
[v0.30.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.30.0...v0.30.1
[v0.30.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.29.1...v0.30.0
[v0.29.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.29.0...v0.29.1
[v0.29.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.28.1...v0.29.0
[v0.28.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.28.0...v0.28.1
[v0.28.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.27.1...v0.28.0
[v0.27.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.27.0...v0.27.1
[v0.27.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.26.0...v0.27.0
[v0.26.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.25.1...v0.26.0
[v0.25.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.25.0...v0.25.1
[v0.25.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.24.1...v0.25.0
[v0.24.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.24.0...v0.24.1
[v0.24.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.4...v0.24.0
[v0.23.4]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.3...v0.23.4
[v0.23.3]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.2...v0.23.3
[v0.23.2]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.1...v0.23.2
[v0.23.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.23.0...v0.23.1
[v0.23.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.22.0...v0.23.0
[v0.22.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.21.1...v0.22.0
[v0.21.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.21.0...v0.21.1
[v0.21.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.20.1...v0.21.0
[v0.20.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.20.0...v0.20.1
[v0.20.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.19.0...v0.20.0
[v0.19.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.18.0...v0.19.0
[v0.18.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.17.1...v0.18.0
[v0.17.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.17.0...v0.17.1
[v0.17.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.16.0...v0.17.0
[v0.16.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.15.0...v0.16.0
[v0.15.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.14.0...v0.15.0
[v0.14.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.13.0...v0.14.0
[v0.13.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.12.0...v0.13.0
[v0.12.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.11.0...v0.12.0
[v0.11.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.10.1...v0.11.0
[v0.10.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.10.0...v0.10.1
[v0.10.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.9.0...v0.10.0
[v0.9.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.7.1...v0.8.0
[v0.7.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.6.0...v0.7.0
[v0.6.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.5.1...v0.6.0
[v0.5.1]: https://github.com/newrelic/newrelic-client-go/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.1.0...v0.2.0
