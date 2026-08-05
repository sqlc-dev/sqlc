-- name: CreateUser :one
INSERT INTO users (name, bio) VALUES ($1, $2) RETURNING *;

-- name: UpdateBio :exec
UPDATE users SET bio = $1 WHERE id = $2;

-- name: UpdateNameReturning :one
UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteUserReturning :one
DELETE FROM users WHERE id = $1 RETURNING id;
