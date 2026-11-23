-- A macro whose body is an expression. It is called before it is defined in
-- the queries below, and after: a macro is expanded where it is called.
-- name: IsRecent :macro
-- @param ts timestamptz
-- @param days int
sqlc.arg(ts) > now() - @days * interval '1 day';

-- A macro whose body is a query.
-- name: ActiveAuthors :macro
-- @param org bigint
SELECT *
FROM authors
WHERE org_id = sqlc.arg(org) -- the tenant
  AND deleted_at IS NULL;

-- name: ListActiveAuthors :many
SELECT a.id, a.name
FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a
ORDER BY a.name;

-- org_id reaches ActiveAuthors through RecentActiveAuthors, and takes its
-- type from the cast ActiveAuthors declares.
-- name: CountRecentAuthors :one
SELECT count(*) FROM sqlc.macro.RecentActiveAuthors(sqlc.arg(org_id)) r;

-- An argument is an expression, and is evaluated as one.
-- name: GetSettledAuthor :one
SELECT * FROM authors
WHERE id = $1 AND NOT sqlc.macro.IsRecent(created_at, 7 + 7);

-- The same macro twice in one query.
-- name: ListSharedNames :many
SELECT a.name
FROM sqlc.macro.ActiveAuthors(sqlc.arg(first_org)) a
JOIN sqlc.macro.ActiveAuthors(sqlc.arg(second_org)) b ON a.name = b.name;

-- A placeholder the query numbers itself keeps its number, and the macro's
-- parameter is numbered around it.
-- name: ListAuthorsWithBio :many
SELECT a.id, a.name, a.bio
FROM authors a
JOIN sqlc.macro.ActiveAuthors(@org_id) active ON active.id = a.id
WHERE a.bio IS NOT NULL AND a.id > $1;

-- A macro that calls other macros. window_days is not one of its parameters,
-- so it becomes a parameter of whichever query the macro lands in.
-- name: RecentActiveAuthors :macro
-- @param org
SELECT * FROM sqlc.macro.ActiveAuthors(sqlc.arg(org)) a
WHERE sqlc.macro.IsRecent(a.created_at, sqlc.arg(window_days));

-- A macro without parameters.
-- name: Unnamed :macro
'unnamed';

-- name: ListNames :many
SELECT coalesce(nullif(name, ''), sqlc.macro.Unnamed())::text AS name FROM authors;
