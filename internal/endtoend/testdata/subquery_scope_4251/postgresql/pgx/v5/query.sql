-- name: GetT1FromT2 :many
-- Unqualified `id` inside the subquery must bind to t2.id (innermost
-- FROM-clause scope), not be flagged as ambiguous against t1.id. See
-- issue #4251.
SELECT id FROM t1
WHERE id IN (
    SELECT t1_id
    FROM t2
    WHERE id = $1
);
