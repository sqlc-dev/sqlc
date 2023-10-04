-- name: CastCoalesce :many
SELECT coalesce(bar, '')::text as login
FROM foo;
