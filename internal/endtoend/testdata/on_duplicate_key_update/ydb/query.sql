-- name: UpsertAuthor :exec
UPSERT INTO authors (name, bio)
VALUES ($name, $bio);

-- name: UpsertAuthorNamed :exec
UPSERT INTO authors (name, bio)
VALUES ($name, $bio);

