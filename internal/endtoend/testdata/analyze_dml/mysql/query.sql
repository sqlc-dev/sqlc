-- name: CreateUser :exec
INSERT INTO users (id, name, bio) VALUES (?, ?, ?);

-- name: UpdateBio :exec
UPDATE users SET bio = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: DeleteUserPosts :exec
DELETE p FROM posts p JOIN users u ON u.id = p.user_id WHERE u.name = ?;
