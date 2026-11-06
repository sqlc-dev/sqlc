-- name: RightJoin :many
SELECT sqlc.embed(p), sqlc.embed(u) FROM posts AS p
RIGHT JOIN users AS u ON p.user_id = u.id;
