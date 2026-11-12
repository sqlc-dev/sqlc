-- name: Trimmed :many
SELECT LEFT(a, 3) AS lf, RIGHT(c, @n) AS rt, TRIM(N'x' FROM a) AS tr, TRIM(LEADING 'x' FROM c) AS tl FROM p;
