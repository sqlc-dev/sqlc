// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
)

const listCampuses = `-- name: ListCampuses :many
SELECT id FROM campus
`

func (q *Queries) ListCampuses(ctx context.Context) ([]string, error) {
	ctx, done := q.observer(ctx, "ListCampuses")
	rows, err := q.db.Query(ctx, listCampuses)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, done(err)
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const listMetadata = `-- name: ListMetadata :many
SELECT id FROM product_meta
`

func (q *Queries) ListMetadata(ctx context.Context) ([]string, error) {
	ctx, done := q.observer(ctx, "ListMetadata")
	rows, err := q.db.Query(ctx, listMetadata)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, done(err)
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const listStudents = `-- name: ListStudents :many
SELECT id FROM students
`

func (q *Queries) ListStudents(ctx context.Context) ([]string, error) {
	ctx, done := q.observer(ctx, "ListStudents")
	rows, err := q.db.Query(ctx, listStudents)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, done(err)
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
