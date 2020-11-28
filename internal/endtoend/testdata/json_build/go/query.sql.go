// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
	"encoding/json"
)

const selectJSONBuildArray = `-- name: SelectJSONBuildArray :one
SELECT 
  json_build_array(),
  json_build_array(1),
  json_build_array(1, 2),
  json_build_array(1, 2, 'foo'),
  json_build_array(1, 2, 'foo', 4)
`

type SelectJSONBuildArrayRow struct {
	JsonBuildArray   json.RawMessage
	JsonBuildArray_2 json.RawMessage
	JsonBuildArray_3 json.RawMessage
	JsonBuildArray_4 json.RawMessage
	JsonBuildArray_5 json.RawMessage
}

func (q *Queries) SelectJSONBuildArray(ctx context.Context) (SelectJSONBuildArrayRow, error) {
	row := q.db.QueryRowContext(ctx, selectJSONBuildArray)
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
	JsonBuildObject   json.RawMessage
	JsonBuildObject_2 json.RawMessage
	JsonBuildObject_3 json.RawMessage
	JsonBuildObject_4 json.RawMessage
	JsonBuildObject_5 json.RawMessage
}

func (q *Queries) SelectJSONBuildObject(ctx context.Context) (SelectJSONBuildObjectRow, error) {
	row := q.db.QueryRowContext(ctx, selectJSONBuildObject)
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
