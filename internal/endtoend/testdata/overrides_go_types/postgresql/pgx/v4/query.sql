-- name: LoadFoo :many
SELECT * FROM foo WHERE id = $1;
