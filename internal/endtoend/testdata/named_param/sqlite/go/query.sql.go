// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const atParams = `-- name: AtParams :many
SELECT name FROM foo WHERE name = ?1
`

func (q *Queries) AtParams(ctx context.Context, slug string, aq ...AdditionalQuery) ([]string, error) {
	query := atParams
	queryParams := []interface{}{slug}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const funcParams = `-- name: FuncParams :many
SELECT name FROM foo WHERE name = ?1
`

func (q *Queries) FuncParams(ctx context.Context, slug string, aq ...AdditionalQuery) ([]string, error) {
	query := funcParams
	queryParams := []interface{}{slug}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertAtParams = `-- name: InsertAtParams :one
INSERT INTO foo(name, bio) values (?1, ?2) returning name
`

type InsertAtParamsParams struct {
	Name string
	Bio  string
}

func (q *Queries) InsertAtParams(ctx context.Context, arg InsertAtParamsParams, aq ...AdditionalQuery) (string, error) {
	query := insertAtParams
	queryParams := []interface{}{arg.Name, arg.Bio}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var name string
	err := row.Scan(&name)
	return name, err
}

const insertFuncParams = `-- name: InsertFuncParams :one
INSERT INTO foo(name, bio) values (?1, ?2) returning name
`

type InsertFuncParamsParams struct {
	Name string
	Bio  string
}

func (q *Queries) InsertFuncParams(ctx context.Context, arg InsertFuncParamsParams, aq ...AdditionalQuery) (string, error) {
	query := insertFuncParams
	queryParams := []interface{}{arg.Name, arg.Bio}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var name string
	err := row.Scan(&name)
	return name, err
}
