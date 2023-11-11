-- name: DefaultSchemaCreate :one
INSERT INTO foo.bar (id, name) VALUES ($1, $2) RETURNING *;

-- name: DefaultSchemaSelect :one
SELECT * FROM foo.bar WHERE id = $1;

-- name: DefaultSchemaSelectAll :many
SELECT * FROM foo.bar;