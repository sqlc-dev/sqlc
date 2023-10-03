// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const getAuthor = `-- name: GetAuthor :one
SELECT a.id, name, b.id, title FROM authors AS a
FULL OUTER JOIN books AS b
 ON a.id = b.id
WHERE a.id = ? LIMIT 1
`

type GetAuthorRow struct {
	ID    sql.NullInt64
	Name  sql.NullString
	ID_2  sql.NullInt64
	Title sql.NullString
}

func (q *Queries) GetAuthor(ctx context.Context, id int64, aq ...AdditionalQuery) (GetAuthorRow, error) {
	query := getAuthor
	queryParams := []interface{}{id}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var i GetAuthorRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ID_2,
		&i.Title,
	)
	return i, err
}
