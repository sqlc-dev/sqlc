CREATE TABLE foo (bar text);

-- name: Coalesce :many
SELECT coalesce(bar, '') as login
FROM foo;
