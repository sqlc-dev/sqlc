-- name: ActiveAuthors :macro
-- @param org bigint
SELECT * FROM authors WHERE org_id = sqlc.arg(org) AND deleted_at IS NULL;

-- An embed is found by where it sits in the query, and a macro that expands
-- before it moves it.
-- name: ListBooksOfActiveAuthors :many
WITH active AS sqlc.macro.ActiveAuthors(sqlc.arg(org_id))
SELECT active.name AS author_name, sqlc.embed(b)
FROM active
JOIN books b ON b.author_id = active.id;

-- name: ListActiveAuthorsAndBooks :many
SELECT sqlc.embed(a), sqlc.embed(b)
FROM authors a
JOIN books b ON b.author_id = a.id
WHERE a.id IN (SELECT active.id FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) active);
