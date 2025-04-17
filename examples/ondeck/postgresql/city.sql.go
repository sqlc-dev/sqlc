// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: city.sql

package ondeck

import (
	"context"
)

const createCity = `-- name: CreateCity :one
INSERT INTO city (
    name,
    slug
) VALUES (
    $1,
    $2
) RETURNING slug, name
`

type CreateCityParams struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Create a new city. The slug must be unique.
// This is the second line of the comment
// This is the third line
func (q *Queries) CreateCity(ctx context.Context, arg CreateCityParams) (City, error) {
	row := q.queryRow(ctx, q.createCityStmt, createCity, arg.Name, arg.Slug)
	var i City
	err := row.Scan(&i.Slug, &i.Name)
	return i, err
}

const getCity = `-- name: GetCity :one
SELECT slug, name
FROM city
WHERE slug = $1
`

func (q *Queries) GetCity(ctx context.Context, slug string) (City, error) {
	row := q.queryRow(ctx, q.getCityStmt, getCity, slug)
	var i City
	err := row.Scan(&i.Slug, &i.Name)
	return i, err
}

const listCities = `-- name: ListCities :many
SELECT slug, name
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
		if err := rows.Scan(&i.Slug, &i.Name); err != nil {
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
SET name = $2
WHERE slug = $1
`

type UpdateCityNameParams struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (q *Queries) UpdateCityName(ctx context.Context, arg UpdateCityNameParams) error {
	_, err := q.exec(ctx, q.updateCityNameStmt, updateCityName, arg.Slug, arg.Name)
	return err
}
