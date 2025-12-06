-- CASE expression
-- name: GetTransactionClassification :many
SELECT
    id,
    user_id,
    amount,
    CASE
        WHEN amount > 1000 THEN 'high'
        WHEN amount > 100 THEN 'medium'
        ELSE 'low'
    END as classification,
    category
FROM transactions
WHERE created_at >= ?
ORDER BY amount DESC;

-- IN operator with subquery style
-- name: GetTransactionsByCategory :many
SELECT
    id,
    user_id,
    amount,
    category
FROM transactions
WHERE category IN ('groceries', 'utilities', 'transportation')
ORDER BY created_at DESC;
