-- name: Lower :many
SELECT bar FROM foo WHERE bar = $1 AND LOWER(bat) = $2;
