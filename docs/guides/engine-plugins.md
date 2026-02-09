# External Database Engines (Engine Plugins)

Engine plugins let you use sqlc with databases that are not built-in. You can add support for other SQL-compatible systems (e.g. CockroachDB, TiDB, or custom engines) by implementing a small external program that parses SQL and returns parameters and result columns.

## Why use an engine plugin?

- Use sqlc with a database that doesn't have native support.
- Reuse an existing SQL parser or dialect in a separate binary.
- Keep engine-specific logic outside the sqlc core.

Data returned by the engine plugin (SQL text, parameters, columns) is passed through to [codegen plugins](plugins.md) without an extra compiler/AST step. The plugin is the single place that defines how queries are interpreted for that engine.

**Limitation:** `sqlc vet` does not support plugin engines. Use vet only with built-in engines (postgresql, mysql, sqlite).

## Overview

An engine plugin is an external process that implements one RPC:

- **Parse** — accepts the **entire contents** of one query file (e.g. `query.sql`) and either schema SQL or connection parameters; returns **one Statement per query block** in that file (each with sql, parameters, columns, and name/cmd).

Process plugins (e.g. written in Go) talk to sqlc over **stdin/stdout** using **Protocol Buffers**. The protocol is defined in `protos/engine/engine.proto`.

## Compatibility

For Go plugins, compatibility is enforced at **compile time** by importing the engine package:

```go
import "github.com/sqlc-dev/sqlc/pkg/engine"
```

- If the plugin builds, it matches this version of the engine API.
- If the API changes in a breaking way, the plugin stops compiling until it's updated.

No version handshake is required; the proto schema defines the contract.

## Configuration

### sqlc.yaml

```yaml
version: "2"

engines:
  - name: external-db
    process:
      # Executable and optional arguments (e.g. --dont-open-wildcard-star).
      # First token is the command; the rest are passed to the plugin before the RPC method name.
      cmd: sqlc-engine-external-db --dont-open-wildcard-star
    env:
      - EXTERNAL_DB_DSN

sql:
  - engine: external-db
    schema: "schema.sql"
    queries: "queries.sql"
    codegen:
      - plugin: go
        out: db
```

### Engine options

| Field | Description |
|-------|-------------|
| `name` | Engine name used in `sql[].engine` |
| `process.cmd` | Command to run: executable path and optional arguments (e.g. `sqlc-engine-external-db --dont-open-wildcard-star`). First token is the executable; remaining tokens are passed as arguments before the RPC method. |
| `env` | Environment variable names passed to the plugin |

Each engine must define either `process` (with `cmd`) or `wasm` (with `url` and `sha256`). See [Configuration reference](../reference/config.md) for the full `engines` schema.

### How sqlc finds the process plugin

For an engine with `process.cmd`, sqlc resolves and runs the plugin as follows:

1. **Command parsing** — `process.cmd` is split on whitespace. The first token is the executable; any further tokens are passed as arguments, and sqlc appends the RPC method name (`parse`) when invoking the plugin.

2. **Executable lookup** — The first token is resolved the same way as in the shell:
   - If it contains a path separator (e.g. `/usr/bin/sqlc-engine-external-db` or `./bin/sqlc-engine-external-db`), it is treated as a path. Absolute paths are used as-is; relative paths are taken relative to the **current working directory of the process running sqlc**.
   - If it has no path separator, the executable is looked up in the **PATH** of the process running sqlc. The plugin binary must be on PATH (e.g. after `go install` or adding its directory to PATH) or `process.cmd` must be an absolute path.

3. **Working directory** — The plugin process is started with its working directory set to the **directory containing the sqlc config file**. That directory is used for resolving relative paths inside the plugin, not for resolving `process.cmd` itself.

If the executable cannot be found or `process.cmd` is empty, sqlc reports an error and refers to this documentation.

## Implementing an engine plugin (Go)

### 1. Dependencies and entrypoint

```go
package main

import "github.com/sqlc-dev/sqlc/pkg/engine"

func main() {
    engine.Run(engine.Handler{
        PluginName:    "external-db",
        PluginVersion: "1.0.0",
        Parse:         handleParse,
    })
}
```

The engine API exposes only **Parse**. There are no separate methods for catalog, keywords, comment syntax, or dialect.

### 2. Parse

sqlc calls Parse **once per query file** (e.g. once for `query.sql`). The plugin receives the full file contents and returns one **Statement** per query block in that file. sqlc then passes each statement to the codegen plugin as a separate query.

**Request**

- `sql` — The **entire contents** of one query file (all query blocks, with `-- name: X :one`-style comments).
- `schema_source` — One of:
  - `schema_sql`: full schema as in schema.sql (for schema-based parsing).
  - `connection_params`: DSN and options for database-only mode.

**Response**

Return `statements`: one `Statement` per query block. Each `Statement` has:

- `name` — Query name (from `-- name: GetUser` etc.).
- `cmd` — Command/type: use the `Cmd` enum (`engine.Cmd_CMD_ONE`, `engine.Cmd_CMD_MANY`, `engine.Cmd_CMD_EXEC`, etc.). See `protos/engine/engine.proto` for the full list.
- `sql` — Processed SQL for that block (as-is or with `*` expanded using schema).
- `parameters` — Parameters for this statement.
- `columns` — Result columns (names, types, nullability, etc.) for this statement.

You may also return **`catalog`** (optional). When present, sqlc passes it to the codegen plugin so codegen can emit model structs from the schema (tables/columns), not only per-query row types. Build an `engine.Catalog` from your schema (e.g. from `schema_sql` or DB metadata): one or more **CatalogSchema**, each with **CatalogTable** (rel = Identifier, columns = **CatalogColumn**). See `protos/engine/engine.proto` for the `Catalog`, `CatalogSchema`, `CatalogTable`, `CatalogColumn`, and `Identifier` messages.

The engine package provides helpers (optional) to split `query.sql` and parse `"-- name: X :cmd"` lines in the same way as the built-in engines:

- `engine.CommentSyntax` — Which comment styles to accept (`Dash`, `SlashStar`, `Hash`).
- `engine.ParseNameAndCmd(line, syntax)` — Parses a single line like `"-- name: ListAuthors :many"` → `(name, cmd, ok)`. `cmd` is `engine.Cmd`.
- `engine.QueryBlocks(content, syntax)` — Splits file content into `[]engine.QueryBlock` (each has `Name`, `Cmd`, `SQL`).
- `engine.StatementMeta(name, cmd, sql)` — Builds a `*engine.Statement` with name/cmd/sql set; you add parameters and columns.

Example handler using helpers:

```go
func handleParse(req *engine.ParseRequest) (*engine.ParseResponse, error) {
    queryFileContent := req.GetSql()
    syntax := engine.CommentSyntax{Dash: true, SlashStar: true, Hash: true}

    var schema *SchemaInfo
    if s := req.GetSchemaSql(); s != "" {
        schema = parseSchema(s)
    }
    // Or use req.GetConnectionParams() for database-only mode.

    blocks, _ := engine.QueryBlocks(queryFileContent, syntax)
    var statements []*engine.Statement
    for _, b := range blocks {
        st := engine.StatementMeta(b.Name, b.Cmd, processSQL(b.SQL, schema))
        st.Parameters = extractParameters(b.SQL)
        st.Columns = extractColumns(b.SQL, schema)
        statements = append(statements, st)
    }
    resp := &engine.ParseResponse{Statements: statements}
    if cat := buildCatalogFromSchema(schema); cat != nil {
        resp.Catalog = cat
    }
    return resp, nil
}
```

Parameter and column types use the `Parameter` and `Column` messages in `engine.proto` (name, position, data_type, nullable, is_array, array_dims; for columns, table_name and schema_name are optional). If you return a **Catalog** (tables/columns from schema or DB), sqlc forwards it to the codegen plugin so generated code can include model types from the schema.

Support for sqlc placeholders (`sqlc.arg()`, `sqlc.narg()`, `sqlc.slice()`, `sqlc.embed()`) is up to the plugin: it can parse and map them into `parameters` (and schema usage) as needed.

#### Catalog (optional)

If the plugin returns a non-empty **`catalog`** in `ParseResponse`, sqlc converts it to the codegen plugin’s Catalog and passes it in the GenerateRequest. Codegen plugins can then emit model structs from the schema (e.g. `type Author struct`) in addition to per-query row types. Build the catalog from `schema_sql` or from live DB metadata when using `connection_params`. The `engine.proto` defines `Catalog`, `CatalogSchema`, `CatalogTable`, `CatalogColumn`, and `Identifier` for this purpose.

### 3. Build and run

```bash
go build -o sqlc-engine-external-db .
# Ensure sqlc-engine-external-db is on PATH or use an absolute path in process.cmd
```

## Protocol

Process plugins use Protocol Buffers on stdin/stdout:

```
sqlc  →  stdin (protobuf)  →  plugin  →  stdout (protobuf)  →  sqlc
```

Invocation:

```bash
sqlc-engine-external-db parse   # stdin: ParseRequest, stdout: ParseResponse
```

The definition lives in `protos/engine/engine.proto` (generated Go in `pkg/engine`). After editing the proto, run `make proto-engine-plugin` to regenerate the Go code.

## Example

The protocol and Go SDK are in this repository: `protos/engine/engine.proto` and `pkg/engine/` (including `sdk.go` with `engine.Run` and `engine.Handler`). Use them to build a binary that implements the Parse RPC; register it under `engines` in sqlc.yaml as shown above.

## Architecture

For each `sql[]` block, `sqlc generate` branches on the configured engine: built-in (postgresql, mysql, sqlite) use the compiler and catalog; any engine listed under `engines:` in sqlc.yaml uses the plugin path (no compiler). For the plugin path, sqlc calls Parse **once per query file**, sending the full file contents and schema (or connection params). The plugin returns **N statements** (one per query block) and optionally a **catalog** (tables/columns from schema or DB). sqlc passes each statement to codegen as a separate query and, when present, passes the catalog so codegen can emit model types from the schema.

```
┌─────────────────────────────────────────────────────────────────┐
│  sqlc generate (plugin engine)                                  │
│  1. Per query file: one Parse(schema_sql|connection_params,      │
│     full query file content)                                    │
│  2. ParseResponse.statements = one Statement per query block     │
│  3. Each statement → one codegen query (N helpers)               │
└─────────────────────────────────────────────────────────────────┘

    sqlc                          sqlc-engine-external-db
      │──── spawn, args: ["parse"] ──────────────────────────────► │
      │──── stdin: ParseRequest{sql=full query.sql, schema_sql|…}  ► │
      │◄─── stdout: ParseResponse{statements: [stmt1, stmt2, …]} ── │
```

## See also

- [Codegen plugins](plugins.md) — Custom code generators that consume engine output.
- [Configuration reference](../reference/config.md)
- Proto schema: `protos/engine/engine.proto`
