-- name: LeftJoin :many
SELECT sqlc.embed(p), sqlc.embed(u) FROM posts AS p
LEFT JOIN users AS u ON p.user_id = u.id;
