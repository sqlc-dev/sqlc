// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package query

import (
	"context"

	"github.com/lib/pq"
)

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio, tags FROM authors
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, id)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		pq.Array(&i.Tags),
	)
	return i, err
}
