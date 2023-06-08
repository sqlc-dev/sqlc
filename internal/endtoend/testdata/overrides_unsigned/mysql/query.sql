CREATE TABLE authors (
  id     SERIAL,
  name   text NOT NULL,
  rating bigint NOT NULL,
  score  int UNSIGNED NOT NULL,
  bio  text
);

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: CreateAuthor :exec
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
);

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;