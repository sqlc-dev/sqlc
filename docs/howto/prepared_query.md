# Preparing queries

```sql
CREATE TABLE records (
  id SERIAL PRIMARY KEY
);

-- name: GetRecord :one
SELECT * FROM records
WHERE id = $1;
```

sqlc has an option to use prepared queries. These prepared queries also work
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

type DBTX interface {
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.getRecordStmt, err = db.PrepareContext(ctx, getRecord); err != nil {
		return nil, err
	}
	return &q, nil
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db            DBTX
	tx            *sql.Tx
	getRecordStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:            tx,
		tx:            tx,
		getRecordStmt: q.getRecordStmt,
	}
}

const getRecord = `-- name: GetRecord :one
SELECT id FROM records
WHERE id = $1
`

func (q *Queries) GetRecord(ctx context.Context, id int) (Record, error) {
	row := q.queryRow(ctx, q.getRecordStmt, getRecord, id)
	var i Record
	err := row.Scan(&i.ID)
	return i, err
}
```

