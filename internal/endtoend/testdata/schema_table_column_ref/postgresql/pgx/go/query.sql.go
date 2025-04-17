// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package querytest

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getTotalSlackQueries = `-- name: GetTotalSlackQueries :one
SELECT
    COUNT(*) AS count
FROM astoria.slack_feedback
WHERE astoria.slack_feedback.workspace_id = $1
AND created_at BETWEEN $2::date AND $3::date
`

type GetTotalSlackQueriesParams struct {
	WorkspaceID int64
	Column2     pgtype.Date
	Column3     pgtype.Date
}

func (q *Queries) GetTotalSlackQueries(ctx context.Context, arg GetTotalSlackQueriesParams) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalSlackQueries, arg.WorkspaceID, arg.Column2, arg.Column3)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalSlackQueriesRequestsCreated = `-- name: GetTotalSlackQueriesRequestsCreated :one
SELECT
    COUNT(*) AS count
FROM astoria.tickets
WHERE astoria.tickets.workspace_id = $1
  AND source = 'RAISED_FROM_BOT'
  AND created_at BETWEEN $2::date AND $3::date
`

type GetTotalSlackQueriesRequestsCreatedParams struct {
	WorkspaceID int64
	Column2     pgtype.Date
	Column3     pgtype.Date
}

func (q *Queries) GetTotalSlackQueriesRequestsCreated(ctx context.Context, arg GetTotalSlackQueriesRequestsCreatedParams) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalSlackQueriesRequestsCreated, arg.WorkspaceID, arg.Column2, arg.Column3)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalSlackQueriesResolved = `-- name: GetTotalSlackQueriesResolved :one
SELECT
    COUNT(*) AS count
FROM astoria.slack_feedback
WHERE astoria.slack_feedback.workspace_id = $1
  AND (issue_raised = false OR issue_raised IS NULL)
  AND created_at BETWEEN $2::date AND $3::date
`

type GetTotalSlackQueriesResolvedParams struct {
	WorkspaceID int64
	Column2     pgtype.Date
	Column3     pgtype.Date
}

func (q *Queries) GetTotalSlackQueriesResolved(ctx context.Context, arg GetTotalSlackQueriesResolvedParams) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalSlackQueriesResolved, arg.WorkspaceID, arg.Column2, arg.Column3)
	var count int64
	err := row.Scan(&count)
	return count, err
}
