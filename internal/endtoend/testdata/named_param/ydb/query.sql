-- name: FuncParams :many
SELECT name FROM foo WHERE name = sqlc.arg(slug) AND CAST(sqlc.arg(filter) AS Bool);

-- name: AtParams :many
SELECT name FROM foo WHERE name = $slug AND CAST($filter AS Bool);

-- name: InsertFuncParams :one
INSERT INTO foo(name, bio) VALUES (sqlc.arg(name), sqlc.arg(bio)) RETURNING name;

-- name: InsertAtParams :one
INSERT INTO foo(name, bio) VALUES ($name, $bio) RETURNING name;

-- name: Update :one
UPDATE foo
SET
  name = CASE WHEN CAST($set_name AS Bool)
    THEN CAST($name AS Text)
    ELSE name
    END
RETURNING *;



