// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
)

const updateAuthor = `-- name: UpdateAuthor :one
update authors
set names[$1] = $2
where id=$3
RETURNING id, names
`

type UpdateAuthorParams struct {
	Names   []string
	Names_2 []string
	ID      int64
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, updateAuthor, arg.Names, arg.Names_2, arg.ID)
	var i Author
	err := row.Scan(&i.ID, &i.Names)
	return i, err
}
