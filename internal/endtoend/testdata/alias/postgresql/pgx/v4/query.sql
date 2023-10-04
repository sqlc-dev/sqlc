-- name: AliasBar :exec
DELETE FROM bar b
WHERE b.id = $1;
