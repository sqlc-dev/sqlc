// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const getAll = `-- name: GetAll :many
SELECT first_name, last_name, age FROM users
`

func (q *Queries) GetAll(ctx context.Context) ([]User, error) {
	ctx, done := q.observer(ctx, "GetAll")
	rows, err := q.db.Query(ctx, getAll)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.FirstName, &i.LastName, &i.Age); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
