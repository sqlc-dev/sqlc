-- name: AliasBar :exec
DELETE FROM bar AS b
WHERE b.id = ?;
