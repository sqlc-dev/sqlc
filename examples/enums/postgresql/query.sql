-- name: UserCreate :exec
INSERT INTO users (first_name, last_name, age, shirt_size)
VALUES ($1, $2, $3, $4);

-- name: ListUsers :many
SELECT * FROM users;

-- name: ListUsersByShirtSizes :many
SELECT * FROM users
WHERE shirt_size = ANY(@shirt_size::size[]);
