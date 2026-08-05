-- name: CreateUser :one
INSERT INTO users (name, bio) VALUES (?, ?) RETURNING *;

-- name: UpdateBio :exec
UPDATE users SET bio = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
