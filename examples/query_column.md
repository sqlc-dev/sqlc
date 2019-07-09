# Returning specific columns

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);

-- name: GetBioForAuthor :one
SELECT bio FROM authors
WHERE id = $1;

-- name: GetInfoForAuthor :one
SELECT bio, birth_year FROM authors
WHERE id = $1;
```


When selecting a single column, only that value that returned. The `GetBioForAuthor`
method takes a single `int` as an argument and returns a `string` and an
`error`.

When selecting multiple columns, a row record (method-specific struct) is
returned. In this case, `GetInfoForAuthor` returns a struct with two fields:
`Bio` and `BirthYear`.

```go
package db

import (
	"context"
	"database/sql"
)

type dbtx interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
}

const getBioForAuthor = `-- name: GetBioForAuthor :one
SELECT bio FROM authors
WHERE id = $1
`

func (q *Queries) GetBioForAuthor(ctx context.Context, id int) (string, error) {
	row := q.db.QueryRowContext(ctx, getBioForAuthor, id)
	var i string
	err := row.Scan(&i)
	return i, err
}

const getInfoForAuthor = `-- name: GetInfoForAuthor :one
SELECT bio, birth_year FROM authors
WHERE id = $1
`

type GetInfoForAuthorRow struct {
	Bio       string
	BirthYear int
}

func (q *Queries) GetBioForAuthor(ctx context.Context, id int) (GetBioForAuthor, error) {
	row := q.db.QueryRowContext(ctx, getInfoForAuthor, id)
	var i GetBioForAuthor
	err := row.Scan(&i.Bio, &i.BirthYear)
	return i, err
}
```
