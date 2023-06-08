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

PostgreSQL supports the Copy Protocol that can insert rows a lot faster than sequential inserts. You can use this easily with sqlc:

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
	return q.db.CopyFrom(ctx, []string{"authors"}, []string{"name", "bio"}, &iteratorForCreateAuthors{rows: arg})
}
```
