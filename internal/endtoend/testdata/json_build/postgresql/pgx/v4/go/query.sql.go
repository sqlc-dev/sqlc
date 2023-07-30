// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgtype"
)

const selectJSONBBuildArray = `-- name: SelectJSONBBuildArray :one
SELECT 
  jsonb_build_array(),
  jsonb_build_array(1),
  jsonb_build_array(1, 2),
  jsonb_build_array(1, 2, 'foo'),
  jsonb_build_array(1, 2, 'foo', 4)
`

type SelectJSONBBuildArrayRow struct {
	JsonbBuildArray   pgtype.JSONB
	JsonbBuildArray_2 pgtype.JSONB
	JsonbBuildArray_3 pgtype.JSONB
	JsonbBuildArray_4 pgtype.JSONB
	JsonbBuildArray_5 pgtype.JSONB
}

func (q *Queries) SelectJSONBBuildArray(ctx context.Context) (SelectJSONBBuildArrayRow, error) {
	row := q.db.QueryRow(ctx, selectJSONBBuildArray)
	var i SelectJSONBBuildArrayRow
	err := row.Scan(
		&i.JsonbBuildArray,
		&i.JsonbBuildArray_2,
		&i.JsonbBuildArray_3,
		&i.JsonbBuildArray_4,
		&i.JsonbBuildArray_5,
	)
	return i, err
}

const selectJSONBBuildObject = `-- name: SelectJSONBBuildObject :one
SELECT
  jsonb_build_object(),
  jsonb_build_object('foo'),
  jsonb_build_object('foo', 1),
  jsonb_build_object('foo', 1, 2),
  jsonb_build_object('foo', 1, 2, 'bar')
`

type SelectJSONBBuildObjectRow struct {
	JsonbBuildObject   pgtype.JSONB
	JsonbBuildObject_2 pgtype.JSONB
	JsonbBuildObject_3 pgtype.JSONB
	JsonbBuildObject_4 pgtype.JSONB
	JsonbBuildObject_5 pgtype.JSONB
}

func (q *Queries) SelectJSONBBuildObject(ctx context.Context) (SelectJSONBBuildObjectRow, error) {
	row := q.db.QueryRow(ctx, selectJSONBBuildObject)
	var i SelectJSONBBuildObjectRow
	err := row.Scan(
		&i.JsonbBuildObject,
		&i.JsonbBuildObject_2,
		&i.JsonbBuildObject_3,
		&i.JsonbBuildObject_4,
		&i.JsonbBuildObject_5,
	)
	return i, err
}

const selectJSONBuildArray = `-- name: SelectJSONBuildArray :one
SELECT 
  json_build_array(),
  json_build_array(1),
  json_build_array(1, 2),
  json_build_array(1, 2, 'foo'),
  json_build_array(1, 2, 'foo', 4)
`

type SelectJSONBuildArrayRow struct {
	JsonBuildArray   pgtype.JSON
	JsonBuildArray_2 pgtype.JSON
	JsonBuildArray_3 pgtype.JSON
	JsonBuildArray_4 pgtype.JSON
	JsonBuildArray_5 pgtype.JSON
}

func (q *Queries) SelectJSONBuildArray(ctx context.Context) (SelectJSONBuildArrayRow, error) {
	row := q.db.QueryRow(ctx, selectJSONBuildArray)
	var i SelectJSONBuildArrayRow
	err := row.Scan(
		&i.JsonBuildArray,
		&i.JsonBuildArray_2,
		&i.JsonBuildArray_3,
		&i.JsonBuildArray_4,
		&i.JsonBuildArray_5,
	)
	return i, err
}

const selectJSONBuildObject = `-- name: SelectJSONBuildObject :one
SELECT
  json_build_object(),
  json_build_object('foo'),
  json_build_object('foo', 1),
  json_build_object('foo', 1, 2),
  json_build_object('foo', 1, 2, 'bar')
`

type SelectJSONBuildObjectRow struct {
	JsonBuildObject   pgtype.JSON
	JsonBuildObject_2 pgtype.JSON
	JsonBuildObject_3 pgtype.JSON
	JsonBuildObject_4 pgtype.JSON
	JsonBuildObject_5 pgtype.JSON
}

func (q *Queries) SelectJSONBuildObject(ctx context.Context) (SelectJSONBuildObjectRow, error) {
	row := q.db.QueryRow(ctx, selectJSONBuildObject)
	var i SelectJSONBuildObjectRow
	err := row.Scan(
		&i.JsonBuildObject,
		&i.JsonBuildObject_2,
		&i.JsonBuildObject_3,
		&i.JsonBuildObject_4,
		&i.JsonBuildObject_5,
	)
	return i, err
}
