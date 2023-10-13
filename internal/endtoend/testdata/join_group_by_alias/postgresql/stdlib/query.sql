-- name: ColumnAsGroupBy :many
SELECT a.email AS id
FROM foo a JOIN foo b ON a.email = b.email
GROUP BY id;
