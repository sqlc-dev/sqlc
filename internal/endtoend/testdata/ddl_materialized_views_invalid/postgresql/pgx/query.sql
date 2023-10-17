-- name: GetTotalEarned :one
SELECT COALESCE(SUM(earned), 0) as total_earned
FROM grouped_kpis
WHERE day = @day;
