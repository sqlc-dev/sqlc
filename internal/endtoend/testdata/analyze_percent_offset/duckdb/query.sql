-- name: Page :many
SELECT id FROM readings ORDER BY id LIMIT 50% OFFSET $1;
