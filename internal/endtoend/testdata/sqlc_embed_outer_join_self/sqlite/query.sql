-- name: SelfJoin :many
SELECT sqlc.embed(a), sqlc.embed(b) FROM users AS a
LEFT JOIN users AS b ON a.id = b.id;
