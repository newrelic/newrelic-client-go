[![Community Project header](https://github.com/newrelic/open-source-office/raw/master/examples/categories/images/Community_Project.png)](https://github.com/newrelic/open-source-office/blob/master/examples/categories/index.md#category-community-project)

# newrelic-client-go

[![Testing](https://github.com/newrelic/newrelic-client-go/workflows/Testing/badge.svg)](https://github.com/newrelic/newrelic-client-go/actions)
[![Security Scan](https://github.com/newrelic/newrelic-client-go/workflows/Security%20Scan/badge.svg)](https://github.com/newrelic/newrelic-client-go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/newrelic/newrelic-client-go?style=flat-square)](https://goreportcard.com/report/github.com/newrelic/newrelic-client-go)
[![GoDoc](https://godoc.org/github.com/newrelic/newrelic-client-go?status.svg)](https://godoc.org/github.com/newrelic/newrelic-client-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/newrelic/newrelic-client-go/blob/main/LICENSE)
[![CLA assistant](https://cla-assistant.io/readme/badge/newrelic/newrelic-client-go)](https://cla-assistant.io/newrelic/newrelic-client-go)
[![Release](https://img.shields.io/github/release/newrelic/newrelic-client-go/all.svg)](https://github.com/newrelic/newrelic-client-go/releases/latest)

The New Relic Client provides the building blocks for tools in the [Developer Toolkit](https://newrelic.github.io/observability-as-code), enabling quick access to the suite of New Relic APIs. As a library, it can also be leveraged within your own custom applications.

- [Getting Started and Example Usage](#getting-started-and-example-usage)
- [Upgrade to the latest version](#upgrade-to-the-latest-version)
- [Community](#community)
- [Development](#development)
	- [Requirements](#requirements)
	- [Building](#building)
	- [Testing](#testing)
		- [Integration tests](#integration-tests)
		- [Go Version Support](#go-version-support)
	- [Commit Messages](#commit-messages)
		- [Format](#format)
		- [Scope](#scope)
	- [Documentation](#documentation)
- [Community Support](#community-support)
- [Issues / Enhancement Requests](#issues--enhancement-requests)
- [Contributing](#contributing)
- [Open Source License](#open-source-license)

<br>

## Getting Started and Example Usage

Follow the steps below to add `github.com/newrelic/newrelic-client-go` as a dependency in your Go project.

1. In the root directory of your project, run `go get github.com/newrelic/newrelic-client-go@latest`. This will update your `go.mod` file with the latest version `newrelic-client-go`.


2. Import `newrelic-client-go` in your project code.
    ```go
    package main

    import "github.com/newrelic/newrelic-client-go/v2/newrelic"

    func main() {
      // Initialize the client.
      client, err := newrelic.New(newrelic.ConfigPersonalAPIKey(os.Getenv("NEW_RELIC_API_KEY")))
      if err != nil {
        // ...
      }
    }
    ```

3. Run `go mod tidy`. This will ensure all your dependencies are in sync with your code.
4. Your module's `go.mod` file should now be updated with the latest version of the client and should look similar to the following example (version number is hypothetical in the example below).
    ```
    module example.com/yourmodule

    go 1.22

    require (
      github.com/newrelic/newrelic-client-go/v2 v2.0.1
    )
    ```
5. The example below demonstrates fetching New Relic entities.
   ```go
    package main

    import (
      "fmt"
      "os"

      log "github.com/sirupsen/logrus"

      "github.com/newrelic/newrelic-client-go/v2/newrelic"
      "github.com/newrelic/newrelic-client-go/v2/pkg/entities"
    )

    func main() {
      // Initialize the client.
      client, err := newrelic.New(newrelic.ConfigPersonalAPIKey(os.Getenv("NEW_RELIC_API_KEY")))
      if err != nil {
        log.Fatal("error initializing client:", err)
      }

      // Search the current account for entities by name and type.
      queryBuilder := entities.EntitySearchQueryBuilder{
        Name: "Example entity",
        Type: entities.EntitySearchQueryBuilderTypeTypes.APPLICATION,
      }

      entitySearch, err := client.Entities.GetEntitySearch(
        entities.EntitySearchOptions{},
        "",
        queryBuilder,
        []entities.EntitySearchSortCriteria{},
      )
      if err != nil {
        log.Fatal("error searching entities:", err)
      }

      fmt.Printf("%+v\n", entitySearch.Results.Entities)
    }
    ```


## Upgrade to the latest version

1. Run the following command to tell Go to download the latest version. You can also check the latest version out in [this page](https://pkg.go.dev/github.com/newrelic/newrelic-client-go/v2/newrelic).
   ```
   go get github.com/newrelic/newrelic-client-go/v2@latest
   ```
2. Run `go mod tidy` to sync your dependencies with your code.
3. Confirm your `go.mod` file is referencing the [latest version](https://github.com/newrelic/newrelic-client-go/releases/latest).


## Community

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices.

- [Issues or Enhancement Requests](https://github.com/newrelic/newrelic-client-go/issues/new/choose) - Issues and enhancement requests can be submitted in the Issues tab of this repository. Please search for and review the existing open issues before submitting a new issue.
- [Contributors Guide](CONTRIBUTING.md) - Contributions are welcome.
- [Community discussion board](https://discuss.newrelic.com/c/build-on-new-relic/developer-toolkit) - Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub.

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

## Development

### Requirements

- Go 1.22+
- GNU Make
- git

### Building

This package does not generate any direct usable assets (it's a library). You can still run the build scripts to validate your code, and generate coverage information.

```bash
# Default target is 'build'
$ make

# Explicitly run build
$ make build

# Locally test the CI build scripts
$ make build-ci
```

### Testing

Before contributing, all linting and tests must pass.

Tests can be run directly via:

```bash
# Tests and Linting
$ make test

# Only unit tests
$ make test-unit

# Only integration tests
$ make test-integration
```

#### Integration tests

Integration tests communicate with the New Relic API, and therefore require proper
account credentials to run properly. The suite makes use of secrets, which will
need to be configured in the environment or else the integraiton tests will be skipped
during a test run. To run the integration test suite, the following secrets will
need to be configured:

```bash
NEW_RELIC_ACCOUNT_ID
NEW_RELIC_API_KEY
NEW_RELIC_INSIGHTS_INSERT_KEY
NEW_RELIC_LICENSE_KEY
NEW_RELIC_REGION
NEW_RELIC_TEST_USER_ID
```

Optional for debugging (defaults to `debug`):

```bash
NEW_RELIC_LOG_LEVEL=trace
```

#### Go Version Support

We'll aim to support the latest supported release of Go, along with the
previous release. This doesn't mean that building with an older version of Go
will not work, but we don't intend to support a Go version in this project that
is not supported by the larger Go community. Please see the [Go
releases][go_releases] page for more details.

### Commit Messages

Using the following format for commit messages allows for auto-generation of
the [CHANGELOG](CHANGELOG.md):

#### Format

`<type>(<scope>): <subject>`

| Type       | Description           | Change log? |
| ---------- | --------------------- | :---------: |
| `chore`    | Maintenance type work |     No      |
| `docs`     | Documentation Updates |     Yes     |
| `feat`     | New Features          |     Yes     |
| `fix`      | Bug Fixes             |     Yes     |
| `refactor` | Code Refactoring      |     No      |

#### Scope

This refers to what part of the code is the focus of the work. For example:

**General:**

- `build` - Work related to the build system (linting, makefiles, CI/CD, etc)
- `release` - Work related to cutting a new release

**Package Specific:**

- `newrelic` - Work related to the New Relic package
- `http` - Work related to the `internal/http` package
- `alerts` - Work related to the `pkg/alerts` package

### Documentation

**Note:** This requires the repo to be in your GOPATH [(godoc issue)](https://github.com/golang/go/issues/26827)

```bash
$ make docs
```

### Releasing

Releases are automated via the Release Github Action on merges to the default branch.  No user interaction is required.

Using [svu](https://github.com/caarlos0/svu), commit messages are parsed to identify if a new release is needed, and to what extent.  Here's the breakdown:

| Commit message                                                                         | Tag increase |
| -------------------------------------------------------------------------------------- | ------------ |
| `fix: fixed something`                                                                 | Patch        |
| `feat: added new button to do X`                                                       | Minor        |
| `fix: fixed thing xyz`<br><br>`BREAKING CHANGE: this will break users because of blah` | Major        |
| `fix!: fixed something`                                                                | Major        |
| `feat!: added blah`                                                                    | Major        |
| `chore: foo`                                                                           | Nothing      |
| `refactor: updated bar`                                                                | Nothing      |


## Community Support

New Relic hosts and moderates an online forum where you can interact with New Relic employees as well as other customers to get help and share best practices. Like all New Relic open source community projects, there's a related topic in the New Relic Explorers Hub. You can find our team's project topic/threads here:

[Developer ToolKit](https://discuss.newrelic.com/t/about-the-developer-toolkit-category/90159)

Please do not report issues with Top to New Relic Global Technical Support. Instead, visit the [`Explorers Hub`](https://discuss.newrelic.com/c/build-on-new-relic) for troubleshooting and best-practices.

## Issues / Enhancement Requests

Issues and enhancement requests can be submitted in te [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.

## Contributing

Contributions are welcome (and if you submit a Enhancement Request, expect to be invited to contribute it yourself :grin:). Please review our [Contributors Guide](CONTRIBUTING.md).

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

## Open Source License

 This project is distributed under the [Apache 2 license](LICENSE).

[go_releases]: https://github.com/golang/go/wiki/Go-Release-Cycle
