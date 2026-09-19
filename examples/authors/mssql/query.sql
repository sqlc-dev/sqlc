-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = @id;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) OUTPUT INSERTED.id, INSERTED.name, INSERTED.bio
VALUES (
  @name, @bio
);

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = @id;
