-- A macro is blanked, and expands where it is called.
-- name: ActiveAuthors :macro
-- @param org bigint
SELECT *
FROM authors
WHERE org_id = sqlc.arg(org) -- the tenant
  AND deleted_at IS NULL;

-- name: ListActiveAuthors :many
SELECT a.id FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a WHERE a.name = sqlc.arg(name);

-- name: CountAuthors :one
SELECT count(*) FROM authors WHERE org_id = sqlc.arg(org_id);
