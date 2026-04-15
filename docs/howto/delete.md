# Deleting rows

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL
);
```

The parameter syntax varies by database engine:

**PostgreSQL:**
```sql
-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;
```

**MySQL and SQLite:**
```sql
-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = ?;
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

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int) error {
	_, err := q.db.ExecContext(ctx, deleteAuthor, id)
	return err
}
```
