-- name: FuncParams :many
SELECT name FROM foo WHERE name = sqlc.arg('slug') AND sqlc.arg(filter)::bool;

-- name: AtParams :many
SELECT name FROM foo WHERE name = @slug AND @filter::bool;

-- name: InsertFuncParams :one
INSERT INTO foo(name, bio) values (sqlc.arg('name'), sqlc.arg('bio')) returning name;

-- name: InsertAtParams :one
INSERT INTO foo(name, bio) values (@name, @bio) returning name;


-- name: Update :one
UPDATE foo
SET
  name = CASE WHEN @set_name::bool
    THEN @name::text
    ELSE name
    END
RETURNING *;
