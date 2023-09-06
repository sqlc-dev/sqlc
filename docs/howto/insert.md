# Inserting rows

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL
);

-- name: CreateAuthor :exec
INSERT INTO authors (bio) VALUES ($1);
```

```go
package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

const createAuthor = `-- name: CreateAuthor :exec
INSERT INTO authors (bio) VALUES ($1)
`

func (q *Queries) CreateAuthor(ctx context.Context, bio string) error {
	_, err := q.db.ExecContext(ctx, createAuthor, bio)
	return err
}
```

## Returning columns from inserted rows

sqlc has full support for the `RETURNING` statement.

```sql
-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateAuthorAndReturnId :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING id;
```

```go
package db

import (
	"context"
	"database/sql"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING id, name, bio
`

type CreateAuthorParams struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRowContext(ctx, createAuthor, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const createAuthorAndReturnId = `-- name: CreateAuthorAndReturnId :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING id
`

type CreateAuthorAndReturnIdParams struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) CreateAuthorAndReturnId(ctx context.Context, arg CreateAuthorAndReturnIdParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createAuthorAndReturnId, arg.Name, arg.Bio)
	var id int64
	err := row.Scan(&id)
	return id, err
}
```

## Using CopyFrom

### PostgreSQL

PostgreSQL supports the [COPY protocol](https://www.postgresql.org/docs/current/sql-copy.html) that can insert rows a lot faster than sequential inserts. You can use this easily with sqlc:

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  name       text   NOT NULL,
  bio        text   NOT NULL
);

-- name: CreateAuthors :copyfrom
INSERT INTO authors (name, bio) VALUES ($1, $2);
```

```go
type CreateAuthorsParams struct {
	Name string
	Bio  string
}

func (q *Queries) CreateAuthors(ctx context.Context, arg []CreateAuthorsParams) (int64, error) {
	...
}
```

The `:copyfrom` command requires either `pgx/v4` or `pgx/v5`.

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "query.sql"
    gen:
      go:
        package: "db"
        sql_package: "pgx/v5"
        out: "db"
```

### MySQL

MySQL supports a similar feature using [LOAD DATA](https://dev.mysql.com/doc/refman/8.0/en/load-data.html).

Errors and duplicate keys are treated as warnings and insertion will
continue, even without an error for some cases. Use this in a transaction
and use SHOW WARNINGS to check for any problems and roll back if necessary.

Check the [error handling](https://dev.mysql.com/doc/refman/8.0/en/load-data.html#load-data-error-handling) documentation for more information.

```sql
CREATE TABLE foo (a text, b integer, c DATETIME, d DATE);

-- name: InsertValues :copyfrom
INSERT INTO foo (a, b, c, d) VALUES (?, ?, ?, ?);
```

```go
func (q *Queries) InsertValues(ctx context.Context, arg []InsertValuesParams) (int64, error) {
	...
}
```

The `:copyfrom` command requires setting the `sql_package` and `sql_driver` options.

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "query.sql"
    schema: "query.sql"
    gen:
      go:
        package: "db"
        sql_package: "database/sql"
        sql_driver: "github.com/go-sql-driver/mysql"
        out: "db"
```