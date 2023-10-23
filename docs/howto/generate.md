# `generate` - Generating code

`sqlc generate` parses SQL, analyzes the results, and outputs code. Your schema and queries are stored in separate SQL files. The paths to these files live in a `sqlc.yaml` configuration file.

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
        sql_package: "pgx/v5"
```

We've written extensive docs on [retrieving](select.md), [inserting](insert.md),
[updating](update.md), and [deleting](delete.md) rows. 

By default, the analysis is run using our built-in query engine. While fast, this engine can't handle some complex queries and type-inference.

## Using a managed database

```{note}
Managed databases are powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

By opting in to [managed database](managed-databases.md), the default analysis is enhanced with metadata from a running database connection. Type inference is improved and query analysis succeeds on a larger set of queries.

```yaml
version: "2"
cloud:
  project: "<PROJECT_ID>"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    database:
      managed: true
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
        sql_package: "pgx/v5"
```

The database analyzer currently supports PostgreSQL, with [MySQL](https://github.com/sqlc-dev/sqlc/issues/2902) and [SQLite](https://github.com/sqlc-dev/sqlc/issues/2903)
support planned in the future.

## Using a database connection

The analyzer uses the configured [database](../reference/config.md#database), whether it be managed or a connection URI.

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    database:
      uri: "postgres://postgres:${PG_PASSWORD}@localhost:5432/postgres"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
        sql_package: "pgx/v5"
```
