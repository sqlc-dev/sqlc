-- name: NotDeleted :macro
-- @param deleted_at timestamptz
sqlc.arg(deleted_at) IS NULL;
