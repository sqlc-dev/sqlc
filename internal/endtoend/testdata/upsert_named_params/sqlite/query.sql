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

-- name: UpsertAuthorDoUpdateWhere :one
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
WHERE excluded.name != @name
RETURNING *;

-- name: UpsertAuthorConflictWhere :one
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
WHERE excluded.bio != sqlc.narg('bio')
RETURNING *;

-- name: UpsertAuthorBothWhere :one
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
WHERE excluded.name != @name
AND excluded.bio != sqlc.narg('bio')
RETURNING *;
