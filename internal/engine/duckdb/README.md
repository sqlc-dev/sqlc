# DuckDB Engine Implementation

This directory contains the DuckDB engine implementation for sqlc.

## Architecture

The DuckDB engine reuses sqlc's PostgreSQL parser and catalog while providing a custom analyzer that connects to an in-memory DuckDB instance. This design leverages DuckDB's PostgreSQL-compatible SQL syntax while enabling accurate type inference through live database analysis.

### Components

1. **Parser**: Reuses `postgresql.NewParser()`
   - DuckDB's SQL syntax is PostgreSQL-compatible
   - No need for a separate parser implementation

2. **Catalog**: Reuses `postgresql.NewCatalog()`
   - Schema metadata is managed using PostgreSQL's catalog structure
   - Compatible with DuckDB's type system

3. **Analyzer**: Custom implementation (`analyzer/analyze.go`)
   - Connects to an in-memory DuckDB instance
   - Extracts column and parameter types by preparing queries
   - Supports schema migrations

## Implementation Details

### Analyzer

The DuckDB analyzer (`analyzer/analyze.go`) implements the `analyzer.Analyzer` interface:

```go
type Analyzer interface {
    Analyze(ctx context.Context, n ast.Node, query string,
            migrations []string, ps *named.ParamSet) (*analysis.Analysis, error)
    Close(ctx context.Context) error
}
```

#### Key Features

- **Lazy Connection**: Database connection is established on first use
- **In-Memory Default**: Uses `:memory:` if no URI is provided
- **Schema Migrations**: Applies migrations before analyzing queries
- **Thread-Safe**: Uses mutex to protect connection initialization
- **Type Inference**: Uses `database/sql.ColumnTypes()` to extract column metadata

#### Database URI Formats

- `:memory:` - In-memory database (default)
- `file.db` - File-based database in current directory
- `/path/to/file.db` - Absolute path to database file

#### Type Extraction

The analyzer extracts type information by:

1. Preparing the query using `sql.PrepareContext()`
2. Querying to get column metadata
3. Using `rows.ColumnTypes()` to extract:
   - Column name
   - Data type
   - Nullability
   - Array dimensions (if applicable)

### Integration Points

#### Engine Registration (`internal/compiler/engine.go`)

```go
case config.EngineDuckDB:
    // Reuse PostgreSQL parser and catalog
    c.parser = postgresql.NewParser()
    c.catalog = postgresql.NewCatalog()
    c.selector = newDefaultSelector()

    // Use DuckDB-specific analyzer
    if conf.Database != nil {
        if conf.Analyzer.Database == nil || *conf.Analyzer.Database {
            c.analyzer = analyzer.Cached(
                duckdbanalyze.New(c.client, *conf.Database),
                combo.Global,
                *conf.Database,
            )
        }
    }
```

#### Vet Support (`internal/cmd/vet.go`)

```go
case config.EngineDuckDB:
    db, err := sql.Open("duckdb", dburl)
    if err != nil {
        return fmt.Errorf("database: connection error: %s", err)
    }
    if err := db.PingContext(ctx); err != nil {
        return fmt.Errorf("database: connection error: %s", err)
    }
    defer db.Close()
    prep = &dbPreparer{db}
    expl = nil // DuckDB supports EXPLAIN, but not enabled yet
```

## Design Decisions

### Why Reuse PostgreSQL Components?

1. **SQL Compatibility**: DuckDB intentionally implements PostgreSQL-compatible SQL
2. **Code Reuse**: Avoid duplicating parser and catalog logic
3. **Maintainability**: Changes to PostgreSQL support automatically benefit DuckDB
4. **Correctness**: Leverage well-tested PostgreSQL parser

### Why Custom Analyzer?

1. **Different Driver**: DuckDB uses a different Go driver (`go-duckdb` vs `pgx`)
2. **Type System**: DuckDB's type system has subtle differences from PostgreSQL
3. **Introspection**: DuckDB's metadata APIs differ from PostgreSQL
4. **In-Memory Focus**: Optimized for in-memory and embedded use cases

## Limitations

1. **Type Inference**: Falls back to `text` type if column types cannot be determined
2. **Parameter Types**: Database/sql doesn't provide standard parameter type introspection
3. **EXPLAIN Support**: Not yet implemented in vet command
4. **Extension System**: DuckDB extensions are not yet supported

## Future Enhancements

1. **Better Type Inference**: Use DuckDB-specific metadata queries
2. **Parameter Type Detection**: Implement DuckDB-specific parameter introspection
3. **EXPLAIN Support**: Add explainer for vet command
4. **Extension Loading**: Support DuckDB extensions
5. **Managed Databases**: Integration with dbmanager for managed DuckDB instances
6. **Performance Optimizations**: Cache prepared statements and metadata

## Testing

### Unit Tests

```bash
go test ./internal/engine/duckdb/analyzer/
```

### Integration Tests

```bash
# With DuckDB driver installed
cd examples/duckdb
sqlc generate
```

### End-to-End Tests

```bash
# Add DuckDB examples to endtoend tests
go test --tags=examples ./internal/endtoend/
```

## Dependencies

- `github.com/marcboeker/go-duckdb` v1.8.5 - DuckDB Go driver

## References

- [DuckDB Documentation](https://duckdb.org/docs/)
- [DuckDB SQL Syntax](https://duckdb.org/docs/sql/introduction)
- [go-duckdb Driver](https://github.com/marcboeker/go-duckdb)
- [sqlc Documentation](https://docs.sqlc.dev/)

## Contributing

When contributing to the DuckDB engine:

1. Maintain PostgreSQL compatibility where possible
2. Document any DuckDB-specific behavior
3. Add tests for new functionality
4. Update examples to demonstrate new features

## License

Same as the parent sqlc project.
