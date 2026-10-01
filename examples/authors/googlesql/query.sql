-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = @id;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  id, name, bio
) VALUES (
  @id, @name, @bio
) THEN RETURN id, name, bio;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = @id;
