-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text,
  country_code CHAR(2) NOT NULL,
  titles TEXT[]
);

-- name: GetAuthor :one
SELECT * FROM authors
WHERE name = $1 AND country_code = $2 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, country_code, titles
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

-- name: DeleteAuthors :exec
DELETE FROM authors
WHERE id IN (sqlc.slice(ids)) AND name = $1;

-- name: CreateAuthorOnlyTitles :one
INSERT INTO authors (name, titles) VALUES ($1, $2) RETURNING *;
