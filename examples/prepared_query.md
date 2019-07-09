# Prepared Queries

```sql
CREATE TABLE records (
  id SERIAL PRIMARY KEY
);

-- name: GetRecord :one
SELECT * FROM records
WHERE id = $1;
```

DinoSQL has an option to use perpared queries. These prepared queries also work
with transactions.

```go
package db

import (
	"context"
	"database/sql"
)

type Record struct {
	ID int
}

type dbtx interface {
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db dbtx) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.getRecord, err = db.PrepareContext(ctx, getRecord); err != nil {
		return nil, err
	}
	return &q, nil
}

type Queries struct {
	db        dbtx
	tx        *sql.Tx
	getRecord *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:        tx,
		tx:        tx,
		getRecord: q.getRecord,
	}
}

const getRecord = `-- name: GetRecord :one
SELECT id FROM records
WHERE id = $1
`

func (q *Queries) GetRecord(ctx context.Context, id int) (Record, error) {
	var row *sql.Row
	switch {
	case q.getRecord != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.getRecord).QueryRowContext(ctx, id)
	case q.getRecord != nil:
		row = q.getRecord.QueryRowContext(ctx, id)
	default:
		row = q.db.QueryRowContext(ctx, getRecord, id)
	}
	var i Record
	err := row.Scan(&i.ID)
	return i, err
}
```

