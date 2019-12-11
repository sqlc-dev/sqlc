# Returning values

sqlc has full support for the `RETURNING` statement.

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL
);

-- name: Delete :exec
DELETE FROM authors WHERE id = $1;

-- name: DeleteAffected :execrows
DELETE FROM authors WHERE id = $1;

-- name: DeleteID :one
DELETE FROM authors WHERE id = $1
RETURNING id;

-- name: DeleteAuthor :one
DELETE FROM authors WHERE id = $1
RETURNING *;
```

```go
package db

import (
	"context"
	"database/sql"
)

type Author struct {
	ID  int
	Bio string
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) error
	QueryRowContext(context.Context, string, ...interface{}) error
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

const delete = `-- name: Delete :exec
DELETE FROM authors WHERE id = $1
`

func (q *Queries) Delete(ctx context.Context, id int) error {
	_, err := q.db.ExecContext(ctx, delete, id)
	return err
}

const deleteAffected = `-- name: DeleteAffected :exec
DELETE FROM authors WHERE id = $1
`

func (q *Queries) DeleteAffected(ctx context.Context, id int) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteAffected, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deleteID = `-- name: DeleteID :one
DELETE FROM authors WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteID(ctx context.Context, id int) (int, error) {
	row := q.db.QueryRowContext(ctx, deleteID, id)
	var i int
	err := row.Scan(&i)
	return i, err
}

const deleteAuhtor = `-- name: DeleteAuthor :one
DELETE FROM authors WHERE id = $1
RETURNING id, bio
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int) (Author, error) {
	row := q.db.QueryRowContext(ctx, deleteAuhtor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Bio)
	return i, err
}
```
