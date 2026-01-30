-- name: SumBaz :many
SELECT bar, CAST(COALESCE(SUM(baz), 0) AS Int64) AS quantity
FROM foo
GROUP BY bar;
