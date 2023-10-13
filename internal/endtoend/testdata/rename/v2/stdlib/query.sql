-- name: ListFoo :many
SELECT id_old as foo_old, id_old as baz_old
FROM bar_old
WHERE ip_old = $1 AND id_old = $2;

-- name: ListBar :many
SELECT * FROM bar_old;

