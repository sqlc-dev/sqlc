// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
	"strings"
)

const funcParamIdent = `-- name: FuncParamIdent :many
SELECT name FROM foo
WHERE name = $1
  AND id IN ($2)
`

type FuncParamIdentParams struct {
	Slug       string
	Favourites []int32
}

func (q *Queries) FuncParamIdent(ctx context.Context, arg FuncParamIdentParams, aq ...AdditionalQuery) ([]string, error) {
	query := funcParamIdent
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Slug)
	if len(arg.Favourites) > 0 {
		for _, v := range arg.Favourites {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:favourites*/?", strings.Repeat(",?", len(arg.Favourites))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:favourites*/?", "NULL", 1)
	}

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

const funcParamString = `-- name: FuncParamString :many
SELECT name FROM foo
WHERE name = $1
  AND id IN ($2)
`

type FuncParamStringParams struct {
	Slug       string
	Favourites []int32
}

func (q *Queries) FuncParamString(ctx context.Context, arg FuncParamStringParams, aq ...AdditionalQuery) ([]string, error) {
	query := funcParamString
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Slug)
	if len(arg.Favourites) > 0 {
		for _, v := range arg.Favourites {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:favourites*/?", strings.Repeat(",?", len(arg.Favourites))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:favourites*/?", "NULL", 1)
	}

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
