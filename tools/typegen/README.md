# typegen

Generate Golang types from GraphQL schema introspection

## Getting Started

1. Include `typegen` in your project as a tool / install in some fashion.
1. Add a `./path/to/package/typegen.yaml` configuration with the type you want generated:

   ```yaml
   ---
   types:
     - name: MyGraphQLTypeName
     - name: AnotherTypeInGraphQL
       createAs: map[string]int
   ```
1. Add a generation command inside the `main.go` (or equivalent)

   ```go
   // Package CoolPackage provides cool stuff, based on generated types
   //go:generate ./path/to/typegen -p $GOPACKAGE
   package CoolPackage
   // ... implementation ...
   ```
1. Run `go generate`
1. Add the `./path/to/package/types.go` file to your repo

# Configuration

## Command Flags

Flags for running the typegen command:

| Flag | Description |
| ---- | ----------- |
| `-p <Package Name>` | Package name used within the generated file. Overrides the configuration file. |
| `-v` | Enable verbose logging |


## Per-package

Configuration on what types to generate, and any overrides from the schema
exist within the package directory in a file named `typegen.yaml`. The file has
a simple configuration format, and includes the following sections:

### types

Types is a list of the types to explicity generate.  Any required sub-type will
also be generated until we hit a Golang type.

| Name | Required | Description |
| ---- | -------- | ----------- |
| name | Yes | The name of the field to search for and create |
| package | No | Name of the package the output file will be part of (see `-p` flag) |
| createAs | No | If you want to override the type that is created, use this to explicitly name the type |

**ORDER MATTERS:** Add types with overrides first, otherwise they might not get
created as you expect. If A => B => gotype, and you want to override B, you
must configure it first.  If you configure A first, B will be generated as a
dependency before we create B via configuration.

**Example:**

```yaml
---
types:
  - name: TheName
    createAs: int
  - name: ComplexType
    createAs: map[string]interface{}
  - name: AnotherName
```
