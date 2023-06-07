CREATE TABLE foo (
    bar text,
    baz bigint
);

-- name: SumBaz :many
SELECT bar, coalesce(sum(baz), 0) as quantity
FROM foo
GROUP BY 1;
