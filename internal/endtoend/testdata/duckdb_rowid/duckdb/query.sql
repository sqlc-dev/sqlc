-- name: ListPlanRowIDs :many
SELECT rowid, id FROM plan;

-- name: ListQualifiedPlanRowIDs :many
SELECT p.rowid, p.id FROM plan p ORDER BY p.rowid;

-- name: GetPlanByRowID :one
SELECT id FROM plan WHERE rowid = $1;

-- name: ListPlans :many
SELECT * FROM plan ORDER BY rowid;

-- name: CreatePlan :exec
INSERT INTO plan VALUES ($1, $2, $3);

-- name: UpdatePlanPrice :exec
UPDATE plan SET price = $1 WHERE rowid = $2;

-- name: DeletePlan :one
DELETE FROM plan WHERE rowid = $1 RETURNING *;

-- name: ListLedger :many
SELECT * FROM ledger;

-- name: GetLedgerAmount :one
SELECT amount FROM ledger WHERE rowid = $1;
