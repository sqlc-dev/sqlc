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

- **Parse** — accepts the query text and either schema SQL or connection parameters, and returns processed SQL, parameter list, and result columns.

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
  - name: mydb
    process:
      cmd: sqlc-engine-mydb
    env:
      - MYDB_DSN

sql:
  - engine: mydb
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
| `process.cmd` | Command to run (PATH or absolute path) |
| `env` | Environment variable names passed to the plugin |

Each engine must define either `process` (with `cmd`) or `wasm` (with `url` and `sha256`). See [Configuration reference](../reference/config.md) for the full `engines` schema.

## Implementing an engine plugin (Go)

### 1. Dependencies and entrypoint

```go
package main

import "github.com/sqlc-dev/sqlc/pkg/engine"

func main() {
    engine.Run(engine.Handler{
        PluginName:    "mydb",
        PluginVersion: "1.0.0",
        Parse:         handleParse,
    })
}
```

The engine API exposes only **Parse**. There are no separate methods for catalog, keywords, comment syntax, or dialect.

### 2. Parse

**Request**

- `sql` — The query text to parse.
- `schema_source` — One of:
  - `schema_sql`: schema as in a schema.sql file (used for schema-based parsing).
  - `connection_params`: DSN and options for database-only mode.

**Response**

- `sql` — Processed query text. Often the same as input; with a schema you may expand `*` into explicit columns.
- `parameters` — List of parameters (position/name, type, nullable, array, etc.).
- `columns` — List of result columns (name, type, nullable, table/schema if known).

Example handler:

```go
func handleParse(req *engine.ParseRequest) (*engine.ParseResponse, error) {
    sql := req.GetSql()

    var schema *SchemaInfo
    if s := req.GetSchemaSql(); s != "" {
        schema = parseSchema(s)
    }
    // Or use req.GetConnectionParams() for database-only mode.

    parameters := extractParameters(sql)
    columns := extractColumns(sql, schema)
    processedSQL := processSQL(sql, schema) // e.g. expand SELECT *

    return &engine.ParseResponse{
        Sql:        processedSQL,
        Parameters: parameters,
        Columns:    columns,
    }, nil
}
```

Parameter and column types use the `Parameter` and `Column` messages in `engine.proto` (name, position, data_type, nullable, is_array, array_dims; for columns, table_name and schema_name are optional).

Support for sqlc placeholders (`sqlc.arg()`, `sqlc.narg()`, `sqlc.slice()`, `sqlc.embed()`) is up to the plugin: it can parse and map them into `parameters` (and schema usage) as needed.

### 3. Build and run

```bash
go build -o sqlc-engine-mydb .
# Ensure sqlc-engine-mydb is on PATH or use an absolute path in process.cmd
```

## Protocol

Process plugins use Protocol Buffers on stdin/stdout:

```
sqlc  →  stdin (protobuf)  →  plugin  →  stdout (protobuf)  →  sqlc
```

Invocation:

```bash
sqlc-engine-mydb parse   # stdin: ParseRequest, stdout: ParseResponse
```

The definition lives in `engine/engine.proto` (and generated Go in `pkg/engine`).

## Example

The protocol and Go SDK are in this repository: `protos/engine/engine.proto` and `pkg/engine/` (including `sdk.go` with `engine.Run` and `engine.Handler`). Use them to build a binary that implements the Parse RPC; register it under `engines` in sqlc.yaml as shown above.

## Architecture

For each `sql[]` block, `sqlc generate` branches on the configured engine: built-in (postgresql, mysql, sqlite) use the compiler and catalog; any engine listed under `engines:` in sqlc.yaml uses the plugin path (no compiler, schema + queries go to the plugin's Parse RPC, then output goes to codegen).

```
┌─────────────────────────────────────────────────────────────────┐
│  sqlc generate                                                  │
│  1. Read sqlc.yaml, find engine for this sql block              │
│  2. If plugin engine: call plugin parse (sql + schema_sql etc.)  │
│  3. Use returned sql, parameters, columns in codegen             │
└─────────────────────────────────────────────────────────────────┘

    sqlc                          sqlc-engine-mydb
      │──── spawn, args: ["parse"] ──────────────────────────────► │
      │──── stdin: ParseRequest{sql, schema_sql|connection_params} ► │
      │◄─── stdout: ParseResponse{sql, parameters, columns} ─────── │
```

## See also

- [Codegen plugins](plugins.md) — Custom code generators that consume engine output.
- [Configuration reference](../reference/config.md)
- Proto schema: `protos/engine/engine.proto`
