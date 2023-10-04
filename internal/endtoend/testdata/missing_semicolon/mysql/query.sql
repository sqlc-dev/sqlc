-- name: SetAuthor :exec
UPDATE  authors
SET     name = ?
WHERE   id = ?
