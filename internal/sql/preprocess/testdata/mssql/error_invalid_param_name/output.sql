-- name: Bad :many
SELECT id FROM t WHERE user_id = sqlc.arg('user-id') - 1;
