-- name: GetAuthor :one
SELECT * FROM authors
WHERE name = $name AND country_code = $country_code LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, country_code, titles
) VALUES (
  $name, $bio, $country_code, $titles
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $id;

-- name: DeleteAuthors :exec
DELETE FROM authors
WHERE id IN sqlc.slice(ids) AND name = $name;

-- name: CreateAuthorOnlyTitles :one
INSERT INTO authors (name, titles) VALUES ($name, $titles) RETURNING *;

-- name: AddNewClient :one
INSERT INTO clients (
  id, name
) VALUES (
  $id, $name
)
RETURNING *;
