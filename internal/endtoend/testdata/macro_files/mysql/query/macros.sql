-- A file of nothing but macros holds no queries, and generates no code.
-- name: NotDeleted :macro
-- @param deleted_at
sqlc.arg(deleted_at) IS NULL;
