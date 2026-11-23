-- The macro these queries call is defined in macros.sql, which is read after
-- this file.
-- name: ListAuthors :many
SELECT * FROM authors a WHERE sqlc.macro.NotDeleted(a.deleted_at) ORDER BY a.name;
