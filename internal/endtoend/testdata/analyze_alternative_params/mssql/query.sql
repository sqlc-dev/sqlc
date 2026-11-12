-- name: Alternatives :many
SELECT GREATEST(1, @n) AS g, LEAST(n, @cap) AS l, COALESCE(@x, name) AS c,
    CASE WHEN id > 0 THEN @label ELSE name END AS k, IIF(id > 0, @alt, at) AS i
FROM t;
