-- name: ActiveAuthors :macro
-- @param org integer
SELECT * FROM authors WHERE org_id = @org AND deleted_at IS NULL;

-- A numbered placeholder keeps its number, and the macro's arguments are
-- numbered around it.
-- name: ListActiveAuthors :many
SELECT a.id
FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a
JOIN sqlc.macro.ActiveAuthors(?1) b ON b.id = a.id
WHERE a.name = @name AND a.id IN (sqlc.slice(ids));
