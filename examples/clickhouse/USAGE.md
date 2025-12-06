# ClickHouse Driver Usage

The ClickHouse example demonstrates how to use sqlc with both the native ClickHouse driver and the standard `database/sql` package.

## Driver Options

### Native Driver (Recommended)

The native `github.com/ClickHouse/clickhouse-go/v2` driver provides better performance and type support.

**Configuration (`sqlc.yaml`):**
```yaml
version: "2"
sql:
  - engine: "clickhouse"
    queries: "queries.sql"
    schema: "schema.sql"
    gen:
      go:
        out: "gen"
        package: "db"
        sql_package: "clickhouse/v2"              # Use native driver
        emit_pointers_for_null_types: true        # Use *T for nullable fields
```

**Usage:**
```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ClickHouse/clickhouse-go/v2"
    "your-project/gen"
)

func main() {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{"localhost:9000"},
        Auth: clickhouse.Auth{
            Database: "default",
            Username: "default",
            Password: "",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    queries := db.New(conn)

    // Use the generated methods
    ctx := context.Background()
    user, err := queries.GetUser(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("User: %s (%s)\n", *user.Name, *user.Email)
}
```

**Generated Code:**
```go
// db.go
type DBTX interface {
    Exec(ctx context.Context, query string, args ...any) error
    Query(ctx context.Context, query string, args ...any) (driver.Rows, error)
    QueryRow(ctx context.Context, query string, args ...any) driver.Row
}

// models.go
type SqlcExampleUser struct {
    ID        *uint32      // Pointer type for nullable field
    Name      *string
    Email     *string
    CreatedAt *time.Time
}
```

### Database/SQL (Standard Library)

Use the standard `database/sql` package for compatibility.

**Configuration (`sqlc.yaml`):**
```yaml
version: "2"
sql:
  - engine: "clickhouse"
    queries: "queries.sql"
    schema: "schema.sql"
    gen:
      go:
        out: "gen"
        package: "db"
        # sql_package defaults to "database/sql" when not specified
```

**Usage:**
```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"

    _ "github.com/ClickHouse/clickhouse-go/v2"
    "your-project/gen"
)

func main() {
    db, err := sql.Open("clickhouse", "clickhouse://localhost:9000/default")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    queries := db.New(db)

    // Use the generated methods
    ctx := context.Background()
    user, err := queries.GetUser(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("User: %s (%s)\n", user.Name.String, user.Email.String)
}
```

**Generated Code:**
```go
// db.go
type DBTX interface {
    ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    PrepareContext(context.Context, string) (*sql.Stmt, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// models.go
type SqlcExampleUser struct {
    ID        sql.NullInt64    // sql.Null* types for nullable fields
    Name      sql.NullString
    Email     sql.NullString
    CreatedAt sql.NullTime
}
```

## Key Differences

| Feature | Native Driver (`clickhouse/v2`) | Database/SQL |
|---------|--------------------------------|--------------|
| **Import** | `github.com/ClickHouse/clickhouse-go/v2` | `database/sql` |
| **Connection** | `clickhouse.Open()` | `sql.Open("clickhouse", dsn)` |
| **Method Signatures** | `Query(ctx, query, args...)` | `QueryContext(ctx, query, args...)` |
| **Null Types** | `*string`, `*int32` (with `emit_pointers_for_null_types`) | `sql.NullString`, `sql.NullInt32` |
| **Performance** | Better (native protocol) | Standard (uses driver internally) |
| **Type Safety** | Better ClickHouse type mapping | Generic sql.Null* types |

## Recommendations

- **Use native driver** (`clickhouse/v2`) for new projects - better performance and type support
- **Enable** `emit_pointers_for_null_types: true` for cleaner nullable field handling
- **Use database/sql** only if you need compatibility with generic SQL tooling

## Example Schema

```sql
CREATE TABLE IF NOT EXISTS sqlc_example.users (
    id UInt32,
    name String,
    email String,
    created_at DateTime
) ENGINE = MergeTree()
ORDER BY id;
```

## Example Query

```sql
-- name: GetUser :one
SELECT id, name, email, created_at
FROM sqlc_example.users
WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO sqlc_example.users (id, name, email, created_at)
VALUES (?, ?, ?, ?);
```

## Testing

The generated code can be tested with the ClickHouse server:

```bash
# Start ClickHouse (using Docker)
docker run -d --name clickhouse -p 9000:9000 clickhouse/clickhouse-server

# Run your application
go run main.go
```
