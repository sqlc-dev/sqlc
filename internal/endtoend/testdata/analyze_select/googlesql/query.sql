-- name: ListUsers :many
SELECT * FROM users;

-- name: CountUsers :one
SELECT COUNT(*) AS total FROM users;

-- name: ListUserPosts :many
SELECT u.name, p.* FROM users u JOIN posts p ON p.user_id = u.id WHERE u.name = @name;
