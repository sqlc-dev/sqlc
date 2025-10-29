# DuckDB Example

This example demonstrates how to use sqlc with DuckDB.

## Overview

DuckDB is an in-process analytical database that supports PostgreSQL-compatible SQL syntax. This integration reuses sqlc's PostgreSQL parser and catalog while providing a DuckDB-specific analyzer that connects to an in-memory DuckDB instance.

## Features

- **PostgreSQL-compatible SQL**: DuckDB uses PostgreSQL-compatible syntax, so you can use familiar SQL constructs
- **In-memory database**: Perfect for testing and development
- **Type-safe Go code**: sqlc generates type-safe Go code from your SQL queries
- **Live database analysis**: The analyzer connects to a DuckDB instance to extract accurate column types

## Configuration

The `sqlc.yaml` file configures sqlc to use the DuckDB engine:

```yaml
version: "2"
sql:
  - name: "duckdb_example"
    engine: "duckdb"  # Use DuckDB engine
    schema:
      - "schema.sql"
    queries:
      - "query.sql"
    database:
      managed: false
      uri: ":memory:"  # Use in-memory database
    analyzer:
      database: true   # Enable live database analysis
    gen:
      go:
        package: "db"
        out: "db"
```

## Database URI

DuckDB supports several URI formats:

- `:memory:` - In-memory database (default if not specified)
- `file.db` - File-based database
- `/path/to/file.db` - Absolute path to database file

## Usage

1. Generate Go code:
   ```bash
   sqlc generate
   ```

2. Use the generated code in your application:
   ```go
   package main

   import (
       "context"
       "database/sql"
       "log"

       _ "github.com/marcboeker/go-duckdb"
       "yourmodule/db"
   )

   func main() {
       // Open DuckDB connection
       conn, err := sql.Open("duckdb", ":memory:")
       if err != nil {
           log.Fatal(err)
       }
       defer conn.Close()

       // Create tables
       schema := `
       CREATE TABLE users (
           id INTEGER PRIMARY KEY,
           name VARCHAR NOT NULL,
           email VARCHAR UNIQUE NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
       );
       `
       if _, err := conn.Exec(schema); err != nil {
           log.Fatal(err)
       }

       // Use generated queries
       queries := db.New(conn)
       ctx := context.Background()

       // Create a user
       user, err := queries.CreateUser(ctx, db.CreateUserParams{
           Name:  "John Doe",
           Email: "john@example.com",
       })
       if err != nil {
           log.Fatal(err)
       }

       log.Printf("Created user: %+v\n", user)

       // Get the user
       fetchedUser, err := queries.GetUser(ctx, user.ID)
       if err != nil {
           log.Fatal(err)
       }

       log.Printf("Fetched user: %+v\n", fetchedUser)
   }
   ```

## Differences from PostgreSQL

While DuckDB supports PostgreSQL-compatible SQL, there are some differences:

1. **Data Types**: DuckDB has its own set of data types, though many are compatible with PostgreSQL
2. **Functions**: Some PostgreSQL functions may not be available or may behave differently
3. **Extensions**: DuckDB uses a different extension system than PostgreSQL

## Benefits of DuckDB

1. **Fast analytical queries**: Optimized for OLAP workloads
2. **Embedded**: No separate server process needed
3. **Portable**: Single file database
4. **PostgreSQL-compatible**: Familiar SQL syntax

## Requirements

- Go 1.24.0 or later
- `github.com/marcboeker/go-duckdb` driver

## Notes

- The DuckDB analyzer uses an in-memory instance to extract query metadata
- Schema migrations are applied to the analyzer instance automatically
- Type inference is done by preparing queries against the DuckDB instance
