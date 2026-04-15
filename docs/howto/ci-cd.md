# Using sqlc in CI/CD

If your project has more than a single developer, we suggest running `sqlc` as
part of your CI/CD pipeline. The four subcommands you'll want to run are `diff`,
`vet`, `verify` and `push`

`sqlc diff` ensures that your generated code is up to date. New developers to a
project may forget to run `sqlc generate` after adding a query or updating a
schema. They also might edit generated code. `sqlc diff` will catch both errors
by comparing the expected output from `sqlc generate` to what's on disk.

```diff
% sqlc diff
--- a/postgresql/query.sql.go
+++ b/postgresql/query.sql.go
@@ -55,7 +55,7 @@

 const listAuthors = `-- name: ListAuthors :many
 SELECT id, name, bio FROM authors
-ORDER BY name
+ORDER BY bio
 `
```

`sqlc vet` runs a set of lint rules against your SQL queries. These rules are
helpful in catching anti-patterns before they make it into production. Please
see the [vet](vet.md) documentation for a complete guide to adding lint rules
for your project.

`sqlc verify` ensures that schema changes do not break production. Existing
queries are checked against new schema changes for correctness. Please see the
[verify](verify.md) documentation for a complete guide.


`sqlc push` pushes your database schema, queries and configuration to sqlc
Cloud. These archives are used by `verify` to catch breaking changes to your
database schema.  Learn more about uploading projects [here](push.md)

## General setup

Install `sqlc` using the [suggested instructions](../overview/install).

Create three steps in your pipeline for `sqlc diff`, `sqlc vet`, and `sqlc
verify`. Run `sqlc push` after merge on your `main` branch.

## GitHub Actions

We provide the [setup-sqlc](https://github.com/marketplace/actions/setup-sqlc)
GitHub Action to install `sqlc`. The action uses the built-in
[tool-cache](https://github.com/actions/toolkit/blob/main/packages/tool-cache/README.md)
to speed up the installation process.

### diff

The following GitHub Workflow configuration runs `sqlc diff` on every push.

```yaml
name: sqlc
on: [push]
jobs:
  diff:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: sqlc-dev/setup-sqlc@v3
      with:
        sqlc-version: '1.30.0'
    - run: sqlc diff
```

### vet

The following GitHub Workflow configuration runs [sqlc vet](vet.md) on every push.
You can use `sqlc vet` without a database connection, but you'll need one if your
`sqlc` configuration references the built-in `sqlc/db-prepare` lint rule.

```yaml
name: sqlc
on: [push]
jobs:
  vet:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: sqlc-dev/setup-sqlc@v3
      with:
        sqlc-version: '1.30.0'
    # Start a PostgreSQL server
    - uses: sqlc-dev/action-setup-postgres@master
      with:
        postgres-version: "16"
      id: postgres
    - run: sqlc vet
      env:
        POSTGRESQL_SERVER_URI: ${{ steps.postgres.outputs.connection-uri }}?sslmode=disable

```

### push

```{note}
Pushing a project is powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

The following GitHub Workflow configuration runs [sqlc push](push.md) on
every push to `main`. Create an auth token via the
[dashboard](https://dashboard.sqlc.dev).

```yaml
name: sqlc
on: [push]
jobs:
  push:
    runs-on: ubuntu-latest
    if: ${{ github.ref == 'refs/heads/main' }}
    steps:
    - uses: actions/checkout@v3
    - uses: sqlc-dev/setup-sqlc@v3
      with:
        sqlc-version: '1.30.0'
    - run: sqlc push
      env:
        SQLC_AUTH_TOKEN: ${{ secrets.SQLC_AUTH_TOKEN }}
```

### verify

```{note}
Verify database migrations is powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

```yaml
name: sqlc
on: [push]
jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: sqlc-dev/setup-sqlc@v3
      with:
        sqlc-version: '1.30.0'
    - uses: sqlc-dev/action-setup-postgres@master
      with:
        postgres-version: "16"
      id: postgres
    - run: sqlc verify
      env:
        POSTGRESQL_SERVER_URI: ${{ steps.postgres.outputs.connection-uri }}?sslmode=disable
        SQLC_AUTH_TOKEN: ${{ secrets.SQLC_AUTH_TOKEN }}
  push:
    runs-on: ubuntu-latest
    if: ${{ github.ref == 'refs/heads/main' }}
    steps:
    - uses: sqlc-dev/setup-sqlc@v3
      with:
        sqlc-version: '1.30.0'
    - run: sqlc push
      env:
        SQLC_AUTH_TOKEN: ${{ secrets.SQLC_AUTH_TOKEN }}
``````
