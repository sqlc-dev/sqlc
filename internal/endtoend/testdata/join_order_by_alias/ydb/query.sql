-- name: ColumnAsOrderBy :many
SELECT a.email AS id
FROM foo a JOIN foo b ON a.email = b.email
ORDER BY id;
