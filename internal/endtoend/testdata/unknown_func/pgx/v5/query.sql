-- name: ListFoos :one
SELECT id FROM foo WHERE id = frobnicate($1);
