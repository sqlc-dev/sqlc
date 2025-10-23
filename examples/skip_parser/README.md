# Skip Parser Example

This example demonstrates the `analyzer.skip_parser` configuration option for PostgreSQL.

## What is skip_parser?

When `analyzer.skip_parser: true` is set in the configuration, sqlc will:

1. **Skip the parser** - No parsing of SQL schema or query files
2. **Skip the catalog** - No building of the internal schema catalog
3. **Use database analyzer only** - Rely entirely on the PostgreSQL database for query analysis

This is useful when:
- Working with complex PostgreSQL syntax not fully supported by the parser
- You want to ensure queries are validated against the actual database schema
- Using database-specific features or extensions

## Configuration

See `sqlc.yaml` for the configuration:

```yaml
version: '2'
sql:
- name: postgresql
  schema: postgresql/schema.sql
  queries: postgresql/query.sql
  engine: postgresql
  database:
    uri: "${SKIP_PARSER_TEST_POSTGRES}"
  analyzer:
    skip_parser: true  # This enables the feature
  gen:
    go:
      package: skipparser
      sql_package: pgx/v5
      out: postgresql
```

## How It Works

1. The schema file (`schema.sql`) is **NOT** parsed by sqlc
2. The schema must be applied to the database separately (tests do this automatically)
3. Query files (`query.sql`) are split using `sqlfile.Split`
4. Each query is sent to PostgreSQL's analyzer for validation
5. Column and parameter types are retrieved from the database
6. Code is generated based on the database analysis

## Running the Tests

### Prerequisites

- PostgreSQL server running and accessible
- Set the `SKIP_PARSER_TEST_POSTGRES` environment variable

### Option 1: Using Docker Compose

```bash
# Start PostgreSQL
docker compose up -d

# Set environment variable
export SKIP_PARSER_TEST_POSTGRES="postgresql://postgres:mysecretpassword@localhost:5432/postgres"

# Run the test
go test -tags=examples ./examples/skip_parser/postgresql
```

### Option 2: Using existing PostgreSQL

```bash
# Set environment variable to your PostgreSQL instance
export SKIP_PARSER_TEST_POSTGRES="postgresql://user:pass@localhost:5432/dbname"

# Run the test
go test -tags=examples ./examples/skip_parser/postgresql
```

### Generating Code

```bash
# Make sure database is running and accessible
export SKIP_PARSER_TEST_POSTGRES="postgresql://postgres:mysecretpassword@localhost:5432/postgres"

# Generate code
cd examples/skip_parser
sqlc generate
```

## Tests Included

The `db_test.go` file includes comprehensive tests:

### TestSkipParser
- Creates a product with arrays and JSON fields
- Tests all CRUD operations (Create, Read, Update, Delete)
- Tests list and search operations
- Tests counting

### TestSkipParserComplexTypes
- Tests PostgreSQL-specific types (arrays, JSONB)
- Tests handling of nil/empty values
- Validates array and JSON handling

## Features Tested

This example tests the following PostgreSQL features with skip_parser:

- **BIGSERIAL** - Auto-incrementing primary keys
- **NUMERIC** - Decimal types
- **TIMESTAMPTZ** - Timestamps with timezone
- **TEXT[]** - Text arrays
- **JSONB** - Binary JSON storage
- **ANY operator** - Array containment queries
- **GIN indexes** - Generalized inverted indexes for arrays
- **RETURNING clause** - Return values from INSERT/UPDATE

All of these are validated directly by PostgreSQL without parser involvement!
