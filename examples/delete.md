# Deleting records

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL
);

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;
```

```go
package db

import (
	"context"
	"database/sql"
)

type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) error
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int) error {
	_, err := q.db.ExecContext(ctx, deleteAuthor, id)
	return err
}
```
