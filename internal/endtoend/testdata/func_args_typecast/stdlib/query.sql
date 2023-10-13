-- name: PlusPositionalCast :one
SELECT plus($1, $2::INTEGER);
