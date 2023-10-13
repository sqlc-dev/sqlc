-- name: UpsertAuthor :exec
INSERT INTO authors (name, bio)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE
SET bio = $2;

-- name: UpsertAuthorNamed :exec
INSERT INTO authors (name, bio)
VALUES ($1, sqlc.arg(bio))
ON CONFLICT (name) DO UPDATE
SET bio = sqlc.arg(bio);
