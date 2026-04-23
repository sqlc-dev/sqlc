-- name: UpsertAuthor :one
INSERT INTO authors (
  id,
  name,
  bio
) VALUES (
  @id,
  @name,
  sqlc.narg('bio')
) ON CONFLICT(id) DO UPDATE SET
  name = @name,
  bio = sqlc.narg('bio')
RETURNING *;
