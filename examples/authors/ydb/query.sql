-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $id LIMIT 1;

-- name: ListAuthors :many 
SELECT * FROM authors
ORDER BY name;

-- name: CreateOrUpdateAuthor :one 
UPSERT INTO authors (name, bio) 
VALUES (
  $name, $bio
)
RETURNING *;

-- name: DeleteAuthor :exec 
DELETE FROM authors WHERE id = $id;

-- name: DropTable :exec
DROP TABLE IF EXISTS authors;
