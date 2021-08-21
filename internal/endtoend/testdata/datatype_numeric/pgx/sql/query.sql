-- name: CreateExample :one
INSERT INTO examples (value) VALUES ($1)
RETURNING *;