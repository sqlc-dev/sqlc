# Managed databases

```{note}
Managed databases are powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

*Added in v1.22.0*

`sqlc` can create and maintain short-lived hosted databases for your project.
These ephemeral databases are immediately useful for powering sqlc's
database-connected query analyzer, an opt-in feature that improves upon sqlc's
built-in query analysis engine. PostgreSQL support is available today, with
MySQL on the way.

Once configured, `sqlc` will also use managed databases when linting queries
with [`sqlc vet`](vet.md) in cases where your lint rules require a connection
to a running database.

Managed databases are under active development, and we're interested in
supporting other use-cases. Outside of sqlc itself, you can use our managed
databases in your tests to quickly stand up a database per test suite or even per test,
providing a real, isolated database for a test run. No cleanup required.

## Configuring managed databases

To configure `sqlc` to use managed databases, remove the `uri` key from your
`database` configuration and replace it with the `managed` key set to `true`.
Set the `project` key in your `cloud` configuration to the value of your
project ID, obtained via the [dashboard](https://dashboard.sqlc.dev).

```yaml
version: '2'
cloud:
  project: '<PROJECT_ID>'
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  database:
    managed: true
```

### Authentication

`sqlc` expects to find a valid auth token in the value of the `SQLC_AUTH_TOKEN`
environment variable. You can create an auth token via the [dashboard](https://dashboard.sqlc.dev).

```shell
export SQLC_AUTH_TOKEN=sqlc_xxxxxxxx
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
cloud:
  project: '<PROJECT_ID>'
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
cloud:
  project: '<PROJECT_ID>'
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  database:
    managed: true
  rules:
  - sqlc/db-prepare
```

## With other tools

With managed databases configured, `sqlc createdb` will create a hosted ephemeral database with your
schema and write the database's connection URI as a string to standard output (stdout). This allows you to use
ephemeral databases with other tools that understand database connection strings.

In the simplest case, you can use psql to poke around:

```shell
psql $(sqlc createdb)
```

Or if you're tired of waiting for us to resolve https://github.com/sqlc-dev/sqlc/issues/296,
you can create databases ad hoc to use with pgtyped:

```shell
DATABASE_URL=$(sqlc createdb) npx pgtyped -c config.json
```

Here's a minimal working configuration if all you need to use is `sqlc createdb`:

```yaml
version: '2'
cloud:
  project: '<PROJECT_ID>'
sql:
- schema: schema.sql
  engine: postgresql
  database:
    managed: true
```
