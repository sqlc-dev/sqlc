CREATE TABLE bar (id integer NOT NULL PRIMARY KEY AUTOINCREMENT);

-- name: AliasBar :exec
DELETE FROM bar AS b
WHERE b.id = ?;
