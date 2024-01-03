# `verify` - Verifying schema changes

```{note}
`verify` is powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

*Added in v1.24.0*

Schema updates and poorly-written queries often bring down production databases. That’s bad.

Out of the box, `sqlc generate` catches some of these issues. Running `sqlc vet` with the `sqlc/db-prepare` rule catches more subtle problems. But there is a large class of issues that sqlc can’t prevent by looking at current schema and queries alone.

For instance, when a schema change is proposed, existing queries and code running in production might fail when the schema change is applied. Enter `sqlc verify`, which analyzes existing queries against new schema changes and errors if there are any issues.

Let's look at an example. Assume you have these two tables in production.

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY
);

CREATE TABLE user_actions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  action TEXT,
  created_at TIMESTAMP
);
```

Your application contains the following query to join user actions against the users table.

```sql
-- name: GetUserActions :many
SELECT * FROM users u
JOIN user_actions ua ON u.id = ua.user_id
ORDER BY created_at;
```

So far, so good. Then assume you propose this schema change:

```sql
ALTER TABLE users ADD COLUMN created_at TIMESTAMP;
```

Running `sqlc generate` fails with this change, returning a `column reference "created_at" is ambiguous` error. You update your query to fix the issue.

```sql
-- name: GetUserActions :many
SELECT * FROM users u
JOIN user_actions ua ON u.id = ua.user_id
ORDER BY u.created_at;
```

While that change fixes the issue, there's a production outage waiting to happen. When the schema change is applied, the existing `GetUserActions` query will begin to fail. The correct way to fix this is to deploy the updated query before applying the schema migration.

It ensures migrations are safe to deploy by sending your current schema and queries to sqlc cloud. There, we run the queries for your latest push against your new schema changes. This check catches backwards incompatible schema changes for existing queries.

Here `sqlc verify` alerts you to the fact that ORDER BY "created_at" is ambiguous.

```sh
$ sqlc verify
FAIL: app query.sql

=== Failed
=== FAIL: app query.sql GetUserActions
    ERROR: column reference "created_at" is ambiguous (SQLSTATE 42702)
```

By the way, this scenario isn't made up! It happened to us a few weeks ago. We've been happily testing early versions of `verify` for the last two weeks and haven't had any issues since.

This type of verification is only the start. If your application is deployed on-prem by your customers, `verify` could tell you if it's safe for your customers to rollback to an older version of your app, even after schema migrations have been run.

Using `verify` requires that you push your queries and schema when you tag a release of your application. We run it on every push to main, as we continuously deploy those commits.

## Authentication

`sqlc` expects to find a valid auth token in the value of the `SQLC_AUTH_TOKEN`
environment variable. You can create an auth token via the [dashboard](https://dashboard.sqlc.dev).

```shell
export SQLC_AUTH_TOKEN=sqlc_xxxxxxxx
```

## Expected workflow

Using `sqlc verify` requires pushing your queries and schema to sqlc Cloud. When
you release a new version of your application, you should push your schema and
queries as well. For example, we run `sqlc push` after any change has been
merged into our `main` branch on Github, as we deploy every commit to
production.

```shell
$ sqlc push --tag main
```

Locally or in pull requests, run `sqlc verify` to check that existing queries
continue to work with your current database schema.

```shell
$ sqlc verify --against main
```

## Picking a tag

Without an `against` argument, `verify` will run its analysis of the provided schema using your most-recently pushed queries. We suggest using the `against` argument to explicitly select a set of queries for comparison.

```shell
$ sqlc verify --against [tag]
```
