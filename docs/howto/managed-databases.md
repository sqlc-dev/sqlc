# Managed databases

*Added in v1.22.0*

`sqlc` can create and maintain hosted databases for your project. These
databases are immediately useful for linting queries with [`sqlc vet`](vet.md)
if your lint rules require a connection to a running database. PostgreSQL
support is available today, with MySQL on the way.

This feature is under active development, and we're interested in supporting
other use-cases. Beyond linting queries, you can use sqlc managed databases
in your tests to quickly stand up a database per test suite or even per test,
providing a real, isolated database for a test run. No cleanup required.

Interested in trying out managed databases? Sign up [here](https://docs.google.com/forms/d/e/1FAIpQLSdxoMzJ7rKkBpuez-KyBcPNyckYV-5iMR--FRB7WnhvAmEvKg/viewform) or send us an email
at [hello@sqlc.dev](mailto:hello@sqlc.dev).

## Configuring managed databases

To configure `sqlc` to use a managed database, remove the `uri` key from your
`database` configuration and replace it with the `managed` key set to `true`.
Set the `project` key in your `cloud` configuration to the value of your
project ID, obtained via the sqlc.dev Dashboard.

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

## Authentication

`sqlc` expects to find a valid auth token in the value of the `SQLC_AUTH_TOKEN`
environment variable. You can create an auth token via the sqlc.dev Dashboard.

```shell
export SQLC_AUTH_TOKEN=sqlc_xxxxxxxx
```

## Linting queries

With managed databases configured, `sqlc vet` will create a database with your
package's schema and use that database when running lint rules that require a
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
