-- name: Rounded :many
SELECT ROUND(price, 1) AS rounded,
    ROUND(price, 1, 1) AS truncated,
    ROUND(ratio, 2, 1) AS ratio_truncated,
    ROUND(qty, -1, 1) AS qty_truncated,
    ROUND(fee, 2, 1) AS fee_truncated
FROM prices;
