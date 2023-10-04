-- name: FuncParams :many
SELECT name FROM foo WHERE name = sqlc.arg('slug');

-- name: AtParams :many
SELECT name FROM foo WHERE name = @slug;

-- name: InsertFuncParams :one
INSERT INTO foo(name, bio) values (sqlc.arg('name'), sqlc.arg('bio')) returning name;

-- name: InsertAtParams :one
INSERT INTO foo(name, bio) values (@name, @bio) returning name;
