-- name: Mixed :many
SELECT id FROM events WHERE name = {name:String} AND id > sqlc.arg(min_id);
