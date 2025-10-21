-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :exec
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
);

-- name: UpdateAuthor :exec
UPDATE authors
SET name = $1, bio = $2
WHERE id = $3;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
