-- name: WithRowID :many
SELECT rowid, id FROM plan;

-- name: WithQualifiedRowID :many
SELECT p.rowid, p.id FROM plan p;

-- name: ByRowID :one
SELECT id FROM plan WHERE rowid = $1;

-- name: OrderByRowID :many
SELECT id FROM plan ORDER BY rowid;

-- name: StarSkipsRowID :many
SELECT * FROM plan;

-- name: DeclaredRowID :many
SELECT rowid, name FROM tag;
