-- https://github.com/sqlc-dev/sqlc/issues/3783
-- name: CausesPgToBeImported :many
SELECT
    id AS author_id,
    name AS author_name,
    bio AS author_bio
FROM
    authors
WHERE
    (sqlc.narg('author_ids') IS NULL OR id IN (sqlc.slice('author_ids')));
