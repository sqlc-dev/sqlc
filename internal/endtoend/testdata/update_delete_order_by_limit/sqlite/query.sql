-- name: DeleteLimitOffset :exec
DELETE FROM t WHERE name = ? ORDER BY score DESC LIMIT ? OFFSET ?;

-- name: DeleteOrderByParam :exec
DELETE FROM t WHERE name = ? ORDER BY abs(score - ?) LIMIT ?;

-- name: UpdateLimitOffset :exec
UPDATE t SET score = ? WHERE name = ? ORDER BY id LIMIT ? OFFSET ?;

-- name: UpdateOrderByParam :exec
UPDATE t SET score = ? WHERE name = ? ORDER BY abs(score - ?) LIMIT 1;
