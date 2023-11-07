-- name: GetTotalSlackQueries :one
SELECT
    COUNT(*) AS count
FROM astoria.slack_feedback
WHERE astoria.slack_feedback.workspace_id = $1
AND created_at BETWEEN $2::date AND $3::date;

-- name: GetTotalSlackQueriesResolved :one
SELECT
    COUNT(*) AS count
FROM astoria.slack_feedback
WHERE astoria.slack_feedback.workspace_id = $1
  AND (issue_raised = false OR issue_raised IS NULL)
  AND created_at BETWEEN $2::date AND $3::date;

-- name: GetTotalSlackQueriesRequestsCreated :one
SELECT
    COUNT(*) AS count
FROM astoria.tickets
WHERE astoria.tickets.workspace_id = $1
  AND source = 'RAISED_FROM_BOT'
  AND created_at BETWEEN $2::date AND $3::date;
