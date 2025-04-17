// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const findByID = `-- name: FindByID :many
SELECT id, name FROM users WHERE ? = id
`

func (q *Queries) FindByID(ctx context.Context, id int32) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, findByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const findByIDAndName = `-- name: FindByIDAndName :many
SELECT id, name FROM users WHERE ? = id AND ? = name
`

type FindByIDAndNameParams struct {
	ID   int32
	Name sql.NullString
}

func (q *Queries) FindByIDAndName(ctx context.Context, arg FindByIDAndNameParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, findByIDAndName, arg.ID, arg.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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
