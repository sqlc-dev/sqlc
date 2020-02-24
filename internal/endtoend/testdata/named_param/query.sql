CREATE TABLE foo (name text not null, bio text not null);

-- name: FuncParams :many
SELECT name FROM foo WHERE name = sqlc.arg('slug') AND sqlc.arg(filter)::bool;

-- name: AtParams :many
SELECT name FROM foo WHERE name = @slug AND @filter::bool;

-- name: Update :one
UPDATE foo
SET
  name = CASE WHEN @set_name::bool
    THEN @name::text
    ELSE name
    END
RETURNING *;
