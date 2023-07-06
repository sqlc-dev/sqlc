# Suggested CI/CD setup

If your project has more than a single developer, we suggest running `sqlc` as
part of your CI/CD pipeline. The two commands you'll want to run are `diff` and `vet`

`sqlc diff` ensures that code is up to date. New developers to a project may
forget to run `sqlc generate`. They also might edit generated code. `diff` will
catch both scenarios.

```diff
% sqlc-dev diff
--- a/postgresql/query.sql.go
+++ b/postgresql/query.sql.go
@@ -55,7 +55,7 @@

 const listAuthors = `-- name: ListAuthors :many
 SELECT id, name, bio FROM authors
-ORDER BY name
+ORDER BY bio
 `
```

`sqlc vet` runs a set of lint checks against your SQL queries. These checks are
helpful in catching anti-patterns before they make it into production. Please
see the [vet](vet.md) documentation for a complete guide on adding checks to your
project.

## General setup

Install `sqlc` using the [suggested instructions](../overview/install).

Create two steps in your pipelines, one for `sqlc diff`and one for `sqlc vet`.

## GitHub Actions

We provide the [setup-sqlc](https://github.com/marketplace/actions/setup-sqlc)
GitHub Action to install `sqlc`. The action uses the built-in
[tool-cache](https://github.com/actions/toolkit/blob/main/packages/tool-cache/README.md)
to speed up the installation process.

The following workflow runs `sqlc diff` on every push.

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
        sqlc-version: '1.19.0'
    - run: sqlc diff
```

We also encourage running [`sqlc vet`](vet.md). To get the most value out of
`vet`, you'll want to set up a running database.

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
        sqlc-version: '1.19.0'
      # Connect and migrate your database here. This is an example which runs
      # commands from a `schema.sql` file.
    - run: psql -h localhost -U postgres -p $PG_PORT -d postgres -f schema.sql
      env:
        PGPASSWORD: postgres
    - run: sqlc vet
```