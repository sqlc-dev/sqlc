# Counting rows

```sql
CREATE TABLE authors (
  id       SERIAL PRIMARY KEY,
  hometown text   NOT NULL
);

-- name: CountAuthors :one
SELECT count(*) FROM authors;

-- name: CountAuthorsByTown :many
SELECT hometown, count(*) FROM authors
GROUP BY 1
ORDER BY 1;
```

```go
package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

const countAuthors = `-- name: CountAuthors :one
SELECT count(*) FROM authors
`

func (q *Queries) CountAuthors(ctx context.Context) (int, error) {
	row := q.db.QueryRowContext(ctx, countAuthors)
	var i int
	err := row.Scan(&i)
	return i, err
}

const countAuthorsByTown = `-- name: CountAuthorsByTown :many
SELECT hometown, count(*) FROM authors
GROUP BY 1
ORDER BY 1
`

type CountAuthorsByTownRow struct {
	Hometown string
	Count    int
}

func (q *Queries) CountAuthorsByTown(ctx context.Context) ([]CountAuthorsByTownRow, error) {
	rows, err := q.db.QueryContext(ctx, countAuthorsByTown)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CountAuthorsByTownRow{}
	for rows.Next() {
		var i CountAuthorsByTownRow
		if err := rows.Scan(&i.Hometown, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
```
