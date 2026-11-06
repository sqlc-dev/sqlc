-- name: LeftJoin :many
SELECT sqlc.embed(posts), sqlc.embed(users) FROM posts
LEFT JOIN users ON posts.user_id = users.id;
