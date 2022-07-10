-- name: InsertUser :one
INSERT INTO users (name) VALUES ($1) RETURNING id;

-- name: UpdateUser :one
UPDATE users SET name = $1
  RETURNING id;

-- name: DeleteUser :one
DELETE FROM users
  WHERE name = $1
  RETURNING id;
