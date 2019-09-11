# The ANY operator, A.K.A. Passing a slice argument to a query
Sometimes, an application will need a function that filters one or more columns, using one or more arguments for each column, for a specific table.
Although the previous statement is a bit wordy, the actual query is can be pretty simple, with the use of the ANY operator.

This SQL script:
```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);

-- name: ListAuthorsByIDs :many
SELECT * FROM authors
WHERE id = ANY($1::int[]);
```

Will turn into this Go code:
```go
package db

import (
	"context"
	"database/sql"
)

type Author struct {
	ID        int
	Bio       string
	BirthYear int
}

type dbtx interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
}

const listAuthors = `-- name: ListAuthorsByIDs :many
SELECT id, bio, birth_year FROM authors
WHERE id = ANY($1::int[])
`

func (q *Queries) ListAuthorsByIDs(ctx context.Context, ids []int) ([]Author, error) {
  rows, err := q.db.QueryContext(ctx, listAuthors, pq.Array(ids))
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var items []Author
  for rows.Next() {
    var i Author
    if err := rows.Scan(&i.ID, &i.Bio, &i.BirthYear); err != nil {
      return nil, err
    }
    items = append(items, i)
  }
  if err := rows.Close(); err != nil {
    return nil, err
  }
  if err := rows.Err(); err != nil {
    return nil, err
  }
  return items, nil
}
```
