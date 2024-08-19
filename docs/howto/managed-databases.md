# Managed databases

*Added in v1.22.0*

`sqlc` can automatically create read-only databases to power query analysis,
linting and verification. These databases are immediately useful for powering
sqlc's database-connected query analyzer, an opt-in feature that improves upon
sqlc's built-in query analysis engine. PostgreSQL support is available today,
with MySQL on the way.

Once configured, `sqlc` will also use managed databases when linting queries
with [`sqlc vet`](vet.md) in cases where your lint rules require a connection
to a running database.

Managed databases are under active development, and we're interested in
supporting other use-cases.

## Configuring managed databases

To configure `sqlc` to use managed databases, remove the `uri` key from your
`database` configuration and replace it with the `managed` key set to `true`.
Access to a running database server is required. Add a connection string to the `servers` mapping.

```yaml
version: '2'
servers:
- engine: postgresql
  uri: "postgres://locahost:5432/postgres?sslmode=disable"
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  database:
    managed: true
```

An environment variable can also be used via the `${}` syntax.

```yaml
version: '2'
servers:
- engine: postgresql
  uri: ${DATABASE_URI}
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  database:
    managed: true
```

## Improving codegen

Without a database connection, sqlc does its best to parse, analyze and compile your queries just using
the schema you pass it and what it knows about the various database engines it supports. In many cases
this works just fine, but for more advanced queries sqlc might not have enough information to produce good code.

With managed databases configured, `sqlc generate` will automatically create a hosted ephemeral database with your
schema and use that database to improve its query analysis. And sqlc will cache its analysis locally
on a per-query basis to speed up future codegen runs. Here's a minimal working configuration:

```yaml
version: '2'
servers:
- engine: postgresql
  uri: "postgres://locahost:5432/postgres?sslmode=disable"
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  database:
    managed: true
  gen:
    go:
      out: "db"
```

## Linting queries

With managed databases configured, `sqlc vet` will automatically create a hosted ephemeral database with your
schema and use that database when running lint rules that require a
database connection, e.g. any [rule relying on `EXPLAIN ...` output](vet.md#rules-using-explain-output).

If you don't yet have any vet rules, the [built-in sqlc/db-prepare rule](vet.md#sqlc-db-prepare)
is a good place to start. It prepares each of your queries against the database
to ensure the query is valid. Here's a minimal working configuration:

```yaml
version: '2'
servers:
- engine: postgresql
  uri: "postgres://locahost:5432/postgres?sslmode=disable"
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  database:
    managed: true
  rules:
  - sqlc/db-prepare
```
