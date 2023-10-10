// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getSubordinates = `-- name: GetSubordinates :many
WITH RECURSIVE subordinates(name, manager) AS (
    SELECT
        NULL, $1::TEXT
    UNION
    SELECT
        s.manager, e.name
    FROM 
        subordinates AS s
    LEFT OUTER JOIN
        employees AS e
    ON
        e.manager = s.manager
    WHERE
        s.manager IS NOT NULL
)
SELECT 
    s.name
FROM
    subordinates AS s
WHERE
    s.name != $1
`

func (q *Queries) GetSubordinates(ctx context.Context, name pgtype.Text) ([]pgtype.Text, error) {
	rows, err := q.db.Query(ctx, getSubordinates, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var name pgtype.Text
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
