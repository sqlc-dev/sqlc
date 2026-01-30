-- name: InsertUserAndReturnID :one
INSERT INTO users (name) VALUES ($name)
  RETURNING id;

-- name: InsertUserAndReturnUser :one
INSERT INTO users (name) VALUES ($name)
  RETURNING *;

-- name: UpdateUserAndReturnID :one
UPDATE users SET name = $name
  WHERE name = $name_2
  RETURNING id;

-- name: UpdateUserAndReturnUser :one
UPDATE users SET name = $name
  WHERE name = $name_2
  RETURNING *;

-- name: DeleteUserAndReturnID :one
DELETE FROM users
  WHERE name = $name
  RETURNING id;

-- name: DeleteUserAndReturnUser :one
DELETE FROM users
  WHERE name = $name
  RETURNING *;
