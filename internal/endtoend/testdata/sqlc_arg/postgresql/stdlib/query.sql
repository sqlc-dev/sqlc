CREATE TABLE foo (name text not null, description text);

-- name: FuncParamIdent :many
SELECT name FROM foo WHERE name = sqlc.arg(slug);

-- name: FuncParamString :many
SELECT name FROM foo WHERE name = sqlc.arg('slug');

-- name: FuncParamStringOptional :exec
UPDATE foo SET name = coalesce(sqlc.arg('slug?'), name);

-- name: FuncParamStringRequired :exec
UPDATE foo SET description = sqlc.arg('slug!');
