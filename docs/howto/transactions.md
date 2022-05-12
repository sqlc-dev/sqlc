# Using transactions

```sql
CREATE TABLE records (
  id SERIAL PRIMARY KEY,
  counter INT NOT NULL
);

-- name: GetRecord :one
SELECT * FROM records
WHERE id = $1;
```

The `WithTx` method allows a `Queries` instance to be associated with a transaction.

```go
package db

import (
	"context"
	"database/sql"
)

type Record struct {
	ID      int
	Counter int
}

type DBTX interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (*Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{db: tx}
}

const getRecord = `-- name: GetRecord :one
SELECT id, counter FROM records
WHERE id = $1
`

func (q *Queries) GetRecord(ctx context.Context, id int) (Record, error) {
	row := q.db.QueryRowContext(ctx, getRecord, id)
	var i Record
	err := row.Scan(&i.ID, &i.Counter)
	return i, err
}
```

With pgx you'd use it like this for example:

```go
func bumpCounter(ctx context.Context, p *pgx.Conn, id int) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	q := db.New(tx)
	r, err := q.GetRecord(ctx, id)
	if err != nil {
		return err
	}
	if err := q.UpdateRecord(ctx, db.UpdateRecordParams{
		ID:      r.ID,
		Counter: r.Counter + 1,
	}); err != nil {
		return err
	}
	return tx.Commit()
}
```
