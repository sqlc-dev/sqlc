// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"
)

const schemaScopedCreate = `-- name: SchemaScopedCreate :one
INSERT INTO foo.bar (id, name) VALUES ($1, $2) RETURNING id
`

type SchemaScopedCreateParams struct {
	ID   int32
	Name string
}

func (q *Queries) SchemaScopedCreate(ctx context.Context, arg SchemaScopedCreateParams, aq ...AdditionalQuery) (int32, error) {
	query := schemaScopedCreate
	queryParams := []interface{}{arg.ID, arg.Name}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var id int32
	err := row.Scan(&id)
	return id, err
}
