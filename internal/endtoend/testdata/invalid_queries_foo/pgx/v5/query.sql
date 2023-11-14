-- name: ListFoos
SELECT id FROM foo;

-- name: ListFoos :one :many
SELECT id FROM foo;

-- name: ListFoos :two
SELECT id FROM foo;

-- name: DeleteFoo :one
DELETE FROM foo WHERE id = $1;

-- name: UpdateFoo :one
UPDATE foo SET id = $2 WHERE id = $1;

-- name: InsertFoo :one
INSERT INTO foo (id) VALUES ($1);

-- name: InsertFoo :batchone
INSERT INTO foo (id) VALUES ($1);
