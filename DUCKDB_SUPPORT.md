# DuckDB Support for sqlc

This document describes the DuckDB engine implementation for sqlc.

## Overview

DuckDB support has been added to sqlc using a database-backed approach, similar to PostgreSQL's analyzer pattern. Unlike MySQL and SQLite which use Go-based catalogs, DuckDB relies entirely on database connections for type inference and schema information.

## Implementation Details

### Core Components

1. **Parser** (`/internal/engine/duckdb/parse.go`)
   - **Minimal pass-through parser** - does not parse SQL into AST
   - All parsing and validation happens in the database via the analyzer
   - Implements the `Parser` interface with `Parse()`, `CommentSyntax()`, and `IsReservedKeyword()` methods
   - Supports `--` and `/* */` comment styles (DuckDB standard)
   - Returns TODO AST nodes - actual parsing done by DuckDB database

2. **Catalog** (`/internal/engine/duckdb/catalog.go`)
   - Minimal catalog implementation
   - Sets "main" as the default schema and "memory" as the default catalog
   - Does not include pre-generated types/functions (database-backed only)

3. **Analyzer** (`/internal/engine/duckdb/analyzer/analyze.go`)
   - **REQUIRED** for DuckDB engine (not optional like PostgreSQL)
   - Connects to DuckDB database via `github.com/marcboeker/go-duckdb`
   - Uses PREPARE and DESCRIBE to analyze queries
   - Handles all SQL parsing and validation via the database
   - Queries column metadata from prepared statements
   - Normalizes DuckDB types to sqlc-compatible types

4. **Reserved Keywords** (`/internal/engine/duckdb/reserved.go`)
   - DuckDB reserved keywords based on official documentation
   - Includes LAMBDA (reserved as of DuckDB 1.3.0)
   - Can be queried from DuckDB using `SELECT * FROM duckdb_keywords()`

## Configuration

### Engine Registration

Added `EngineDuckDB` constant to `/internal/config/config.go`:
```go
const (
    EngineDuckDB     Engine = "duckdb"
    EngineMySQL      Engine = "mysql"
    EnginePostgreSQL Engine = "postgresql"
    EngineSQLite     Engine = "sqlite"
)
```

### Compiler Integration

Registered in `/internal/compiler/engine.go` with required database analyzer:
```go
case config.EngineDuckDB:
    c.parser = duckdb.NewParser()
    c.catalog = duckdb.NewCatalog()
    c.selector = newDefaultSelector()
    // DuckDB requires database analyzer
    if conf.Database == nil {
        return nil, fmt.Errorf("duckdb engine requires database configuration")
    }
    if conf.Analyzer.Database == nil || *conf.Analyzer.Database {
        c.analyzer = analyzer.Cached(
            duckdbanalyze.New(c.client, *conf.Database),
            combo.Global,
            *conf.Database,
        )
    }
```

## Usage Example

### sqlc.yaml Configuration

```yaml
version: "2"
sql:
  - name: "basic"
    engine: "duckdb"
    schema: "schema/"
    queries: "query/"
    database:
      uri: ":memory:"  # or path to .db file
    gen:
      go:
        out: "db"
        package: "db"
        emit_json_tags: true
        emit_interface: true
```

### Schema Example

```sql
CREATE TABLE authors (
  id   INTEGER PRIMARY KEY,
  name VARCHAR NOT NULL,
  bio  TEXT
);
```

### Query Example

```sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :exec
INSERT INTO authors (name, bio)
VALUES ($1, $2);
```

## Key Differences from Other Engines

### vs PostgreSQL
- **PostgreSQL**: Optional database analyzer, rich Go-based catalog with pg_catalog, full AST parsing
- **DuckDB**: Required database analyzer, minimal catalog, no AST parsing (database validates SQL)

### vs MySQL/SQLite
- **MySQL/SQLite**: Go-based catalog with built-in functions, TiDB/ANTLR parser with full AST
- **DuckDB**: Database-backed only, no Go-based catalog, minimal parser (database parses SQL)

### Unique Approach
DuckDB is the only engine that doesn't parse SQL in Go. All SQL parsing, validation, and type checking happens directly in the DuckDB database. This ensures 100% compatibility with DuckDB's SQL syntax without needing to maintain a separate parser.

## Type Mapping

DuckDB types are normalized to sqlc-compatible types:

| DuckDB Type | sqlc Type |
|-------------|-----------|
| INTEGER, INT, INT4 | integer |
| BIGINT, INT8, LONG | bigint |
| SMALLINT, INT2, SHORT | smallint |
| TINYINT, INT1 | tinyint |
| DOUBLE, FLOAT8 | double |
| REAL, FLOAT4, FLOAT | real |
| VARCHAR, TEXT, STRING | varchar |
| BOOLEAN, BOOL | boolean |
| DATE | date |
| TIME | time |
| TIMESTAMP | timestamp |
| TIMESTAMPTZ | timestamptz |
| BLOB, BYTEA, BINARY | bytea |
| UUID | uuid |
| JSON | json |
| DECIMAL, NUMERIC | decimal |

## Dependencies

Added to `go.mod`:
```go
github.com/marcboeker/go-duckdb v1.8.5
```

## Setup Instructions

1. **Install dependencies** (requires network access):
   ```bash
   go mod tidy
   ```

2. **Build sqlc**:
   ```bash
   go build ./cmd/sqlc
   ```

3. **Run code generation**:
   ```bash
   ./sqlc generate
   ```

## Testing

An example project is provided in `/examples/duckdb/basic/` with:
- Schema definitions
- Sample queries
- sqlc.yaml configuration

To test:
```bash
cd examples/duckdb/basic
sqlc generate
```

## Database Requirements

DuckDB engine **requires** a database connection. You must configure:
```yaml
database:
  uri: "path/to/database.db"  # or ":memory:" for in-memory
```

Without this configuration, the compiler will return an error:
```
duckdb engine requires database configuration
```

## Limitations

1. **Network dependency**: Requires network access to download go-duckdb initially
2. **Parameter type inference**: DuckDB doesn't provide parameter types without execution, so parameters are typed as "any" by the analyzer
3. **Database required**: Unlike other engines, DuckDB cannot generate code without a database connection (no offline mode)

## Future Enhancements

1. Improve parameter type inference by analyzing query patterns
2. Add support for DuckDB-specific types (STRUCT, LIST, UNION, MAP)
3. Support DuckDB extensions and extension-specific functions
4. Add DuckDB-specific selector for custom column handling
5. Improve error messages with DuckDB-specific error codes
6. Cache database connections for better performance
7. Support managed databases via database manager

## Files Modified/Created

### Created:
- `/internal/engine/duckdb/parse.go` - Minimal pass-through parser
- `/internal/engine/duckdb/catalog.go` - Minimal catalog
- `/internal/engine/duckdb/reserved.go` - Reserved keywords
- `/internal/engine/duckdb/analyzer/analyze.go` - Database analyzer
- `/examples/duckdb/basic/schema/schema.sql` - Example schema
- `/examples/duckdb/basic/query/query.sql` - Example queries
- `/examples/duckdb/basic/sqlc.yaml` - Example configuration

### Modified:
- `/internal/config/config.go` - Added `EngineDuckDB` constant
- `/internal/compiler/engine.go` - Registered DuckDB engine with analyzer
- `/go.mod` - Added `github.com/marcboeker/go-duckdb v1.8.5`

## Notes

- **No SQL parsing in Go**: DuckDB engine validates all SQL via the database, not in Go code
- DuckDB uses "main" as the default schema (different from PostgreSQL's "public")
- DuckDB uses "memory" as the default catalog name
- Comment syntax supports only `--` and `/* */`, not `#`
- Reserved keyword LAMBDA was added in DuckDB 1.3.0
- Reserved keyword GRANT was removed in DuckDB 1.3.0
- 100% compatibility with DuckDB syntax since the database itself parses SQL
