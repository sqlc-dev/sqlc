-- name: In :many
SELECT *
FROM bar
WHERE id IN ($p1, $p2);
