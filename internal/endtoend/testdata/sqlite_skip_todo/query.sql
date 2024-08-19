-- name: PragmaForeignKeysEnable :exec
PRAGMA foreign_keys = 1;

-- name: ListFoo :many
SELECT * FROM foo;

-- name: PragmaForeignKeysDisable :exec
PRAGMA foreign_keys = 0;

-- name: PragmaForeignKeysGet :one
PRAGMA foreign_keys;

-- name: GetFoo :many
SELECT * FROM foo
WHERE bar = ?;

