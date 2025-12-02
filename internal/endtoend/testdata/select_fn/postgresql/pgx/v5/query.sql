-- name: SelectFoos :many
SELECT * FROM foo_fn();

-- name: SelectFooEmbed :many
SELECT sqlc.embed(foo), 1 AS one FROM foo_fn() AS foo;

-- name: SelectSingleColumn :many
SELECT baz FROM foo_fn();
