// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createMemories = `-- name: CreateMemories :many
INSERT INTO memories (vampire_id)
SELECT
    unnest($1::uuid[]) AS vampire_id
RETURNING
    id, vampire_id, created_at, updated_at
`

func (q *Queries) CreateMemories(ctx context.Context, vampireID []uuid.UUID, aq ...AdditionalQuery) ([]Memory, error) {
	query := createMemories
	queryParams := []interface{}{pq.Array(vampireID)}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Memory
	for rows.Next() {
		var i Memory
		if err := rows.Scan(
			&i.ID,
			&i.VampireID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
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

const getVampireIDs = `-- name: GetVampireIDs :many
SELECT vampires.id::uuid FROM unnest($1::uuid[]) AS vampires (id)
`

func (q *Queries) GetVampireIDs(ctx context.Context, vampireID []uuid.UUID, aq ...AdditionalQuery) ([]uuid.UUID, error) {
	query := getVampireIDs
	queryParams := []interface{}{pq.Array(vampireID)}

	if len(aq) > 0 {
		query += " " + aq[0].SQL
		queryParams = append(queryParams, aq[0].Args...)
	}

	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var vampires_id uuid.UUID
		if err := rows.Scan(&vampires_id); err != nil {
			return nil, err
		}
		items = append(items, vampires_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
