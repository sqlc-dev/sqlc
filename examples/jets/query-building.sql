-- name: CountPilots :one
SELECT COUNT(*) FROM pilots;

-- name: ListPilots :many
SELECT * FROM pilots LIMIT 5;

-- name: DeletePilot :exec
DELETE FROM pilots WHERE id = $1;
