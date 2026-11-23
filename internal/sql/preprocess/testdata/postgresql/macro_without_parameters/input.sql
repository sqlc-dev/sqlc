-- name: Unnamed :macro
'unnamed';

-- name: Names :many
SELECT coalesce(name, sqlc.macro.Unnamed()) FROM authors;
