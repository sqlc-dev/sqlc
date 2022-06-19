// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const joinWhereClause = `-- name: JoinWhereClause :many
SELECT foo.barid
FROM foo
JOIN bar ON bar.id = barid
WHERE owner = $1
`

func (q *Queries) JoinWhereClause(ctx context.Context, owner string) ([]int32, error) {
	ctx, done := q.observer(ctx, "JoinWhereClause")
	rows, err := q.db.Query(ctx, joinWhereClause, owner)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var barid int32
		if err := rows.Scan(&barid); err != nil {
			return nil, done(err)
		}
		items = append(items, barid)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
