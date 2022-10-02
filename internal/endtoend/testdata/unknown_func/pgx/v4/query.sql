CREATE TABLE foo (id text not null);
-- name: ListFoos :one
SELECT id FROM foo WHERE id = frobnicate($1);
