# A macro is blanked, and expands where it is called.
-- name: ActiveAuthors :macro
-- @param org signed
SELECT /*+ MAX_EXECUTION_TIME(1000) */ *
FROM `authors`
WHERE org_id = sqlc.arg(org) # the tenant
  AND deleted_at IS NULL;

-- Every ? is an argument of its own, numbered in the order the expanded query
-- sends them.
-- name: ListActiveAuthors :many
SELECT a.id
FROM authors x
JOIN sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a ON a.id = x.id
JOIN sqlc.macro.ActiveAuthors(?) b ON b.id = x.id
WHERE x.name = ? AND x.bio <> sqlc.arg(bio) AND @uservar IS NULL;
