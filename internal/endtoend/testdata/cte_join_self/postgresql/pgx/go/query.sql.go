// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"
)

const listItems = `-- name: ListItems :one
WITH
    items1 AS (SELECT 'id'::TEXT AS id, 'name'::TEXT AS name),
    items2 AS (SELECT 'id'::TEXT AS id, 'name'::TEXT AS name)
SELECT
    i1.id AS id1,
    i2.id AS id2
FROM
    items1 i1
        JOIN items1 i2 ON 1 = 1
`

type ListItemsRow struct {
	Id1 string
	Id2 string
}

func (q *Queries) ListItems(ctx context.Context) (ListItemsRow, error) {
	row := q.db.QueryRow(ctx, listItems)
	var i ListItemsRow
	err := row.Scan(&i.Id1, &i.Id2)
	return i, err
}
