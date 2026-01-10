# Database Engine Plugins

sqlc supports adding custom database backends through engine plugins. This allows you to use sqlc with databases that aren't natively supported (like MyDB, CockroachDB, or other SQL-compatible databases).

## Overview

Engine plugins are external programs that implement the sqlc engine interface:
- **Process plugins** (Go): Communicate via **Protocol Buffers** over stdin/stdout
- **WASM plugins** (any language): Communicate via **JSON** over stdin/stdout

## Compatibility Guarantee

For Go process plugins, compatibility is guaranteed at **compile time**:

```go
import "github.com/sqlc-dev/sqlc/pkg/engine"
```

When you import this package:
- If your plugin compiles successfully → it's compatible with this version of sqlc
- If types change incompatibly → your plugin won't compile until you update it

The Protocol Buffer schema ensures binary compatibility. No version negotiation needed.

## Configuration

### sqlc.yaml

```yaml
version: "2"

# Define engine plugins
engines:
  - name: mydb
    process:
      cmd: sqlc-engine-mydb
    env:
      - MYDB_CONNECTION_STRING

sql:
  - engine: mydb             # Use the MyDB engine
    schema: "schema.sql"
    queries: "queries.sql"
    gen:
      go:
        package: db
        out: db
```

### Configuration Options

| Field | Description |
|-------|-------------|
| `name` | Unique name for the engine (used in `sql[].engine`) |
| `process.cmd` | Command to run (must be in PATH or absolute path) |
| `wasm.url` | URL to download WASM module (`file://` or `https://`) |
| `wasm.sha256` | SHA256 checksum of the WASM module |
| `env` | Environment variables to pass to the plugin |

## Creating a Go Engine Plugin

### 1. Import the SDK

```go
import "github.com/sqlc-dev/sqlc/pkg/engine"
```

### 2. Implement the Handler

```go
package main

import (
    "github.com/sqlc-dev/sqlc/pkg/engine"
)

func main() {
    engine.Run(engine.Handler{
        PluginName:        "mydb",
        PluginVersion:     "1.0.0",
        Parse:             handleParse,
        GetCatalog:        handleGetCatalog,
        IsReservedKeyword: handleIsReservedKeyword,
        GetCommentSyntax:  handleGetCommentSyntax,
        GetDialect:        handleGetDialect,
    })
}
```

### 3. Implement Methods

#### Parse

Parses SQL text into statements with AST.

```go
func handleParse(req *engine.ParseRequest) (*engine.ParseResponse, error) {
    sql := req.GetSql()
    // Parse SQL using your database's parser
    
    return &engine.ParseResponse{
        Statements: []*engine.Statement{
            {
                RawSql:       sql,
                StmtLocation: 0,
                StmtLen:      int32(len(sql)),
                AstJson:      astJSON, // AST encoded as JSON bytes
            },
        },
    }, nil
}
```

#### GetCatalog

Returns the initial catalog with built-in types and functions.

```go
func handleGetCatalog(req *engine.GetCatalogRequest) (*engine.GetCatalogResponse, error) {
    return &engine.GetCatalogResponse{
        Catalog: &engine.Catalog{
            DefaultSchema: "public",
            Name:          "mydb",
            Schemas: []*engine.Schema{
                {
                    Name: "public",
                    Functions: []*engine.Function{
                        {Name: "now", ReturnType: &engine.DataType{Name: "timestamp"}},
                    },
                },
            },
        },
    }, nil
}
```

#### IsReservedKeyword

Checks if a string is a reserved keyword.

```go
func handleIsReservedKeyword(req *engine.IsReservedKeywordRequest) (*engine.IsReservedKeywordResponse, error) {
    reserved := map[string]bool{
        "select": true, "from": true, "where": true,
    }
    return &engine.IsReservedKeywordResponse{
        IsReserved: reserved[strings.ToLower(req.GetKeyword())],
    }, nil
}
```

#### GetCommentSyntax

Returns supported SQL comment syntax.

```go
func handleGetCommentSyntax(req *engine.GetCommentSyntaxRequest) (*engine.GetCommentSyntaxResponse, error) {
    return &engine.GetCommentSyntaxResponse{
        Dash:      true,  // -- comment
        SlashStar: true,  // /* comment */
        Hash:      false, // # comment
    }, nil
}
```

#### GetDialect

Returns SQL dialect information for formatting.

```go
func handleGetDialect(req *engine.GetDialectRequest) (*engine.GetDialectResponse, error) {
    return &engine.GetDialectResponse{
        QuoteChar:   "`",             // Identifier quoting character
        ParamStyle:  "dollar",        // $1, $2, ...
        ParamPrefix: "$",             // Parameter prefix
        CastSyntax:  "cast_function", // CAST(x AS type) or "double_colon" for ::
    }, nil
}
```

### 4. Build and Install

```bash
go build -o sqlc-engine-mydb .
mv sqlc-engine-mydb /usr/local/bin/
```

## Protocol

### Process Plugins (Go)

Process plugins use **Protocol Buffers** for serialization:

```
sqlc → stdin (protobuf) → plugin → stdout (protobuf) → sqlc
```

The proto schema is published at `buf.build/sqlc/sqlc` in `engine/engine.proto`.

Methods are invoked as command-line arguments:
```bash
sqlc-engine-mydb parse        # stdin: ParseRequest, stdout: ParseResponse
sqlc-engine-mydb get_catalog  # stdin: GetCatalogRequest, stdout: GetCatalogResponse
```

### WASM Plugins

WASM plugins use **JSON** for broader language compatibility:

```
sqlc → stdin (JSON) → wasm module → stdout (JSON) → sqlc
```

## Full Example

See `examples/plugin-based-codegen/` for a complete engine plugin implementation.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│  sqlc generate                                                  │
│                                                                 │
│  1. Read sqlc.yaml                                              │
│  2. Find engine: mydb → look up in engines[]                    │
│  3. Run: sqlc-engine-mydb parse < schema.sql                    │
│  4. Get AST via protobuf on stdout                              │
│  5. Generate Go code                                            │
└─────────────────────────────────────────────────────────────────┘

Process Plugin Communication (Protobuf):

    sqlc                          sqlc-engine-mydb
    ────                          ────────────────
      │                                  │
      │──── spawn process ─────────────► │
      │     args: ["parse"]              │
      │                                  │
      │──── protobuf on stdin ─────────► │
      │     ParseRequest{sql: "..."}     │
      │                                  │
      │◄─── protobuf on stdout ───────── │
      │     ParseResponse{statements}    │
      │                                  │
```

## See Also

- [Codegen Plugins](plugins.md) - For custom code generators
- [Configuration Reference](../reference/config.md)
- Proto schema: `protos/engine/engine.proto`
