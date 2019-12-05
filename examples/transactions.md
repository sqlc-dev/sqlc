# Transactions

```sql
CREATE TABLE records (
  id SERIAL PRIMARY KEY
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
	ID int
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
SELECT id FROM records
WHERE id = $1
`

func (q *Queries) GetRecord(ctx context.Context, id int) (Record, error) {
	row := q.db.QueryRowContext(ctx, getRecord, id)
	var i Record
	err := row.Scan(&i.ID)
	return i, err
}
```
