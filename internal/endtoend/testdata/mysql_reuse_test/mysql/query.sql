-- name: ReuseParam :exec
UPDATE foo SET name = sqlc.arg(name) WHERE id = sqlc.arg(id) OR name = sqlc.arg(name);
