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

By default, sqlc runs its analysis using a built-in query analysis engine. While fast, this engine can't handle some complex queries and type-inference.

You can configure sqlc to use a database connection for enhanced analysis using metadata from that database.

The database-backed analyzer currently supports PostgreSQL, with [MySQL](https://github.com/sqlc-dev/sqlc/issues/2902) and [SQLite](https://github.com/sqlc-dev/sqlc/issues/2903)
support planned in the future.

## Enhanced analysis with managed databases

```{note}
Managed databases are powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

With [managed databases](managed-databases.md) configured, `generate` will automatically create a hosted ephemeral database with your
schema and use that database to improve its query analysis. And sqlc will cache its analysis locally
on a per-query basis to speed up future `generate` runs. This saves you the trouble of running and maintaining a database with
an up-to-date schema. Here's a minimal working configuration:

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
        out: "db"
        sql_package: "pgx/v5"
```

## Enhanced analysis using your own database

You can opt-in to database-backed analysis using your own database, by providing a `uri` in your sqlc
[database](../reference/config.md#database) configuration.

The `uri` string can contain references to environment variables using the `${...}`
syntax. In the following example, the connection string will have the value of
the `PG_PASSWORD` environment variable set as its password.

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
        out: "db"
        sql_package: "pgx/v5"
```

Databases configured with a `uri` must have an up-to-date schema for query analysis to work correctly, and `sqlc` does not apply schema migrations your database. Use your migration tool of choice to create the necessary
tables and objects before running `sqlc generate`.
