-- name: UpsertAuthor :exec
INSERT INTO authors (name, bio)
VALUES (?, ?)
ON DUPLICATE KEY
    UPDATE bio = ?;

-- name: UpsertAuthorNamed :exec
INSERT INTO authors (name, bio)
VALUES (?, sqlc.arg(bio))
ON DUPLICATE KEY
    UPDATE bio = sqlc.arg(bio);
