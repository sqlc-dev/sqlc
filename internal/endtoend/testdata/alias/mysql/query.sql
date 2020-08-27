CREATE TABLE bar (id serial not null);

-- name: AliasBar :exec
DELETE FROM bar b
WHERE b.id = ?;
