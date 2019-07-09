# Inserting records

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

type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) error
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
}

const createAuthor = `-- name: CreateAuthor :exec
INSERT INTO authors (bio) VALUES ($1)
`

func (q *Queries) CreateAuthor(ctx context.Context, bio string) error {
	_, err := q.db.ExecContext(ctx, createAuthor, bio)
	return err
}
```
