-- name: GetAll :many
SELECT * FROM users;

-- name: GetFirst :one
SELECT first_name FROM users LIMIT 1;
