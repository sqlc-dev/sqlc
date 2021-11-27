// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
	"fmt"
	"strings"
)

const funcParamIdent = `-- name: FuncParamIdent :many
SELECT name FROM foo
WHERE name = ?
  AND id IN (/*REPLACE:favourites*/?)
`

type FuncParamIdentParams struct {
	Slug       string
	Favourites []int32
}

func (q *Queries) FuncParamIdent(ctx context.Context, arg FuncParamIdentParams) ([]string, error) {
	sql := funcParamIdent
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Slug)
	if len(arg.Favourites) == 0 {
		return nil, fmt.Errorf("slice FuncParamIdentParams.Favourites must have at least one element")
	}
	for _, v := range arg.Favourites {
		queryParams = append(queryParams, v)
	}
	sql = strings.Replace(sql, "/*REPLACE:favourites*/?", strings.Repeat(",?", len(arg.Favourites))[1:], 1)
	rows, err := q.db.QueryContext(ctx, sql, queryParams...)
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

const funcParamSoloArg = `-- name: FuncParamSoloArg :many
SELECT name FROM foo
WHERE id IN (/*REPLACE:favourites*/?)
`

func (q *Queries) FuncParamSoloArg(ctx context.Context, favourites []int32) ([]string, error) {
	sql := funcParamSoloArg
	var queryParams []interface{}
	if len(favourites) == 0 {
		return nil, fmt.Errorf("slice favourites must have at least one element")
	}
	for _, v := range favourites {
		queryParams = append(queryParams, v)
	}
	sql = strings.Replace(sql, "/*REPLACE:favourites*/?", strings.Repeat(",?", len(favourites))[1:], 1)
	rows, err := q.db.QueryContext(ctx, sql, queryParams...)
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

const funcParamString = `-- name: FuncParamString :many
SELECT name FROM foo
WHERE name = ?
  AND id IN (/*REPLACE:favourites*/?)
`

type FuncParamStringParams struct {
	Slug       string
	Favourites []int32
}

func (q *Queries) FuncParamString(ctx context.Context, arg FuncParamStringParams) ([]string, error) {
	sql := funcParamString
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Slug)
	if len(arg.Favourites) == 0 {
		return nil, fmt.Errorf("slice FuncParamStringParams.Favourites must have at least one element")
	}
	for _, v := range arg.Favourites {
		queryParams = append(queryParams, v)
	}
	sql = strings.Replace(sql, "/*REPLACE:favourites*/?", strings.Repeat(",?", len(arg.Favourites))[1:], 1)
	rows, err := q.db.QueryContext(ctx, sql, queryParams...)
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
