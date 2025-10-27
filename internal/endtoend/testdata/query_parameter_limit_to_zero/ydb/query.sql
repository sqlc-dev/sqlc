-- name: GetAuthor :one
SELECT * FROM authors
WHERE name = $name AND country_code = $country_code LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, country_code
) VALUES (
  $name, $bio, $country_code
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $id;
