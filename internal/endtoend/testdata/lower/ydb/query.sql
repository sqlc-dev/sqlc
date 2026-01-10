-- name: Lower :many
SELECT bar FROM foo WHERE bar = $bar AND Unicode::ToLower(bat) = $bat_lower;
