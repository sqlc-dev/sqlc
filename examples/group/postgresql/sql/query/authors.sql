-- name: Get :one
SELECT * FROM authors
WHERE author_id = $1;

-- name: Create :one
INSERT INTO authors (name) VALUES ($1)
RETURNING *;
