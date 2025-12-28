# Plugin-Based Code Generation Example

This example demonstrates how to use **custom database engine plugins** and **custom code generation plugins** with sqlc.

## Overview

```
┌─────────────────────────────────────────────────────────────────┐
│  sqlc generate                                                  │
│                                                                 │
│  1. Read schema.sql & queries.sql                               │
│  2. Send to sqlc-engine-sqlite3 (custom DB engine)              │
│  3. Get AST & catalog                                           │
│  4. Send to sqlc-gen-rust (custom codegen)                      │
│  5. Get generated Rust code                                     │
└─────────────────────────────────────────────────────────────────┘
```

## Structure

```
plugin-based-codegen/
├── go.mod                 # This module depends on sqlc
├── sqlc.yaml              # Configuration
├── schema.sql             # Database schema (SQLite3)
├── queries.sql            # SQL queries
├── plugin_test.go         # Integration test
├── plugins/
│   ├── sqlc-engine-sqlite3/   # Custom database engine plugin
│   │   └── main.go
│   └── sqlc-gen-rust/         # Custom code generator plugin
│       └── main.go
└── gen/
    └── rust/
        └── queries.rs     # ✅ Generated Rust code
```

## Quick Start

### 1. Build the plugins

```bash
cd plugins/sqlc-engine-sqlite3 && go build -o sqlc-engine-sqlite3 .
cd ../sqlc-gen-rust && go build -o sqlc-gen-rust .
cd ../..
```

### 2. Run tests

```bash
go test -v ./...
```

### 3. Generate code (requires sqlc with plugin support)

```bash
SQLCDEBUG=processplugins=1 sqlc generate
```

## How It Works

### Database Engine Plugin (`sqlc-engine-sqlite3`)

The engine plugin implements the `pkg/engine.Handler` interface:

```go
import "github.com/sqlc-dev/sqlc/pkg/engine"

func main() {
    engine.Run(engine.Handler{
        Parse:             handleParse,      // Parse SQL
        GetCatalog:        handleGetCatalog, // Return initial catalog
        IsReservedKeyword: handleIsReservedKeyword,
        GetCommentSyntax:  handleGetCommentSyntax,
        GetDialect:        handleGetDialect,
    })
}
```

Communication: **Protobuf over stdin/stdout**

### Code Generation Plugin (`sqlc-gen-rust`)

The codegen plugin uses the `pkg/plugin.Run` helper:

```go
import "github.com/sqlc-dev/sqlc/pkg/plugin"

func main() {
    plugin.Run(func(req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
        // Generate Rust code from req.Queries and req.Catalog
        return &plugin.GenerateResponse{
            Files: []*plugin.File{{Name: "queries.rs", Contents: rustCode}},
        }, nil
    })
}
```

Communication: **Protobuf over stdin/stdout**

## Compatibility

Both plugins import public packages from sqlc:

- `github.com/sqlc-dev/sqlc/pkg/engine` - Engine plugin SDK
- `github.com/sqlc-dev/sqlc/pkg/plugin` - Codegen plugin SDK

**Compile-time compatibility**: If the plugin compiles, it's compatible with this version of sqlc.

## Configuration

```yaml
version: "2"

engines:
  - name: sqlite3
    process:
      cmd: ./plugins/sqlc-engine-sqlite3/sqlc-engine-sqlite3

plugins:
  - name: rust
    process:
      cmd: ./plugins/sqlc-gen-rust/sqlc-gen-rust

sql:
  - engine: sqlite3      # Use custom engine
    schema: "schema.sql"
    queries: "queries.sql"
    codegen:
      - plugin: rust     # Use custom codegen
        out: gen/rust
```

## Generated Code Example

The `sqlc-gen-rust` plugin generates type-safe Rust code from SQL:

**Input (`queries.sql`):**
```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO users (id, name, email) VALUES (?, ?, ?);
```

**Output (`gen/rust/queries.rs`):**
```rust
use sqlx::{FromRow, SqlitePool};
use anyhow::Result;

#[derive(Debug, FromRow)]
pub struct Users {
    pub id: i32,
    pub name: String,
    pub email: String,
}

pub async fn get_user(pool: &SqlitePool, id: i32) -> Result<Option<Users>> {
    const QUERY: &str = "SELECT * FROM users WHERE id = ?";
    let row = sqlx::query_as(QUERY)
        .bind(id)
        .fetch_optional(pool)
        .await?;
    Ok(row)
}

pub async fn create_user(pool: &SqlitePool, id: i32, name: String, email: String) -> Result<()> {
    const QUERY: &str = "INSERT INTO users (id, name, email) VALUES (?, ?, ?)";
    sqlx::query(QUERY)
        .bind(id)
        .bind(name)
        .bind(email)
        .execute(pool)
        .await?;
    Ok(())
}
```

## See Also

- [Engine Plugins Documentation](../../docs/howto/engine-plugins.md)
- [Codegen Plugins Documentation](../../docs/howto/plugins.md)

