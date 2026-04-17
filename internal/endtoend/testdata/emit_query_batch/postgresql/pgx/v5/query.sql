-- name: GetUser :one
SELECT * FROM myschema.users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM myschema.users ORDER BY id;

-- name: CreateUser :one
INSERT INTO myschema.users (name, email) VALUES ($1, $2) RETURNING *;

-- name: UpdateUser :exec
UPDATE myschema.users SET name = $1, email = $2 WHERE id = $3;

-- name: DeleteUser :execrows
DELETE FROM myschema.users WHERE id = $1;

-- name: ArchiveUser :execresult
UPDATE myschema.users SET name = 'archived' WHERE id = $1;
