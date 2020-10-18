// Code generated by sqlc. DO NOT EDIT.
// source: city.sql

package ondeck

import (
	"context"
	"database/sql"
	"encoding/json"
)

const createCity = `-- name: CreateCity :execresult
INSERT INTO city (
    name,
    slug,
    data
) VALUES (
    ?,
    ?,
    ?
)
`

type CreateCityParams struct {
	Name string          `json:"name"`
	Slug string          `json:"slug"`
	Data json.RawMessage `json:"data"`
}

func (q *Queries) CreateCity(ctx context.Context, arg CreateCityParams) (sql.Result, error) {
	return q.exec(ctx, q.createCityStmt, createCity, arg.Name, arg.Slug, arg.Data)
}

const getCity = `-- name: GetCity :one
SELECT slug, name, data
FROM city
WHERE slug = ?
`

func (q *Queries) GetCity(ctx context.Context, slug string) (City, error) {
	row := q.queryRow(ctx, q.getCityStmt, getCity, slug)
	var i City
	err := row.Scan(&i.Slug, &i.Name, &i.Data)
	return i, err
}

const listCities = `-- name: ListCities :many
SELECT slug, name, data
FROM city
ORDER BY name
`

func (q *Queries) ListCities(ctx context.Context) ([]City, error) {
	rows, err := q.query(ctx, q.listCitiesStmt, listCities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []City
	for rows.Next() {
		var i City
		if err := rows.Scan(&i.Slug, &i.Name, &i.Data); err != nil {
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

const updateCityName = `-- name: UpdateCityName :exec
UPDATE city
SET name = ?
WHERE slug = ?
`

type UpdateCityNameParams struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (q *Queries) UpdateCityName(ctx context.Context, arg UpdateCityNameParams) error {
	_, err := q.exec(ctx, q.updateCityNameStmt, updateCityName, arg.Name, arg.Slug)
	return err
}
