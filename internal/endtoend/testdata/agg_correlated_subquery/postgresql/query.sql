-- name: Correlated :many
SELECT (SELECT ARRAY_AGG(o.*) FROM b) FROM o;
