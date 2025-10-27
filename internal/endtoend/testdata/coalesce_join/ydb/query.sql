-- name: GetBar :many
SELECT foo.*, COALESCE(bar.id, 0) AS bar_id
FROM foo
LEFT JOIN bar ON foo.id = bar.id;
