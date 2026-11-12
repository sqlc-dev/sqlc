-- name: Weights :many
SELECT a.id, a.name, v.weight
FROM authors a
JOIN (VALUES (1, 10), (2, 20)) AS v(id, weight) ON v.id = a.id;

-- name: Correlated :many
SELECT a.id, v.n
FROM authors a
CROSS APPLY (VALUES (1), (2)) v(n)
WHERE v.n < a.id;

-- name: Lateral :many
SELECT a.id
FROM authors a
CROSS APPLY (VALUES (a.id), (a.id + 1)) v(n)
WHERE v.n > 1;
