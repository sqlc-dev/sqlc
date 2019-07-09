# Updating records

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL
);

-- name: UpdateAuthor :exec
UPDATE authors SET bio = $2
WHERE id = $1;
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

const updateAuthor = `-- name: UpdateAuthor :exec
UPDATE authors SET bio = $2
WHERE id = $1
`

func (q *Queries) UpdateAuthor(ctx context.Context, id int, bio string) error {
	_, err := q.db.ExecContext(ctx, updateAuthor, id, bio)
	return err
}
```
