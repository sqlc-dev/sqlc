# A macro whose body is an expression.
-- name: IsRecent :macro
-- @param ts
-- @param days signed
TIMESTAMPDIFF(DAY, sqlc.arg(ts), now()) < sqlc.arg(days);

/* A macro whose body is a query. */
-- name: ActiveAuthors :macro
-- @param org signed
SELECT *
FROM authors
WHERE org_id = sqlc.arg(org) # the tenant
  AND deleted_at IS NULL;

-- name: ListActiveAuthors :many
SELECT a.id, a.name
FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a
ORDER BY a.name;

-- name: ListRecentAuthors :many
SELECT r.id, r.name
FROM sqlc.macro.RecentActiveAuthors(sqlc.arg(org_id)) r
ORDER BY r.name;

-- Every ? is an argument of its own, numbered in the order the expanded
-- query sends them.
-- name: ListLaterAuthors :many
SELECT x.id, x.name
FROM authors x
JOIN sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a ON a.id = x.id
WHERE x.name <> ? AND x.id > ?;

-- name: GetSettledAuthor :one
SELECT * FROM authors
WHERE id = ? AND NOT sqlc.macro.IsRecent(created_at, 7 + 7);

-- The same macro twice in one query.
-- name: ListSharedNames :many
SELECT a.name
FROM sqlc.macro.ActiveAuthors(sqlc.arg(first_org)) a
JOIN sqlc.macro.ActiveAuthors(sqlc.arg(second_org)) b ON a.name = b.name;

-- A macro that calls other macros. window_days is not one of its parameters,
-- so it becomes a parameter of whichever query the macro lands in.
-- name: RecentActiveAuthors :macro
-- @param org
SELECT * FROM sqlc.macro.ActiveAuthors(sqlc.arg(org)) a
WHERE sqlc.macro.IsRecent(a.created_at, sqlc.arg(window_days));
