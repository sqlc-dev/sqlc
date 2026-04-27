-- name: BrokenQuery :many
SELECT *
FROM mytable
WHERE
    typ IN (sqlc.slice(types))
    AND (sqlc.arg(allnames) OR (name IN (sqlc.slice(names))));
