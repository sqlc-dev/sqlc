CREATE TABLE bar (id serial not null);

-- name: baz :one
SELECT foo FROM bar WHERE baz = $4;

-- name: bar
SELECT foo FROM bar WHERE baz = $1 AND baz = $3;

-- name: foo :one
SELECT foo FROM bar;

-- name: Named :many
SELECT id FROM bar WHERE id = $1 AND sqlc.arg(named) = true AND id = $5;
