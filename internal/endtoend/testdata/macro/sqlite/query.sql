-- A macro whose body is an expression.
-- name: HasBio :macro
-- @param bio
sqlc.arg(bio) IS NOT NULL AND @bio <> '';

-- A macro whose body is a query.
-- name: ActiveAuthors :macro
-- @param org integer
SELECT *
FROM authors
WHERE org_id = sqlc.arg(org) -- the tenant
  AND deleted_at IS NULL;

-- name: ListActiveAuthors :many
SELECT a.id, a.name
FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a
ORDER BY a.name;

-- name: ListAuthorsWithBio :many
SELECT w.id, w.name
FROM sqlc.macro.AuthorsWithBio(sqlc.arg(org_id)) w
ORDER BY w.name;

-- A numbered placeholder keeps its number, and the macro's parameters are
-- numbered around it.
-- name: ListLaterAuthors :many
SELECT x.id, x.name
FROM authors x
JOIN sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a ON a.id = x.id
WHERE x.id > ?1;

-- The same macro twice in one query.
-- name: ListSharedNames :many
SELECT a.name
FROM sqlc.macro.ActiveAuthors(@first_org) a
JOIN sqlc.macro.ActiveAuthors(@second_org) b ON a.name = b.name;

-- A macro that calls other macros.
-- name: AuthorsWithBio :macro
-- @param org
SELECT * FROM sqlc.macro.ActiveAuthors(sqlc.arg(org)) a
WHERE sqlc.macro.HasBio(a.bio);
