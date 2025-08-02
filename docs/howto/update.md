# Updating rows

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL
);
```

## Single parameter

If your query has a single parameter, your Go method will also have a single
parameter.

The parameter syntax varies by database engine:

**PostgreSQL:**
```sql
-- name: UpdateAuthorBios :exec
UPDATE authors SET bio = $1;
```

**MySQL and SQLite:**
```sql
-- name: UpdateAuthorBios :exec
UPDATE authors SET bio = ?;
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

const updateAuthorBios = `-- name: UpdateAuthorBios :exec
UPDATE authors SET bio = $1
`

func (q *Queries) UpdateAuthorBios(ctx context.Context, bio string) error {
	_, err := q.db.ExecContext(ctx, updateAuthorBios, bio)
	return err
}
```

## Multiple parameters

If your query has more than one parameter, your Go method will accept a
`Params` struct.

**PostgreSQL:**
```sql
-- name: UpdateAuthor :exec
UPDATE authors SET bio = $2
WHERE id = $1;
```

**MySQL and SQLite:**
```sql
-- name: UpdateAuthor :exec
UPDATE authors SET bio = ?
WHERE id = ?;
```

Note: For MySQL and SQLite, parameters are bound in the order they appear in the query, regardless of the order in the function signature.

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

const updateAuthor = `-- name: UpdateAuthor :exec
UPDATE authors SET bio = $2
WHERE id = $1
`

type UpdateAuthorParams struct {
	ID  int32
	Bio string
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) error {
	_, err := q.db.ExecContext(ctx, updateAuthor, arg.ID, arg.Bio)
	return err
}
```

