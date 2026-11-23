-- A macro is never checked on its own. This one expands, and the query it
-- leaves does not analyze: the error is reported against the query that
-- called it.
-- name: Archived :macro
-- @param id
SELECT * FROM foo_archive WHERE id = sqlc.arg(id);

-- name: ListFoo :many
SELECT * FROM foo;

-- name: ListArchived :many
SELECT * FROM foo
WHERE id IN (SELECT a.id FROM sqlc.macro.Archived(sqlc.arg(id)) a);
