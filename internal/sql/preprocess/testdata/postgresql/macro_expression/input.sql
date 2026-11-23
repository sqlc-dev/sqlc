-- name: IsRecent :macro
-- @param ts timestamptz
-- @param days
@ts > now() - sqlc.arg(days) * interval '1 day';

-- name: ListRecent :many
SELECT * FROM authors WHERE NOT sqlc.macro.IsRecent(created_at, 7) OR id = $1;
