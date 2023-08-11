# Preparing queries

If you're using `pgx/v5` you get its
[implicit support](https://github.com/jackc/pgx/wiki/Automatic-Prepared-Statement-Caching)
for prepared statements. No additional `sqlc` configuration is required.

For other drivers, `sqlc` can give you the option to explicitly use prepared
queries. These prepared queries also work with transactions.

You'll need to set `emit_prepared_queries` to `true` in your `sqlc` configuration
to generate code similar to the example below.

```sql
CREATE TABLE records (
  id SERIAL PRIMARY KEY
);

-- name: GetRecord :one
SELECT * FROM records
WHERE id = $1;
```

```go
package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Record struct {
	ID int32
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
		return nil, fmt.Errorf("error preparing query GetRecord: %w", err)
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

func (q *Queries) GetRecord(ctx context.Context, id int32) (int32, error) {
	row := q.queryRow(ctx, q.getRecordStmt, getRecord, id)
	err := row.Scan(&id)
	return id, err
}
```
