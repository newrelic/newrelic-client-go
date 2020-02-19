<a name="unreleased"></a>
## [Unreleased]


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


<a name="v0.7.0"></a>
## [v0.7.0] - 2020-01-23
### Features
- **newrelic:** add ConfigOptions for logging
- **newrelic:** add the ability to configure base URLs per API


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


<a name="v0.4.0"></a>
## [v0.4.0] - 2020-01-15
### Bug Fixes
- retry HTTP requests on 429 status codes

### Features
- **entities:** add entities search and entity tagging


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


[Unreleased]: https://github.com/newrelic/newrelic-client-go/compare/v0.10.0...HEAD
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
