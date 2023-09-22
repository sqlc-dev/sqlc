# Managed databases

*Added in v1.22.0*

`sqlc` can create and maintain hosted databases for your project. These
databases can be used for linting queries. Right now, only PostgreSQL is
supported, with MySQL on the way.

This feature is under active development. Beyond linting queries, managed
databases can be created per test suite or even per test, providing a real,
isolated PostgreSQL database for a test run, no cleanup required.

Interested in trying out managed databases? Sign up [here](https://docs.google.com/forms/d/e/1FAIpQLSdxoMzJ7rKkBpuez-KyBcPNyckYV-5iMR--FRB7WnhvAmEvKg/viewform) or send us an email
at <mailto:hello@sqlc.dev>.

## Configuring managed databases

To configured a managed database, remove the `uri` key, replacing it with the
`managed` key set to `true`. Set the `project` key to your project ID, obtained
via the sqlc Dashboard.

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

You'll also need to create an access token and make it available via the
`SQLC_AUTH_TOKEN` environment variable.

```shell
export SQLC_AUTH_TOKEN=sqlc_xxxxxxxx
```

## Linting queries

In managed mode, `sqlc vet` will create a database with the provided schema and
use that database when running lint rules. If you don't currently have any
rules, the [built-in sqlc/db-prepare] rule verifies each of your queries against
the database by creating a prepared statement.

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

