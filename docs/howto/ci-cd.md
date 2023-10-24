# Using sqlc in CI/CD

If your project has more than a single developer, we suggest running `sqlc` as
part of your CI/CD pipeline. The three subcommands you'll want to run are `diff`, `vet` and `upload`

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

`sqlc upload` pushes your database schema and queries to sqlc Cloud. Once
uploaded, we ensure that future releases of sqlc do not break your code. Learn
more about uploading projects [here](upload.md)

## General setup

Install `sqlc` using the [suggested instructions](../overview/install).

Create two steps in your pipeline, one for `sqlc diff` and one for `sqlc vet`. Run `sqlc upload` after merge on your `main` branch.

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
        sqlc-version: '1.23.0'
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
    services:
      postgres:
        image: "postgres:15"
        env:
          POSTGRES_DB: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_USER: postgres
        ports:
        - 5432:5432
        # needed because the postgres container does not provide a healthcheck
        options: --health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5
    env:
      PG_PORT: ${{ job.services.postgres.ports['5432'] }}

    steps:
    - uses: actions/checkout@v3
    - uses: sqlc-dev/setup-sqlc@v3
      with:
        sqlc-version: '1.23.0'
      # Connect and migrate your database here. This is an example which runs
      # commands from a `schema.sql` file.
    - run: psql -h localhost -U postgres -p $PG_PORT -d postgres -f schema.sql
      env:
        PGPASSWORD: postgres
    - run: sqlc vet
```

#### Managed databases

```{note}
Managed databases are powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

If you're using [managed databases](managed-databases.md), the `services` block
in the previous workflow isn't required.

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
        sqlc-version: '1.23.0'
    - run: sqlc vet
```

### upload

```{note}
Project uploads are powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

The following GitHub Workflow configuration runs [sqlc upload](upload.md) on
every push to `main`. Create an auth token via the
[dashboard](https://dashboard.sqlc.dev).

```yaml
name: sqlc
on: [push]
jobs:
  upload:
    runs-on: ubuntu-latest
    if: ${{ github.ref == 'refs/heads/main' }}
    steps:
    - uses: actions/checkout@v3
    - uses: sqlc-dev/setup-sqlc@v3
      with:
        sqlc-version: '1.23.0'
    - run: sqlc upload
      env:
        SQLC_AUTH_TOKEN: ${{ secrets.SQLC_AUTH_TOKEN }}
```
