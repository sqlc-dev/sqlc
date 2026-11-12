-- name: Trimmed :many
SELECT LEFT(a, 3) AS lf, RIGHT(c, @n) AS rt, TRIM(N'x' FROM a) AS tr, TRIM(LEADING 'x' FROM c) AS tl FROM p;

-- name: Repeated :many
SELECT REPLICATE(a, @r) AS r, CHARINDEX(N'x', a, @ci) AS ci FROM p;
