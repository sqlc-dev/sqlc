# Returning a single row 

To generate a database access method, annotate a query with a specific comment.

```sql
CREATE TABLE records (
  id SERIAL PRIMARY KEY
);

-- name: GetRecord :one
SELECT * FROM records
WHERE id = $1;
```

A few new pieces of code are generated beyond the `Record` struct. An interface
for the underlying database is generated. The `*sql.DB` and `*sql.Tx` types
satisty this interface.

The database access methods are added to a `Queries` struct, which is created
using the `New` method.

Note that the `*` in our query has been replaced with explicit column names.
This change ensures that the query will never return unexpected data.

Our query was annotated with `:one`, meaning that it should only return a
single row. We scan the data from that one into a `Record` struct.

Since the query has a single parameter, the `GetRecord` method takes a single
`int` as an argument.

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
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
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
